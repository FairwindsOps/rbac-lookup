#!/bin/bash

set -e

if [ -z "$CI_SHA1" ]; then
    echo "CI_SHA1 not set. Something is wrong"
    exit 1
else
    echo "CI_SHA1: $CI_SHA1"
fi

apk add go=1.12.12-r0
/usr/bin/go build -o ./rbac-lookup
docker cp ./rbac-lookup e2e-command-runner:rbac-lookup
