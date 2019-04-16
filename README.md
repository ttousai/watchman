# watchman
[![Go Report Card](https://goreportcard.com/badge/github.com/ttousai/watchman)](https://goreportcard.com/report/github.com/ttousai/watchman)

`watchman` is a lightweight NGINX configuration management tool for Docker Swarm
focused on:

* starting NGINX
* keeping NGINX service configuration files up-to-date using [Docker](https://docs.docker.com/)
  Swarm service labels
* reloading NGINX to pick up new config file changes


## Getting Started

As of this writing watchman must be scheduled on the Swarm manager instance
to run correctly.

### Create overlay network
```
docker network create --driver overlay shared-net
```

### Create watchman service
Creates a service called nginx, running watchman in the default configuration.
Replace `watchman.service.url=guestbook.example.com` with a domain you control.

```
docker service create --name nginx \
  --mount type=bind,source=/var/run/docker.sock,destination=/var/run/docker.sock \
  --constraint 'node.role==manager' \
  --network shared-net --publish 80:80 \
  ttousai/watchman:v0.1.0-alpha.1

docker service ls
ID                  NAME                  MODE                REPLICAS            IMAGE                             PORTS
tnuebu0w6dfu        nginx                 replicated          1/1                 ttousai/watchman:v0.1.0-alpha.1   *:80->80/tcp
```

### Publish service
```
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

### Verify
Verify by visiting http://guestbook.example.com in your browser.

## Next steps
Read more about the [configuring](docs/configuration-guide.md) watchman,
and check out the [docs directory](docs) for more docs.

## Contributing
Please read [CONTRIBUTING.md](docs/contributing.md) for details on our code of conduct,
and the process for submitting pull requests to us.

## Versioning
We use [SemVer](http://semver.org/) for versioning. For the versions available,
see the [tags on this repository](https://github.com/ttousai/watchman/tags). 

## Authors
* **Abubakr-Sadik Nii Nai Davis** - [ttousai](https://github.com/ttousai)

See also the list of [contributors](https://github.com/ttousai/watchman/contributors)
who participated in this project.

## License
This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file
for details

## Acknowledgments
* Inspired by [confd](https://github.com/kelseyhightower/confd)
* PurpleBooth [README template](https://gist.github.com/PurpleBooth/109311bb0361f32d87a2)
