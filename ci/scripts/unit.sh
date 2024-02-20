#!/bin/bash -eux

export cwd=$(pwd)

# XXX go install github.com/kevinburke/go-bindata/v4/...@v4.0.2
pushd $cwd/dp-frontend-router
  make test
popd
