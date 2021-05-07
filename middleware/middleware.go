// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/6

package middleware

import (
	"net/http"
)

type Handler interface {
	NameSpace() string
	http.Handler
}
