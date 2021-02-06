// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/6

package balance

import (
	"context"
	"strconv"
	"testing"

	"proxy-go/model"
)

func TestBalancer_Update(t *testing.T) {
	b := New(context.Background())
	tt := []struct {
		Add     bool
		handler *model.NamedHandler
	}{
		{
			Add: false,
			handler: &model.NamedHandler{
				Handler:  nil,
				Host:     nil,
				Weight:   0,
				Deadline: 0,
			},
		},
		{
			Add: true,
			handler: &model.NamedHandler{
				Handler: nil,
				Host: &model.Host{
					Scheme: "http",
					Host:   "127.0.0.1:8150",
				},
				Weight:   1,
				Deadline: 0,
			},
		},
		{
			Add: true,
			handler: &model.NamedHandler{
				Handler: nil,
				Host: &model.Host{
					Scheme: "http",
					Host:   "localhost:8150",
				},
				Weight:   1,
				Deadline: 0,
			},
		},
		{
			Add: true,
			handler: &model.NamedHandler{
				Handler: nil,
				Host: &model.Host{
					Scheme: "http",
					Host:   "127.0.0.1:8150",
				},
				Weight:   2,
				Deadline: 0,
			},
		},
		{
			Add: false,
			handler: &model.NamedHandler{
				Handler: nil,
				Host: &model.Host{
					Scheme: "http",
					Host:   "127.0.0.1:8150",
				},
				Weight:   2,
				Deadline: 0,
			},
		},
	}
	for index, hand := range tt {
		t.Run("index:"+strconv.Itoa(index), func(t *testing.T) {
			b.Update(hand.Add, hand.handler)
			for _, handler := range b.handlers {
				t.Logf("%+v", handler)
			}
		})
	}
}
