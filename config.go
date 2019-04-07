package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/kelseyhightower/confd/log"
	"gopkg.in/yaml.v2"
)

// A Config structure is used to configure watchman.
type Config struct {
	Interval            int    `yaml:"interval"`
	LogLevel            string `yaml:"log-level"`
	PrintVersion        bool
	ConfigFile          string
	ConfDir             string
	ServiceTemplateFile string
	OutputDir           string
}

var config Config

func init() {
	flag.StringVar(&config.ConfDir, "conf-dir", "/etc/watchman", "watchman conf directory")
	flag.StringVar(&config.ConfigFile, "config-file", "/etc/watchman/watchman.yaml", "the watchman config file")
	flag.IntVar(&config.Interval, "interval", 60, "backend polling interval")
	flag.StringVar(&config.LogLevel, "log-level", "", "level which watchman should log messages")
	flag.BoolVar(&config.PrintVersion, "version", false, "print version and exit")
	flag.StringVar(&config.ServiceTemplateFile, "service-template-file", "/etc/watchman/nginx/service.tmpl", "the NGINX service template file")
	flag.StringVar(&config.OutputDir, "output-dir", "/etc/nginx/conf.d", "the NGINX configuration directory")
}

// initConfig initializes the watchman configuration by first setting defaults,
// then overriding settings from the watchman config file, then overriding,
// and finally overriding settings from flags set on the command line.
// It returns an error if any.
func initConfig() error {
	_, err := os.Stat(config.ConfigFile)
	if os.IsNotExist(err) {
		log.Debug("Skipping watchman config file.")
	} else {
		log.Debug("Loading " + config.ConfigFile)
		configBytes, err := ioutil.ReadFile(config.ConfigFile)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(configBytes, &config)
		if err != nil {
			return err
		}
	}

	if config.LogLevel != "" {
		log.SetLevel(config.LogLevel)
	}

	return nil
}
