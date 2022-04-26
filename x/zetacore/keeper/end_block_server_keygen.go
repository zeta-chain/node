package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func (k Keeper) InitializeGenesisKeygen(goCtx context.Context) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if ctx.BlockHeight() == 100 {
		accts := k.GetAllNodeAccount(ctx)
		var pubkeys []string
		for _, acct := range accts {
			if isBondedValidator(acct.Creator, validators) {
				pubkeys = append(pubkeys, acct.PubkeySet.Secp256k1.String())
			}
		}
		kg := types.Keygen{
			Creator:     "genesis keygen",
			Status:      0, // to keygen
			Pubkeys:     pubkeys,
			BlockNumber: 110,
		}
		k.SetKeygen(ctx, kg)
	}
}
