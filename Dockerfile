FROM golang:1.14

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/mirobertod/drone-cloudformation