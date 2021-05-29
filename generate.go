// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/25

package proxy_go

//go:generate go install github.com/securego/gosec/v2/cmd/gosec@v2.7.0
//go:generate gosec -fmt=json -out=results.json .\...
//go:generate go test -cover .\...
