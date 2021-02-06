// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/6

package balance

import (
	"context"
	"reflect"
	"strconv"
	"testing"
)

func TestBalancer_Update(t *testing.T) {
	b := New(context.Background(), NewRandom(), nil)
	tt := []struct {
		Add  bool
		Node *Node
	}{
		{
			Add: false,
			Node: &Node{
				Scheme:   "http",
				Host:     "127.0.0.1:8150",
				Metadata: nil,
			},
		},
		{
			Add: true,
			Node: &Node{
				Scheme:   "http",
				Host:     "192.168.31.62:8150",
				Metadata: nil,
			},
		},
		{
			Add: true,
			Node: &Node{
				Scheme:   "http",
				Host:     "localhost:8150",
				Metadata: nil,
			},
		},
		{
			Add: true,
			Node: &Node{
				Scheme:   "http",
				Host:     "127.0.0.1:8150",
				Metadata: nil,
			},
		},
		{
			Add: false,
			Node: &Node{
				Scheme:   "http",
				Host:     "127.0.0.1:8150",
				Metadata: nil,
			},
		},
	}
	for index, hand := range tt {
		t.Run("index:"+strconv.Itoa(index), func(t *testing.T) {
			b.Update(hand.Add, hand.Node, 1)
			for _, handler := range b.selector.List() {
				t.Logf("%+v", handler)
			}
		})
	}
}

func TestBalancer_Name(t *testing.T) {
	b := New(nil, nil, nil)
	t.Log(reflect.TypeOf(b))
}
