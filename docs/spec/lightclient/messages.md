# Messages

## MsgUpdateVerificationFlags

UpdateVerificationFlags updates the light client verification flags.
This disables/enables blocks verification of the light client for the specified chain.
Emergency group can disable flags while operation group can enable/disable

```proto
message MsgUpdateVerificationFlags {
	string creator = 1;
	VerificationFlags verification_flags = 2;
}
```

