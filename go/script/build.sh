#!bin/bash
set -ex
go build -trimpath -ldflags="-s -w" ./cmd/proxy
