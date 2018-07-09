.SILENT:

APP_VERSION=`cat VERSION`
GOPRINCE_PRODUCTION_PORT=80

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
	docker run -it \
	    --name goprince_1 \
	    --net=host \
	    --volume $(pwd)/src/goprince:/go/src/goprince \
      	--volume $(pwd)/public:/public \
		--volume $(pwd)/logs:/var/log/goprince \
		--publish 80:80 \
	    goprince:$(APP_VERSION) \
	    sh -c "goprince --port 80"
