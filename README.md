# GoatNickels - The Shitcoin

I'm building this to learn Go and Blockchain concepts

## Block Data

Each block includes a data element. Data is made of two arrays: accounts and transactions.

### Accounts

[{
  "account": "HASH",
  "balance": 123792167
}]

### Transactions

[{
  "from": "HASH",
  "to": "HASH",
  "amount": 32134967423
}]

## Consensus Algorithm - TO BE IMPLEMENTED

Block consensus is reached through a Byzantine Fault Tolerant Proof of Stake algorithm. It combines elements of the [Ripple Consensus Algorithm](https://ripple.com/files/ripple_consensus_whitepaper.pdf) and [Ethereum Proof of Stake](https://github.com/ethereum/wiki/wiki/Proof-of-Stake-FAQ) model. The intent is to process transactions as quickly as Ripple without requiring a trusted group of non-collaborating nodes.

Votes scale by the amount staked as a percentage of total stakes and by successful block and transaction contribution history(?).

## GoatSatoshis

Account balances are in Pellets, the smallest unit of goat money.

1 Pellet = .00000001 GoatNickels

## Using GoatNickels...

TBD