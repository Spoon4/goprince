FROM golang:1.10-alpine3.7 AS build

ARG app_env
ENV APP_ENV $app_env

RUN apk add --update --no-cache \
    make musl git g++ \
    && rm -rf /var/cache/apk/*

ENV GOPATH /go

# Add Go dependencies
RUN go get github.com/gin-gonic/gin
RUN go get github.com/stretchr/testify/assert

ADD ./ /$GOPATH
WORKDIR $GOPATH

RUN go install goprince

FROM alpine:3.7

RUN apk add --update --no-cache \
    libxml2 pixman tiff giflib libpng lcms2 libjpeg-turbo libcurl libgomp \
    msttcorefonts-installer fontconfig freetype \
    && rm -rf /var/cache/apk/* \
    && update-ms-fonts

# Install Prince
ENV PRINCE_VERSION 12

RUN mkdir /src
RUN wget -qO- https://www.princexml.com/download/prince-${PRINCE_VERSION}-alpine3.7-x86_64.tar.gz \
        | tar xvz --strip-components=1 -C /src \
    && printf "/usr\n" | ./src/install.sh \
    && rm -Rf src

RUN mkdir -p /public /var/log/goprince
VOLUME /public
VOLUME /var/log/goprince

WORKDIR /
COPY --from=build /go/bin/goprince /usr/local/bin

ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD ["goprince"]