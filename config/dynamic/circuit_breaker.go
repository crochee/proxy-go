// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/15

package dynamic

type CircuitBreaker struct {
	Expression string `json:"expression,omitempty" yaml:"expression,omitempty"`
}
