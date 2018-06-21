# goprince

[![Docker Layers](https://images.microbadger.com/badges/image/spoon4/goprince.svg)][microbadger]
[![Docker Build Status](https://img.shields.io/docker/build/spoon4/goprince.svg)][dockerstore]

REST API in Go to use [Prince][prince].

## How to use?

### In Development
 
To build the application's docker image run:
```bash
$ make build-dev
```

To launch a docker container run:
```bash
$ make run-dev
```

To access the API browse at
```text
http://localhost:8080/
```

### In Production

We can build a docker image for our go application by running:
```bash
$ make build-prod
```

And to launch a docker container for the image created above we run:
```bash
$ make run-prod
```

## References

* [Dockerized Development and Production Environment For Go (GoLang)][tarkan-article]
* [Gin documentation][gin-doc]

[microbadger]:      https://microbadger.com/images/ardeveloppement/node
[dockerstore]:      https://store.docker.com/community/images/ardeveloppement/node
[prince]:           http://www.princexml.com
[tarkan-article]:   https://www.surenderthakran.com/articles/tech/dockerized-development-and-production-environment-golang
[gin-doc]:          https://github.com/gin-gonic/gin/blob/master/README.md
