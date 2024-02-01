package storage

import (
	"bytes"
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/haupc/foundryutils/client"
	"github.com/haupc/foundryutils/hash"
	"github.com/haupc/foundryutils/helper"
)

var (
	maxUint128, _ = new(big.Int).SetString("ffffffffffffffffffffffffffffffff", 16)
)

// Get slot for account balance
// when know slot index of the
// `balanceOf` map (map address => uint256)
func Erc20AccountBalanceSlotIdx(balanceOfSlotIdx *big.Int, addr common.Address) common.Hash {
	// https://dev.kit.eco/ethereum-simple-deep-dive-into-evm-storage
	mappingTyp, _ := abi.NewType("tuple", "mappings", []abi.ArgumentMarshaling{
		{
			Name: "key",
			Type: "address",
		},
		{
			Name: "slotNum",
			Type: "uint256",
		},
	})
	arg := abi.Arguments{
		{
			Type: mappingTyp,
			Name: "mappingsArg",
		},
	}
	packed, _ := arg.Pack(struct {
		Key     common.Address
		SlotNum *big.Int
	}{addr, balanceOfSlotIdx})
	return hash.Keccak256(packed)
}

// Set storage slot at slot index
func SetStorageAt(contractAddress common.Address, idx, data common.Hash) error {
	return client.GlobalClient.RpcClient.Call(nil, "anvil_setStorageAt", contractAddress.Hex(), idx.Hex(), data.Hex())
}

func FindErc20BalanceOfSlotIdx(contractAddress common.Address) (*big.Int, error) {
	var randomAddress = common.HexToAddress("0xDeaDBeef") // 0x00000000000000000000000000000000DeaDBeef
	randomBalance := common.BytesToHash(maxUint128.Bytes())
	for i := 0; i < 100; i++ {
		balanceOfSlotIdx := big.NewInt(int64(i))
		accountBalanceSlotIdx := Erc20AccountBalanceSlotIdx(balanceOfSlotIdx, randomAddress)
		b, err := client.GlobalClient.GethClient.CallContract(context.Background(), ethereum.CallMsg{
			To:   &contractAddress,
			Data: helper.BalanceOfCallData(randomAddress),
		}, nil, &map[common.Address]gethclient.OverrideAccount{
			contractAddress: {
				StateDiff: map[common.Hash]common.Hash{
					accountBalanceSlotIdx: randomBalance,
				},
			},
		})

		if err != nil {
			return nil, err
		}

		if bytes.Equal(b, randomBalance.Bytes()) {
			return balanceOfSlotIdx, nil
		}
	}
	return nil, errors.New("not found")
}
