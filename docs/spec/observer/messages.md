# Messages

## MsgAddObserver

```proto
message MsgAddObserver {
	string creator = 1;
	string observer_address = 2;
	string zetaclient_grantee_pubkey = 3;
	bool add_node_account_only = 4;
}
```

## MsgUpdateObserver

```proto
message MsgUpdateObserver {
	string creator = 1;
	string old_observer_address = 2;
	string new_observer_address = 3;
	ObserverUpdateReason update_reason = 4;
}
```

## MsgUpdateChainParams

```proto
message MsgUpdateChainParams {
	string creator = 1;
	ChainParams chainParams = 2;
}
```

## MsgRemoveChainParams

```proto
message MsgRemoveChainParams {
	string creator = 1;
	int64 chain_id = 2;
}
```

## MsgVoteBlame

```proto
message MsgVoteBlame {
	string creator = 1;
	int64 chain_id = 2;
	Blame blame_info = 3;
}
```

## MsgUpdateKeygen

```proto
message MsgUpdateKeygen {
	string creator = 1;
	int64 block = 2;
}
```

## MsgVoteBlockHeader

```proto
message MsgVoteBlockHeader {
	string creator = 1;
	int64 chain_id = 2;
	bytes block_hash = 3;
	int64 height = 4;
	pkg.proofs.HeaderData header = 5;
}
```

## MsgResetChainNonces

```proto
message MsgResetChainNonces {
	string creator = 1;
	int64 chain_id = 2;
	int64 chain_nonce_low = 3;
	int64 chain_nonce_high = 4;
}
```

## MsgVoteTSS

```proto
message MsgVoteTSS {
	string creator = 1;
	string tss_pubkey = 2;
	int64 keygen_zeta_height = 3;
	pkg.chains.ReceiveStatus status = 4;
}
```

## MsgEnableCCTX

```proto
message MsgEnableCCTX {
	string creator = 1;
	bool enableInbound = 2;
	bool enableOutbound = 3;
}
```

## MsgDisableCCTX

```proto
message MsgDisableCCTX {
	string creator = 1;
	bool disableInbound = 2;
	bool disableOutbound = 3;
}
```

## MsgUpdateGasPriceIncreaseFlags

```proto
message MsgUpdateGasPriceIncreaseFlags {
	string creator = 1;
	GasPriceIncreaseFlags gasPriceIncreaseFlags = 2;
}
```

