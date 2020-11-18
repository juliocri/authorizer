# ********************************************************************
# * Dockerfile                                                       *
# *                                                                  *
# * 2020-03-16 First Version, JR                                     *
# * 2020-03-17 Change build instructions to make rule, JR            *
# * 2020-11-18 Changes golang to lightweight alpine version, JR      *
# *                                                                  *
# * File with Instructions to build the authorizer docker image.     *
# *                                                                  *
# * Usage:                                                           *
# * $ docker image build -t authorizer:go .                          *
# ********************************************************************

FROM golang:alpine

RUN apk update && \
    apk --no-cache add build-base git make && \
    rm -rf /var/cache/apk/*

RUN mkdir /go/src/authorizer

COPY . /go/src/authorizer
WORKDIR /go/src/authorizer/

RUN make build

ENTRYPOINT ["/go/bin/authorizer"]
