#!/usr/bin/env bash

set -euo pipefail

if [[ -d $PWD/go-module-cache && ! -d ${GOPATH}/pkg/mod ]]; then
  mkdir -p ${GOPATH}/pkg
  ln -s $PWD/go-module-cache ${GOPATH}/pkg/mod
fi

GOOS="linux" go build -ldflags='-s -w' -o bin/build build/main.go
GOOS="linux" go build -ldflags='-s -w' -o bin/detect detect/main.go
GOOS="linux" go build -ldflags='-s -w' -o bin/class-counter memcalc/helper/main.go
GOOS="linux" go build -ldflags='-s -w' -o bin/link-local-dns dns/helper/main.go
GOOS="linux" go build -ldflags='-s -w' -o bin/security-provider-configurer provider/helper/main.go
