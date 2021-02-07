// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/6

package dynamic

type BalanceNode struct {
	ServiceName string            `json:"service_name"`
	Add         bool              `json:"add"`
	Scheme      string            `json:"scheme"`
	Host        string            `json:"host"`
	Metadata    map[string]string `json:"metadata"`
	Weight      float64           `json:"weight"`
}
