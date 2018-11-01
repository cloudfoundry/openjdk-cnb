#!/usr/bin/env bash

set -euo pipefail

if [[ -d $PWD/go-module-cache && ! -d $GOPATH/pkg/mod ]]; then
  mkdir -p $GOPATH/pkg
  ln -s $PWD/go-module-cache $GOPATH/pkg/mod
fi

OUTPUT="$PWD/artifactory"

cd openjdk-buildpack
go build -i -ldflags='-s -w' -o bin/package package/main.go
bin/package $OUTPUT
