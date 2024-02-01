package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupClient_Success(t *testing.T) {
	assert.NotPanics(t, func() {
		SetupClient("http://localhost:8545")
	})
	assert.NotNil(t, GlobalClient)
}

func TestSetupClient_Panic(t *testing.T) {
	assert.Panics(t, func() {
		SetupClient("")
	})
}
