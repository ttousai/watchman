package main

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/kelseyhightower/confd/log"
)

type intervalProcessor struct {
	client   *dockerClient
	config   Config
	stopChan chan bool
	doneChan chan bool
	errChan  chan error
	interval int
}

var serviceMap = make(map[string]service)

func newIntervalProcessor(client *dockerClient, config Config, stopChan, doneChan chan bool, errChan chan error, interval int) *intervalProcessor {
	return &intervalProcessor{client, config, stopChan, doneChan, errChan, interval}
}

func (p *intervalProcessor) process() {
	defer close(p.doneChan)
	for {
		log.Info("Generating NGINX service configs")
		restart, err := p.generateServiceConfig(p.config)
		if err != nil {
			p.errChan <- err
		}

		if restart {
			log.Info("Reloading NGINX configuration")
			// TODO: need to handle case of restart failing.
			// conf change -> restart -> fail restart. leads to
			// no change -> no restart -> new services above don't get loaded.
			err = p.reloadServerConfig()
			if err != nil {
				p.errChan <- err
			}
		}

		select {
		case <-p.stopChan:
			break
		case <-time.After(time.Duration(p.interval) * time.Second):
			continue
		}
	}
}

func (p *intervalProcessor) generateServiceConfig(config Config) (bool, error) {
	src := config.ServiceTemplateFile
	restart := false

	services, err := p.client.getServices()
	if err != nil {
		return restart, err
	}

	for _, s := range services {
		if _, ok := serviceMap[s.ServiceName]; !ok {
			log.Debug("processing %s", s.ServiceName)
			serviceMap[s.ServiceName] = s

			dst := fmt.Sprintf("%s/%s.conf", config.OutputDir, s.ServiceName)
			b, err := os.Create(dst)
			if err != nil {
				return restart, err
			}

			t := template.Must(template.ParseFiles(src))
			err = t.Execute(b, s)
			if err != nil {
				return restart, err
			}

			restart = true
		} else {
			log.Debug("skipping %s", s.ServiceName)
		}
	}

	return restart, nil
}

// startServer runs NGINX command and does not wait for completion.
func (p *intervalProcessor) startServer() error {
	cmd := exec.Command("nginx")
	err := cmd.Start()
	return err
}

func (p *intervalProcessor) reloadServerConfig() error {
	cmd := exec.Command("nginx", "-s", "reload")
	err := cmd.Run()
	return err
}
