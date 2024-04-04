package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrBlockAlreadyExist               = errorsmod.Register(ModuleName, 1101, "block already exists")
	ErrNoParentHash                    = errorsmod.Register(ModuleName, 1102, "no parent hash")
	ErrInvalidTimestamp                = errorsmod.Register(ModuleName, 1103, "invalid timestamp")
	ErrBlockHeaderVerificationDisabled = errorsmod.Register(ModuleName, 1104, "block header verification is disabled")
	ErrorVerificationFlagsNotFound     = errorsmod.Register(ModuleName, 1105, "verification flags not found")
)
