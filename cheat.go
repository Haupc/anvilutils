package anvilutils

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/haupc/anvilutils/client"
	"github.com/haupc/anvilutils/contracts"
	"github.com/haupc/anvilutils/helper"
	"github.com/haupc/anvilutils/storage"
)

type Cheat struct {
	client *client.Client
}

func NewCheat(client *client.Client) *Cheat {
	return &Cheat{client: client}
}

// write balance of erc20 token for an account
func (c *Cheat) WriteErc20Balance(contract, user common.Address, amount *big.Int) error {
	balanceOfSlotIdx, err := storage.FindErc20BalanceOfSlotIdx(c.client, contract)
	if err != nil {
		return err
	}
	accountBalanceSlotIdx := storage.Erc20AccountBalanceSlotIdx(c.client, balanceOfSlotIdx, user)
	return c.SetStorageAt(contract, accountBalanceSlotIdx, common.BytesToHash(amount.Bytes()))
}

// write native balance for an account
func (c *Cheat) WriteNativeBalance(user common.Address, amount *big.Int) error {
	return c.client.RpcClient.Call(nil, "anvil_setBalance", user.Hex(), hexutil.Encode(amount.Bytes()))
}

// set code for an address
func (c *Cheat) SetCode(account common.Address, code []byte) error {
	return c.client.RpcClient.Call(nil, "anvil_setCode", account.Hex(), hexutil.Encode(code))
}

// set storge for an contract
func (c *Cheat) SetStorageAt(contractAddress common.Address, slotIdx, data common.Hash) error {
	return storage.SetStorageAt(c.client, contractAddress, slotIdx, data)
}

// after this call, all transactions from specified account
// can be executed without signing
func (c *Cheat) StartImpersonateAccount(account common.Address) error {
	return c.client.RpcClient.Call(nil, "anvil_impersonateAccount", account.Hex())
}

// stop impersonating account
func (c *Cheat) StopImpersonateAccount(account common.Address) error {
	return c.client.RpcClient.Call(nil, "anvil_stopImpersonatingAccount", account.Hex())
}

// make an impersonate txn
func (c *Cheat) SendImpersonateTxn(from, to common.Address, value *big.Int, data []byte, gas uint64, gasPrice *big.Int) (txHash common.Hash, err error) {
	if value == nil {
		value = big.NewInt(0)
	}
	if gas == 0 {
		gas, err = c.client.EthClient.EstimateGas(context.Background(), ethereum.CallMsg{
			From:  from,
			To:    &to,
			Value: value,
			Data:  data,
		})
		if err != nil {
			return
		}
	}
	if gasPrice == nil || gasPrice.Cmp(big.NewInt(0)) == 0 {
		gasPrice, err = c.client.EthClient.SuggestGasPrice(context.Background())
		if err != nil {
			return
		}
	}

	err = c.client.RpcClient.Call(&txHash, "eth_sendTransaction",
		struct {
			From     string `json:"from"`
			To       string `json:"to"`
			Value    string `json:"value"`
			Gas      string `json:"gas"`
			GasPrice string `json:"gasPrice"`
			Data     string `json:"data"`
		}{
			From:     from.Hex(),
			To:       to.Hex(),
			Value:    fmt.Sprintf("0x%s", value.Text(16)),
			Gas:      fmt.Sprintf("0x%x", gas),
			GasPrice: hexutil.Encode(gasPrice.Bytes()),
			Data:     hexutil.Encode(data),
		},
	)
	return
}

func (c *Cheat) ImpersonateAccountAndSendTransaction(from, to common.Address, value *big.Int, data []byte, gas uint64, gasPrice *big.Int) (txHash common.Hash, err error) {
	if err = c.StartImpersonateAccount(from); err != nil {
		return
	}
	defer c.StopImpersonateAccount(from)
	return c.SendImpersonateTxn(from, to, value, data, gas, gasPrice)
}

func (c *Cheat) SetApprovalErc20(owner, token, spender common.Address, amount *big.Int) error {
	callData := helper.Erc20ApproveCallData(spender, amount)
	if _, err := c.ImpersonateAccountAndSendTransaction(owner, token, nil, callData, 0, nil); err != nil {
		return err
	}
	return nil
}

func (c *Cheat) TakeErc721Token(tokenAddress common.Address, tokenId *big.Int, receiver common.Address) error {
	erc721Contract, _ := contracts.NewErc721(tokenAddress, c.client.EthClient)
	tokenOwner, err := erc721Contract.OwnerOf(nil, tokenId)
	if err != nil {
		return err
	}
	// function transferFrom(address _from, address _to, uint256 _tokenId)
	abi, _ := contracts.Erc721MetaData.GetAbi()
	callData, _ := abi.Pack("transferFrom", tokenOwner, receiver, tokenId)
	_, err = c.ImpersonateAccountAndSendTransaction(tokenOwner, tokenAddress, nil, callData, 0, nil)
	return err
}
