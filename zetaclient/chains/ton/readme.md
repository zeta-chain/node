# TON Observer-Signer

This document provides an overview for fellow Zeta contributors on how the TON integration is implemented 🙌.

The core smart contracts are located in the [zeta-chain/protocol-contracts-ton](https://github.com/zeta-chain/protocol-contracts-ton) repository. 
For local dev and testing, TON is also supported in our [localnet](https://github.com/zeta-chain/localnet).

> ⚠️ Before proceeding, please familiarize yourself with TON's [basic concepts](https://docs.ton.org/v3/concepts/dive-into-ton/introduction). 
> Afterward, it's highly recommended to review the contracts repository and its test specs for practical examples.

## `Gateway`

The Gateway (GW) wrapper and transaction encoder/decoder are implemented in `pkg/contracts/ton`.

It is very similar to the TypeScript Gateway wrapper in the contracts repository.

## Observer-Signer

The logic is similar to other ZetaChain's observer-signers.

`ton.go` contains a high-level orchestrator that schedules different tasks for inbound and outbound transactions.

- `check_rpc_status` - checks RPC health
- `post_gas_price` - posts gas prices. They rarely change (only by governance proposals).
- `observe_inbound` - scrolls Gateway transactions, detects inbounds, and converts them into votes.
- `process_inbound_trackers` - same as other inbound trackers, but note that the hash here is `$logical_time:$hash`.
- `process_outbound_trackers` - same as other outbound trackers, but note that the hash here is `$logical_time:$hash`.
- `schedule_cctx` - lists pending outbound CCTXs, then constructs, signs, and broadcasts the corresponding TON transactions.

## Notes and Pitfalls

### Accounts

- All TON accounts are smart contracts. Wallets just store the key and sign hashes that are broadcasted to a contract.
- TON uses different addresses for mainnet/testnet and bouncable/non-bouncable addresses. A converter is available here: https://ton.org/address/
- We rely only on **raw** addresses (`$workchain:$bytes`). Example: `0:87115e4a012e747d9bce013ce2244010c6d5e3b0f88ddbc63420519b8619e5a0`

### Transactions

Technically, each "action" in TON is async, following the Actor Model. 
One logical action (e.g., swapping TON for USDT) is actually a chain of multiple messages between smart contracts, 
resulting in multiple transactions across different blocks and shards. 
This process might take a while to fully execute (sometimes 30+ seconds).

Even a simple TON transfer from Alice to Bob consists of two transactions: 
"Alice sends a message with 1 TON to Bob" and "Bob receives a message from Alice with 1 TON".

A cool property of this is that each "physical transaction" that mutates a single account has **instant finality**.

### Transaction Retrieval

**It's not possible to retrieve a transaction by its hash from RPC!** 
A transaction can only be uniquely identified and retrieved with 
a combination of `account_address` + `lt` (Logical Time) + `tx_hash`.

Another implication from TON's async nature is that because there are multiple in/out messages between
different accounts, we don't know the destination's transaction hash. So it's not possible to 
perform a classical flow of `tx = build(); send(tx); rpc.get(tx.hash())`.

As of now, we only rely on the Gateway's address (so the account is always known), but instead of the transaction hash, 
we use `$logical_time:$hash` in the cross-chain module to have the full arguments for retrieving a transaction.

Also, because EACH account is considered a "shard-chain", we can't filter by inbound or outbound transactions. 
This can be implemented only in the runtime, i.e., we simply scroll all new Gateway transactions, then parse them 
and determine whether this is an inbound, outbound, or other transaction.

This is why `observe_inbound` might also process a finalized withdrawal that was invoked
by another goroutine in `schedule_cctx`.

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

### RPC

For more details on the RPC API, 
see the official documentation: https://docs.ton.org/v3/guidelines/dapps/apis-sdks/api-types

TON provides a binary `lite-server` protocol for data access and transaction submission. 
The initial integration was implemented using a lite-client, but we eventually encountered infrastructure issues: 
not many RPC providers support lite-server and know how to operate it.

That's why we switched to the widely adopted toncenter-v2 RPC implementation. Technically, it is a Python webserver 
with a lite-client wrapper and a JSON API hosted by node providers.

The pruning period differs from provider to provider. Most non-archival nodes prune transactions 
after ~14 days, which is an important consideration for the observer 
(specifically for `process_inbound_trackers` and `process_outbound_trackers`)

### Gas

TON has a complex [gas model](https://docs.ton.org/v3/documentation/smart-contracts/transaction-fees/fees) 
with dynamic pricing. Also, all outbound operations are paid by the **contract**, i.e., the Gateway pays 
the gas fee for withdrawing to end recipients.

In order to keep gas calculation manageable, we measured the gas cost for all operations empirically and 
placed a "ceiling" that we treat as the gas fee. This is also a suggested approach by the TON team.

Also, the TON gas price can only be changed via a governance proposal, and all VM operations
have a predefined gas price, so it should not be an issue. This approach makes transactions slightly "overpriced," but the actual difference in tx cost is typically less than a cent ($0.01).

The Gateway has a getter `int calculate_gas_fee(int op) method_id` for retrieving the transaction gas cost in nanoTON.

See [gateway.fc](https://github.com/zeta-chain/protocol-contracts-ton/blob/main/contracts/gateway.fc#L34).

### Event Logs

There are no event logs in TON like in EVM. For `deposit` and `deposit_and_call`, 
we need an equivalent mechanism. This is achieved using a feature of TON called an "external message" 
— an outbound message with no destination. The Gateway wrapper's parsing logic relies on this to detect inbound events.

This is documented as a "raw message" or "external message" in the official TON documentation: 
https://docs.ton.org/v3/documentation/smart-contracts/message-management/sending-messages#types-of-messages

Let's take an example from `depositAndCall` operation:

```go
// taken from pkg/contracts/ton/gateway_parse.go

type Deposit struct {
  // Parsed from an incoming message (IM / internal message)
  Sender ton.AccountID 

  // Equals to depositLog.Amount
  Amount math.Uint

  // Parsed from message payload;
  // Payload is just a part of the same IM
  Recipient eth.Address 
}

type DepositAndCall struct {
  Deposit

  // Parsed from the payload using snakeCell encoding
  CallData []byte 
}

// represents outbound external message that acts as an event log
// It contains FACTUAL TON amount that we treat as CCTX deposit 
// and a deposit fee based on the operation (not used by zetaclient as of now)
type depositLog struct {
  Amount     math.Uint
  DepositFee math.Uint
}
```

### Signing Outbound Transactions

TON uses EdDSA for cryptography, but TSS uses ECDSA. TVM provides the `ECRECOVER` opcode that performs 
ECDSA recovery of the TSS signature. The `protocol-contracts-ton` repository explains how this 
signature scheme works in detail. The encoding is implemented in the Gateway wrapper within this repository.

- See [`func (w *Withdrawal) AsBody() (*boc.Cell, error)`](https://github.com/zeta-chain/node/blob/9fbdb7767674e12add976b1f25b2fba94c4361c3/pkg/contracts/ton/gateway_msg.go#L176)
- See [`async Gateway.sendTSSCommand`](https://github.com/zeta-chain/protocol-contracts-ton/blob/5898b04fffc937864f441623b45834a5250aa2e6/wrappers/Gateway.ts#L143)

We also implemented a sequential `seqno` in the Gateway contract that acts similarly to an EVM nonce.

