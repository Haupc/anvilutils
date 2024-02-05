[![Go Reference](https://pkg.go.dev/badge/github.com/Haupc/foundryutils.svg)](https://pkg.go.dev/github.com/Haupc/foundryutils)
[![Go Report Card](https://goreportcard.com/badge/github.com/haupc/foundryutils)](https://goreportcard.com/report/github.com/haupc/foundryutils)
[![Go Coverage Badge](https://raw.githubusercontent.com/Haupc/foundryutils/badges/.badges/master/coverage.svg)](https://raw.githubusercontent.com/Haupc/foundryutils/badges/.badges/master/coverage.svg)
# foundry utils

## cheat code:

- [x] Write Erc20 balance
- [x] Write native balance
- [ ] Write Erc20 approval
- [x] Set code for address
- [x] Start impersonate an address
- [x] Stop impersonate an address
- [ ] Impersonate and make only 1 txn
- [ ] Write Erc721 balance

## How to Use:
```
    forkRpcEndpoint, err := foundryutils.StartFork("https://rpc.ankr.com/eth", nil)
    if err != nil {
        // do something with err
    }
    client.SetupClient(forkRpcEndpoint)
    cheat.WriteErc20Balance(helper.DummyContract, helper.DummyAccount, big.NewInt(1234567890123))
    
    // call api, do everything for testing
    // do something with onchain data using client.GlobalClient
    // to sop fork
    foundryutils.StopFork()
```