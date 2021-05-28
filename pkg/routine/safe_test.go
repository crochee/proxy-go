// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/16

package routine

import (
	"context"
	"log"
	"testing"
)

func TestNewPool(t *testing.T) {
	testList := []struct {
		name    string
		recover func(interface{})
		input   []func(ctx context.Context)
	}{
		{
			name: "default",
			input: []func(ctx context.Context){
				func(ctx context.Context) {
				},
				func(ctx context.Context) {
					panic("test0")
				},
				func(ctx context.Context) {
					panic("test1")
				},
				func(ctx context.Context) {
					panic("test2")
				},
				func(ctx context.Context) {
					panic("test3")
				}, func(ctx context.Context) {
					panic("test4")
				}, func(ctx context.Context) {
					panic("test5")
				}, func(ctx context.Context) {
					panic("test6")
				}, func(ctx context.Context) {
					panic("test7")
				},
			},
		},
		{
			name: "recover",
			recover: func(i interface{}) {
				log.Println(i)
			},
			input: []func(ctx context.Context){
				func(ctx context.Context) {
					panic("op")
				},
			},
		},
	}
	for _, data := range testList {
		t.Run(data.name, func(t *testing.T) {
			var p *pool
			if data.recover != nil {
				p = NewPool(context.Background(), Recover(data.recover))
			} else {
				p = NewPool(context.Background())
			}
			for _, f := range data.input {
				p.Go(f)
			}
			p.Stop()
		})
	}
}
