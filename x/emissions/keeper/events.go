package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/emissions/types"
)

func EmitValidatorEmissions(
	ctx sdk.Context,
	bondFactor, reservesFactor, durationsFactor, validatorRewards, observerRewards, tssRewards string,
) {
	err := ctx.EventManager().EmitTypedEvents(&types.EventBlockEmissions{
		MsgTypeUrl:               "/zetachain.zetacore.emissions.internal.BlockEmissions",
		BondFactor:               bondFactor,
		DurationFactor:           durationsFactor,
		ReservesFactor:           reservesFactor,
		ValidatorRewardsForBlock: validatorRewards,
		ObserverRewardsForBlock:  observerRewards,
		TssRewardsForBlock:       tssRewards,
	})
	if err != nil {
		ctx.Logger().Error("Error emitting ValidatorEmissions :", err)
	}
}

func EmitObserverEmissions(ctx sdk.Context, em []*types.ObserverEmission) {
	err := ctx.EventManager().EmitTypedEvents(&types.EventObserverEmissions{
		MsgTypeUrl: "/zetachain.zetacore.emissions.internal.ObserverEmissions",
		Emissions:  em,
	})
	if err != nil {
		ctx.Logger().Error("Error emitting ObserverEmissions :", err)
	}
}
