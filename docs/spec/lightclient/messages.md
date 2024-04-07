# Messages

## MsgUpdateVerificationFlags

UpdateVerificationFlags updates the crosschain related flags.
Emergency group can disable flags while operation group can enable/disable

```proto
message MsgUpdateVerificationFlags {
	string creator = 1;
	VerificationFlags verification_flags = 2;
}
```

