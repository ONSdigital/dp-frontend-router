#!/bin/bash -eux

export BINPATH=$(pwd)/build
export GOPATH=$(pwd)/go

pushd $GOPATH/src/github.com/ONSdigital/dp-frontend-router
  BINPATH=$BINPATH make build
  cp Dockerfile.concourse $BINPATH
popd
