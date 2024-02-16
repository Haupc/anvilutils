[![Go Reference](https://pkg.go.dev/badge/github.com/Haupc/anvilutils.svg)](https://pkg.go.dev/github.com/Haupc/anvilutils)
[![Go Report Card](https://goreportcard.com/badge/github.com/haupc/anvilutils)](https://goreportcard.com/report/github.com/haupc/anvilutils)
[![Go Coverage Badge](https://raw.githubusercontent.com/Haupc/anvilutils/badges/.badges/master/coverage.svg)](https://raw.githubusercontent.com/Haupc/anvilutils/badges/.badges/master/coverage.svg)
# anvil utils
anvil utils is a library that provides cheating functionality executed in anvil foundry for interacting with account, transaction especially erc20 and erc721 token

## installation

To use this library, you need to install anvil foundry first. ([see this link](https://book.getfoundry.sh/getting-started/installation))

TL;DR:
```bash
curl -L https://foundry.paradigm.xyz | bash
foundryup
```

Import:
```go
import github.com/haupc/anvilutils
```

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
- [ ] Set approve Erc721 token

## How to Use:

You can run anvil separately or use forkCmd to run an anvil chain.
```go
    // if you want to fork on test code, you can use fork command
    // otherwise, you can run anvil separately then call setup cheats
    forkCmd, err := anvilutils.NewForkCommand(ForkOpts{})
    if err != nil {
        // do something with err
    }
    forkUrl, err := forkCmd.Start()
    if err != nil {
        // do something with err
    }

    // setup client for anvil node
    client := client.NewClient(forkUrl)
    cheat := anvilutils.NewCheat(client)
    cheat.WriteErc20Balance(helper.DummyContract, helper.DummyAccount, big.NewInt(1234567890123))
    
    // call api, do everything for testing
    // do something with onchain data using client.GlobalClient
    // to sop fork
    forkCmd.Stop()
```
