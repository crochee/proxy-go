// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/30

package internal

import (
	"io"
	"net/http"
	"net/http/httptest"
)

// PerformRequest unit test function httptest.ResponseRecorder
func PerformRequest(r http.Handler, method, path string, body io.Reader,
	headers http.Header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	req.Header = headers
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
