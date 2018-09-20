#!/usr/bin/env bash

set -euo pipefail

GOCACHE="$PWD/go-build"

go build -i -ldflags='-s -w' -o bin/build build/main.go
go build -i -ldflags='-s -w' -o bin/detect detect/main.go