#!/usr/bin/env bash

set -euo pipefail

if [[ -d $PWD/go-cache ]]; then
  export GOPATH=$PWD/go-cache
fi

GOOS="linux" go build -ldflags='-s -w' -o bin/build cmd/build/main.go
GOOS="linux" go build -ldflags='-s -w' -o bin/class-counter cmd/class-counter/main.go
GOOS="linux" go build -ldflags='-s -w' -o bin/detect cmd/detect/main.go
GOOS="linux" go build -ldflags='-s -w' -o bin/link-local-dns cmd/link-local-dns/main.go
GOOS="linux" go build -ldflags='-s -w' -o bin/security-provider-configurer cmd/security-provider-configurer/main.go
