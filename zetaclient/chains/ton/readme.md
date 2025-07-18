# TON Observer-Signer

This document contains some notes for fellow Zeta contributors on how this integration is implemented 🙌

Contracts are located here: [zeta-chain/protocol-contracts-ton](https://github.com/zeta-chain/protocol-contracts-ton)

> ⚠️ Please check contracts repo first with the test spec as it contains real-world examples.

> ⚠️ Please read some TON basics before moving to this document

## `Gateway{}`

Gateway (GW) wrapper and transaction encoder/decoder are implemented in `pkg/contracts/ton`. 

It is very similar to TypeScript Gateway wrapper in contracts repo.

## Observer-Signer

The logic is similar to other Zeta's observer-signers. 

`ton.go` contains high-level orchestrator that scheduled different tasks for inbound/outbound txs.

- `check_rpc_status` - checks RPC health
- `post_gas_price` - posts gas prices. They rarely change (only by gov. proposals)
- `observe_inbound` - scroll GW transaction, detect inbounds, and convert them to votes 
- `process_inbound_trackers` - same as other inbound trackers, but note that hash here is `$logical_time:$hash`
- `process_outbound_trackers` - same as other outbound tracker, but note that hash here is `$logical_time:$hash`
- `schedule_cctx` - lists outbound pending CCTX, constructs, signs, and broadcasts TON txs

# Notes/Pitfalls

## Accounts

- All TON accounts are smart-contracts. Wallet's just store the key & sign hashes that are broadcasted to a contract
- TON uses different addresses for mainnet/testnet & bouncable/nonbouncable. Converter: https://ton.org/address/

## Transactions

Technically, each "action" in TON is async. One logical action (eg swap TON for USDT) is actually multiple transactions in different blocks & shards that might take a while to execute (sometimes 30+ seconds)

Even a simple TON transfer from Alice to Bob is two txs: "Alice sends a message with TON to Bob" and "Bob receives a message from Alice with TON".

The implication is that each "physical tx" that mutates a single account has **instant finality**

## Transaction Retrieval

**It's not possible to get tx by its hash from RPC!** A tx can be retrieved with a combination of `account` + `tx_hash` + `logical_time` (think of some virtual time in TON's distributed system)

Another implication from TON's async nature is that because there are multiple in/out messages between different accounts, we don't know the destination's tx hash! So it's not possible to perform a classical flow of `tx = build(); send(tx); rpc.get(tx.hash())`.

As of now, we only rely on Gateway's address (so the account is always known) but instead of tx hash we use `$logical_time:$hash` in cross-chain module to have a full args for retrieving a tx.

Also, because EACH account is considered as a "shard-chain", we can't filter by inbound or outbound transactions. This can be implemented only in the runtime i.e. we simply scroll all new Gateway transactions, then parse then and determine whether this is an inbound / outbound / other tx...

That is why in `observe_inbound` we might also process finalized withdrawal that was invoked by another goroutine in `schedule_cctx` 

```sh
# https://athens.explorer.zetachain.com/cc/tx/0x865e8bf2292872a5b5cc7dacc45739812ee37b8db03fa4dcc5b1765c6b48c17f
 ▸ ~ ▸ zq crosschain show-cctx 0x865e8bf2292872a5b5cc7dacc45739812ee37b8db03fa4dcc5b1765c6b48c17f -o json  | jq '.CrossChainTx.outbound_params[0].hash'
"36905041000001:a0a423d045e3b4099e28957eff5a9f893f69edd649a82fb35154b76770961812"
```

```sh
curl -s \
  'https://testnet.toncenter.com/api/v2/getTransactions?address=0%3A87115e4a012e747d9bce013ce2244010c6d5e3b0f88ddbc63420519b8619e5a0&limit=1&lt=36905041000001&hash=a0a423d045e3b4099e28957eff5a9f893f69edd649a82fb35154b76770961812&to_lt=0&archival=false' \
  -H 'accept: application/json' | jq '.result[0]'
```

