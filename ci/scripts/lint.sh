#!/bin/bash -eux

go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0

pushd dp-frontend-router
  make lint
popd
