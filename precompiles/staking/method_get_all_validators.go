package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

func (c *Contract) GetAllValidators(
	ctx sdk.Context,
	method *abi.Method,
) ([]byte, error) {
	validators := c.stakingKeeper.GetAllValidators(ctx)

	validatorsRes := make([]Validator, len(validators))
	for i, v := range validators {
		validatorsRes[i] = Validator{
			OperatorAddress: v.OperatorAddress,
			ConsensusPubKey: v.ConsensusPubkey.String(),
			// Safe casting from int32 to uint8, as BondStatus is an enum.
			//nolint:gosec
			BondStatus: uint8(v.Status),
			Jailed:     v.Jailed,
		}
	}

	return method.Outputs.Pack(validatorsRes)
}
