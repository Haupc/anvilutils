package anvilutils

import (
	"context"
	"fmt"

	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/haupc/anvilutils/client"
	"github.com/haupc/anvilutils/contracts"
	"github.com/haupc/anvilutils/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CheatSuite struct {
	suite.Suite
	cheat *Cheat
}

func (s *CheatSuite) TestSetCode() {
	b, err := s.cheat.client.EthClient.CodeAt(context.Background(), helper.DummyErc20Contract, nil)
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
	s.Assert().NoError(s.cheat.WriteErc20Balance(helper.DummyErc20Contract, helper.DummyAccount, big.NewInt(1234567890123)))
	balance, err := s.cheat.client.EthClient.CallContract(context.Background(), ethereum.CallMsg{
		To:   &helper.DummyErc20Contract,
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
	s.Assert().NoError(s.cheat.SetApprovalErc20(helper.DummyAccount, helper.DummyErc20Contract, spender, helper.MaxUint128))
	erc20Contract, _ := contracts.NewErc20(helper.DummyErc20Contract, s.cheat.client.EthClient)
	allowance, err := erc20Contract.Allowance(nil, helper.DummyAccount, spender)
	s.Assert().NoError(err)
	s.Assert().Equal(helper.MaxUint128.String(), allowance.String())
}

func (s *CheatSuite) TestTakeErc721Token() {

	tokenAddress := helper.DummyErc721Contract
	tokenId := big.NewInt(986148)
	receiver := common.HexToAddress("0xdeaddead")

	mockErc721Balance(s.cheat)
	err := s.cheat.TakeErc721Token(tokenAddress, tokenId, receiver)
	s.Assert().NoError(err)

	erc721Contract, _ := contracts.NewErc721(tokenAddress, s.cheat.client.EthClient)
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
	err = cheatSuite.cheat.SetCode(helper.DummyErc20Contract, helper.WETH_Code)
	assert.NoError(t, err)
	suite.Run(t, cheatSuite)
}

func mockErc721Balance(cheat *Cheat) {
	cheat.SetCode(helper.DummyErc721Contract, helper.UniV3_NFPM_CODE)
	cheat.SetStorageAt(helper.DummyErc721Contract, common.HexToHash("0x2"), common.HexToHash("0xf8562"))
	cheat.SetStorageAt(helper.DummyErc721Contract, common.HexToHash("0x405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3d6b81f"), common.HexToHash("0xdeadbeef"))
	cheat.SetStorageAt(helper.DummyErc721Contract, common.HexToHash("0x52d1192e38ae8dd4312020815b03b1d6249abb2c370c9f29b1bf1901c490ea58"), common.HexToHash("0xf"))
	cheat.SetStorageAt(helper.DummyErc721Contract, common.HexToHash("0x8003b79ccb3357ef955735bf2d79c8eec95faebc6572403c109c4d6635d6ed76"), common.HexToHash("0x1"))
	cheat.SetStorageAt(helper.DummyErc721Contract, common.HexToHash("0x9ae6d7f0a910acb25dc64a1b43dd42976c51bcec30ecde275f5a3a9988788ca9"), common.HexToHash("0xdaea9"))
	cheat.SetStorageAt(helper.DummyErc721Contract, common.HexToHash("0xb4c69c9f05a8ee32626026766df0c1c20b3cea795556d6c90198edfbd4803f6f"), common.HexToHash("0x1"))
	cheat.SetStorageAt(helper.DummyErc721Contract, common.HexToHash("0xeea09047952cbb2a75aa7bc824be08b7e980734c90435c281119c68beb49b725"), common.HexToHash("0xf0c24"))
}
