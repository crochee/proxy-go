// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package ratelimit

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/crochee/proxy-go/config/dynamic"
	"github.com/crochee/proxy-go/internal"
	"github.com/crochee/proxy-go/logger"
	"github.com/crochee/proxy-go/middleware"
)

type rateLimiter struct {
	limiter *rate.Limiter
	mux     sync.RWMutex
	next    http.Handler

	maxDelay time.Duration
	every    time.Duration
	burst    int
	mode     int
}

// New returns a rate limiter middleware.
func New(next http.Handler) *rateLimiter {
	rateLimiter := &rateLimiter{
		next:  next,
		every: 100 * time.Microsecond,
		burst: 1000 * 1000 * 1000,
		mode:  1,
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

func (rl *rateLimiter) Name() middleware.HandlerName {
	return middleware.RateLimiter
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
			http.Error(rw, "No bursty traffic allowed", http.StatusTooManyRequests)
			return
		}
	case 2:
		res := rl.limiter.Reserve()
		if !res.OK() {
			http.Error(rw, "No bursty traffic allowed", http.StatusTooManyRequests)
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
	rl.mode = limit.Mode
	rl.mux.Unlock()
}

func (rl *rateLimiter) Get() *dynamic.RateLimit {
	rl.mux.RLock()
	rlValue := &dynamic.RateLimit{
		Every: rl.every,
		Burst: rl.burst,
		Mode:  rl.mode,
	}
	rl.mux.RUnlock()
	return rlValue
}
