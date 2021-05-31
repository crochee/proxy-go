package trace

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/crochee/proxy-go/pkg/logger"
	"github.com/crochee/proxy-go/pkg/tracex"
	"github.com/crochee/proxy-go/pkg/writer"
)

type entryPoint struct {
	*tracex.Tracer
	entryPoint string
	next       http.Handler
}

// NewEntryPoint creates a new middleware that the incoming request.
func NewTraceEntryPoint(t *tracex.Tracer, entryPointName string, next http.Handler) http.Handler {
	return &entryPoint{
		Tracer:     t,
		entryPoint: entryPointName,
		next:       next,
	}
}

func (e *entryPoint) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	spanCtx, err := e.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(request.Header))
	if err != nil {
		logger.FromContext(request.Context()).Errorf("Failed to extract the context: %v", err)
	}

	span, req, finish := e.StartSpanf(request, ext.SpanKindRPCServerEnum, "EntryPoint",
		[]string{e.entryPoint, request.Host}, " ", ext.RPCServerOption(spanCtx))
	defer finish()

	ext.Component.Set(span, e.ServiceName)
	tracex.RecordRequest(span, req)

	req = req.WithContext(tracex.WithTracer(req.Context(), e.Tracer))

	e.next.ServeHTTP(rw, req)

	if recorder, ok := rw.(writer.Capture); ok {
		tracex.RecordResponseCode(span, recorder.Status())
	}
}
