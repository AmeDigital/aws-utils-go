# This dockerfile creates a docker image having the aws-utils-go lib copied to /go/src/stash.b2w/asp/aws-utils-go.git

FROM golang:1.12.9

RUN mkdir -p /go/src/stash.b2w/asp/aws-utils-go.git

COPY ./ /go/src/stash.b2w/asp/aws-utils-go.git/
