.SILENT:

APP_VERSION=0.1

include  ./Makefile.dist

build-dev:
	@docker-compose build

run-dev:
	@docker-compose up

sh:
	@docker-compose exec app ash

build-prod:
	@docker build \
	    -t goprince:v$(APP_VERSION) . \
	    --build-arg app_env=production

run-prod:
	@docker run -it -d \
	    --name goprince_1 \
	    --net=host \
	    goprince:v$(APP_VERSION)
