// Package proxy
package proxygo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	assert.NoError(t, Server())
}
