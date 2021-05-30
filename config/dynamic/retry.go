// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/30

package dynamic

import (
	"time"
)

// Retry holds the retry configuration.
type Retry struct {
	Attempts        int           `json:"attempts,omitempty" yaml:"attempts,omitempty"`
	InitialInterval time.Duration `json:"initialInterval,omitempty" yaml:"initialInterval,omitempty"`
}
