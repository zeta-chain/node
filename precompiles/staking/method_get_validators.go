package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	precompiletypes "github.com/zeta-chain/node/precompiles/types"
)

// getValidators queries the list of validators for a given delegator.
func (c *Contract) getDelegatorValidators(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 1 {
		return nil, &precompiletypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 1,
		}
	}

	delegatorAddr, err := unpackGetValidatorsArgs(args)
	if err != nil {
		return nil, err
	}

	// Get the cosmos address of the caller.
	delegatorCosmosAddr, err := precompiletypes.GetCosmosAddress(c.bankKeeper, delegatorAddr)
	if err != nil {
		return nil, err
	}

	// Query the validator list of the given delegator.
	dstrQuerier := distrkeeper.NewQuerier(c.distributionKeeper)

	res, err := dstrQuerier.DelegatorValidators(ctx, &distrtypes.QueryDelegatorValidatorsRequest{
		DelegatorAddress: delegatorCosmosAddr.String(),
	})
	if err != nil {
		return nil, precompiletypes.ErrUnexpected{
			When: "DelegatorValidators",
			Got:  err.Error(),
		}
	}

	// Return immediately, no need to check the slice.
	// If there are no validators we simply return an empty array to calling contracts.
	return method.Outputs.Pack(res.Validators)
}

func unpackGetValidatorsArgs(args []interface{}) (delegator common.Address, err error) {
	delegator, ok := args[0].(common.Address)
	if !ok {
		return common.Address{}, &precompiletypes.ErrInvalidAddr{
			Got: delegator.String(),
		}
	}

	return delegator, nil
}
