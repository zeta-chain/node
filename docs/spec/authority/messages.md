# Messages

## MsgUpdatePolicies

UpdatePolicies updates policies

```proto
message MsgUpdatePolicies {
	string signer = 1;
	Policies policies = 2;
}
```

## MsgUpdateChainInfo

UpdateChainInfo updates the chain inffo structure that adds new static chain info or overwrite existing chain info
on the hard-coded chain info

```proto
message MsgUpdateChainInfo {
	string signer = 1;
}
```

