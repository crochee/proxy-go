package circuitbreaker

import (
	"net/http"

	"github.com/vulcand/oxy/cbreaker"

	"github.com/crochee/proxy-go/config/dynamic"
	"github.com/crochee/proxy-go/internal"
	"github.com/crochee/proxy-go/pkg/logger"
	"github.com/crochee/proxy-go/pkg/tracex"
)

func New(cfg dynamic.CircuitBreaker) http.Handler {
	return &circuitBreaker{
		expression: cfg.Expression,
	}
}

type circuitBreaker struct {
	expression     string
	next           http.Handler
	circuitBreaker *cbreaker.CircuitBreaker
}

func (c *circuitBreaker) Name() string {
	return "CIRCUIT_BREAKER"
}

func (c *circuitBreaker) Level() int {
	return 2
}

func (c *circuitBreaker) Next(handler http.Handler) {
	c.next = handler
	oxyCircuitBreaker, err := cbreaker.New(c.next, c.expression, createCircuitBreakerOptions(c.expression))
	if err != nil {
		logger.Error(err.Error())
		return
	}
	c.circuitBreaker = oxyCircuitBreaker
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
