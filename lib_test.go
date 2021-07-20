// Package proxy_go
package proxy_go

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	assert.NoError(t, Server())
}
