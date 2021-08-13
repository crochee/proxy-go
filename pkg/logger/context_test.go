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
