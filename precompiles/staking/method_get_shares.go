package staking

import (
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	precompiletypes "github.com/zeta-chain/node/precompiles/types"
)

func (c *Contract) GetShares(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 2 {
		return nil, &(precompiletypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		})
	}
	stakerAddress, ok := args[0].(common.Address)
	if !ok {
		return nil, precompiletypes.ErrInvalidArgument{
			Got: args[0],
		}
	}

	validatorAddress, ok := args[1].(string)
	if !ok {
		return nil, precompiletypes.ErrInvalidArgument{
			Got: args[1],
		}
	}

	validator, err := sdk.ValAddressFromBech32(validatorAddress)
	if err != nil {
		return nil, err
	}
	shares := big.NewInt(0)
	delegation, err := c.stakingKeeper.Delegation(ctx, sdk.AccAddress(stakerAddress.Bytes()), validator)
	if err != nil {
		if errors.Is(err, stakingtypes.ErrNoDelegation) {
			return method.Outputs.Pack(shares)
		}
		return nil, err
	}

	if delegation != nil {
		shares = delegation.GetShares().BigInt()
	}

	return method.Outputs.Pack(shares)
}
