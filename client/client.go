package client

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var GlobalClient *Client

type Client struct {
	RpcClient  *rpc.Client
	EthClient  *ethclient.Client
	GethClient *gethclient.Client
}

// init 3 types of client to
// interact with rpc node endpoint
func SetupClient(endpoint string) {
	rpcClient, err := rpc.Dial(endpoint)
	if err != nil {
		panic(err)
	}
	ethClient := ethclient.NewClient(rpcClient)
	gethClient := gethclient.New(rpcClient)
	GlobalClient = &Client{
		RpcClient:  rpcClient,
		EthClient:  ethClient,
		GethClient: gethClient,
	}
}
