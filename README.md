# goprince

[![Docker Layers](https://images.microbadger.com/badges/image/spoon4/goprince.svg)][microbadger]
[![Docker Build Status](https://img.shields.io/docker/build/spoon4/goprince.svg)][dockerstore]

REST API in Go to use [Prince][prince].

## How to use?

### Routes

#### `POST` /generate/{filename}

##### Parameters

* `filename` _string_: Name of the output PDF file

##### Body

* `input_file` _file_: HTML file to convert (**required**)
* `stylesheet` _file_: CSS file to upload to pass to Prince

##### Query parameters

* `output` _string_ (optional): 
    * `stream`: returns bytes of the file
    * `file`: serves PDF file to download
    * _not present_: returns file output path

### Env vars

|Var|Description|
|---|---|
|`APP_ENV`| `dev` or `production`|
|`LICENSE_FILE`|Path to Prince license file|
|`LICENSE_FILE`|Prince license hash key|

## Deployment

### In Development environment
 
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

To access the conteneur with shell
```bash
$ make sh
```

### In Production environment

We can build a docker image for our go application by running:
```bash
$ make build-prod
```

And to launch a docker container for the image created above we run:
```bash
$ make run-prod
```

## Logs

Both Prince and Gin are logged but in separated files.
Logs are written in the container, in `/var/log/goprince` folder.

In **development** environment, the Docker log directory is mapped on `logs/` folder by a volume.

## References

* [Dockerized Development and Production Environment For Go (GoLang)][tarkan-article]
* [Gin documentation][gin-doc]
* [How to debug Golang applications inside Docker containers using Delve][go-remote-debug]

[microbadger]:      https://microbadger.com/images/ardeveloppement/node
[dockerstore]:      https://store.docker.com/community/images/ardeveloppement/node
[prince]:           http://www.princexml.com
[tarkan-article]:   https://www.surenderthakran.com/articles/tech/dockerized-development-and-production-environment-golang
[gin-doc]:          https://github.com/gin-gonic/gin/blob/master/README.md
[go-remote-debug]:  https://mikemadisonweb.github.io/2018/06/14/go-remote-debug/