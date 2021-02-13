// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package dynamic

import "time"

type RateLimit struct {
	Every time.Duration
	Burst int
	Mode  int
}
