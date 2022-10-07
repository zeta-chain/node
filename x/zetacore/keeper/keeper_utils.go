package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func (k Keeper) isAuthorized(ctx sdk.Context, address string, senderChain string, observationType string) (bool, error) {
	observerMapper, found := k.zetaObserverKeeper.GetObserverMapper(ctx, senderChain, observationType)
	if !found {
		return false, errors.Wrap(types.ErrNotAuthorized, fmt.Sprintf("Chain/Observation type not supported Chain : %s , Observation type : %s", senderChain, observationType))
	}
	for _, obs := range observerMapper.ObserverList {
		if obs == address {
			return true, nil
		}
	}
	return false, errors.Wrap(types.ErrNotAuthorized, fmt.Sprintf("Adress: %s", address))
}

func (k Keeper) hasSuperMajorityValidators(ctx sdk.Context, signers []string) bool {
	numSigners := len(signers)
	validators := k.StakingKeeper.GetAllValidators(ctx)
	numValidValidators := 0
	for _, v := range validators {
		if v.IsBonded() {
			numValidValidators++
		}
	}
	threshold := numValidValidators * 2 / 3
	if threshold < 2 {
		threshold = 2
	}
	if numValidValidators == 1 {
		threshold = 1
	}
	return numSigners == threshold
}
