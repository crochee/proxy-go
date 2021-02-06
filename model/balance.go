// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/6

package model

type Node struct {
	Scheme   string            `json:"scheme"`
	Host     string            `json:"host"`
	Metadata map[string]string `json:"metadata"`
}
