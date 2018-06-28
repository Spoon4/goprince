FROM golang:1.10-alpine3.7

ARG app_env
ENV APP_ENV $app_env

RUN apk add --update --no-cache \
    make musl git g++ \
    libxml2 pixman tiff giflib libpng lcms2 libjpeg-turbo libcurl libgomp \
    msttcorefonts-installer fontconfig freetype \
    && rm -rf /var/cache/apk/* \
    && update-ms-fonts

# Install Prince
ENV PRINCE_VERSION 12

RUN wget -qO- https://www.princexml.com/download/prince-${PRINCE_VERSION}-alpine3.7-x86_64.tar.gz \
        | tar xvz --strip-components=1 \
    && printf "/usr\n" | ./install.sh \
    && rm -Rf *

RUN mkdir -p /public
VOLUME /public

ENV GOPATH /go
ENV PATH "$PATH:$GOPATH/bin"

RUN mkdir -p $(GOPATH)/bin $(GOPATH)/pkg

# Add Go dependencies
RUN go get github.com/gin-gonic/gin
RUN go get github.com/stretchr/testify/assert
RUN go get github.com/cespare/reflex

ADD ./ /$GOPATH
WORKDIR $GOPATH

RUN go get ./...
RUN go install goprince
