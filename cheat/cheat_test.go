package cheat

import (
	"context"
	"io"
	"math/big"
	"os/exec"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/haupc/foundryutils/client"
	"github.com/haupc/foundryutils/helper"
	"github.com/stretchr/testify/suite"
)

type CheatSuite struct {
	suite.Suite
	anvilCmd *exec.Cmd
}

func (s *CheatSuite) TearDownTest() {
	err := s.anvilCmd.Process.Kill()
	s.Assert().NoError(err)
}

func (s *CheatSuite) SetupTest() {
	client.SetupClient("http://127.0.0.1:8545")
	_, err := exec.LookPath("anvil")
	s.Assert().NoError(err)
	s.anvilCmd = exec.Command("anvil")
	s.T().Log(s.anvilCmd.String())

	out, err := s.anvilCmd.StdoutPipe()
	s.Assert().NoError(err)

	err = s.anvilCmd.Start()
	s.Assert().NoError(err)
	s.Assert().NoError(waitForAnvilReady(out))
	s.Assert().NoError(SetCode(helper.DummyContract, helper.WETH_Code))
}

func (s *CheatSuite) TestSetCode() {
	b, err := client.GlobalClient.EthClient.CodeAt(context.Background(), helper.DummyContract, nil)
	s.Assert().NoError(err)
	s.Assert().Equal(hexutil.Encode(helper.WETH_Code), hexutil.Encode(b))
}

func (s *CheatSuite) TestWriteNativeBalance() {
	s.Assert().NoError(WriteNativeBalance(helper.DummyAccount, big.NewInt(1000000000000)))
	balance, err := client.GlobalClient.EthClient.BalanceAt(context.Background(), helper.DummyAccount, nil)
	s.Assert().NoError(err)
	s.Assert().Equal("1000000000000", balance.String())
}

func (s *CheatSuite) TestWriteErc20Balance() {
	s.Assert().NoError(WriteErc20Balance(helper.DummyContract, helper.DummyAccount, big.NewInt(1234567890123)))
	balance, err := client.GlobalClient.EthClient.CallContract(context.Background(), ethereum.CallMsg{
		To:   &helper.DummyContract,
		Data: helper.BalanceOfCallData(helper.DummyAccount),
	}, nil)
	s.Assert().NoError(err)
	s.Assert().Equal("1234567890123", new(big.Int).SetBytes(balance).String())
}

func (s *CheatSuite) TestStartImpersonateAccount() {
	WriteNativeBalance(helper.DummyAccount, big.NewInt(1234567890987654321))
	s.Assert().NoError(
		StartImpersonateAccount(helper.DummyAccount),
	)
	s.Assert().NoError(
		SendImpersonateTxn(helper.DummyAccount, common.HexToAddress("0x08081999"), 21000, big.NewInt(1000000000000000000), big.NewInt(10000000000000), nil),
	)
	b, err := client.GlobalClient.EthClient.BalanceAt(context.Background(), common.HexToAddress("0x08081999"), nil)
	s.Assert().NoError(err)
	s.Assert().Equal(hexutil.Encode(big.NewInt(1000000000000000000).Bytes()), hexutil.Encode(b.Bytes()))
}

func waitForAnvilReady(out io.ReadCloser) error {
	readyString := "Listening on 127.0.0.1:8545"
	for {
		buff := make([]byte, 1024)
		if _, err := out.Read(buff); err != nil {
			return err
		}
		if strings.Contains(string(buff), readyString) {
			return nil
		}
	}
}

func TestCheatSuite(t *testing.T) {
	suite.Run(t, new(CheatSuite))
}
