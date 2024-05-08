package app

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
)

const releaseVersion = "v17"

func SetupHandlers(app *App) {
	// Set param key table for params module migration
	for _, subspace := range app.ParamsKeeper.GetSubspaces() {
		var keyTable paramstypes.KeyTable
		switch subspace.Name() {
		case authtypes.ModuleName:
			keyTable = authtypes.ParamKeyTable() //nolint:staticcheck
		case banktypes.ModuleName:
			keyTable = banktypes.ParamKeyTable() //nolint:staticcheck
		case stakingtypes.ModuleName:
			keyTable = stakingtypes.ParamKeyTable()
		case distrtypes.ModuleName:
			keyTable = distrtypes.ParamKeyTable() //nolint:staticcheck
		case slashingtypes.ModuleName:
			keyTable = slashingtypes.ParamKeyTable() //nolint:staticcheck
		case govtypes.ModuleName:
			keyTable = govv1.ParamKeyTable() //nolint:staticcheck
		case crisistypes.ModuleName:
			keyTable = crisistypes.ParamKeyTable() //nolint:staticcheck
		case emissionstypes.ModuleName:
			keyTable = emissionstypes.ParamKeyTable()
		default:
			continue
		}
		if !subspace.HasKeyTable() {
			subspace.WithKeyTable(keyTable)
		}
	}
	baseAppLegacySS := app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())
	app.UpgradeKeeper.SetUpgradeHandler(releaseVersion, func(ctx sdk.Context, plan types.Plan, vm module.VersionMap) (module.VersionMap, error) {
		app.Logger().Info("Running upgrade handler for " + releaseVersion)
		// Migrate Tendermint consensus parameters from x/params module to a dedicated x/consensus module.
		baseapp.MigrateParams(ctx, baseAppLegacySS, &app.ConsensusParamsKeeper)
		// Updated version map to the latest consensus versions from each module
		for m, mb := range app.mm.Modules {
			if module, ok := mb.(module.HasConsensusVersion); ok {
				vm[m] = module.ConsensusVersion()
			}
		}

		VersionMigrator{v: vm}.TriggerMigration(authtypes.ModuleName)
		VersionMigrator{v: vm}.TriggerMigration(banktypes.ModuleName)
		VersionMigrator{v: vm}.TriggerMigration(stakingtypes.ModuleName)
		VersionMigrator{v: vm}.TriggerMigration(distrtypes.ModuleName)
		VersionMigrator{v: vm}.TriggerMigration(slashingtypes.ModuleName)
		VersionMigrator{v: vm}.TriggerMigration(govtypes.ModuleName)
		VersionMigrator{v: vm}.TriggerMigration(crisistypes.ModuleName)

		VersionMigrator{v: vm}.TriggerMigration(emissionstypes.ModuleName)

		return app.mm.RunMigrations(ctx, app.configurator, vm)
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}
	if upgradeInfo.Name == releaseVersion && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{consensustypes.ModuleName, crisistypes.ModuleName},
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
