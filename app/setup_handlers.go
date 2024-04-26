package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
)

const releaseVersion = "v15"

func SetupHandlers(app *App) {
	app.UpgradeKeeper.SetUpgradeHandler(releaseVersion, func(ctx sdk.Context, _ types.Plan, vm module.VersionMap) (module.VersionMap, error) {
		app.Logger().Info("Running upgrade handler for " + releaseVersion)
		return app.mm.RunMigrations(ctx, app.configurator, vm)
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}
	if upgradeInfo.Name == releaseVersion && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{authoritytypes.ModuleName, lightclienttypes.ModuleName},
		}
		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(types.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
