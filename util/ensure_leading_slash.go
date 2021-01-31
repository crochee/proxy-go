// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package util

// EnsureLeadingSlash makes sure str has lead slash
func EnsureLeadingSlash(str string) string {
	if str == "" {
		return str
	}

	if str[0] == '/' {
		return str
	}

	return "/" + str
}
