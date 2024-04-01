package types

import errorsmod "cosmossdk.io/errors"

var ErrUnauthorized = errorsmod.Register(ModuleName, 1102, "sender not authorized")
