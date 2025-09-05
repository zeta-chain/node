package types

// DONTCOVER

import (
	cosmoserrors "cosmossdk.io/errors"
)

// x/fungible module sentinel errors
var (
	ErrABIPack                 = cosmoserrors.Register(ModuleName, 1101, "abi pack error")
	ErrABIGet                  = cosmoserrors.Register(ModuleName, 1102, "abi get error")
	ErrABIUnpack               = cosmoserrors.Register(ModuleName, 1104, "abi unpack error")
	ErrContractNotFound        = cosmoserrors.Register(ModuleName, 1107, "contract not found")
	ErrContractCall            = cosmoserrors.Register(ModuleName, 1109, "contract call error")
	ErrSystemContractNotFound  = cosmoserrors.Register(ModuleName, 1110, "system contract not found")
	ErrInvalidAddress          = cosmoserrors.Register(ModuleName, 1111, "invalid address")
	ErrStateVariableNotFound   = cosmoserrors.Register(ModuleName, 1112, "state variable not found")
	ErrEmitEvent               = cosmoserrors.Register(ModuleName, 1114, "emit event error")
	ErrInvalidDecimals         = cosmoserrors.Register(ModuleName, 1115, "invalid decimals")
	ErrInvalidGasLimit         = cosmoserrors.Register(ModuleName, 1118, "invalid gas limit")
	ErrSetBytecode             = cosmoserrors.Register(ModuleName, 1119, "set bytecode error")
	ErrInvalidContract         = cosmoserrors.Register(ModuleName, 1120, "invalid contract")
	ErrPausedZRC20             = cosmoserrors.Register(ModuleName, 1121, "ZRC20 is paused")
	ErrForeignCoinNotFound     = cosmoserrors.Register(ModuleName, 1122, "foreign coin not found")
	ErrForeignCoinCapReached   = cosmoserrors.Register(ModuleName, 1123, "foreign coin cap reached")
	ErrCallNonContract         = cosmoserrors.Register(ModuleName, 1124, "cannot call a non-contract address")
	ErrForeignCoinAlreadyExist = cosmoserrors.Register(ModuleName, 1125, "foreign coin already exist")
	ErrNilGasPrice             = cosmoserrors.Register(ModuleName, 1127, "nil gas price")
	ErrAccountNotFound         = cosmoserrors.Register(ModuleName, 1128, "account not found")
	ErrGatewayContractNotSet   = cosmoserrors.Register(ModuleName, 1129, "gateway contract not set")
	ErrZRC20ZeroAddress        = cosmoserrors.Register(ModuleName, 1130, "ZRC20 address cannot be zero")
	ErrZRC20NotWhiteListed     = cosmoserrors.Register(ModuleName, 1131, "ZRC20 is not whitelisted")
	ErrZRC20NilABI             = cosmoserrors.Register(ModuleName, 1132, "ZRC20 ABI is nil")
	ErrZeroAddress             = cosmoserrors.Register(ModuleName, 1133, "address cannot be zero")
	ErrInvalidAmount           = cosmoserrors.Register(ModuleName, 1134, "invalid amount")
	ErrMaxSupplyReached        = cosmoserrors.Register(ModuleName, 1135, "max supply reached")
	ErrCallEvmWithData         = cosmoserrors.Register(
		ModuleName,
		1136,
		"contract call failed when calling EVM with data",
	)
	ErrDepositZetaToEvmAccount = cosmoserrors.Register(
		ModuleName,
		1137,
		"error depositing ZETA to users EVM account",
	)
	ErrDepositZetaToFungibleAccount = cosmoserrors.Register(
		ModuleName,
		1138,
		"error depositing ZETA to fungible module account",
	)
	ErrProcessDeposit = cosmoserrors.Register(ModuleName, 1139, "error processing deposit")
	ErrZeroBalance    = cosmoserrors.Register(ModuleName, 1140, "balance is zero")
	ErrFailedToBurn   = cosmoserrors.Register(ModuleName, 1141, "failed to burn coins")
	ErrGasLimitNotSet = cosmoserrors.Register(ModuleName, 1142, "gas limit not set")
)
