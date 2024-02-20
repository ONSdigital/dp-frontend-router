#!/bin/bash -eux

export BINPATH=$(pwd)/build
export cwd=$(pwd)

pushd $cwd/dp-frontend-router
  BINPATH=$BINPATH make build
  cp Dockerfile.concourse $BINPATH
popd
