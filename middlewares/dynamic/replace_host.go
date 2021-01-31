// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package dynamic

type ReplaceHost struct {
	Name string
	*Host
}

type Host struct {
	Scheme string
	Host   string
}
