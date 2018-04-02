# GoatNickels - The Shitcoin

GoatNickels is a blockchain-based cryptocurrency implemented in Go. This is the server and mining application. The wallet app is called [GoatBucket](https://github.com/seanmclane/goatbucket).

## Quick Start Guide

See the GoatNickels Quick Start Guide [here](goatnickels.com/gn-quick-start-guide).

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

