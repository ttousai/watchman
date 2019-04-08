[![Go Report Card](https://goreportcard.com/badge/github.com/ttousai/watchman)](https://goreportcard.com/report/github.com/ttousai/watchman)

# Watchman

`watchman` is a lightweight tool to generate NGINX configs from Docker Service labels.

Features:
* starting NGINX service
* keeping NGINX service configuration files up-to-date using [Docker Swarm](https://docs.docker.com/) service labels
* reloading NGINX to pick up new config file changes

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. 

Go 1.10+ is required to build watchman, which uses gomod.

```
$ git clone https://github.com/ttousai/watchman.git
$ cd watchman 
$ make
```

You should now have watchman in your `bin/` directory:

```
$ ls bin/
watchman
```

### Prerequisites

What things you need to install the software and how to install them

* [Go 1.10+](https://golang.org/download)
* [Docker](https://docs.docker.com/install/)


### Installing

A step by step series of examples that tell you how to get a development env running

`watchman` is designed to run in a Docker Swarm and the recommended installation is
via a Docker Swarm service.

You can build the Docker image locally with the commmand
```
make release
docker images | grep ttousai/watchman
```

or directly with Docker
```
docker build -f Dockerfile.alpine -t <sometag> .
docker images | grep <sometag>
```

The following example installs `watchman` in a Docker Swarm using the `docker stack deploy`
command. Please run these commands on the swarm manager.

First create a stack.yml file with the following content
```
version: "3"
services:

  api:
    image: <some-api>
    volumes:
      - <volume>:/data
    networks:
      - shared-net
    deploy:
      labels:
        watchman.service.enable: "true"
        watchman.service.port: "8080"
        watchman.service.url: "api.somedomain.com"
      placement:
        constraints:
          - node.role == worker
      replicas: 1
      restart_policy:
        condition: on-failure
          
  nginx:
    image: ttousai/watchman:<tag>
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    networks:
      - shared-net
    ports:
      - "80:80"
      - "443:443"
    deploy:
      placement:
        constraints:
          - node.role == manager
      replicas: 1
      restart_policy:
        condition: on-failure

networks:
  shared-net:
    driver: overlay
```

On the Docker Swarm manager run
```
docker stack deploy -c stack.yml
```

Check that all the services are running
```
docker service ls
```

Access the API service at http://api.somedomain.com.

## Running the tests

```
make test
```

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/ttousai/watchman/tags). 

## Authors

* **Abubakr-Sadik Nii Nai Davis** - *Initial work* - [ttousai](https://github.com/ttousai)

See also the list of [contributors](https://github.com/ttousai/watchman/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Inspired by [confd](https://github.com/kelseyhightower/confd)
* PurpleBooth [README template](https://gist.github.com/PurpleBooth/109311bb0361f32d87a2)
