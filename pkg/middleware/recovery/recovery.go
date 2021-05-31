package recovery

import (
	"net/http"
	"runtime/debug"

	"github.com/crochee/proxy-go/internal"
	"github.com/crochee/proxy-go/pkg/logger"
)

type recovery struct {
	next http.Handler
}

// New creates recovery middleware
func New(next http.Handler) http.Handler {
	return &recovery{
		next: next,
	}
}

func (re *recovery) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
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
	re.next.ServeHTTP(rw, req)
}
