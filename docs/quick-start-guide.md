# Quick Start Guide
`watchman` is intended to be run in a Docker Swarm, so this guide will focus
on running `watchman` in Swarm mode using the watchman Docker image.

As of this writing watchman must be scheduled on the Swarm manager instance
to run correctly.

## Create overlay network
```shell
docker network create --driver overlay shared-net
```

## Create watchman service
```shell
docker service create --name nginx \
  --mount type=bind,source=/var/run/docker.sock,destination=/var/run/docker.sock \
  --constraint 'node.role==manager' \
  --network shared-net --publish 80:80 \
  ttousai/watchman:v0.1.0-alpha.1

docker service ls
ID                  NAME                  MODE                REPLICAS            IMAGE                             PORTS
tnuebu0w6dfu        nginx                 replicated          1/1                 ttousai/watchman:v0.1.0-alpha.1   *:80->80/tcp
```

Creates a service called nginx, running watchman in the default configuration.
Replace `watchman.service.url=guestbook.example.com` with a domain you control.

## Publish service
```shell
docker service create --name guestbook \
  --label watchman.service.enable=true \
  --label watchman.service.port=80 \
  --label watchman.service.url=guestbook.example.com \
  --network shared-net \
  gcr.io/google-samples/gb-frontend:v4

docker service ls                                                                                                                     
ID                  NAME                  MODE                REPLICAS            IMAGE                                  PORTS                                                     
zajwgv6zghdw        guestbook             replicated          1/1                 gcr.io/google-samples/gb-frontend:v4                                                             
tnuebu0w6dfu        nginx                 replicated          1/1                 ttousai/watchman:v0.1.0-alpha.1        *:80->80/tcp 
```

## Verify
Verify by visiting http://guestbook.example.com in your browser.

Read the [configuration guide](docs/configuration-guide.md) for details on how to configure watchman.
