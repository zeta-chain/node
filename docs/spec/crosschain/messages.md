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
reached, the ballot is finalized. When a ballot is finalized, if the amount
of zeta minted does not match the outbound transaction amount an error is
thrown. If the amounts match, the outbound transaction hash and the "last
updated" timestamp are updated.

The transaction is proceeded to be finalized:

If the observation was successful, the status is changed from "pending
revert/outbound" to "reverted/mined". The difference between zeta burned
and minted is minted by the bank module and deposited into the module
account.

If the observation was unsuccessful, and if the status is "pending outbound",
prices and nonce are updated and the status is changed to "pending revert".
If the status was "pending revert", the status is changed to "aborted".

If there's an error in the finalization process, the CCTX status is set to
'aborted'.

After finalization the outbound transaction tracker and pending nonces are
removed, and the CCTX is updated in the store.

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

If the receiver chain is a ZetaChain, the EVM deposit is handled and the
status of CCTX is changed to "outbound mined". If EVM deposit handling fails,
the status of CCTX is chagned to 'aborted'.

If the receiver chain is a connected chain, the inbound CCTX is finalized
(prices and nonce are updated) and status changes to "pending outbound". If
the finalization fails, the status of CCTX is changed to 'aborted'.

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

