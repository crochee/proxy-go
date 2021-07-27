// Package proxy
package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	assert.NoError(t, Server())
}
