# Messages

## MsgAddOutboundTracker

```proto
message MsgAddOutboundTracker {
	string creator = 1;
	int64 chain_id = 2;
	uint64 nonce = 3;
	string tx_hash = 4;
	pkg.proofs.Proof proof = 5;
	string block_hash = 6;
	int64 tx_index = 7;
}
```

## MsgAddInboundTracker

```proto
message MsgAddInboundTracker {
	string creator = 1;
	int64 chain_id = 2;
	string tx_hash = 3;
	pkg.coin.CoinType coin_type = 4;
	pkg.proofs.Proof proof = 5;
	string block_hash = 6;
	int64 tx_index = 7;
}
```

## MsgRemoveOutboundTracker

```proto
message MsgRemoveOutboundTracker {
	string creator = 1;
	int64 chain_id = 2;
	uint64 nonce = 3;
}
```

## MsgVoteGasPrice

```proto
message MsgVoteGasPrice {
	string creator = 1;
	int64 chain_id = 2;
	uint64 price = 3;
	uint64 block_number = 4;
	string supply = 5;
}
```

## MsgVoteOutbound

```proto
message MsgVoteOutbound {
	string creator = 1;
	string cctx_hash = 2;
	string observed_outbound_hash = 3;
	uint64 observed_outbound_block_height = 4;
	uint64 observed_outbound_gas_used = 10;
	string observed_outbound_effective_gas_price = 11;
	uint64 observed_outbound_effective_gas_limit = 12;
	string value_received = 5;
	pkg.chains.ReceiveStatus status = 6;
	int64 outbound_chain = 7;
	uint64 outbound_tss_nonce = 8;
	pkg.coin.CoinType coin_type = 9;
}
```

## MsgVoteInbound

```proto
message MsgVoteInbound {
	string creator = 1;
	string sender = 2;
	int64 sender_chain_id = 3;
	string receiver = 4;
	int64 receiver_chain = 5;
	string amount = 6;
	string message = 8;
	string inbound_hash = 9;
	uint64 inbound_block_height = 10;
	uint64 gas_limit = 11;
	pkg.coin.CoinType coin_type = 12;
	string tx_origin = 13;
	string asset = 14;
	uint64 event_index = 15;
}
```

## MsgWhitelistERC20

```proto
message MsgWhitelistERC20 {
	string creator = 1;
	string erc20_address = 2;
	int64 chain_id = 3;
	string name = 4;
	string symbol = 5;
	uint32 decimals = 6;
	int64 gas_limit = 7;
}
```

## MsgUpdateTssAddress

```proto
message MsgUpdateTssAddress {
	string creator = 1;
	string tss_pubkey = 2;
}
```

## MsgMigrateTssFunds

```proto
message MsgMigrateTssFunds {
	string creator = 1;
	int64 chain_id = 2;
	string amount = 3;
}
```

## MsgAbortStuckCCTX

```proto
message MsgAbortStuckCCTX {
	string creator = 1;
	string cctx_index = 2;
}
```

## MsgRefundAbortedCCTX

```proto
message MsgRefundAbortedCCTX {
	string creator = 1;
	string cctx_index = 2;
	string refund_address = 3;
}
```

## MsgUpdateRateLimiterFlags

```proto
message MsgUpdateRateLimiterFlags {
	string creator = 1;
	RateLimiterFlags rate_limiter_flags = 2;
}
```

