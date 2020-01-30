#!/bin/bash -eux

export GOMODPATH=$(pwd)

pushd $GOMODPATH/dp-frontend-router
  make generate-build test
popd
