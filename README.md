# GoatNickels - The Shitcoin

GoatNickels is a blockchain-based cryptocurrency implemented in Go. This is the server and mining application. The wallet app is called [GoatBucket](https://github.com/seanmclane/goatbucket).

## Quick Start Guide

### Building the binary
Pull or download this repo to the appropriate folder in your 'src' directory under your GOPATH then run the following.
```
go install GOPATH/src/github.com/seanmclane/goatnickels
```

### Starting the server
GoatNickels ships with a default configuration file including a directory for storing the blockchain data, bootstrap nodes running a GoatNickels server, and an account value (empty). The config lives in your home directory at `~/.goatnickels/config.json`.

Joining the network and downloading the blockchain is as easy as running `goatnickels -serve y`.

### Creating an account
```
// Make an account to store GoatNickels in
goatnickels -generate-acct y -save y
```
The `-save` flag will write the account to a keystore file at `~/.goatnickels/keystore.json` and update the account in `config.json`. **Make a backup of the keystore file!** If you lose your key, you will not be able to recover GoatNickels in the account.

All accounts start with `goat_` and are followed by a long hash. Raw account balances are in Pellets, the smallest unit of goat money. 1 Pellet = .00000001 GoatNickels.

### Mining GoatNickels
During each consensus voting round (every ten seconds), miners are rewarded GoatNickels based on the amount they have "staked" as a long-term deposit for the right to vote.

*Staking to be implemented*

## Block Data

Each block includes a data element. Data is made of two arrays: state (accounts) and transactions.

### Accounts

[{
  "account": {"balance": 123792167, "sequence": 3, "stake": 123793}
}]

### Transaction

[{
  "from": "HASH",
  "to": "HASH",
  "amount": 32134967423
  "sequence": 1,
  "r": "HASH",
  "s": "HASH"
}]

## Consensus Algorithm - PROBABLY SHIT

Block consensus is reached through a Byzantine Fault Tolerant Proof of Stake algorithm. It combines elements of the [Ripple Consensus Algorithm](https://ripple.com/files/ripple_consensus_whitepaper.pdf) and [Ethereum Proof of Stake](https://github.com/ethereum/wiki/wiki/Proof-of-Stake-FAQ) model. The intent is to process transactions as quickly as Ripple without requiring a trusted group of non-collaborating nodes.

### Voting Process

