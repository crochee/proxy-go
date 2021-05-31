// Package tcpx
package tcpx

import (
	"errors"
	"io"
	"net"
	"net/http/httputil"
	"strings"
	"sync"
	"time"

	"github.com/crochee/proxy-go/config/dynamic"
	"github.com/crochee/proxy-go/internal"
	"github.com/crochee/proxy-go/pkg/logger"
	"github.com/crochee/proxy-go/pkg/selector"
)

type Proxy struct {
	log              logger.Builder
	BufferPool       httputil.BufferPool
	TargetSelector   map[string]selector.Selector
	terminationDelay time.Duration
	rw               sync.RWMutex
}

func New(log logger.Builder, cfg *dynamic.Config) Handler {
	p := &Proxy{
		log:              log,
		BufferPool:       internal.BufPool,
		TargetSelector:   make(map[string]selector.Selector),
		terminationDelay: 0,
	}
	for _, balance := range cfg.Balance.Transfers {
		p.TargetSelector[balance.ServiceName] = createSelector(&balance.Balance)
	}
	return p
}

func (p *Proxy) ServeTCP(conn WriteCloser) {
	p.log.Debugf("Handling connection from %s", conn.RemoteAddr())

	// needed because of e.g. server.trackedConnection
	defer conn.Close()

	connBackend, err := p.dialBackend(conn.RemoteAddr().String())
	if err != nil {
		p.log.Errorf("Error while connecting to backend: %v", err)
		return
	}

	// maybe not needed, but just in case
	defer connBackend.Close()
	errChan := make(chan error)

	go p.connCopy(conn, connBackend, errChan)
	go p.connCopy(connBackend, conn, errChan)

	err = <-errChan
	if err != nil {
		p.log.Errorf("Error during connection: %v", err)
	}

	<-errChan
}

func (p *Proxy) dialBackend(addr string) (*net.TCPConn, error) {
	p.rw.RLock()
	s, ok := p.TargetSelector[addr]
	p.rw.RUnlock()
	if !ok {
		return nil, errors.New("no backend tcp addr")
	}
	node, err := s.Next()
	if err != nil {
		return nil, err
	}
	if strings.ToUpper(node.Scheme) != "TCP" {
		return nil, errors.New("scheme is not tcp")
	}
	var tcp *net.TCPAddr
	if tcp, err = net.ResolveTCPAddr("tcp", node.Host); err != nil {
		return nil, err
	}
	return net.DialTCP("tcp", nil, tcp)
}

func (p *Proxy) connCopy(dst, src WriteCloser, errCh chan error) {
	var buf []byte
	if p.BufferPool != nil {
		buf = p.BufferPool.Get()
		defer p.BufferPool.Put(buf)
	}
	_, err := io.CopyBuffer(dst, src, buf)
	errCh <- err

	errClose := dst.CloseWrite()
	if errClose != nil {
		p.log.Debugf("Error while terminating connection: %v", errClose)
		return
	}

	if p.terminationDelay >= 0 {
		err = dst.SetReadDeadline(time.Now().Add(p.terminationDelay))
		if err != nil {
			p.log.Debugf("Error while setting deadline: %v", err)
		}
	}
}

func createSelector(balance *dynamic.Balance) selector.Selector {
	var s selector.Selector
	switch strings.Title(balance.Selector) {
	case "Random":
		s = selector.NewRandom()
	case "RoundRobin":
		s = selector.NewRoundRobin()
	case "Heap":
		s = selector.NewHeap()
	case "Wrr":
		fallthrough
	default:
		s = selector.NewWeightRoundRobin()
	}
	for _, node := range balance.Nodes {
		s.AddNode(&selector.Node{
			Scheme:   node.Scheme,
			Host:     node.Host,
			Metadata: node.Metadata,
			Weight:   node.Weight,
		})
	}
	return s
}
