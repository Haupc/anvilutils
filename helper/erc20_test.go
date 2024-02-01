package helper

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestBalanceOfCallData(t *testing.T) {
	address := common.HexToAddress("0xdeadbeef")
	b := BalanceOfCallData(address)
	assert.Equal(t, "0x70a0823100000000000000000000000000000000000000000000000000000000deadbeef", hexutil.Encode(b))
}
