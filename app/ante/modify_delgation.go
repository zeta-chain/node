package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	zetaObserverModuleKeeper "github.com/zeta-chain/zetacore/x/observer/keeper"
)

type ValidateModifyDelegationDecorator struct {
	zok *zetaObserverModuleKeeper.Keeper
}

func NewValidateModifyDelegationDecorator(keeper *zetaObserverModuleKeeper.Keeper) ValidateModifyDelegationDecorator {
	return ValidateModifyDelegationDecorator{zok: keeper}
}

func (vcd ValidateModifyDelegationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		if err := vcd.validateMsg(ctx, msg); err != nil {
			return ctx, err
		}
	}
	return next(ctx, tx, simulate)
}

func (vcd ValidateModifyDelegationDecorator) validateMsg(ctx sdk.Context, msg sdk.Msg) error {
	switch msg := msg.(type) {
	case *stakingtypes.MsgUndelegate:
		return validateAddress(msg.DelegatorAddress, vcd.zok.GetAllObserverAddresses(ctx))
	case *stakingtypes.MsgBeginRedelegate:
		return validateAddress(msg.DelegatorAddress, vcd.zok.GetAllObserverAddresses(ctx))
	}
	return nil
}

func validateAddress(address string, observerList []string) error {
	//for _, observer := range observerList {
	//	if address == observer {
	//		return sdkerrors.Wrapf(
	//			sdkerrors.ErrInvalidRequest,
	//			"cannot change delegation for observer address %s", address)
	//	}
	//}
	return nil
}
