// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/6

package model

import "net/http"

type NamedHandler struct {
	http.Handler
	*Host
	Weight   float64
	Deadline float64
}
