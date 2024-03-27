# Messages

## MsgDeploySystemContracts

DeploySystemContracts deploy new instances of the system contracts

Authorized: admin policy group 2.

```proto
message MsgDeploySystemContracts {
	string creator = 1;
}
```

## MsgDeployFungibleCoinZRC20

DeployFungibleCoinZRC20 deploys a fungible coin from a connected chains as a ZRC20 on ZetaChain.

If this is a gas coin, the following happens:

* ZRC20 contract for the coin is deployed
* contract address of ZRC20 is set as a token address in the system
contract
* ZETA tokens are minted and deposited into the module account
* setGasZetaPool is called on the system contract to add the information
about the pool to the system contract
* addLiquidityETH is called to add liquidity to the pool

If this is a non-gas coin, the following happens:

* ZRC20 contract for the coin is deployed
* The coin is added to the list of foreign coins in the module's state

Authorized: admin policy group 2.

```proto
message MsgDeployFungibleCoinZRC20 {
	string creator = 1;
	string ERC20 = 2;
	int64 foreign_chain_id = 3;
	uint32 decimals = 4;
	string name = 5;
	string symbol = 6;
	coin.CoinType coin_type = 7;
	int64 gas_limit = 8;
}
```

## MsgRemoveForeignCoin

RemoveForeignCoin removes a coin from the list of foreign coins in the
module's state.

Authorized: admin policy group 2.

```proto
message MsgRemoveForeignCoin {
	string creator = 1;
	string name = 2;
}
```

## MsgUpdateSystemContract

UpdateSystemContract updates the system contract

```proto
message MsgUpdateSystemContract {
	string creator = 1;
	string new_system_contract_address = 2;
}
```

## MsgUpdateContractBytecode

UpdateContractBytecode updates the bytecode of a contract from the bytecode
of an existing contract Only a ZRC20 contract or the WZeta connector contract
can be updated IMPORTANT: the new contract bytecode must have the same
storage layout as the old contract bytecode the new contract can add new
variable but cannot remove any existing variable

Authozied: admin policy group 2

```proto
message MsgUpdateContractBytecode {
	string creator = 1;
	string contract_address = 2;
	string new_code_hash = 3;
}
```

## MsgUpdateZRC20WithdrawFee

UpdateZRC20WithdrawFee updates the withdraw fee and gas limit of a zrc20 token

```proto
message MsgUpdateZRC20WithdrawFee {
	string creator = 1;
	string zrc20_address = 2;
	string new_withdraw_fee = 6;
	string new_gas_limit = 7;
}
```

## MsgUpdateZRC20PausedStatus

UpdateZRC20PausedStatus updates the paused status of a ZRC20
The list of ZRC20s are either paused or unpaused

Authorized: admin policy group 1 (pausing), group 2 (pausing & unpausing)

```proto
message MsgUpdateZRC20PausedStatus {
	string creator = 1;
	string zrc20_addresses = 2;
	UpdatePausedStatusAction action = 3;
}
```

## MsgUpdateZRC20LiquidityCap

UpdateZRC20LiquidityCap updates the liquidity cap for a ZRC20 token.

Authorized: admin policy group 2.

```proto
message MsgUpdateZRC20LiquidityCap {
	string creator = 1;
	string zrc20_address = 2;
	string liquidity_cap = 3;
}
```

