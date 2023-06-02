# Messages

## MsgAddObserver

```proto
message MsgAddObserver {
	string creator = 1;
	int64 chain_id = 2;
	ObservationType observationType = 3;
}
```

## MsgUpdateCoreParams

```proto
message MsgUpdateCoreParams {
	string creator = 1;
	CoreParams coreParams = 2;
}
```

