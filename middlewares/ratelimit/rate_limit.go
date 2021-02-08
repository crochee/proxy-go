// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package ratelimit

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"proxy-go/config/dynamic"
	"proxy-go/internal"
	"proxy-go/logger"
	"proxy-go/middlewares"
)

type rateLimiter struct {
	limiter *rate.Limiter
	mux     sync.Mutex
	next    http.Handler
	ctx     context.Context

	maxDelay time.Duration
	every    time.Duration
	burst    int
}

// New returns a rate limiter middleware.
func New(ctx context.Context, next http.Handler) *rateLimiter {
	rateLimiter := &rateLimiter{
		next:  next,
		ctx:   ctx,
		every: 100 * time.Microsecond,
		burst: 1000 * 1000 * 1000,
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

func (rl *rateLimiter) Name() string {
	return middlewares.RateLimiter
}

func (rl *rateLimiter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	res := rl.limiter.Reserve()
	if !res.OK() {
		http.Error(rw, "No bursty traffic allowed", http.StatusTooManyRequests)
		return
	}

	delay := res.Delay()
	if delay > rl.maxDelay {
		res.Cancel()
		rl.serveDelayError(rw, delay)
		return
	}
	time.Sleep(delay)

	rl.next.ServeHTTP(rw, req)
}

func (rl *rateLimiter) serveDelayError(w http.ResponseWriter, delay time.Duration) {
	w.Header().Set("Retry-After", fmt.Sprintf("%.0f", delay.Seconds()))
	w.Header().Set("X-Retry-In", delay.String())
	w.WriteHeader(http.StatusTooManyRequests)

	if _, err := w.Write(internal.Bytes(internal.StatusText(http.StatusTooManyRequests))); err != nil {
		logger.FromContext(rl.ctx).Errorf("could not serve 429: %v", err)
	}
}

func (rl *rateLimiter) Update(limit *dynamic.RateLimit) {
	rl.mux.Lock()
	if rl.every != limit.Every || rl.burst != limit.Burst {
		rl.every, rl.burst = limit.Every, limit.Burst

		every := rate.Every(rl.every)
		if every < 1 {
			rl.maxDelay = 500 * time.Millisecond
		} else {
			rl.maxDelay = time.Second / (time.Duration(every) * 2)
		}
		rl.limiter = rate.NewLimiter(every, rl.burst)
	}
	rl.mux.Unlock()
}
