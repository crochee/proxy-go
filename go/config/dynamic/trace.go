package dynamic

import "github.com/crochee/proxy/pkg/tracex/jaeger"

type TraceInfo struct {
	Jaeger *jaeger.Config `description:"Settings for Jaeger." json:"jaeger,omitempty" yaml:"jaeger,omitempty"`
}
