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

