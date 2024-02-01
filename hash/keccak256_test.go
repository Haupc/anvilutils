package hash

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestKeccak256(t *testing.T) {
	b := Keccak256([]byte("balanceOf(address)"))
	assert.Equal(t, "0x70a08231", hexutil.Encode(b[:4]))
}
