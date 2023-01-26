package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ValidatorRewardsAmount = "ValidatorRewardsAmount"
	ReservesFactor         = "ReservesFactor"
	BondFactor             = "BondFactor"
	DurationFactor         = "DurationFactor"
	ValidatorRewardsLeft   = "ValidatorRewardsLeft"
)

const (
	ValidatorEmissons = "emissions/ValidatorEmissions"
)

func EmitValidatorEmissions(ctx sdk.Context, bondFactor, reservesFactor, durationsFactor, blockrewards, blockrewardsLeft string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(ValidatorEmissons,
			sdk.NewAttribute(BondFactor, blockrewards),
			sdk.NewAttribute(DurationFactor, blockrewards),
			sdk.NewAttribute(ReservesFactor, blockrewards),
			sdk.NewAttribute(ValidatorRewardsAmount, blockrewards),
			sdk.NewAttribute(ValidatorRewardsLeft, blockrewardsLeft),
		),
	)
}
