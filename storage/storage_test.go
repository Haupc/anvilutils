package storage

import (
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/haupc/foundryutils/client"
	"github.com/haupc/foundryutils/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type StorageSuite struct {
	suite.Suite
	client *client.Client
}

func (s *StorageSuite) TestFindErc20BalanceOfSlotIdx() {
	slotIdx, err := FindErc20BalanceOfSlotIdx(s.client, helper.DummyErc20Contract)
	s.Assert().NoError(err)
	s.Assert().NotNil(slotIdx)
	s.Assert().Equal(int64(3), slotIdx.Int64())
}

func (s *StorageSuite) TestFindErc20BalanceOfSlotIdx_ErrNotFound() {
	slotIdx, err := FindErc20BalanceOfSlotIdx(s.client, helper.DummyAccount) // dummy address
	s.Assert().EqualError(err, "not found")
	s.Assert().Nil(slotIdx)
}

func (s *StorageSuite) TestFindErc20BalanceOfSlotIdx_Err() {
	client := client.NewClient("http://dead.link")
	slotIdx, err := FindErc20BalanceOfSlotIdx(client, helper.DummyAccount) // dummy address
	s.Assert().Error(err)
	s.Assert().Nil(slotIdx)
}

func (s *StorageSuite) TestSetStorageAt() {
	err := SetStorageAt(s.client, helper.DummyErc20Contract, common.HexToHash("0x1"), helper.DummyHash)
	s.Assert().NoError(err)
	b, err := s.client.EthClient.StorageAt(context.Background(), helper.DummyErc20Contract, common.HexToHash("0x1"), nil)
	s.Assert().NoError(err)
	s.Assert().Equal(helper.DummyHash, common.BytesToHash(b))
}

func TestStorageSuite(t *testing.T) {
	storageSuite := new(StorageSuite)

	cmd, endpoint, err := startAnvil()
	assert.NoError(t, err)
	defer cmd.Process.Kill()
	client := client.NewClient(endpoint)
	err = client.RpcClient.Call(nil, "anvil_setCode", helper.DummyErc20Contract, hexutil.Encode(helper.WETH_Code))
	assert.NoError(t, err)

	storageSuite.client = client
	suite.Run(t, storageSuite)
}

func startAnvil() (*exec.Cmd, string, error) {
	anvilPort := fmt.Sprint(rand.Int31n(65535-1024) + 1024)
	if _, err := exec.LookPath("anvil"); err != nil {
		return nil, "", err
	}
	cmd := exec.Command("anvil", "--port", anvilPort)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, "", err
	}
	if err := cmd.Start(); err != nil {
		return nil, "", err
	}

	// wait for anvil to start
	readyString := "Listening on 127.0.0.1:" + anvilPort
	for {
		buff := make([]byte, 1024)
		if _, err := out.Read(buff); err != nil {
			return nil, "", err
		}
		if strings.Contains(string(buff), readyString) {
			return cmd, "http://127.0.0.1:" + anvilPort, nil
		}
	}
}
