#!/bin/bash

set -e

go build -ldflags "-s -w" -o rbac-lookup

docker cp rbac-lookup e2e-command-runner:/rbac-lookup

docker cp e2e/deploy e2e-command-runner:/
