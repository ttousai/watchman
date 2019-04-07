package main

import (
	"context"

	docker "docker.io/go-docker"
	"docker.io/go-docker/api/types"
)

type service struct {
	ServiceName string
	ServicePort string
	ServiceURL  string
}

// Client provides a wrapper around the Docker client
type dockerClient struct {
	client *docker.Client
}

// NewDockerClient returns a new client to for the given address
func newDockerClient() (*dockerClient, error) {
	// TODO: works only when there is direct access to /var/run/docker.sock
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}
	return &dockerClient{cli}, nil
}

func (c *dockerClient) getServices() ([]service, error) {
	svcOpts := types.ServiceListOptions{}

	ctx := context.Background()
	swarmSvcs, err := c.client.ServiceList(ctx, svcOpts)
	if err != nil {
		return nil, err
	}

	services := make([]service, 0)
	for _, svc := range swarmSvcs {
		lbls := svc.Spec.Annotations.Labels
		if lbls["watchman.service.enable"] == "true" {
			s := service{}
			s.ServicePort = lbls["watchman.service.port"]
			s.ServiceURL = lbls["watchman.service.url"]
			s.ServiceName = svc.Spec.Annotations.Name
			services = append(services, s)
		}
	}

	return services, nil
}
