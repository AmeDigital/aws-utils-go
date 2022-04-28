#!/bin/bash

function error () {
    message=$1
    echo "ERROR: $message"
    exit 1
}

env=$1

folders=(cognitoutils dynamodbutils localstack s3utils sessionutils snsutils sqsutils)

echo "Installing aws-utils-go..."

for folder in "${folders[@]}"; do
    cd "$folder" || error "$folder does not exists!"
    echo -n "Installing $folder..."
    go get -d -v || error "$folder installation failed!"
    if [[ $env == "local" ]]; then
        if [[ $folder == "localstack" ]] || [[ $folder == "dynamodbutils" ]]; then
            echo -n " Testing $folder..."
            go test || error "$folder tests failed!"
        else
            echo "👌"
        fi
    else
        echo "👌"
    fi
    cd ..
done

echo "aws-utils-go successfully installed"
