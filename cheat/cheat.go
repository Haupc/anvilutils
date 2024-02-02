package cheat

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/haupc/foundryutils/client"
	"github.com/haupc/foundryutils/storage"
)

// write balance of erc20 token for an account
func WriteErc20Balance(contract, user common.Address, amount *big.Int) error {
	balanceOfSlotIdx, err := storage.FindErc20BalanceOfSlotIdx(contract)
	if err != nil {
		return err
	}
	accountBalanceSlotIdx := storage.Erc20AccountBalanceSlotIdx(balanceOfSlotIdx, user)
	return storage.SetStorageAt(contract, accountBalanceSlotIdx, common.BytesToHash(amount.Bytes()))
}

// write native balance for an account
func WriteNativeBalance(user common.Address, amount *big.Int) error {
	return client.GlobalClient.RpcClient.Call(nil, "anvil_setBalance", user.Hex(), hexutil.Encode(amount.Bytes()))
}

// set code for an address
func SetCode(account common.Address, code []byte) error {
	return client.GlobalClient.RpcClient.Call(nil, "anvil_setCode", account.Hex(), hexutil.Encode(code))
}
