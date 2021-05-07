@echo off
golangci-lint run -c ./build/ci/golangci-lint.yml --tests=false  --out-format=json  > golangci-lint.json 2>&1