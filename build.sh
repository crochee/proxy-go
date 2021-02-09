#!bin/bash
set -ex

go build -tags=jsoniter ./cmd/proxy
