#!/bin/bash -eux

export BINPATH=$(pwd)/build
export GOMODPATH=$(pwd)

pushd $GOMODPATH/dp-frontend-router
  BINPATH=$BINPATH make build
  cp Dockerfile.concourse $BINPATH
popd
