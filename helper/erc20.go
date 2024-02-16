package helper

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/haupc/anvilutils/contracts"
)

// pack data for balanceOf method
// solidity: balanceOf(address)
func Erc20BalanceOfCallData(userAddress common.Address) []byte {
	abi, _ := contracts.Erc20MetaData.GetAbi()
	packedData, _ := abi.Pack("balanceOf", userAddress)
	return packedData
}

// pack data for approve method
// solidity: approve(address spender, uint256 amount) returns(bool)
func Erc20ApproveCallData(spender common.Address, amount *big.Int) []byte {
	abi, _ := contracts.Erc20MetaData.GetAbi()
	packedData, _ := abi.Pack("approve", spender, amount)
	return packedData
}
