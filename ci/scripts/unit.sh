#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-frontend-router
  make test
popd
