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
)

type rateLimiter struct {
	limiterMap *sync.Map
	mux        sync.Mutex
	next       http.Handler
	ctx        context.Context

	maxDelay time.Duration
	every    time.Duration
	burst    int
}

// New returns a rate limiter middleware.
func New(ctx context.Context, next http.Handler) *rateLimiter {
	rateLimiter := &rateLimiter{
		limiterMap: new(sync.Map),
		next:       next,
		ctx:        ctx,
		every:      10 * time.Microsecond,
		burst:      1,
	}
	every := rate.Every(rateLimiter.every)
	if every < 1 {
		rateLimiter.maxDelay = 500 * time.Millisecond
	} else {
		rateLimiter.maxDelay = time.Second / (time.Duration(every) * 2)
	}
	return rateLimiter
}

func (rl *rateLimiter) Name() string {
	return "RateLimiter"
}

func (rl *rateLimiter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if rl.limit(rw, req) {
		rl.next.ServeHTTP(rw, req)
	}
}

func (rl *rateLimiter) serveDelayError(w http.ResponseWriter, delay time.Duration) {
	w.Header().Set("Retry-After", fmt.Sprintf("%.0f", delay.Seconds()))
	w.Header().Set("X-Retry-In", delay.String())
	w.WriteHeader(http.StatusTooManyRequests)

	if _, err := w.Write(internal.Bytes(internal.StatusText(http.StatusTooManyRequests))); err != nil {
		logger.FromContext(rl.ctx).Errorf("could not serve 429: %v", err)
	}
}

func (rl *rateLimiter) Update(limit dynamic.RateLimit) {
	rl.mux.Lock()
	rl.every = limit.Every
	rl.burst = limit.Burst
	rl.mux.Unlock()
}

func (rl *rateLimiter) limit(rw http.ResponseWriter, req *http.Request) bool {
	var limiter *rate.Limiter
	value, ok := rl.limiterMap.Load(req.RemoteAddr)
	if !ok {
		rl.mux.Lock()
		every := rate.Every(rl.every)
		if every < 1 {
			rl.maxDelay = 500 * time.Millisecond
		} else {
			rl.maxDelay = time.Second / (time.Duration(every) * 2)
		}
		limiter = rate.NewLimiter(every, rl.burst)
		rl.mux.Unlock()

		rl.limiterMap.Store(req.RemoteAddr, limiter)
	}
	if limiter, ok = value.(*rate.Limiter); !ok {
		time.Sleep(rl.maxDelay)
	}

	res := limiter.Reserve()
	if !res.OK() {
		http.Error(rw, "No bursty traffic allowed", http.StatusTooManyRequests)
		return false
	}

	delay := res.Delay()
	if delay > rl.maxDelay {
		res.Cancel()
		rl.serveDelayError(rw, delay)
		return false
	}
	time.Sleep(delay)
	return true
}
