#!/bin/bash -eux

export BINPATH=$(pwd)/build
export cwd=$(pwd)

go install github.com/kevinburke/go-bindata/v4/...@v4.0.2
pushd $cwd/dp-frontend-router
  BINPATH=$BINPATH make build
  cp Dockerfile.concourse $BINPATH
popd
