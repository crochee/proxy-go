@echo off
go build -trimpath -ldflags="-s -w" ./cmd/proxy