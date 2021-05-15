// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/15

package dynamic

import "github.com/crochee/proxy-go/service/tracex/jaeger"

type TraceInfo struct {
	Jaeger *jaeger.Config `description:"Settings for Jaeger." json:"jaeger,omitempty" yaml:"jaeger,omitempty"`
}
