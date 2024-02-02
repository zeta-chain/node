package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.True(t, cfg.JSONRPC.Enable)
	assert.Equal(t, cfg.JSONRPC.Address, DefaultJSONRPCAddress)
	assert.Equal(t, cfg.JSONRPC.WsAddress, DefaultJSONRPCWsAddress)
}
