package proxy_go

// 火焰图
//go:generate go get -u github.com/google/pprof
//go:generate pprof -http=:8080 cpu.prof

//go:generate go install github.com/securego/gosec/v2/cmd/gosec@v2.7.0
//go:generate gosec -fmt=json -out=results.json .\...
//go:generate go test -cover .\...
