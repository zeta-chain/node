# Messages

## MsgEnableHeaderVerification

EnableHeaderVerification enables the verification flags for the given chain IDs
Enabled chains allow the submissions of block headers and using it to verify the correctness of proofs

```proto
message MsgEnableHeaderVerification {
	string creator = 1;
	int64 chain_id_list = 2;
}
```

## MsgDisableHeaderVerification

DisableHeaderVerification disables the verification flags for the given chain IDs
Disabled chains do not allow the submissions of block headers or using it to verify the correctness of proofs

```proto
message MsgDisableHeaderVerification {
	string creator = 1;
	int64 chain_id_list = 2;
}
```

