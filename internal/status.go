// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/4

package internal

import "net/http"

// StatusClientClosedRequest non-standard HTTP status code for client disconnection.
const StatusClientClosedRequest = 499

// StatusClientClosedRequestText non-standard HTTP status for client disconnection.
const StatusClientClosedRequestText = "Client Closed Request"

func StatusText(statusCode int) string {
	if statusCode == StatusClientClosedRequest {
		return StatusClientClosedRequestText
	}
	return http.StatusText(statusCode)
}

var (
	JsonContentType      = "application/json; charset=utf-8"
	JsonAsciiContentType = "application/json"
)
