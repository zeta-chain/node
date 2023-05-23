# Messages

## MsgDeployFungibleCoinZRC20

```proto
message MsgDeployFungibleCoinZRC20 {
	string creator = 1;
	string ERC20 = 2;
	string foreignChain = 3;
	uint32 decimals = 4;
	string name = 5;
	string symbol = 6;
	common.CoinType coinType = 7;
	int64 gasLimit = 8;
}
```

## MsgRemoveForeignCoin

```proto
message MsgRemoveForeignCoin {
	string creator = 1;
	string name = 2;
}
```

