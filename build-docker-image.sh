#!/bin/bash
# This script builds a docker image and publishes it to "registry.b2w.io/b2wbuild"

function die {
    declare MSG="$@"
    echo -e "$0: Error: $MSG">&2
    exit 1
}

function tag_exists {
    declare tag="$1"
    curl http://registry.b2w.io/repository/docker-private/v2/b2wbuild/golang-aws-utils-go/tags/list 2>/dev/null| jq -r ".tags" | grep -q $tag
}

TAG="$1"

[ -z "$TAG" ] && die "Parameter 'TAGNAME' cannot be empty. You must specify a tag name to be associated to the docker image that will be created.\nExample:\n$0 '1.0.0'."

tag_exists $TAG && die "TAG $TAG already exists in the repository"


IMAGE_NAME="registry.b2w.io/b2wbuild/golang-aws-utils-go:${TAG}"

SEPARATOR="#######################################################################################"

echo $SEPARATOR
echo "creating image $IMAGE_NAME"

docker build -t $IMAGE_NAME . || die "failed to build docker image $IMAGE_NAME"

echo $SEPARATOR
echo "build completed successfully."

echo $SEPARATOR
echo "publishing image $IMAGE_NAME to the repository..."

docker push $IMAGE_NAME || die "failed to publish image to the repository."

echo $SEPARATOR
echo "Image $IMAGE_NAME was published successfully."
