package ratelimit

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"github.com/crochee/proxy/internal"
	"github.com/crochee/proxy/pkg/logger"
	"github.com/crochee/proxy/pkg/middleware"
)

type option struct {
	every time.Duration
	burst int
	mode  int
}

func Every(t time.Duration) func(*option) {
	return func(o *option) { o.every = t }
}

func Burst(b int) func(*option) {
	return func(o *option) { o.burst = b }
}

func Mode(mode int) func(*option) {
	return func(o *option) { o.mode = mode }
}

// New returns a rate limiter middleware.
func New(opts ...func(*option)) *rateLimiter {
	l := &rateLimiter{
		option: option{
			every: 500 * time.Millisecond,
			burst: 1000 * 1000,
			mode:  1,
		},
	}
	for _, opt := range opts {
		opt(&l.option)
	}
	every := rate.Every(l.every)
	if every < 1 {
		l.maxDelay = 500 * time.Millisecond
	} else {
		l.maxDelay = time.Second / (time.Duration(every) * 2)
	}
	l.limiter = rate.NewLimiter(every, l.burst)
	return l
}

type rateLimiter struct {
	limiter *rate.Limiter
	next    middleware.Handler
	option
	maxDelay time.Duration
}

func (rl *rateLimiter) Name() string {
	return "RATE_LIMITER"
}

func (rl *rateLimiter) Level() int {
	return 5
}

func (rl *rateLimiter) Next(handler middleware.Handler) middleware.Handler {
	rl.next = handler
	return rl
}

func (rl *rateLimiter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch rl.mode {
	case 0:
		if err := rl.limiter.Wait(req.Context()); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	case 1:
		if !rl.limiter.Allow() {
			http.Error(rw, "No bursty allowed", http.StatusTooManyRequests)
			return
		}
	case 2:
		res := rl.limiter.Reserve()
		if !res.OK() {
			http.Error(rw, "No bursty allowed", http.StatusTooManyRequests)
			return
		}
		delay := res.Delay()
		if delay > rl.maxDelay {
			res.Cancel()
			rl.serveDelayError(rw, req, delay)
			return
		}
		time.Sleep(delay)
	default:
	}
	rl.next.ServeHTTP(rw, req)
}

func (rl *rateLimiter) serveDelayError(w http.ResponseWriter, req *http.Request, delay time.Duration) {
	w.Header().Set("Retry-After", fmt.Sprintf("%.0f", delay.Seconds()))
	w.Header().Set("X-Retry-In", delay.String())
	w.WriteHeader(http.StatusTooManyRequests)

	if _, err := w.Write(internal.Bytes(internal.StatusText(http.StatusTooManyRequests))); err != nil {
		logger.FromContext(req.Context()).Errorf("could not serve 429: %v", err)
	}
}
