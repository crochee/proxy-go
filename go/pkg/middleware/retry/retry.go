package retry

import (
	"io"
	"math"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/cenkalti/backoff/v4"

	"github.com/crochee/proxy/config/dynamic"
	"github.com/crochee/proxy/internal"
	"github.com/crochee/proxy/pkg/middleware"
)

// nexter returns the duration to wait before retrying the operation.
type nexter interface {
	NextBackOff() time.Duration
}

// New create a middleware that retries requests.
func New(rt dynamic.Retry) *retry {
	return &retry{
		initialInterval: rt.InitialInterval,
		attempts:        rt.Attempts,
	}
}

// retry is a middleware that retries requests.
type retry struct {
	next            middleware.Handler
	initialInterval time.Duration
	attempts        int
}

func (r *retry) Name() string {
	return "RETRY"
}

func (r *retry) Level() int {
	return 1
}

func (r *retry) Next(handler middleware.Handler) middleware.Handler {
	r.next = handler
	return r
}

func (r *retry) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if r.attempts > 1 {
		// 一般情况发送失败的时候，body会关闭并返回err,导致重试时，数据被破坏
		body := req.Body
		defer internal.Close(body)
		req.Body = io.NopCloser(body)
	}

	var attempts int
	backOff := r.newBackOff() // 退避算法 保证时间间隔为指数级增长
	currentInterval := 0 * time.Millisecond
	t := time.NewTimer(currentInterval)
	for {
		select {
		case <-t.C:
			shouldRetry := attempts < r.attempts
			retryResponseWriter := newResponseWriter(rw, shouldRetry)
			// Disable retries when the backend already received request data
			trace := &httptrace.ClientTrace{
				WroteRequest: func(httptrace.WroteRequestInfo) {
					retryResponseWriter.DisableRetries()
				},
			}
			newCtx := httptrace.WithClientTrace(req.Context(), trace)

			r.next.ServeHTTP(retryResponseWriter, req.WithContext(newCtx))

			if !retryResponseWriter.ShouldRetry() {
				t.Stop()
				return
			}
			// 计算下一次
			currentInterval = backOff.NextBackOff()
			attempts++
			// 定时器重置
			t.Reset(currentInterval)
		case <-req.Context().Done():
			t.Stop()
			return
		}
	}
}

func (r *retry) newBackOff() nexter {
	if r.attempts < 2 || r.initialInterval <= 0 {
		return &backoff.ZeroBackOff{}
	}

	b := backoff.NewExponentialBackOff()
	b.InitialInterval = r.initialInterval

	// calculate the multiplier for the given number of attempts
	// so that applying the multiplier for the given number of attempts will not exceed 2 times the initial interval
	// it allows to control the progression along the attempts
	b.Multiplier = math.Pow(2, 1/float64(r.attempts-1))

	// according to docs, b.Reset() must be called before using
	b.Reset()
	return b
}
