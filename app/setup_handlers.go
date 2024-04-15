package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

const releaseVersion = "v15"

func SetupHandlers(app *App) {
	// Set param key table for params module migration
	for _, subspace := range app.ParamsKeeper.GetSubspaces() {
		subspace := subspace

		switch subspace.Name() {
		// TODO: add all modules when cosmos-sdk is updated
		case emissionstypes.ModuleName:
			subspace.WithKeyTable(emissionstypes.ParamKeyTable())
		case observertypes.ModuleName:
			subspace.WithKeyTable(observertypes.ParamKeyTable())
		}
	}
	app.UpgradeKeeper.SetUpgradeHandler(releaseVersion, func(ctx sdk.Context, plan types.Plan, vm module.VersionMap) (module.VersionMap, error) {
		app.Logger().Info("Running upgrade handler for " + releaseVersion)
		// Updated version map to the latest consensus versions from each module
		for m, mb := range app.mm.Modules {
			if module, ok := mb.(module.HasConsensusVersion); ok {
				vm[m] = module.ConsensusVersion()
			}
		}
		VersionMigrator{v: vm}.TriggerMigration(observertypes.ModuleName)
		VersionMigrator{v: vm}.TriggerMigration(emissionstypes.ModuleName)

		return app.mm.RunMigrations(ctx, app.configurator, vm)
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}
	if upgradeInfo.Name == releaseVersion && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{authoritytypes.ModuleName},
		}
		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(types.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

type VersionMigrator struct {
	v module.VersionMap
}

func (v VersionMigrator) TriggerMigration(moduleName string) module.VersionMap {
	v.v[moduleName] = v.v[moduleName] - 1
	return v.v
}
