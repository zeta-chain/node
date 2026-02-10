package ante

import (
	"slices"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

var _ sdk.AnteDecorator = VestingAccountDecorator{}

// VestingAccountDecorator blocks vesting messages from reaching the mempool
type VestingAccountDecorator struct {
	disabledMsgTypeURLs []string
}

// NewVestingAccountDecorator creates a decorator to block vesting messages from reaching the mempool
func NewVestingAccountDecorator() VestingAccountDecorator {
	return VestingAccountDecorator{
		disabledMsgTypeURLs: []string{
			sdk.MsgTypeURL(&vesting.MsgCreateVestingAccount{}),
			sdk.MsgTypeURL(&vesting.MsgCreatePermanentLockedAccount{}),
			sdk.MsgTypeURL(&vesting.MsgCreatePeriodicVestingAccount{}),
		},
	}
}

// AnteHandle implements AnteDecorator
func (vad VestingAccountDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		typeURL := sdk.MsgTypeURL(msg)

		if slices.Contains(vad.disabledMsgTypeURLs, typeURL) {
			return ctx, errorsmod.Wrapf(
				sdkerrors.ErrUnauthorized,
				"MsgTypeURL %s not supported",
				typeURL,
			)
		}
	}

	return next(ctx, tx, simulate)
}
