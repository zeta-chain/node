package app

import (
	"os"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/zeta-chain/node/pkg/constant"
)

func SetupHandlers(app *App) {
	allUpgrades := upgradeTracker{
		upgrades: []upgradeTrackerItem{
			// TODO: enable back IBC
			// these commented lines allow for the IBC modules to be added to the upgrade tracker
			// https://github.com/zeta-chain/node/issues/2573
			//{
			//	index: <CURRENT TIMESTAMP>,
			//	storeUpgrade: &storetypes.StoreUpgrades{
			//		Added: []string{
			//			capabilitytypes.ModuleName,
			//			ibcexported.ModuleName,
			//			ibctransfertypes.ModuleName,
			//		},
			//	},
			//},
			//{
			//	index: <CURRENT TIMESTAMP>,
			//	storeUpgrade: &storetypes.StoreUpgrades{
			//		Added: []string{ibccrosschaintypes.ModuleName},
			//	},
			//},
		},
		stateFileDir: DefaultNodeHome,
	}

	var upgradeHandlerFns []upgradeHandlerFn
	var storeUpgrades *storetypes.StoreUpgrades
	var err error
	_, useIncrementalTracker := os.LookupEnv("ZETACORED_USE_INCREMENTAL_UPGRADE_TRACKER")
	if useIncrementalTracker {
		upgradeHandlerFns, storeUpgrades, err = allUpgrades.getIncrementalUpgrades()
		if err != nil {
			panic(err)
		}
	} else {
		upgradeHandlerFns, storeUpgrades = allUpgrades.mergeAllUpgrades()
	}

	app.UpgradeKeeper.SetUpgradeHandler(
		constant.Version,
		func(ctx sdk.Context, _ types.Plan, vm module.VersionMap) (module.VersionMap, error) {
			app.Logger().Info("Running upgrade handler for " + constant.Version)

			var err error
			for _, upgradeHandler := range upgradeHandlerFns {
				vm, err = upgradeHandler(ctx, vm)
				if err != nil {
					return vm, err
				}
			}

			return app.mm.RunMigrations(ctx, app.configurator, vm)
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}
	if upgradeInfo.Name == constant.Version && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(types.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
