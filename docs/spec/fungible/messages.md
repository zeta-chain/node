# Messages

## MsgDeploySystemContracts

```proto
message MsgDeploySystemContracts {
	string creator = 1;
}
```

## MsgDeployFungibleCoinZRC20

```proto
message MsgDeployFungibleCoinZRC20 {
	string creator = 1;
	string ERC20 = 2;
	int64 foreign_chain_id = 3;
	uint32 decimals = 4;
	string name = 5;
	string symbol = 6;
	pkg.coin.CoinType coin_type = 7;
	int64 gas_limit = 8;
}
```

## MsgRemoveForeignCoin

```proto
message MsgRemoveForeignCoin {
	string creator = 1;
	string name = 2;
}
```

## MsgUpdateSystemContract

```proto
message MsgUpdateSystemContract {
	string creator = 1;
	string new_system_contract_address = 2;
}
```

## MsgUpdateContractBytecode

```proto
message MsgUpdateContractBytecode {
	string creator = 1;
	string contract_address = 2;
	string new_code_hash = 3;
}
```

## MsgUpdateZRC20WithdrawFee

```proto
message MsgUpdateZRC20WithdrawFee {
	string creator = 1;
	string zrc20_address = 2;
	string new_withdraw_fee = 6;
	string new_gas_limit = 7;
}
```

## MsgUpdateZRC20LiquidityCap

```proto
message MsgUpdateZRC20LiquidityCap {
	string creator = 1;
	string zrc20_address = 2;
	string liquidity_cap = 3;
}
```

## MsgPauseZRC20

```proto
message MsgPauseZRC20 {
	string creator = 1;
	string zrc20_addresses = 2;
}
```

## MsgUnpauseZRC20

```proto
message MsgUnpauseZRC20 {
	string creator = 1;
	string zrc20_addresses = 2;
}
```

