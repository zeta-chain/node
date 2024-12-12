package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/crypto"
	"github.com/zeta-chain/node/x/fungible/types"
)

type fungibleKeeper interface {
	GetAllForeignCoins(ctx sdk.Context) (list []types.ForeignCoins)
	SetForeignCoins(ctx sdk.Context, foreignCoins types.ForeignCoins)
}

// MigrateStore migrates the x/fungible module state from the consensus version 2 to 3
// It updates all existing address in ForeignCoin to use checksum format if the address is EVM type
func MigrateStore(ctx sdk.Context, fungibleKeeper fungibleKeeper) error {
	fcs := fungibleKeeper.GetAllForeignCoins(ctx)
	for _, fc := range fcs {
		if fc.Asset != "" && crypto.IsEVMAddress(fc.Asset) && !crypto.IsChecksumAddress(fc.Asset) {
			checksumAddress := crypto.ToChecksumAddress(fc.Asset)
			ctx.Logger().Info("Patching zrc20 asset", "zrc20", fc.Symbol, "old", fc.Asset, "new", checksumAddress)

			fc.Asset = checksumAddress
			fungibleKeeper.SetForeignCoins(ctx, fc)
		}
	}

	return nil
}
