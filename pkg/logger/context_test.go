// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/29

package logger

import (
	"context"
	"testing"
)

func TestContext(t *testing.T) {
	ctx := Context(context.Background(), NewLogger())
	t.Log(ctx)
	l := FromContext(ctx)
	t.Log(l)
}
