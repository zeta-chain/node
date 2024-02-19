package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/fungible module sentinel errors
var (
	ErrABIPack                 = sdkerrors.Register(ModuleName, 1101, "abi pack error")
	ErrABIGet                  = sdkerrors.Register(ModuleName, 1102, "abi get error")
	ErrABIUnpack               = sdkerrors.Register(ModuleName, 1104, "abi unpack error")
	ErrContractNotFound        = sdkerrors.Register(ModuleName, 1107, "contract not found")
	ErrContractCall            = sdkerrors.Register(ModuleName, 1109, "contract call error")
	ErrSystemContractNotFound  = sdkerrors.Register(ModuleName, 1110, "system contract not found")
	ErrInvalidAddress          = sdkerrors.Register(ModuleName, 1111, "invalid address")
	ErrStateVariableNotFound   = sdkerrors.Register(ModuleName, 1112, "state variable not found")
	ErrEmitEvent               = sdkerrors.Register(ModuleName, 1114, "emit event error")
	ErrInvalidDecimals         = sdkerrors.Register(ModuleName, 1115, "invalid decimals")
	ErrInvalidGasLimit         = sdkerrors.Register(ModuleName, 1118, "invalid gas limit")
	ErrSetBytecode             = sdkerrors.Register(ModuleName, 1119, "set bytecode error")
	ErrInvalidContract         = sdkerrors.Register(ModuleName, 1120, "invalid contract")
	ErrPausedZRC20             = sdkerrors.Register(ModuleName, 1121, "ZRC20 is paused")
	ErrForeignCoinNotFound     = sdkerrors.Register(ModuleName, 1122, "foreign coin not found")
	ErrForeignCoinCapReached   = sdkerrors.Register(ModuleName, 1123, "foreign coin cap reached")
	ErrCallNonContract         = sdkerrors.Register(ModuleName, 1124, "can't call a non-contract address")
	ErrForeignCoinAlreadyExist = sdkerrors.Register(ModuleName, 1125, "foreign coin already exist")
	ErrInvalidHash             = sdkerrors.Register(ModuleName, 1126, "invalid hash")
	ErrNilGasPrice             = sdkerrors.Register(ModuleName, 1127, "nil gas price")
)
