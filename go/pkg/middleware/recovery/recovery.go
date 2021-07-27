package recovery

import (
	"net/http"
	"runtime/debug"

	"github.com/crochee/proxy/internal"
	"github.com/crochee/proxy/pkg/logger"
	"github.com/crochee/proxy/pkg/middleware"
)

// New creates recovery middleware
func New() *recovery {
	return &recovery{}
}

type recovery struct {
	next middleware.Handler
}

func (r *recovery) Name() string {
	return "RECOVERY"
}

func (r *recovery) Level() int {
	return 5
}

func (r *recovery) Next(handler middleware.Handler) middleware.Handler {
	r.next = handler
	return r
}

func (r *recovery) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log := logger.FromContext(req.Context())
			if r == http.ErrAbortHandler { // nolint:errorlint
				log.Debugf("Request has been aborted [%s - %s]: %v", req.RemoteAddr, req.URL, r)
				return
			}
			log.Errorf("[Recovery] from panic in HTTP handler [%s - %s]: %+v\nStack:\n%s",
				req.RemoteAddr, req.URL, r, debug.Stack())

			http.Error(rw, internal.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}()
	r.next.ServeHTTP(rw, req)
}
