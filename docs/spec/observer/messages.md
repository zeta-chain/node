# Messages

## MsgAddObserver

Authorized: admin policy group 2.

```proto
message MsgAddObserver {
	string creator = 1;
	string observer_address = 2;
	string zetaclient_grantee_pubkey = 3;
	bool add_node_account_only = 4;
}
```

## MsgUpdateObserver

Authorized: admin policy group 2 (admin update), old observer address (if the
reason is that the observer was tombstoned).

```proto
message MsgUpdateObserver {
	string creator = 1;
	string old_observer_address = 2;
	string new_observer_address = 3;
	ObserverUpdateReason update_reason = 4;
}
```

## MsgUpdateCoreParams

UpdateCoreParams updates core parameters for a specific chain. Core parameters include
confirmation count, outbound transaction schedule interval, ZETA token,
connector and ERC20 custody contract addresses, etc.

Throws an error if the chain ID is not supported.

Only the admin policy account is authorized to broadcast this message.

```proto
message MsgUpdateCoreParams {
	string creator = 1;
	CoreParams coreParams = 2;
}
```

## MsgAddBlameVote

```proto
message MsgAddBlameVote {
	string creator = 1;
	int64 chain_id = 2;
	Blame blame_info = 3;
}
```

## MsgUpdateCrosschainFlags

UpdateCrosschainFlags updates the crosschain related flags.

Aurthorized: admin policy group 1 (except enabling/disabled
inbounds/outbounds and gas price increase), admin policy group 2 (all).

```proto
message MsgUpdateCrosschainFlags {
	string creator = 1;
	bool isInboundEnabled = 3;
	bool isOutboundEnabled = 4;
	GasPriceIncreaseFlags gasPriceIncreaseFlags = 5;
	BlockHeaderVerificationFlags blockHeaderVerificationFlags = 6;
}
```

## MsgUpdateKeygen

UpdateKeygen updates the block height of the keygen and sets the status to
"pending keygen".

Authorized: admin policy group 1.

```proto
message MsgUpdateKeygen {
	string creator = 1;
	int64 block = 2;
}
```

## MsgAddBlockHeader

AddBlockHeader handles adding a block header to the store, through majority voting of observers

```proto
message MsgAddBlockHeader {
	string creator = 1;
	int64 chain_id = 2;
	bytes block_hash = 3;
	int64 height = 4;
	common.HeaderData header = 5;
}
```

