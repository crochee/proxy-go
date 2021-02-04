// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package ratelimit

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"proxy-go/internal"
	"proxy-go/logger"
	"proxy-go/model"
)

type rateLimiter struct {
	limiter  *rate.Limiter // reqs/s
	next     http.Handler
	maxDelay time.Duration
	ctx      context.Context
}

// New returns a rate limiter middleware.
func New(ctx context.Context, next http.Handler, limit *model.RateLimit) http.Handler {
	rateLimiter := &rateLimiter{
		next: next,
		ctx:  ctx,
	}
	every := rate.Every(limit.Every)
	if every < 1 {
		rateLimiter.maxDelay = 500 * time.Millisecond
	} else {
		rateLimiter.maxDelay = time.Second / (time.Duration(every) * 2)
	}
	rateLimiter.limiter = rate.NewLimiter(every, limit.Burst)
	return rateLimiter
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
