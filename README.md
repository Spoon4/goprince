# goprince

[![Docker Layers](https://images.microbadger.com/badges/image/spoon4/goprince.svg)][microbadger]
[![Docker Build Status](https://img.shields.io/docker/build/spoon4/goprince.svg)][dockerstore]

Convert HTML files to PDF with [Prince][prince] through an API.

## How to use?

```bash
$ goprince [--help] [--log-dir=/my/directory/] [--stdout] [--port 8080]
```

### Command line options

* `--help`: show help
* `--log-dir`: set Gin and Prince log directory. Default `/var/log/goprince`
* `--stdout`: if given, Gin and Prince also log on stdout
* `--port`: set port listening on. Default `8080`

## API

### `POST` /generate/{filename}

#### Parameters

* `filename` _string_: Name of the output PDF file

#### Body

* `input_file` _file_: HTML file to convert (**required**)
* `stylesheet` _file_: CSS file to upload to pass to Prince

#### Query parameters

* `output` _string_ (optional): 
    * `stream`: returns bytes of the file
    * `file`: serves PDF file to download
    * _not present_: returns file output path

## Env vars

|Var|Description|
|---|---|
|`APP_ENV`| `dev` / `production`|
|`LICENSE_FILE`|Path to Prince license file|
|`LICENSE_KEY`|Prince license hash key|

## Logs

Both Prince and Gin are logged in separated files.
Logs are written in the container, in `/var/log/goprince` folder.

In **development** environment, the Docker log directory is mapped on `logs/` folder by a volume in [docker-compose.yml](docker-compose.yml).

## Deployment

PHPStorm doesn't handle Docker packages..so you have to install all packages on your host system :vomit: ...
Use Golang [applications structure][app-structure] to manage as well a as possible package requirements.

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

You can build a Docker image for your go application by running:

```bash
$ make build-prod
```
:warning: You need to increment number in [VERSION](version) file to increment Docker's image tag.

And to launch a docker container for the image created above we run:
```bash
$ make run-prod
```


## References

* [Dockerized Development and Production Environment For Go (GoLang)][tarkan-article]
* [Gin documentation][gin-doc]
* [How to debug Golang applications inside Docker containers using Delve][go-remote-debug]
* [Golang: Gracefully stop application][kpbird-graceful]

[microbadger]:      https://microbadger.com/images/spoon4/goprince
[dockerstore]:      https://store.docker.com/community/images/spoon4/goprince
[prince]:           http://www.princexml.com
[app-structure]:    https://golang.org/doc/code.html
[tarkan-article]:   https://www.surenderthakran.com/articles/tech/dockerized-development-and-production-environment-golang
[gin-doc]:          https://github.com/gin-gonic/gin/blob/master/README.md
[go-remote-debug]:  https://mikemadisonweb.github.io/2018/06/14/go-remote-debug/
[kpbird-graceful]:  https://medium.com/@kpbird/golang-gracefully-stop-application-23c2390bb212
