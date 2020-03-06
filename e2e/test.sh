#!/bin/bash

set -e

kubectl create -f deploy/

./rbac-lookup  e2e-test |grep -v "No RBAC Bindings found"
