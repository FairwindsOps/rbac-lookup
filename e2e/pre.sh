#!/bin/bash

set -e

CGO_ENABLED=0 go build -ldflags "-s -w" -o rbac-lookup

docker cp ./ e2e-command-runner:/rbac-lookup
