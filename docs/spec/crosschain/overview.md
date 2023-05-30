# Overview

The `crosschain` module tracks inbound and outbound cross-chain transactions
(CCTX).

The main actors interacting with the Crosschain module are observer validators
(or "observers"). Observers are running an off-chain program (called
`zetaclient`) that watches connected blockchains for inbound transactions and
watches ZetaChain for pending outbound transactions and watches connected chains
for outbound transactions.

After observing an inbound or an outbound transaction, an observer participates
in a voting process.

## Voting

When an observer submits a vote for a transaction a `ballot` is created (if it
wasn't created before). Observers are allowed to cast votes that will be
associated with this ballot. Based on `BallotThreshold`, when enough votes are
cast ballot is considered to be "finalized".

The last vote that moves the ballot to the "finalized" state triggers execution
and pays the gas costs of the cross-chain transaction.

Any votes cast after the ballot has been finalized are discarded.

## Inbound Transaction

Inbound transactions are cross-chain transactions observed on connected chains.
To vote on an inbound transaction an observer broadcasts
`MsgVoteOnObservedInboundTx`.

The last vote that moves the ballot to the "finalized" state triggers execution
of the cross-chain transaction.

If the destination chain is ZetaChain and the CCTX does not contain a message,
ZRC20 tokens are deposited into an account on ZetaChain.

If the destination chain is ZetaChain and the CCTX contains a message, ZRC20
tokens are deposited and a contract on ZetaChain is called. Contract address and
the arguments for the contract call are contained within the message.

If the destination chain is not ZetaChain, the status of a transaction is
changed to "pending outbound" and the CCTX to be processed as an outbound
transaction.

## Outbound Transaction

### Pending Outbound

Observers watch ZetaChain for pending outbound transactions. To process a
pending outbound transactions observers enter into a TSS keysign ceremony to
sign the transaction, and then broadcast the signed transaction to the connected
blockchains.

### Observed Outbound

Observers watch connected blockchains for the broadcasted outbound transactions.
Once a transaction is "confirmed" (or "mined") on a connected blockchains,
observers vote on ZetaChain by sending a `VoteOnObservedOutboundTx` message.

After the vote passes the threshold, the voting is finalized and a transaction's
status is changed to final.

## Permissions

| Message                     | Admin policy account | Observer validator |
| --------------------------- | -------------------- | ------------------ |
| MsgCreateTSSVoter           |                      | ✅                 |
| MsgGasPriceVoter            |                      | ✅                 |
| MsgVoteOnObservedOutboundTx |                      | ✅                 |
| MsgVoteOnObservedInboundTx  |                      | ✅                 |
| MsgAddToOutTxTracker        | ✅                   | ✅                 |
| MsgRemoveFromOutTxTracker   | ✅                   |                    |
| MsgUpdatePermissionFlags    | ✅                   |                    |
| MsgUpdateKeygenPermission   | ✅                   |                    |

## State

The module stores the following information in the state:

- List of outbound transactions
- List of chain nonces
- List of last chain heights
- List of cross-chain transactions
- List of
- Mapping between inbound transactions and cross-chain transactions
- Keygen
- TSS key
- Gas prices on connected chains submitted by observers
