#!bin/bash
set -ex
go build -tags jsoniter -o proxy ./cmd/proxy
