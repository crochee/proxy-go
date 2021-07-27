package circuitbreaker

import (
	"net/http"

	"github.com/vulcand/oxy/cbreaker"

	"github.com/crochee/proxy/config/dynamic"
	"github.com/crochee/proxy/internal"
	"github.com/crochee/proxy/pkg/logger"
	"github.com/crochee/proxy/pkg/middleware"
	"github.com/crochee/proxy/pkg/tracex"
)

func New(cfg dynamic.CircuitBreaker) *circuitBreaker {
	return &circuitBreaker{
		expression: cfg.Expression,
	}
}

type circuitBreaker struct {
	expression     string
	next           middleware.Handler
	circuitBreaker *cbreaker.CircuitBreaker
}

func (c *circuitBreaker) Name() string {
	return "CIRCUIT_BREAKER"
}

func (c *circuitBreaker) Level() int {
	return 2
}

func (c *circuitBreaker) Next(handler middleware.Handler) middleware.Handler {
	c.next = handler
	oxyCircuitBreaker, err := cbreaker.New(c.next, c.expression, createCircuitBreakerOptions(c.expression))
	if err != nil {
		logger.Error(err.Error())
		return c
	}
	c.circuitBreaker = oxyCircuitBreaker
	return c
}

func (c *circuitBreaker) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if c.circuitBreaker == nil {
		c.next.ServeHTTP(writer, request)
		return
	}
	c.circuitBreaker.ServeHTTP(writer, request)
}

// createCircuitBreakerOptions returns a new cbreaker.CircuitBreakerOption.
func createCircuitBreakerOptions(expression string) cbreaker.CircuitBreakerOption {
	return cbreaker.Fallback(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		tracex.SetErrorWithEvent(req, "blocked by circuit-breaker (%q)", expression)
		rw.WriteHeader(http.StatusServiceUnavailable)

		if _, err := rw.Write(internal.Bytes(internal.StatusText(http.StatusServiceUnavailable))); err != nil {
			logger.FromContext(req.Context()).Error(err.Error())
		}
	}))
}
