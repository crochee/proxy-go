// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/6

package balance

import (
	"container/heap"
	"errors"
	"math/rand"
	"reflect"
	"sync"
	"time"
)

type Node struct {
	Scheme   string            `json:"scheme"`
	Host     string            `json:"host"`
	Metadata map[string]string `json:"metadata"`
}

var ErrNoneAvailable = errors.New("none available")

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Selector strategy algorithm
type Selector interface {
	Update(bool, *Node, float64)
	Next() (*Node, error)
	List() []*Node
}

type Random struct {
	mux  sync.RWMutex
	list []*Node
}

func NewRandom() *Random {
	return &Random{
		list: make([]*Node, 0, 4),
	}
}

func (r *Random) Update(add bool, node *Node, weight float64) {
	r.mux.Lock()
	defer r.mux.Unlock()
	var equal bool
	for index, list := range r.list {
		if reflect.DeepEqual(list, node) {
			if !add {
				if index == len(r.list)-1 {
					r.list = r.list[:index]
					return
				}
				r.list = append(r.list[:index], r.list[index+1:]...)
				return
			}
			equal = true
		}
	}
	if !equal {
		r.list = append(r.list, node)
	}
}

func (r *Random) Next() (*Node, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	length := len(r.list)
	if length == 0 {
		return nil, ErrNoneAvailable
	}
	i := rand.Int() % length
	return r.list[i], nil
}

func (r *Random) List() []*Node {
	return r.list
}

type RoundRobin struct {
	randIndex int
	mux       sync.Mutex
	list      []*Node
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		randIndex: rand.Int(),
		list:      make([]*Node, 0, 4),
	}
}

func (r *RoundRobin) Update(add bool, node *Node, weight float64) {
	r.mux.Lock()
	defer r.mux.Unlock()
	var equal bool
	for index, list := range r.list {
		if reflect.DeepEqual(list, node) {
			if !add {
				if index == len(r.list)-1 {
					r.list = r.list[:index]
					return
				}
				r.list = append(r.list[:index], r.list[index+1:]...)
				return
			}
			equal = true
		}
	}
	if !equal {
		r.list = append(r.list, node)
	}
}

func (r *RoundRobin) Next() (*Node, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	length := len(r.list)
	if length == 0 {
		return nil, ErrNoneAvailable
	}
	r.randIndex %= length
	node := r.list[r.randIndex]
	r.randIndex++
	return node, nil
}

func (r *RoundRobin) List() []*Node {
	return r.list
}

type weightHandler struct {
	node     *Node
	weight   float64
	deadline float64
}

type Heap struct {
	mutex       sync.RWMutex
	handlers    []*weightHandler
	curDeadline float64
}

func NewHeap() *Heap {
	return &Heap{
		handlers: make([]*weightHandler, 0, 4),
	}
}

func (h *Heap) Len() int {
	return len(h.handlers)
}

func (h *Heap) Less(i, j int) bool {
	return h.handlers[i].deadline < h.handlers[j].deadline
}

func (h *Heap) Swap(i, j int) {
	h.handlers[i], h.handlers[j] = h.handlers[j], h.handlers[i]
}

func (h *Heap) Push(x interface{}) {
	handler, ok := x.(*weightHandler)
	if !ok {
		return
	}
	h.handlers = append(h.handlers, handler)
}

func (h *Heap) Pop() interface{} {
	if h.Len() < 1 {
		return nil
	}
	handler := h.handlers[len(h.handlers)-1]
	h.handlers = h.handlers[:len(h.handlers)-1]
	return handler
}

func (h *Heap) Update(add bool, node *Node, weight float64) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	var equal bool
	for index, list := range h.handlers {
		if reflect.DeepEqual(list, node) {
			if !add {
				if index == len(h.handlers)-1 {
					h.handlers = h.handlers[:index]
					return
				}
				h.handlers = append(h.handlers[:index], h.handlers[index+1:]...)
				return
			}
			equal = true
		}
	}
	if !equal {
		w := &weightHandler{node: node, weight: weight}
		h.Push(w)
		// use RWLock to protect b.curDeadline
		w.deadline = h.curDeadline + 1/w.weight
	}
}

func (h *Heap) Next() (*Node, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.Len() == 0 {
		return nil, ErrNoneAvailable
	}
	handler, ok := heap.Pop(h).(*weightHandler)
	if !ok {
		return nil, ErrNoneAvailable
	}
	// curDeadline should be handler's deadline so that
	// new added entry would have a fair competition environment with the old ones.
	h.curDeadline = handler.deadline
	handler.deadline += 1 / handler.weight
	heap.Push(h, handler)

	return handler.node, nil
}

func (h *Heap) List() []*Node {
	list := make([]*Node, 0, len(h.handlers))
	h.mutex.RLock()
	for _, handler := range h.handlers {
		list = append(list, handler.node)
	}
	return list
}
