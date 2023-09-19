# Messages

## MsgAddObserver

AddObserver adds in a new observer to the store.It can be executed using an admin policy account
Once added, the function also resets keygen and pauses inbound so that a new TSS can be generated.

```proto
message MsgAddObserver {
	string creator = 1;
	string observer_address = 2;
	string zetaclient_grantee_pubkey = 3;
	bool add_node_account_only = 4;
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
	bool isOutboundEnabled = 4;
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

## MsgAddBlockHeader

MsgAddBlockHeader handles adding a block header to the store, through majority voting of observers

```proto
message MsgAddBlockHeader {
	string creator = 1;
	int64 chain_id = 2;
	bytes block_hash = 3;
	int64 height = 4;
	HeaderData header = 5;
}
```

