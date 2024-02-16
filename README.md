[![Go Reference](https://pkg.go.dev/badge/github.com/Haupc/foundryutils.svg)](https://pkg.go.dev/github.com/Haupc/foundryutils)
[![Go Report Card](https://goreportcard.com/badge/github.com/haupc/foundryutils)](https://goreportcard.com/report/github.com/haupc/foundryutils)
[![Go Coverage Badge](https://raw.githubusercontent.com/Haupc/foundryutils/badges/.badges/master/coverage.svg)](https://raw.githubusercontent.com/Haupc/foundryutils/badges/.badges/master/coverage.svg)
# foundry utils

## cheat code:

- [x] Write Erc20 balance
- [x] Write native balance
- [x] Set Erc20 approval
- [x] Set code for address
- [x] Start impersonate an address
- [x] Stop impersonate an address
- [x] Impersonate a txn(require impersonate account)
- [x] Impersonate and make only 1 txn
- [x] Take Erc721 token from another account
- [ ] Approve Erc721 token

## How to Use:
```go
    // if you want to fork on test code, you can use fork command
    // otherwise, you can run anvil separately then call setup cheats
    forkCmd, err := foundryutils.NewForkCommand(ForkOpts{})
    if err != nil {
        // do something with err
    }
    forkUrl, err := forkCmd.Start()
    if err != nil {
        // do something with err
    }

    // setup client for anvil node
    client := client.NewClient(forkUrl)
    cheat := foundryutils.NewCheat(client)
    cheat.WriteErc20Balance(helper.DummyContract, helper.DummyAccount, big.NewInt(1234567890123))
    
    // call api, do everything for testing
    // do something with onchain data using client.GlobalClient
    // to sop fork
    forkCmd.Stop()
```