#!bin/bash
set -ex
go build -trimpath -ldflags="-s -w" -o proxy ./cmd/proxy
