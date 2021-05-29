// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/30

package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnable(t *testing.T) {
	Enable.Store(true)
	require.Equal(t, true, Enable.Load())
}
