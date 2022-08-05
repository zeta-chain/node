package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/mirror module sentinel errors
var (
	ErrSample                 = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrABIPack                = sdkerrors.Register(ModuleName, 1101, "abi pack error")
	ErrABIUnpack              = sdkerrors.Register(ModuleName, 1102, "abi unpack error")
	ErrUnexpectedEvent        = sdkerrors.Register(ModuleName, 1103, "unexpected event")
	ErrTOkenPairAlreadyExists = sdkerrors.Register(ModuleName, 1104, "token pair already exists")
	ErrDeployERC20Mirror      = sdkerrors.Register(ModuleName, 1105, "deploy erc20 mirror error")
	ErrTokenPairNotFound      = sdkerrors.Register(ModuleName, 1106, "token pair not found")
	ErrInvalidAmount          = sdkerrors.Register(ModuleName, 1107, "invalid amount")
	ErrEVMCall                = sdkerrors.Register(ModuleName, 1108, "evm call error")
	ErrBalanceInvariance      = sdkerrors.Register(ModuleName, 1109, "balance invariance error")
)
