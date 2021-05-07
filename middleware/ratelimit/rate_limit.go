// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package ratelimit

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"github.com/crochee/proxy-go/internal"
	"github.com/crochee/proxy-go/logger"
	"github.com/crochee/proxy-go/middleware"
)

type option struct {
	maxDelay time.Duration
	every    time.Duration
	burst    int
	mode     int
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

type rateLimiter struct {
	limiter *rate.Limiter
	next    http.Handler
	option
}

// New returns a rate limiter middleware.
func New(next http.Handler, opts ...func(*option)) middleware.Handler {
	rateLimiter := &rateLimiter{
		next: next,
		option: option{
			every: 500 * time.Millisecond,
			burst: 1000 * 1000,
			mode:  1,
		},
	}
	for _, opt := range opts {
		opt(&rateLimiter.option)
	}
	every := rate.Every(rateLimiter.every)
	if every < 1 {
		rateLimiter.maxDelay = 500 * time.Millisecond
	} else {
		rateLimiter.maxDelay = time.Second / (time.Duration(every) * 2)
	}
	rateLimiter.limiter = rate.NewLimiter(every, rateLimiter.burst)
	return rateLimiter
}

func (rl *rateLimiter) NameSpace() string {
	return "RateLimiter"
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
