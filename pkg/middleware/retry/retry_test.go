package retry

import (
	"net/http"
	"testing"
	"time"

	"github.com/crochee/proxy-go/config/dynamic"
	"github.com/crochee/proxy-go/internal"
)

func TestNew(t *testing.T) {
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/proxy", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(504)
	})
	mux := New(dynamic.Retry{
		Attempts:        5,
		InitialInterval: time.Second,
	})
	mux.Next(mux1)
	resp := internal.PerformRequest(mux, http.MethodGet, "/proxy", nil, nil)
	t.Logf("%v", resp)
}

func BenchmarkNew(b *testing.B) {
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/proxy", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(503)
	})
	mux := New(dynamic.Retry{
		Attempts:        5,
		InitialInterval: time.Second,
	})
	mux.Next(mux1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		internal.PerformRequest(mux, http.MethodGet, "/proxy", nil, nil)
	}
}
