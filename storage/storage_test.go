package storage

import (
	"context"
	"io"
	"os/exec"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/haupc/foundryutils/client"
	"github.com/haupc/foundryutils/helper"
	"github.com/stretchr/testify/suite"
)

type StorageSuite struct {
	suite.Suite
	anvilCmd *exec.Cmd
}

func (s *StorageSuite) TearDownTest() {
	err := s.anvilCmd.Process.Kill()
	s.Assert().NoError(err)
}

func (s *StorageSuite) SetupTest() {
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
	s.Assert().NoError(client.GlobalClient.RpcClient.Call(nil, "anvil_setCode", helper.DummyContract, hexutil.Encode(helper.WETH_Code)))
}

func (s *StorageSuite) TestFindErc20BalanceOfSlotIdx() {
	slotIdx, err := FindErc20BalanceOfSlotIdx(helper.DummyContract)
	s.Assert().NoError(err)
	s.Assert().NotNil(slotIdx)
	s.Assert().Equal(int64(3), slotIdx.Int64())
}

func (s *StorageSuite) TestFindErc20BalanceOfSlotIdx_ErrNotFound() {
	slotIdx, err := FindErc20BalanceOfSlotIdx(helper.DummyAccount) // dummy address
	s.Assert().EqualError(err, "not found")
	s.Assert().Nil(slotIdx)
}

func (s *StorageSuite) TestFindErc20BalanceOfSlotIdx_Err() {
	client.SetupClient("http://dead.link")
	slotIdx, err := FindErc20BalanceOfSlotIdx(helper.DummyAccount) // dummy address
	s.Assert().Error(err)
	s.Assert().Nil(slotIdx)
}

func (s *StorageSuite) TestSetStorageAt() {
	err := SetStorageAt(helper.DummyContract, common.HexToHash("0x1"), helper.DummyHash)
	s.Assert().NoError(err)
	b, err := client.GlobalClient.EthClient.StorageAt(context.Background(), helper.DummyContract, common.HexToHash("0x1"), nil)
	s.Assert().NoError(err)
	s.Assert().Equal(helper.DummyHash, common.BytesToHash(b))
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

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageSuite))
}
