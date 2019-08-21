# This dockerfile creates a docker image having the aws-utils-go lib copied to /go/src/stash.b2w/asp/aws-utils-go.git
# there is also a symbolic link from /go/src/stash.b2w/asp/aws-utils-go.git to /go/src/aws-utils-go 
# to let aws-utils-go internal files be able to build with each other.

FROM golang:1.12.9

RUN mkdir -p /go/src/stash.b2w/asp/aws-utils-go.git && ln -s /go/src/stash.b2w/asp/aws-utils-go.git /go/src/aws-utils-go

COPY ./ /go/src/stash.b2w/asp/aws-utils-go.git/
