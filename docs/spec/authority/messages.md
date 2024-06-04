# Messages

## MsgUpdatePolicies

```proto
message MsgUpdatePolicies {
	string creator = 1;
	Policies policies = 2;
}
```

## MsgUpdateChainInfo

```proto
message MsgUpdateChainInfo {
	string creator = 1;
	ChainInfo chain_info = 2;
}
```

## MsgAddAuthorization

```proto
message MsgAddAuthorization {
	string creator = 1;
	string msg_url = 2;
	PolicyType authorized_policy = 3;
}
```

## MsgRemoveAuthorization

```proto
message MsgRemoveAuthorization {
	string creator = 1;
	string msg_url = 2;
}
```

