# Watchman
[![Go Report Card](https://goreportcard.com/badge/github.com/ttousai/watchman)](https://goreportcard.com/report/github.com/ttousai/watchman)

`watchman` is a lightweight NGINX configuration management tool, heavily inspired by [confd](https://github.com/kelseyhightower/confd) focused on:

* starting NGINX service
* keeping NGINX service configuration files up-to-date using [Docker](https://docs.docker.com/) Swarm service labels
* reloading NGINX to pick up new config file changes

## Building

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

## Getting Started

Before we begin be sure to [download and install watchman](docs/installation.md).

* [quick start guide](docs/quick-start-guide.md)

## Next steps

Check out the [docs directory](docs) for more docs.
