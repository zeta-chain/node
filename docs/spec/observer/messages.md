# Messages

## MsgAddObserver

Not implemented.

```proto
message MsgAddObserver {
	string creator = 1;
	int64 chain_id = 2;
	ObservationType observationType = 3;
}
```

## MsgUpdateCoreParams

Updates core parameters for a specific chain. Core parameters include
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

