# Developer Guide

## Prerequisites

What things you need to install the software and how to install them

* [Go 1.10+](https://golang.org/download)
* [Docker](https://docs.docker.com/install/)
* [Vagrant](https://www.vagrantup.com/downloads.html)

## Set up a development environment
TODO: write Vagrantfile

## Build

```
git clone https://github.com/ttousai/watchman.git
cd watchman 
make
```

You should now have watchman in your `bin/` directory:

```
ls bin/
watchman
```

## Running the tests
```
make test
```
