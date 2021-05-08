// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/8

package dynamic

type ApiService struct {
	Value map[string]*ApiInfo
}

type ApiInfo struct {
	Method string
	Path   string
}

// 路由树
