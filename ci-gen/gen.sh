#!bash

set -euo pipefail

##
# since this script is meant to run for the user inside the container, there's
# a bit of hackery we have to do around the fact we're changing to the host
# user's UID. This creates a temporary home dir for storing the module cache
# and a gopath for storing go downloads.
##

if [ ! -w "${GOPATH}" ]
then
  mkdir -p /tmp/gopath /tmp/home
  export HOME=/tmp/home
  export GOPATH=/tmp/gopath:${GOPATH}
fi

echo $GOPATH

oapi-codegen -package uisvc -o ${PWD}/ci-gen/openapi/services/uisvc/uisvc.gen.go ${PWD}/ci-gen/openapi/spec.yaml

protoc -I/usr/include:/go/src /go/src/github.com/tinyci/ci-agents/ci-gen/grpc/types/*.proto --go_out=plugins=grpc:/go/src

for i in $(find ci-gen/grpc/services -maxdepth 1 -type d -name '*' | tail -n +2)
do 
  SPEC=$(basename $i .proto)
  protoc -I/usr/include:/go/src ${PWD}/ci-gen/grpc/services/${SPEC}/server.proto --go_out=plugins=grpc:/go/src
  protoc -I/usr/include:${PWD}/ci-gen:/go/src ${PWD}/ci-gen/grpc/services/${SPEC}/server.proto --go_out=plugins=grpc:/go/src
done
