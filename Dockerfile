# This dockerfile creates a docker image having the aws-utils-go lib copied to /go/src/github.com/AmeDigital/aws-utils-go

FROM golang:1.17.9

RUN mkdir -p /go/src/github.com/AmeDigital/aws-utils-go

COPY ./ /go/src/github.com/AmeDigital/aws-utils-go/
