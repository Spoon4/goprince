FROM spoon4/goprince:0.1

RUN apk add --update --no-cache \
    ca-certificates \
    dpkg \
    musl-dev \
    libc6-compat \
    && rm -rf /var/cache/apk/* \
    && update-ms-fonts

ENV PATH "$PATH:/usr/local/go/bin"

# Add Delve
RUN go get -u github.com/derekparker/delve/cmd/dlv

WORKDIR $GOPATH/src/goprince

CMD ["dlv", "debug", "--headless", "--listen=:2345", "--log=true", "--api-version=2", "goprince"]
