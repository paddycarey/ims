#!/bin/sh
set -e

# install dependencies
apk add --update git libjpeg-turbo-utils
apk add --update-cache --repository http://dl-3.alpinelinux.org/alpine/edge/main/ go
apk add --update-cache --repository http://dl-3.alpinelinux.org/alpine/edge/testing/ gifsicle

# copy source into temporary GOPATH
export GOPATH=/go
mkdir -p /go/src/github.com/paddycarey
cp -r /src /go/src/github.com/paddycarey/ims
cd /go/src/github.com/paddycarey/ims

# build ims inside the container
go get
go build -o /bin/ims

# remove build-only dependencies and clean up build artifacts
apk del go git wget
rm -rf /go
rm -rf /var/cache/apk/*
