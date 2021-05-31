// Package main
package proxy_go

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/crochee/proxy-go/cmd"
)

func TestServer(t *testing.T) {
	assert.NoError(t, cmd.Server())
}
