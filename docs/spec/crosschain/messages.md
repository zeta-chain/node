# Messages

## MsgAddToOutTxTracker

Adds a new record to the outbound transaction tracker.

Only the admin policy account and the observer validators are authorized to
broadcast this message.

```proto
message MsgAddToOutTxTracker {
	string creator = 1;
	int64 chain_id = 2;
	uint64 nonce = 3;
	string tx_hash = 4;
}
```

## MsgRemoveFromOutTxTracker

Removes a record from the outbound transaction tracker by chain ID and nonce.

Only the admin policy account is authorized to broadcast this message.

```proto
message MsgRemoveFromOutTxTracker {
	string creator = 1;
	int64 chain_id = 2;
	uint64 nonce = 3;
}
```

## MsgCreateTSSVoter

Vote on creating a TSS key and recording the information about it (public
key, participant and operator addresses, finalized and keygen heights).

If the vote passes, the information about the TSS key is recorded on chain
and the status of the keygen is set to "success".

Fails if the keygen does not exist, the keygen has been already
completed, or the keygen has failed.

Only node accounts are authorized to broadcast this message.

```proto
message MsgCreateTSSVoter {
	string creator = 1;
	string tss_pubkey = 2;
	int64 keyGenZetaHeight = 3;
	common.ReceiveStatus status = 4;
}
```

## MsgGasPriceVoter

Submit information about the connected chain's gas price at a specific block
height. Gas price submitted by each validator is recorded separately and a
median index is updated.

Only observer validators are authorized to broadcast this message.

```proto
message MsgGasPriceVoter {
	string creator = 1;
	int64 chain_id = 2;
	uint64 price = 3;
	uint64 block_number = 4;
	string supply = 5;
}
```

## MsgNonceVoter

Deprecated.

```proto
message MsgNonceVoter {
	string creator = 1;
	int64 chain_id = 2;
	uint64 nonce = 3;
}
```

## MsgVoteOnObservedOutboundTx

Casts a vote on an outbound transaction observed on a connected chain (after
it has been broadcasted to and finalized on a connected chain). If this is
the first vote, a new ballot is created. When a threshold of votes is
reached, the ballot is finalized. When a ballot is finalized, the outbound
transaction is processed.

If the observation is successful, the difference between zeta burned
and minted is minted by the bank module and deposited into the module
account.

If the observation was unsuccessful, the logic depends on the previous
status.

If the previous status was `PendingOutbound`, a new revert transaction is
created. To cover the revert transaction fee, the required amount of tokens
submitted with the CCTX are swapped using a Uniswap pool on ZetaChain for the
ZRC20 of the gas token of the receiver chain. The ZRC20 tokens are then
burned. The nonce is updated. If everything is successful, the CCTX status is
changed to `PendingRevert`.

If the previous status was `PendingRevert`, the CCTX is aborted.

```mermaid
stateDiagram-v2

	state finalize_outbound <<choice>>
	state observation <<choice>>
	state success_old_status <<choice>>
	state fail_old_status <<choice>>
	[*] --> finalize_outbound
	finalize_outbound --> observation: Finalize outbound
	observation --> success_old_status: Observation succeeded
	success_old_status --> Reverted: Old status is PendingRevert
	success_old_status --> OutboundMined: Old status is PendingOutbound
	observation --> fail_old_status: Observation failed
	fail_old_status --> PendingRevert: Old status is PendingOutbound
	fail_old_status --> Aborted: Old status is PendingRevert
	finalize_outbound --> Aborted: Finalize outbound error
	Aborted --> [*]
	Reverted --> [*]
	OutboundMined --> [*]
	PendingRevert --> [*]

```
Only observer validators are authorized to broadcast this message.

```proto
message MsgVoteOnObservedOutboundTx {
	string creator = 1;
	string cctx_hash = 2;
	string observed_outTx_hash = 3;
	uint64 observed_outTx_blockHeight = 4;
	string zeta_minted = 5;
	common.ReceiveStatus status = 6;
	int64 outTx_chain = 7;
	uint64 outTx_tss_nonce = 8;
	common.CoinType coin_type = 9;
}
```

## MsgVoteOnObservedInboundTx

Casts a vote on an inbound transaction observed on a connected chain. If this
is the first vote, a new ballot is created. When a threshold of votes is
reached, the ballot is finalized. When a ballot is finalized, a new CCTX is
created.

If the receiver chain is ZetaChain, `HandleEVMDeposit` is called. If the
tokens being deposited are ZETA, `MintZetaToEVMAccount` is called and the
tokens are minted to the receiver account on ZetaChain. If the tokens being
deposited are gas tokens or ERC20 of a connected chain, ZRC20's `deposit`
method is called and the tokens are deposited to the receiver account on
ZetaChain. If the message is not empty, system contract's `depositAndCall`
method is also called and an omnichain contract on ZetaChain is executed.
Omnichain contract address and arguments are passed as part of the message.
If everything is successful, the CCTX status is changed to `OutboundMined`.

If the receiver chain is a connected chain, the `FinalizeInbound` method is
called to prepare the CCTX to be processed as an outbound transaction. To
cover the outbound transaction fee, the required amount of tokens submitted
with the CCTX are swapped using a Uniswap pool on ZetaChain for the ZRC20 of
the gas token of the receiver chain. The ZRC20 tokens are then burned. The
nonce is updated. If everything is successful, the CCTX status is changed to
`PendingOutbound`.

```mermaid
stateDiagram-v2

	state is_zetachain <<choice>>
	state evm_deposit_success <<choice>>
	state finalize_inbound <<choice>>
	[*] --> PendingInbound: Create New CCTX
	PendingInbound --> is_zetachain
	is_zetachain --> evm_deposit_success: Receiver is ZetaChain
	evm_deposit_success --> OutboundMined: EVM deposit success
	evm_deposit_success --> Aborted: EVM deposit error
	is_zetachain --> finalize_inbound: Receiver is connected chain
	finalize_inbound --> Aborted: Finalize inbound error
	finalize_inbound --> PendingOutbound: Finalize inbound success
	Aborted --> [*]
	PendingOutbound --> [*]
	OutboundMined --> [*]

```

Only observer validators are authorized to broadcast this message.

```proto
message MsgVoteOnObservedInboundTx {
	string creator = 1;
	string sender = 2;
	int64 sender_chain_id = 3;
	string receiver = 4;
	int64 receiver_chain = 5;
	string amount = 6;
	string message = 8;
	string in_tx_hash = 9;
	uint64 in_block_height = 10;
	uint64 gas_limit = 11;
	common.CoinType coin_type = 12;
	string tx_origin = 13;
	string asset = 14;
}
```

## MsgSetNodeKeys

Not implemented yet.

```proto
message MsgSetNodeKeys {
	string creator = 1;
	common.PubKeySet pubkeySet = 2;
	string tss_signer_Address = 3;
}
```

## MsgUpdatePermissionFlags

Updates permissions. Currently, this is only used to enable/disable the
inbound transactions.

Only the admin policy account is authorized to broadcast this message.

```proto
message MsgUpdatePermissionFlags {
	string creator = 1;
	bool isInboundEnabled = 3;
}
```

## MsgUpdateKeygen

Updates the block height of the keygen and sets the status to "pending
keygen".

Only the admin policy account is authorized to broadcast this message.

```proto
message MsgUpdateKeygen {
	string creator = 1;
	int64 block = 2;
}
```

