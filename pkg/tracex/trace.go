package tracex

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/crochee/proxy-go/pkg/logger"
)

type tracerKey struct{}

type Finish func()

// WithTracing Adds Tracer into the context.
func WithTracer(ctx context.Context, tracer *Tracer) context.Context {
	return context.WithValue(ctx, tracerKey{}, tracer)
}

// FromContext Gets Tracing from context.
func FromContext(ctx context.Context) (*Tracer, error) {
	tracer, ok := ctx.Value(tracerKey{}).(*Tracer)
	if !ok {
		return nil, errors.New("unable to find tracing in the context")
	}
	return tracer, nil
}

// Backend is an abstraction for tracker backend (Jaeger, Zipkin, ...).
type Backend interface {
	Setup(componentName string) (opentracing.Tracer, io.Closer, error)
}

// Tracer middleware.
type Tracer struct {
	ServiceName string `description:"Set the name for this service" export:"true"`
	// nolint:lll
	SpanNameLimit int `description:"Set the maximum character limit for Span names (default 0 = no limit)" export:"true"`

	tracer opentracing.Tracer
	closer io.Closer
}

// NewTracer Creates a Tracer.
func NewTracer(serviceName string, spanNameLimit int, tracingBackend Backend) (*Tracer, error) {
	tracing := &Tracer{
		ServiceName:   serviceName,
		SpanNameLimit: spanNameLimit,
	}

	var err error
	tracing.tracer, tracing.closer, err = tracingBackend.Setup(serviceName)
	if err != nil {
		return nil, err
	}
	return tracing, nil
}

func (t *Tracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return t.tracer.StartSpan(operationName, opts...)
}

func (t *Tracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return t.tracer.Inject(sm, format, carrier)
}

func (t *Tracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return t.tracer.Extract(format, carrier)
}

// IsEnabled determines if tracer was successfully activated.
func (t *Tracer) IsEnabled() bool {
	return t != nil && t.tracer != nil
}

// Close tracer.
func (t *Tracer) Close() {
	if t.closer != nil {
		err := t.closer.Close()
		if err != nil {
			logger.Warn(err.Error())
		}
	}
}

// StartSpanf delegates to StartSpan.
func (t *Tracer) StartSpanf(r *http.Request, spanKind ext.SpanKindEnum, opPrefix string, opParts []string,
	separator string, opts ...opentracing.StartSpanOption) (opentracing.Span, *http.Request, Finish) {
	operationName := generateOperationName(opPrefix, opParts, separator, t.SpanNameLimit)

	return StartSpan(r, operationName, spanKind, opts...)
}

// RecordRequest used to create span tags from the request.
func RecordRequest(span opentracing.Span, r *http.Request) {
	if span != nil && r != nil {
		ext.HTTPMethod.Set(span, r.Method)
		ext.HTTPUrl.Set(span, r.URL.String())
		span.SetTag("http.host", r.Host)
	}
}

// RecordResponseCode used to log response code in span.
func RecordResponseCode(span opentracing.Span, code int) {
	if span != nil {
		ext.HTTPStatusCode.Set(span, uint16(code))
		if code >= http.StatusInternalServerError {
			ext.Error.Set(span, true)
		}
	}
}

// GetSpan used to retrieve span from request context.
func GetSpan(r *http.Request) opentracing.Span {
	return opentracing.SpanFromContext(r.Context())
}

// InjectRequestHeaders used to inject OpenTracing headers into the request.
func InjectRequestHeaders(r *http.Request) {
	if span := GetSpan(r); span != nil {
		err := opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			logger.FromContext(r.Context()).Error(err.Error())
		}
	}
}

// RecordEventf logs an event to the span in the request context.
func RecordEventf(r *http.Request, format string, args ...interface{}) {
	if span := GetSpan(r); span != nil {
		span.LogKV("event", fmt.Sprintf(format, args...))
	}
}

// StartSpan starts a new span from the one in the request context.
func StartSpan(r *http.Request, operationName string, spanKind ext.SpanKindEnum,
	opts ...opentracing.StartSpanOption) (opentracing.Span, *http.Request, Finish) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), operationName, opts...)

	switch spanKind {
	case ext.SpanKindRPCClientEnum:
		ext.SpanKindRPCClient.Set(span)
	case ext.SpanKindRPCServerEnum:
		ext.SpanKindRPCServer.Set(span)
	case ext.SpanKindProducerEnum:
		ext.SpanKindProducer.Set(span)
	case ext.SpanKindConsumerEnum:
		ext.SpanKindConsumer.Set(span)
	default:
		// noop
	}

	return span, r.WithContext(ctx), func() { span.Finish() }
}

// SetError flags the span associated with this request as in error.
func SetError(r *http.Request) {
	if span := GetSpan(r); span != nil {
		ext.Error.Set(span, true)
	}
}

// SetErrorWithEvent flags the span associated with this request as in error and log an event.
func SetErrorWithEvent(r *http.Request, format string, args ...interface{}) {
	SetError(r)
	RecordEventf(r, format, args...)
}
