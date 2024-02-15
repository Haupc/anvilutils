package helper

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestErc20BalanceOfCallData(t *testing.T) {
	address := common.HexToAddress("0xdeadbeef")
	b := Erc20BalanceOfCallData(address)
	assert.Equal(t, "0x70a0823100000000000000000000000000000000000000000000000000000000deadbeef", hexutil.Encode(b))
}

func TestErc20ApproveCallData(t *testing.T) {
	amount, _ := new(big.Int).SetString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	b := Erc20ApproveCallData(common.HexToAddress("0xdeadbeef"), amount)
	assert.Equal(t, "0x095ea7b300000000000000000000000000000000000000000000000000000000deadbeefffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", hexutil.Encode(b))
}
