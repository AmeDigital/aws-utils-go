# This dockerfile creates a docker image having the aws-utils-go lib copied to /go/src/stash.b2w/asp/aws-utils-go.git
# it also packages aws-sdk-go and aws-lambda-go

FROM golang:1.12.9

RUN mkdir -p /go/src/github.com/aws/aws-lambda-go
COPY ./aws-lambda-go/ /go/src/github.com/aws/aws-lambda-go/

RUN mkdir -p /go/src/github.com/aws/aws-sdk-go 
COPY ./aws-sdk-go/ /go/src/github.com/aws/aws-sdk-go/
