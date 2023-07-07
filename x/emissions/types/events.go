package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ValidatorRewardsAmount = "ValidatorRewardsAmount"
	ReservesFactor         = "ReservesFactor"
	BondFactor             = "BondFactor"
	DurationFactor         = "DurationFactor"
	ObserverRewardsAmount  = "ObserverRewardsLeft"
	TssRewardsAmount       = "ObserverRewardsLeft"
)

const (
	ValidatorEmissons = "ValidatorEmissions"
)

func EmitValidatorEmissions(ctx sdk.Context, bondFactor, reservesFactor, durationsFactor, validatorRewards, observerRewards, tssRewards string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(ValidatorEmissons,
			sdk.NewAttribute(BondFactor, bondFactor),
			sdk.NewAttribute(DurationFactor, durationsFactor),
			sdk.NewAttribute(ReservesFactor, reservesFactor),
			sdk.NewAttribute(ValidatorRewardsAmount, validatorRewards),
			sdk.NewAttribute(ObserverRewardsAmount, observerRewards),
			sdk.NewAttribute(TssRewardsAmount, tssRewards),
		),
	)
}
