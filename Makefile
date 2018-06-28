.SILENT:

APP_VERSION=`cat VERSION`

include  ./Makefile.dist

version:
	@echo $(APP_VERSION)

build-dev:
	@docker-compose build app

run-dev:
	@docker-compose up app

sh:
	@docker-compose exec app ash

build-prod:
	@docker build \
	    -t goprince:$(APP_VERSION) \
	    --build-arg app_env=production \
	    --no-cache .

run-prod:
	@docker run -it -d \
	    --name goprince_1 \
	    --net=host \
	    goprince:$(APP_VERSION)
