package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/kelseyhightower/confd/log"
)

func main() {
	flag.Parse()
	if config.PrintVersion {
		fmt.Printf("watchman %s (Git SHA: %s, Go Version: %s)\n", Version, GitSHA, runtime.Version())
		os.Exit(0)
	}

	if err := initConfig(); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Starting watchman")

	cli, err := newDockerClient()
	if err != nil {
		log.Fatal(err.Error())
	}

	interval := config.Interval

	stopChan := make(chan bool)
	doneChan := make(chan bool)
	errChan := make(chan error, 10)

	processor := newIntervalProcessor(cli, config, stopChan, doneChan, errChan, interval)

	// TODO: worry about NGINX server startup and shutdown, is this the best approach?
	log.Info("Starting NGINX server")
	if err = processor.startServer(); err != nil {
		log.Fatal(err.Error())
	}

	go processor.process()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case err := <-errChan:
			log.Error(err.Error())
		case s := <-signalChan:
			log.Info(fmt.Sprintf("Captured %v. Exiting...", s))
			close(doneChan)
		case <-doneChan:
			// TODO: worry about NGINX server startup and shutdown
			os.Exit(0)
		}
	}
}
