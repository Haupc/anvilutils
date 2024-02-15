package client

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	RpcClient  *rpc.Client
	EthClient  *ethclient.Client
	GethClient *gethclient.Client
}

// init 3 types of client to
// interact with rpc node endpoint
func NewClient(endpoint string) *Client {
	rpcClient, err := rpc.Dial(endpoint)
	if err != nil {
		panic(err)
	}
	ethClient := ethclient.NewClient(rpcClient)
	gethClient := gethclient.New(rpcClient)
	return &Client{
		RpcClient:  rpcClient,
		EthClient:  ethClient,
		GethClient: gethClient,
	}
}
