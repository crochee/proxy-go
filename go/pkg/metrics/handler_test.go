package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnable(t *testing.T) {
	Enable.Store(true)
	require.Equal(t, true, Enable.Load())
}
