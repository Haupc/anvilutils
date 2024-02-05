package cheat

import (
	"fmt"
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

// after this call, all transactions from specified account
// can be executed without signing
func StartImpersonateAccount(account common.Address) error {
	return client.GlobalClient.RpcClient.Call(nil, "anvil_impersonateAccount", account.Hex())
}

// stop impersonating account
func StopImpersonateAccount(account common.Address) error {
	return client.GlobalClient.RpcClient.Call(nil, "anvil_stopImpersonatingAccount", account.Hex())
}

// make an impersonate txn
func SendImpersonateTxn(from, to common.Address, gas uint64, value, gasPrice *big.Int, data []byte) error {
	return client.GlobalClient.RpcClient.Call(nil, "eth_sendTransaction",
		struct {
			From     string `json:"from"`
			To       string `json:"to"`
			Value    string `json:"value"`
			Gas      string `json:"gas"`
			GasPrice string `json:"gasPrice"`
		}{
			From:     from.Hex(),
			To:       to.Hex(),
			Value:    hexutil.Encode(value.Bytes()),
			Gas:      fmt.Sprintf("0x%x", gas),
			GasPrice: hexutil.Encode(gasPrice.Bytes()),
		},
	)
}
