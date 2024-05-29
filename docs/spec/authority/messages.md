# Messages

## MsgUpdatePolicies

UpdatePolicies updates policies

```proto
message MsgUpdatePolicies {
	string signer = 1;
	Policies policies = 2;
}
```

## MsgUpdateAuthorizations

```proto
message MsgUpdateAuthorizations {
	string signer = 1;
	AuthorizationList add_authorization_list = 2;
	AuthorizationList remove_authorization_list = 3;
}
```

