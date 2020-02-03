#!/bin/bash -eux

export GOMODPATH=$(pwd)

pushd $GOMODPATH/dp-frontend-router
  make test
popd
