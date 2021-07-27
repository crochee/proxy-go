package internal

import "go.uber.org/zap/buffer"

var (
	BufferPool = buffer.NewPool()
	// Get retrieves a buffer from the pool, creating one if necessary.
	GetBuffer = BufferPool.Get
)
