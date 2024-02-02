package helper

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/haupc/foundryutils/contracts"
)

// pack data for balanceOf method
// solidity: balanceOf(address)
func BalanceOfCallData(userAddress common.Address) []byte {
	abi, _ := contracts.Erc20MetaData.GetAbi()
	packedData, _ := abi.Pack("balanceOf", userAddress)
	return packedData
}
