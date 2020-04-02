# ********************************************************************
# * Dockerfile                                                       *
# *                                                                  *
# * 2020-03-16 First Version, JR                                     *
# * 2020-03-17 Change build instructions to make rule, JR            *
# *                                                                  *
# * File with Instructions to build the authorizer docker image.     *
# *                                                                  *
# * Usage:                                                           *
# * $ docker image build -t authorizer:go .                          *
# ********************************************************************

FROM golang

RUN mkdir /go/src/authorizer
COPY . /go/src/authorizer
WORKDIR /go/src/authorizer/
RUN make build
