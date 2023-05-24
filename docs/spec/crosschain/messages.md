# Messages

## MsgAddToOutTxTracker

```proto
message MsgAddToOutTxTracker {
	string creator = 1;
	int64 chain_id = 2;
	uint64 nonce = 3;
	string tx_hash = 4;
}
```

## MsgRemoveFromOutTxTracker

```proto
message MsgRemoveFromOutTxTracker {
	string creator = 1;
	int64 chain_id = 2;
	uint64 nonce = 3;
}
```

## MsgCreateTSSVoter

```proto
message MsgCreateTSSVoter {
	string creator = 1;
	string tss_pubkey = 2;
	int64 keyGenZetaHeight = 3;
	common.ReceiveStatus status = 4;
}
```

## MsgGasPriceVoter

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

Should be removed

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
'Aborted'.

After finalization the outbound transaction tracker and pending nonces are
removed, and the CCTX is updated in the store.

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
created. If the receiver chain is a ZetaChain, the EVM deposit is handled.
If the receiver chain is a connected chain, a swap is performed and the
CCTX is finalized.

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

```proto
message MsgSetNodeKeys {
	string creator = 1;
	common.PubKeySet pubkeySet = 2;
	string tss_signer_Address = 3;
}
```

## MsgUpdatePermissionFlags

```proto
message MsgUpdatePermissionFlags {
	string creator = 1;
	bool isInboundEnabled = 3;
}
```

## MsgUpdateKeygen

```proto
message MsgUpdateKeygen {
	string creator = 1;
	int64 block = 2;
}
```

