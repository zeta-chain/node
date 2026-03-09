package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrUnableToAddVote = errorsmod.Register(ModuleName, 1100, "unable to add vote ")
	ErrParamsThreshold = errorsmod.Register(ModuleName, 1101, "threshold cannot be more than 1")
	ErrSupportedChains = errorsmod.Register(ModuleName, 1102, "chain not supported")
	ErrInvalidStatus   = errorsmod.Register(ModuleName, 1103, "invalid Voting Status")

	ErrNotValidator = errorsmod.Register(
		ModuleName,
		1106,
		"user needs to be a validator before applying to become an observer",
	)
	ErrValidatorStatus = errorsmod.Register(
		ModuleName,
		1107,
		"corresponding validator needs to be bonded and not jailed",
	)
	ErrInvalidAddress = errorsmod.Register(ModuleName, 1108, "invalid Address")
	ErrSelfDelegation = errorsmod.Register(ModuleName, 1109, "self Delegation for operator not found")
	ErrKeygenNotFound = errorsmod.Register(
		ModuleName,
		1113,
		"Keygen not found, Keygen block can only be updated,New keygen cannot be set",
	)
	ErrKeygenBlockTooLow = errorsmod.Register(
		ModuleName,
		1114,
		"please set a block number at-least 10 blocks higher than the current block number",
	)
	ErrKeygenCompleted = errorsmod.Register(ModuleName, 1115, "keygen already completed")

	ErrLastObserverCountNotFound   = errorsmod.Register(ModuleName, 1123, "last observer count not found")
	ErrUpdateObserver              = errorsmod.Register(ModuleName, 1124, "unable to update observer")
	ErrNodeAccountNotFound         = errorsmod.Register(ModuleName, 1125, "node account not found")
	ErrInvalidChainParams          = errorsmod.Register(ModuleName, 1126, "invalid chain params")
	ErrChainParamsNotFound         = errorsmod.Register(ModuleName, 1127, "chain params not found")
	ErrParamsMinObserverDelegation = errorsmod.Register(ModuleName, 1128, "min observer delegation cannot be nil")
	ErrMinDelegationNotFound       = errorsmod.Register(ModuleName, 1129, "min delegation not found")
	ErrObserverSetNotFound         = errorsmod.Register(ModuleName, 1130, "observer set not found")
	ErrTssNotFound                 = errorsmod.Register(ModuleName, 1131, "tss not found")

	ErrInboundDisabled = errorsmod.Register(
		ModuleName,
		1132,
		"inbound tx processing is disabled",
	)
	ErrInvalidZetaCoinTypes                  = errorsmod.Register(ModuleName, 1133, "invalid zeta coin types")
	ErrNotObserver                           = errorsmod.Register(ModuleName, 1134, "sender is not an observer")
	ErrDuplicateObserver                     = errorsmod.Register(ModuleName, 1135, "observer already exists")
	ErrObserverNotFound                      = errorsmod.Register(ModuleName, 1136, "observer not found")
	ErrInvalidObserverAddress                = errorsmod.Register(ModuleName, 1137, "invalid observer address")
	ErrOperationalFlagsRestartHeightNegative = errorsmod.Register(
		ModuleName,
		1138,
		"restart height cannot be negative",
	)
	ErrOperationalFlagsSignerBlockTimeOffsetNegative = errorsmod.Register(
		ModuleName,
		1139,
		"signer block time offset cannot be negative",
	)
	ErrOperationalFlagsSignerBlockTimeOffsetLimit = errorsmod.Register(
		ModuleName,
		1140,
		"signer block time offset exceeds limit",
	)
	ErrOperationalFlagsInvalidMinimumVersion = errorsmod.Register(
		ModuleName,
		1141,
		"minimum version is not a valid semver string",
	)
	ErrValidatorJailed               = errorsmod.Register(ModuleName, 1142, "validator is jailed")
	ErrValidatorTombstoned           = errorsmod.Register(ModuleName, 1143, "validator is tombstoned")
	ErrParamsStabilityPoolPercentage = errorsmod.Register(
		ModuleName,
		1144,
		"stability pool percentage cannot be more than 100")
	ErrUnauthorizedGrantee = errorsmod.Register(
		ModuleName,
		1145,
		"grantee is not the registered hotkey for the observer")
)
