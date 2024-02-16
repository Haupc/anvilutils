package foundryutils

import (
	"context"
	"fmt"

	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/haupc/foundryutils/client"
	"github.com/haupc/foundryutils/contracts"
	"github.com/haupc/foundryutils/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CheatSuite struct {
	suite.Suite
	cheat *Cheat
}

func (s *CheatSuite) TestSetCode() {
	b, err := s.cheat.client.EthClient.CodeAt(context.Background(), helper.DummyContract, nil)
	s.Assert().NoError(err)
	s.Assert().Equal(hexutil.Encode(helper.WETH_Code), hexutil.Encode(b))
}

func (s *CheatSuite) TestWriteNativeBalance() {
	s.Assert().NoError(s.cheat.WriteNativeBalance(helper.DummyAccount, big.NewInt(1000000000000)))
	balance, err := s.cheat.client.EthClient.BalanceAt(context.Background(), helper.DummyAccount, nil)
	s.Assert().NoError(err)
	s.Assert().Equal("1000000000000", balance.String())
}

func (s *CheatSuite) TestWriteErc20Balance() {
	s.Assert().NoError(s.cheat.WriteErc20Balance(helper.DummyContract, helper.DummyAccount, big.NewInt(1234567890123)))
	balance, err := s.cheat.client.EthClient.CallContract(context.Background(), ethereum.CallMsg{
		To:   &helper.DummyContract,
		Data: helper.Erc20BalanceOfCallData(helper.DummyAccount),
	}, nil)
	s.Assert().NoError(err)
	s.Assert().Equal("1234567890123", new(big.Int).SetBytes(balance).String())
}

func (s *CheatSuite) TestStartImpersonateAccount() {
	s.cheat.WriteNativeBalance(helper.DummyAccount, big.NewInt(1234567890987654321))
	s.Assert().NoError(
		s.cheat.StartImpersonateAccount(helper.DummyAccount),
	)
	txHash, err := s.cheat.SendImpersonateTxn(helper.DummyAccount, common.HexToAddress("0x08081999"), big.NewInt(1000000000000000000), nil, 0, nil)
	s.Assert().NoError(err)
	s.Assert().NotEqual(txHash, common.Hash{})
	b, err := s.cheat.client.EthClient.BalanceAt(context.Background(), common.HexToAddress("0x08081999"), nil)
	s.Assert().NoError(err)
	s.Assert().Equal(hexutil.Encode(big.NewInt(1000000000000000000).Bytes()), hexutil.Encode(b.Bytes()))
}

func (s *CheatSuite) TestSetApprovalErc20() {
	s.cheat.WriteNativeBalance(helper.DummyAccount, big.NewInt(1234567890987654321))
	spender := common.HexToAddress("0xdeaddead")
	s.Assert().NoError(s.cheat.SetApprovalErc20(helper.DummyAccount, helper.DummyContract, spender, helper.MaxUint128))
	erc20Contract, _ := contracts.NewErc20(helper.DummyContract, s.cheat.client.EthClient)
	allowance, err := erc20Contract.Allowance(nil, helper.DummyAccount, spender)
	s.Assert().NoError(err)
	s.Assert().Equal(helper.MaxUint128.String(), allowance.String())
}

func (s *CheatSuite) TestTakeErc721Token() {

	tokenAddress := common.HexToAddress("0xc36442b4a4522e871399cd717abdd847ab11fe88")
	tokenId := big.NewInt(986148)
	receiver := helper.DummyAccount

	forkCmd, err := NewForkCommand(ForkOpts{
		ForkUrl: "https://rpc.ankr.com/arbitrum",
		Port:    fmt.Sprint(rand.Int31n(65535-1024) + 1024),
	})
	s.Assert().NoError(err)
	endpoint, err := forkCmd.Start()
	s.Assert().NoError(err)
	defer forkCmd.Stop()
	cheat := NewCheat(client.NewClient(endpoint))
	err = cheat.TakeErc721Token(tokenAddress, tokenId, receiver)
	s.Assert().NoError(err)

	erc721Contract, _ := contracts.NewErc721(tokenAddress, cheat.client.EthClient)
	tokenOwner, err := erc721Contract.OwnerOf(nil, tokenId)
	s.Assert().NoError(err)
	s.Assert().Equal(receiver, tokenOwner)
}

func TestCheatSuite(t *testing.T) {
	cheatSuite := new(CheatSuite)
	anvilPort := fmt.Sprint(rand.Int31n(65535-1024) + 1024)

	forkCmd, err := NewForkCommand(ForkOpts{
		Port: anvilPort,
	})
	assert.NoError(t, err)

	forkUrl, err := forkCmd.Start()
	defer forkCmd.Stop()
	assert.NoError(t, err)

	assert.NotPanics(
		t, func() {
			cheatSuite.cheat = NewCheat(client.NewClient(forkUrl))
		},
	)
	err = cheatSuite.cheat.SetCode(helper.DummyContract, helper.WETH_Code)
	assert.NoError(t, err)
	suite.Run(t, cheatSuite)
}
