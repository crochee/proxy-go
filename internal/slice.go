// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/4

package internal

import "sync"

const _poolSize = 32 * 1024

var (
	BufPool = NewBufPool()
)

func NewBufPool() *bufferPool {
	return &bufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, _poolSize)
			},
		},
	}
}

type bufferPool struct {
	pool sync.Pool
}

func (b *bufferPool) Get() []byte {
	return b.pool.Get().([]byte)
}

func (b *bufferPool) Put(bytes []byte) {
	b.pool.Put(bytes)
}
