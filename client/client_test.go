package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupClient_Success(t *testing.T) {
	var client *Client
	assert.NotPanics(t, func() {
		client = NewClient("http://localhost:8545")
	})
	assert.NotNil(t, client)
}

func TestSetupClient_Panic(t *testing.T) {
	assert.Panics(t, func() {
		_ = NewClient("")
	})
}
