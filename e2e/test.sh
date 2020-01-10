#!/bin/bash



printf "\n\n"
echo "**************************"
echo "** Begin E2E Test Setup **"
echo "**************************"
printf "\n\n"

set -e

printf "\n\n"
echo "********************************************************************"
echo "** Create Test Account **"
echo "********************************************************************"
printf "\n\n"

kubectl create serviceaccount rbac-lookup -n default
kubectl create clusterrole test-rbac-lookup --verb="create" --resource=deployment
kubectl create clusterrolebinding test-rbac-lookup --clusterrole=test-rbac-lookup --serviceaccount=default:rbac-lookup

printf "\n\n"
echo "********************************************************************"
echo "** Test rbac-lookup **"
echo "********************************************************************"
printf "\n\n"
ls -al /usr/local/bin/
chmod+x /usr/local/bin/rbac-lookup
/usr/local/bin/rbac-lookup "rbac-lookup" --kind service | grep "ClusterRole/test-rbac-lookup"

if [ $? != 0 ]; then
  echo "Cluster Role not found.  Did rbac-lookup fail?"
  exit 1
fi


