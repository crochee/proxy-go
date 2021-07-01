package proxy_go

//go:generate go install github.com/swaggo/swag/cmd/swag
//go:generate swag i -g pkg/router/api.go

//go:generate go install github.com/securego/gosec/v2/cmd/gosec@v2.7.0
//go:generate gosec -fmt=json -out=results.json .\...

//go:generate go test -cover .\...

// 火焰图
//go:generate go get -u github.com/google/pprof
//go:generate pprof -http=:8080 cpu.prof

// mock
//go:generate go get -u github.com/golang/mock/mockgen