```json
{
  "@type": "raw.transaction",
  "address": {
    "@type": "accountAddress",
    "account_address": "EQCHEV5KAS50fZvOATziJEAQxtXjsPiN28Y0IFGbhhnloNo_"
  },
  "utime": 1752761575,
  "data": "te6cckECCwEAAjYAA7V4cRXkoBLnR9m84BPOIkQBDG1eOw+I3bxjQgUZuGGeWgAAAhkJ/b9kFgGCyvAl9jda3ItoTycX39lCJe2GwVPFslGWm0IZm6qAAAIZCXg7tDaHkE5wADRqRBGoAQIDAgHgBAUAgnKcldazix1TmERB1tiV0A+WMYVRApEowcQ7WArpgXYfqV0RG1kNUijQrrfLnkkxEXVEnZ96THWaSa6xlhz/xooNAg8MZoYaGzWEQAkKAUWIAQ4ivJQCXOj7N5wCecRIgCGNq8dh8Ru3jGhAozcMM8tADAYBAd8IAYIAIoOFv9Gx0Koa2WizcV0st+93XMChZpWkgpwcMIkRmFVk7Oo3qtvpbBHV10s1eUB2MYDy+6me9bFSy+V9eD2ySgcAWQAAAMiACQEvj6WVYf53b3u07RPkruiGPJGOdn2I7DWiSRmO8YUnMS0AAAAA5wCvSAEOIryUAlzo+zecAnnESIAhjavHYfEbt4xoQKM3DDPLQQASAl8fSyrD/O7e92naJ8ld0Qx5Ixzs+xHYa0SSMx3jCk5iWgAGCCNaAABDIT+37ITQ8gnOQACdRWRjE4gAAAAAAAAAADmAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIABvyYMNQEwII0wAAAAAAAIAAAAAAALuX5dqZdWjN6Ey2poojJI0wR15GakQiu5sVlZczAYyoEBQFczBfHAj",
  "transaction_id": {
    "@type": "internal.transactionId",
    "lt": "36905041000001",
    "hash": "oKQj0EXjtAmeKJV+/1qfiT9p7dZJqC+zUVS3Z3CWGBI="
  },
  "fee": "5648954",
  "storage_fee": "154",
  "other_fee": "5648800",
  "in_msg": {
    "@type": "raw.message",
    "hash": "nx4neUgZFUqNkKalyE1u8G7ROtRjSLCF3tLcReRFeYE=",
    "source": "",
    "destination": "EQCHEV5KAS50fZvOATziJEAQxtXjsPiN28Y0IFGbhhnloNo_",
    "value": "0",
    "extra_currencies": [],
    "fwd_fee": "0",
    "ihr_fee": "0",
    "created_lt": "0",
    "body_hash": "MzaTC6Ttqr2BfSSefzv7OWxFD7gH4qm2c4egDAAQQvY=",
    "msg_data": {
      "@type": "msg.dataRaw",
      "body": "te6cckEBAgEAcwABggAig4W/0bHQqhrZaLNxXSy373dcwKFmlaSCnBwwiRGYVWTs6jeq2+lsEdXXSzV5QHYxgPL7qZ71sVLL5X14PbJKAQBZAAAAyIAJAS+PpZVh/ndve7TtE+Su6IY8kY52fYjsNaJJGY7xhScxLQAAAADn/KKzlw==",
      "init_state": ""
    },
    "message": "ACKDhb/RsdCqGtlos3FdLLfvd1zAoWaVpIKcHDCJEZhVZOzqN6rb6WwR1ddLNXlAdjGA8vupnvWxUsvlfXg9sko="
  },
  "out_msgs": [
    {
      "@type": "raw.message",
      "hash": "FM39MwV9b+0kSTGnXtx6e0vf3DLvKSNoSGdl9N44pzw=",
      "source": "EQCHEV5KAS50fZvOATziJEAQxtXjsPiN28Y0IFGbhhnloNo_",
      "destination": "EQBICXx9LKsP87t73adonyV3RDHkjHOz7EdhrRJIzHeMKcF4",
      "value": "10000000",
      "extra_currencies": [],
      "fwd_fee": "266669",
      "ihr_fee": "0",
      "created_lt": "36905041000002",
      "body_hash": "lqKW0iTyhcZ77pPDD4owkVfw2qNdxbh+QQt4YwoJz8c=",
      "msg_data": {
        "@type": "msg.dataRaw",
        "body": "te6cckEBAQEAAgAAAEysuc0=",
        "init_state": ""
      },
      "message": ""
    }
  ]
}
```

## RPC

https://docs.ton.org/v3/guidelines/dapps/apis-sdks/api-types

TON provides a binary `lite-server` protocol for data access and tx submission. Initial integration was indeed implemented using lite-client, but eventually we encountered infra issues: not many RPC providers support lite-server & know how to operate one.

That's why we switched to RPC (toncenter-v2) implementation what is widely adopted. Technically, it's just a python webserver with lite-client wrapper & json api hosted by node providers.

Pruning period differs from provider to provider. 

I assume most non-archival nodes prune txs after ~14 days.

## Gas

TON has a complex [gas model](https://docs.ton.org/v3/documentation/smart-contracts/transaction-fees/fees) with dynamic pricing. Also, all outbound operations are paid by the **contract** i.e. Gateway pays gas fee for withdrawing to end recipients.

In order to make gas calculation sane, we measured gas cost for all operations empirically and placed a "ceiling" that we treat as gas fee. This is also a suggested approach by the TON team.

Also, TON gas price can be only changed via a gov. proposal, and all VM operations have predefined gas price, so it should not be an issue. Yes, we make txs slightly overpriced due to such approximation, but the difference is less than a cent.

GW has a getter `int calculate_gas_fee(int op) method_id` for retrieving tx gas cost in nanoTON.

See [gateway.fc](https://github.com/zeta-chain/protocol-contracts-ton/blob/main/contracts/gateway.fc#L34)

## Event Logs

There are no event logs in TON like in EVM. But for `deposit` and `deposit_and_call` we need something similar, thus we use "external log message". Technically, it is just an outbound message w/o destination. GW wrapper depends on it.

https://docs.ton.org/v3/documentation/smart-contracts/message-management/sending-messages#types-of-messages

## Withdrawals

TON uses EdDSA for cryptography, but TSS uses ECDSA. TVM provides `ECRECOVER` op code that perform ECDSA recovery of TSS signature. Protocol-contract-ton repo explains how exactly such signature works. Encoding is implement in in Gateway's wrapper in this repo

We also implemented sequential `seqno` in the Gateway contract that acts similarly to EVM's nonce.
