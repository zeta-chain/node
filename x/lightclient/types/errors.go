package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrBlockAlreadyExist               = errorsmod.Register(ModuleName, 1101, "block already exists")
	ErrNoParentHash                    = errorsmod.Register(ModuleName, 1102, "no parent hash")
	ErrInvalidTimestamp                = errorsmod.Register(ModuleName, 1103, "invalid timestamp")
	ErrBlockHeaderVerificationDisabled = errorsmod.Register(ModuleName, 1104, "block header verification is disabled")
	ErrVerificationFlagsNotFound       = errorsmod.Register(ModuleName, 1105, "verification flags not found")
	ErrChainNotSupported               = errorsmod.Register(ModuleName, 1106, "chain not supported")
	ErrInvalidBlockHash                = errorsmod.Register(ModuleName, 1107, "invalid block hash")
	ErrBlockHeaderNotFound             = errorsmod.Register(ModuleName, 1108, "block header not found")
	ErrProofVerificationFailed         = errorsmod.Register(ModuleName, 1109, "proof verification failed")
	ErrInvalidHeight                   = errorsmod.Register(ModuleName, 1110, "invalid height")
	ErrInvalidBlockHeader              = errorsmod.Register(ModuleName, 1111, "invalid block header")
)
