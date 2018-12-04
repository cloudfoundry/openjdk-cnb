#!/usr/bin/env bash

set -euo pipefail

if [[ -d $PWD/go-module-cache && ! -d ${GOPATH}/pkg/mod ]]; then
  mkdir -p ${GOPATH}/pkg
  ln -s $PWD/go-module-cache ${GOPATH}/pkg/mod
fi

go build -ldflags='-s -w' -o bin/build build/main.go
go build -ldflags='-s -w' -o bin/detect detect/main.go
