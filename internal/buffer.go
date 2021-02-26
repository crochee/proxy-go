// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/4

package internal

import "go.uber.org/zap/buffer"

var (
	BufferPool = buffer.NewPool()
	// Get retrieves a buffer from the pool, creating one if necessary.
	GetBuffer = BufferPool.Get
)
