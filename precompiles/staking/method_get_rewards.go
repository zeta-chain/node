package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	dstrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	precompiletypes "github.com/zeta-chain/node/precompiles/types"
)

// getRewards returns the list of ZRC20 cosmos coins, available for withdrawal by the delegator.
func (c *Contract) getRewards(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 2 {
		return nil, &precompiletypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		}
	}

	delegatorAddr, validatorAddr, err := unpackGetRewardsArgs(args)
	if err != nil {
		return nil, err
	}

	// Get delegator Cosmos address.
	delegatorCosmosAddr, err := precompiletypes.GetCosmosAddress(c.bankKeeper, delegatorAddr)
	if err != nil {
		return nil, err
	}

	// Query the delegation rewards through the distribution keeper querier.
	dstrQuerier := distrkeeper.NewQuerier(c.distributionKeeper)

	res, err := dstrQuerier.DelegationRewards(ctx, &dstrtypes.QueryDelegationRewardsRequest{
		DelegatorAddress: delegatorCosmosAddr.String(),
		ValidatorAddress: validatorAddr,
	})
	if err != nil {
		return nil, precompiletypes.ErrUnexpected{
			When: "DelegationRewards",
			Got:  err.Error(),
		}
	}

	coins := res.GetRewards()
	if !coins.IsValid() {
		return nil, precompiletypes.ErrUnexpected{
			When: "GetRewards",
			Got:  "invalid coins",
		}
	}

	zrc20Coins := make([]sdk.DecCoin, 0)
	for _, coin := range coins {
		if precompiletypes.CoinIsZRC20(coin.Denom) {
			zrc20Coins = append(zrc20Coins, coin)
		}
	}

	return method.Outputs.Pack(zrc20Coins)
}

func unpackGetRewardsArgs(args []interface{}) (delegator common.Address, validator string, err error) {
	delegator, ok := args[0].(common.Address)
	if !ok {
		return common.Address{}, "", &precompiletypes.ErrInvalidAddr{
			Got: delegator.String(),
		}
	}

	validator, ok = args[1].(string)
	if !ok {
		return common.Address{}, "", &precompiletypes.ErrInvalidAddr{
			Got: validator,
		}
	}

	return delegator, validator, nil
}
