package selector

import (
	"container/heap"
	"crypto/rand"
	"errors"
	"math/big"
	"sync"
)

type Node struct {
	Scheme   string            `json:"scheme"`
	Host     string            `json:"host"`
	Metadata map[string]string `json:"metadata"`
	Weight   float64           `json:"weight"`
}

var ErrNoneAvailable = errors.New("none available")

// selector strategy algorithm
type Selector interface {
	// 只能在初始化阶段使用
	AddNode(*Node)
	Next() (*Node, error)
}

type Random struct {
	list []*Node
}

func NewRandom() *Random {
	return &Random{
		list: make([]*Node, 0, 4),
	}
}

func (r *Random) AddNode(node *Node) {
	r.list = append(r.list, node)
}

func (r *Random) Next() (*Node, error) {
	length := len(r.list)
	if length == 0 {
		return nil, ErrNoneAvailable
	}
	bInt, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
	if err != nil {
		return nil, err
	}
	return r.list[bInt.Int64()], nil
}

type RoundRobin struct {
	randIndex int
	mux       sync.Mutex
	list      []*Node
}

func NewRoundRobin() *RoundRobin {
	bInt, err := rand.Prime(rand.Reader, 32)
	if err != nil {
		bInt = new(big.Int)
	}
	return &RoundRobin{
		randIndex: int(bInt.Int64()),
		list:      make([]*Node, 0, 4),
	}
}

func (r *RoundRobin) AddNode(node *Node) {
	r.list = append(r.list, node)
}

func (r *RoundRobin) Next() (*Node, error) {
	length := len(r.list)
	if length == 0 {
		return nil, ErrNoneAvailable
	}
	r.mux.Lock()
	r.randIndex %= length
	node := r.list[r.randIndex]
	r.randIndex++
	r.mux.Unlock()
	return node, nil
}

type deadlineNode struct {
	Node     *Node
	deadline float64
}

type Heap struct {
	mutex       sync.RWMutex
	handlers    []*deadlineNode
	curDeadline float64
}

func NewHeap() *Heap {
	return &Heap{
		handlers: make([]*deadlineNode, 0, 4),
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
	handler, ok := x.(*deadlineNode)
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

func (h *Heap) AddNode(node *Node) {
	h.Push(&deadlineNode{
		Node: node,
	})
}

func (h *Heap) Next() (*Node, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.Len() == 0 {
		return nil, ErrNoneAvailable
	}
	handler, ok := heap.Pop(h).(*deadlineNode)
	if !ok {
		return nil, ErrNoneAvailable
	}
	// curDeadline should be handler's deadline so that
	// new added entry would have a fair competition environment with the old ones.
	h.curDeadline = handler.deadline
	handler.deadline += 1 / handler.Node.Weight
	heap.Push(h, handler)

	return handler.Node, nil
}

type WeightNode struct {
	*Node
	currentWeight float64 //当前权重
}

type WeightRoundRobin struct {
	list []*WeightNode
	mux  sync.RWMutex
}

func NewWeightRoundRobin() *WeightRoundRobin {
	return &WeightRoundRobin{
		list: make([]*WeightNode, 0, 4),
	}
}

func (w *WeightRoundRobin) AddNode(node *Node) {
	w.list = append(w.list, &WeightNode{
		Node: node,
	})
}

func (w *WeightRoundRobin) Next() (*Node, error) {
	w.mux.Lock()
	defer w.mux.Unlock()
	var best *WeightNode
	var total float64
	for _, node := range w.list {
		node.currentWeight += node.Weight // 将当前权重与有效权重相加
		total += node.Weight              //累加总权重
		if best == nil || node.currentWeight > best.currentWeight {
			best = node
		}
	}
	if best == nil {
		return nil, ErrNoneAvailable
	}
	best.currentWeight -= total //将当前权重改为当前权重-总权重
	return best.Node, nil
}
