// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/15

package circuitbreaker

import (
	"net/http"

	"github.com/vulcand/oxy/cbreaker"

	"github.com/crochee/proxy-go/config/dynamic"
	"github.com/crochee/proxy-go/internal"
	"github.com/crochee/proxy-go/logger"
	"github.com/crochee/proxy-go/service/tracex"
)

type circuitBreaker struct {
	circuitBreaker *cbreaker.CircuitBreaker
}

func New(cfg dynamic.CircuitBreaker, next http.Handler) (http.Handler, error) {
	oxyCircuitBreaker, err := cbreaker.New(next, cfg.Expression, createCircuitBreakerOptions(cfg.Expression))
	if err != nil {
		return nil, err
	}
	return &circuitBreaker{
		circuitBreaker: oxyCircuitBreaker,
	}, nil
}
func (c *circuitBreaker) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c.circuitBreaker.ServeHTTP(writer, request)
}

// NewCircuitBreakerOptions returns a new CircuitBreakerOption.
func createCircuitBreakerOptions(expression string) cbreaker.CircuitBreakerOption {
	return cbreaker.Fallback(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		tracex.SetErrorWithEvent(req, "blocked by circuit-breaker (%q)", expression)
		rw.WriteHeader(http.StatusServiceUnavailable)

		if _, err := rw.Write([]byte(internal.StatusText(http.StatusServiceUnavailable))); err != nil {
			logger.FromContext(req.Context()).Error(err.Error())
		}
	}))
}
