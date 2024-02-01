package storage

import (
	"context"
	"io"
	"os/exec"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/haupc/foundryutils/client"
	"github.com/stretchr/testify/suite"
)

// anvil --fork-url $OPTIMISM_NODE_ENDPOINT
type StorageSuite struct {
	suite.Suite
	anvilCmd *exec.Cmd
}

func (s *StorageSuite) TearDownTest() {
	err := s.anvilCmd.Process.Kill()
	s.Assert().NoError(err)
}

func (s *StorageSuite) SetupTest() {
	_, err := exec.LookPath("anvil")
	s.Assert().NoError(err)
	s.anvilCmd = exec.Command("anvil", "--fork-url", "https://rpc.ankr.com/optimism")
	s.T().Log(s.anvilCmd.String())

	out, err := s.anvilCmd.StdoutPipe()
	s.Assert().NoError(err)

	err = s.anvilCmd.Start()
	s.Assert().NoError(waitForAnvilReady(out))
	s.Assert().NoError(err)
}

func (s *StorageSuite) TestFindErc20BalanceOfSlotIdx() {
	client.SetupClient("http://127.0.0.1:8545")
	slotIdx, err := FindErc20BalanceOfSlotIdx(common.HexToAddress("0x4200000000000000000000000000000000000006")) // WETH
	s.Assert().NoError(err)
	s.Assert().NotNil(slotIdx)
	s.Assert().Equal(int64(3), slotIdx.Int64())
}

func (s *StorageSuite) TestFindErc20BalanceOfSlotIdx_ErrNotFound() {
	client.SetupClient("http://127.0.0.1:8545")
	slotIdx, err := FindErc20BalanceOfSlotIdx(common.HexToAddress("0xDeadBeef")) // dummy address
	s.Assert().EqualError(err, "not found")
	s.Assert().Nil(slotIdx)
}

func (s *StorageSuite) TestFindErc20BalanceOfSlotIdx_ErrCallContract() {
	client.SetupClient("http://127.0.0.1:8545")
	slotIdx, err := FindErc20BalanceOfSlotIdx(common.HexToAddress("0xDeadBeef")) // dummy address
	s.Assert().Error(err)
	s.Assert().Nil(slotIdx)

}

func (s *StorageSuite) TestSetStorageAt() {
	contractAddress := common.HexToAddress("0x4200000000000000000000000000000000000006") // WETH
	err := SetStorageAt(contractAddress, common.HexToHash("0x1"), common.HexToHash("0xDeadBeef"))
	s.Assert().NoError(err)
	b, err := client.GlobalClient.EthClient.StorageAt(context.Background(), contractAddress, common.HexToHash("0x1"), nil)
	s.Assert().NoError(err)
	s.Assert().Equal(common.HexToHash("0xDeadBeef"), common.BytesToHash(b))
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
