#!/bin/bash

set -e


printf "\n\n"
echo "***************************"
echo "** Install and Run Venom **"
echo "***************************"
printf "\n\n"

curl -LO https://github.com/ovh/venom/releases/download/v1.1.0/venom.linux-amd64
mv venom.linux-amd64 /usr/local/bin/venom
chmod +x /usr/local/bin/venom

mv /rbac-lookup/rbac-lookup /rbac-lookup/e2e/tests

cd /rbac-lookup/e2e
mkdir -p /tmp/test-results
venom run tests/* --output-dir=/tmp/test-results
exit $?
