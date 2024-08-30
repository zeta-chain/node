package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/emissions/types"
)

func (k Keeper) GetBlockRewardComponents(ctx sdk.Context, params types.Params) (sdk.Dec, sdk.Dec, sdk.Dec) {
	reservesFactor := k.GetReservesFactor(ctx)
	if reservesFactor.LTE(sdk.ZeroDec()) {
		return sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()
	}
	bondFactor := params.GetBondFactor(k.stakingKeeper.BondedRatio(ctx))
	durationFactor := params.GetDurationFactor(ctx.BlockHeight())
	return reservesFactor, bondFactor, durationFactor
}

func (k Keeper) GetReservesFactor(ctx sdk.Context) sdk.Dec {
	reserveAmount := k.GetBankKeeper().GetBalance(ctx, types.EmissionsModuleAddress, config.BaseDenom)
	return sdk.NewDecFromInt(reserveAmount.Amount)
}

func (k Keeper) GetFixedBlockRewards() (sdk.Dec, error) {
	return CalculateFixedValidatorRewards(types.AvgBlockTime)
}

func CalculateFixedValidatorRewards(avgBlockTimeString string) (sdk.Dec, error) {
	azetaAmountTotalRewards, err := coin.GetAzetaDecFromAmountInZeta(types.BlockRewardsInZeta)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	avgBlockTime, err := sdk.NewDecFromStr(avgBlockTimeString)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	numberOfBlocksInAMonth := sdk.NewDec(types.SecsInMonth).Quo(avgBlockTime)
	numberOfBlocksTotal := numberOfBlocksInAMonth.Mul(sdk.NewDec(12)).Mul(sdk.NewDec(types.EmissionScheduledYears))
	constantRewardPerBlock := azetaAmountTotalRewards.Quo(numberOfBlocksTotal)
	return constantRewardPerBlock, nil
}
