#!/usr/bin/env bash

set -euo pipefail

if [[ -d $PWD/go-module-cache && ! -d $GOPATH/pkg/mod ]]; then
  mkdir -p $GOPATH/pkg
  ln -s $PWD/go-module-cache $GOPATH/pkg/mod
fi

cd openjdk-buildpack
go test ./...
