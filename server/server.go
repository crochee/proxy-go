// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"proxy-go/config"
	"proxy-go/logger"
	"proxy-go/safe"
)

type Server struct {
	config   *config.Config
	pool     *safe.Pool
	list     []*http.Server
	handler  http.Handler
	stopChan chan struct{}
	ctx      context.Context
	sync.RWMutex
}

// New returns an initialized Server.
func New(ctx context.Context, routinesPool *safe.Pool, cf *config.Config, handler http.Handler) *Server {
	return &Server{
		config:   cf,
		pool:     routinesPool,
		list:     make([]*http.Server, 0, len(cf.Server.Medata)),
		handler:  handler,
		stopChan: make(chan struct{}, 1),
		ctx:      ctx,
	}
}

func (s *Server) Start() {
	go func() {
		<-s.ctx.Done()
		logger.FromContext(s.ctx).Info("Stopping server gracefully")
		s.Stop()
	}()
	if s.config.Server == nil {
		s.Stop()
	}
	for _, medata := range s.config.Server.Medata {
		s.pool.GoCtx(func(ctx context.Context) {
			srv := &http.Server{
				Addr:      fmt.Sprintf(":%d", medata.Port),
				Handler:   s.handler,
				TLSConfig: nil,
			}
			s.Lock()
			s.list = append(s.list, srv)
			s.Unlock()
			log := logger.FromContext(ctx)
			switch medata.Scheme {
			case "http":
				log.Info("http server running...")
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Error(err.Error())
				}
			case "https":
				log.Info("https server running...")
				if err := srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
					log.Error(err.Error())
				}
			default:
				return
			}
			s.Lock()
			s.list = s.list[:len(s.list)-1]
			s.Unlock()
		})
	}
}

// Wait blocks until the server shutdown.
func (s *Server) Wait() {
	<-s.stopChan
}

// Stop stops the server.
func (s *Server) Stop() {
	for _, srv := range s.list {
		s.pool.GoCtx(func(ctx context.Context) {
			Shutdown(ctx, srv)
		})
	}
	s.stopChan <- struct{}{}
	logger.FromContext(s.ctx).Info("server stopped")
}

// Close destroys the server.
func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)

	go func(ctx context.Context) {
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.Canceled) {
			return
		}
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			panic("timeout while stopping proxy, killing instance ✝")
		}
	}(ctx)

	s.pool.Stop()

	close(s.stopChan)
	cancel()
}

func Shutdown(ctx context.Context, server *http.Server) {
	err := server.Shutdown(ctx)
	if err == nil {
		return
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		logger.FromContext(ctx).Debugf("server failed to shutdown within deadline because: %s", err)
		if err = server.Close(); err != nil {
			logger.Error(err.Error())
		}
		return
	}
	logger.FromContext(ctx).Error(err.Error())
	// We expect Close to fail again because Shutdown most likely failed when trying to close a listener.
	// We still call it however, to make sure that all connections get closed as well.
	_ = server.Close()
}
