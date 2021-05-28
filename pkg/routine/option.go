// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/16

package routine

type option struct {
	recoverFunc func(interface{})
}

// Recover register to pool
func Recover(f func(interface{})) func(*option) {
	return func(o *option) { o.recoverFunc = f }
}
