#!/bin/bash

set -e

go build -ldflags "-s -w" -o rbac-lookup

docker cp ./ e2e-command-runner:/rbac-lookup
