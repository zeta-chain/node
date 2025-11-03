package app

import (
	"context"
	"os"

	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	"golang.org/x/mod/semver"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
)

// GetDefaultUpgradeHandlerVersion prints the default upgrade handler version
//
// There may be multiple upgrade handlers configured on some releases if different
// migrations needto be run in different environment
func GetDefaultUpgradeHandlerVersion() string {
	// semver must have v prefix, but we store without prefix
	vVersion := constant.GetNormalizedVersion()

	// development builds always use the full version in the release handlers
	if semver.Build(vVersion) != "" || semver.Prerelease(vVersion) != "" {
		return constant.Version
	}

	// release builds use just the major version (v22.0.0 -> v22)
	return semver.Major(vVersion)
}

func createUpgrades(chainID string) []upgradeTrackerItem {
	addErc20ModuleUpgrade := upgradeTrackerItem{
		index: 1752528615,
		storeUpgrade: &storetypes.StoreUpgrades{
			Added: []string{erc20types.ModuleName},
		},
		// TODO: enable back IBC
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
	}
	upgrades := make([]upgradeTrackerItem, 0)
	if chainID != "" {
		evmChaindID, err := chains.CosmosToEthChainID(chainID)
		if err != nil {
			panic("invalid chain ID: " + chainID + ", error: " + err.Error())
		}
		if evmChaindID == chains.ZetaChainMainnet.ChainId {
			return append(upgrades, addErc20ModuleUpgrade)
		}
	}
	return upgrades
}

func SetupHandlers(app *App) {
	allUpgrades := upgradeTracker{
		upgrades:     createUpgrades(app.ChainID()),
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

	upgradeHandlerVersion := GetDefaultUpgradeHandlerVersion()

	app.UpgradeKeeper.SetUpgradeHandler(
		upgradeHandlerVersion,
		func(ctx context.Context, _ types.Plan, vm module.VersionMap) (module.VersionMap, error) {
			app.Logger().Info("Running upgrade handler for " + upgradeHandlerVersion)

			var err error
			for _, upgradeHandler := range upgradeHandlerFns {
				vm, err = upgradeHandler(ctx, vm)
				if err != nil {
					return vm, err
				}
			}

			// TODO: are these fields ok?
			app.BankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
				Description: "The native staking token for zetacored",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    "azeta",
						Exponent: 0,
						Aliases:  nil,
					},
					{
						Denom:    "zeta",
						Exponent: 18,
						Aliases:  nil,
					},
				},
				Base:    "azeta",
				Display: "zeta",
				Name:    "Zeta Token",
				Symbol:  "ZETA",
				URI:     "",
				URIHash: "",
			})

			// (Required for NON-18 denom chains *only)
			// Update EVM params to add Extended denom options
			// Ensure that this corresponds to the EVM denom
			// (tyically the bond denom)
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			// evmParams := app.EvmKeeper.GetParams(sdkCtx)
			// evmParams.ExtendedDenomOptions = &types.ExtendedDenomOptions{ExtendedDenom: "atest"}
			// err = app.EvmKeeper.SetParams(sdkCtx, evmParams)
			// if err != nil {
			// 	return nil, err
			// }
			// Initialize EvmCoinInfo in the module store
			if err := app.EvmKeeper.InitEvmCoinInfo(sdkCtx); err != nil {
				return nil, err
			}

			return app.mm.RunMigrations(ctx, app.configurator, vm)
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}
	if upgradeInfo.Name == upgradeHandlerVersion && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(types.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
