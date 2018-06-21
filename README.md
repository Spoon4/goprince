# goprince

REST API in Go to use [Prince][prince].
## How to use?

### In Development

Running our application in development is a very easy process thanks to the `docker-compose` command. 
To build our application's docker image we run:
```bash
$ make build-dev
```

And to launch a docker container, we run:
```bash
$ make run-dev
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

[prince]:           http://www.princexml.com
[tarkan-article]:   https://www.surenderthakran.com/articles/tech/dockerized-development-and-production-environment-golang
[gin-doc]:          https://github.com/gin-gonic/gin/blob/master/README.md
