FROM golang:1.10-alpine

ARG app_env
ENV APP_ENV $app_env

ENV GOPATH /workspace
ENV PATH "$PATH:$GOPATH/bin"

RUN apk add --update \
    make musl git \
    && rm -rf /var/cache/apk/*

ADD ./ /workspace
WORKDIR $GOPATH

RUN make --no-print-directory install
#CMD make --no-print-directory run

#RUN go get ./
#RUN go build -o goprince

#CMD if [ ${APP_ENV} = production ]; \
#	then \
#	app; \
#	else \
#	go get github.com/pilu/fresh && \
#	fresh; \
#	fi
	
#EXPOSE 8080

