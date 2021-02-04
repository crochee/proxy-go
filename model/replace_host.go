// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package model

type ReplaceHost struct {
	Name string `json:"name" binding:"required"`
	*Host
}

type Host struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host" binding:"required"`
}
