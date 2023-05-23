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

FIXME: use more specific error types & codes

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
