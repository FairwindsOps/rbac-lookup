#!/bin/bash

set -e

if [ -z "$CI_SHA1" ]; then
    echo "CI_SHA1 not set. Something is wrong"
    exit 1
else
    echo "CI_SHA1: $CI_SHA1"
fi

curl -O https://storage.googleapis.com/golang/go1.12.9.linux-amd64.tar.gz
tar -xvf go1.12.9.linux-amd64.tar.gz
chown -R root:root ./go
mv go /usr/local


/usr/local/go/bin/go build -o ./rbac-lookup
docker cp ./rbac-lookup e2e-command-runner:rbac-lookup
