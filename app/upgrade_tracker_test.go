package app

import (
	"os"
	"path"
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/stretchr/testify/require"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
)

func TestUpgradeTracker(t *testing.T) {
	r := require.New(t)

	tmpdir, err := os.MkdirTemp("", "storeupgradetracker-*")
	r.NoError(err)
	defer os.RemoveAll(tmpdir)

	allUpgrades := upgradeTracker{
		upgrades: []upgradeTrackerItem{
			{
				index: 1000,
				storeUpgrade: &storetypes.StoreUpgrades{
					Added: []string{authoritytypes.ModuleName},
				},
			},
			{
				index: 2000,
				storeUpgrade: &storetypes.StoreUpgrades{
					Added: []string{lightclienttypes.ModuleName},
				},
				upgradeHandler: func(ctx sdk.Context, vm module.VersionMap) (module.VersionMap, error) {
					return vm, nil
				},
			},
			{
				index: 3000,
				upgradeHandler: func(ctx sdk.Context, vm module.VersionMap) (module.VersionMap, error) {
					return vm, nil
				},
			},
		},
		stateFileDir: tmpdir,
	}

	upgradeHandlers, storeUpgrades := allUpgrades.mergeAllUpgrades()
	r.Len(storeUpgrades.Added, 2)
	r.Len(storeUpgrades.Renamed, 0)
	r.Len(storeUpgrades.Deleted, 0)
	r.Len(upgradeHandlers, 2)

	// should return all migrations on first call
	upgradeHandlers, storeUpgrades, err = allUpgrades.getDevelopUpgrades()
	r.NoError(err)
	r.Len(storeUpgrades.Added, 2)
	r.Len(storeUpgrades.Renamed, 0)
	r.Len(storeUpgrades.Deleted, 0)
	r.Len(upgradeHandlers, 2)

	// should return no upgrades on second call
	upgradeHandlers, storeUpgrades, err = allUpgrades.getDevelopUpgrades()
	r.NoError(err)
	r.Len(storeUpgrades.Added, 0)
	r.Len(storeUpgrades.Renamed, 0)
	r.Len(storeUpgrades.Deleted, 0)
	r.Len(upgradeHandlers, 0)

	// now add a upgrade and ensure that it gets run without running
	// the other upgrades
	allUpgrades.upgrades = append(allUpgrades.upgrades, upgradeTrackerItem{
		index: 4000,
		storeUpgrade: &storetypes.StoreUpgrades{
			Deleted: []string{"example"},
		},
	})

	upgradeHandlers, storeUpgrades, err = allUpgrades.getDevelopUpgrades()
	r.NoError(err)
	r.Len(storeUpgrades.Added, 0)
	r.Len(storeUpgrades.Renamed, 0)
	r.Len(storeUpgrades.Deleted, 1)
	r.Len(upgradeHandlers, 0)
}

func TestUpgradeTrackerBadState(t *testing.T) {
	r := require.New(t)

	tmpdir, err := os.MkdirTemp("", "storeupgradetracker-*")
	r.NoError(err)
	defer os.RemoveAll(tmpdir)

	stateFilePath := path.Join(tmpdir, developUpgradeTrackerStateFile)

	err = os.WriteFile(stateFilePath, []byte("badstate"), 0o600)
	r.NoError(err)

	allUpgrades := upgradeTracker{
		upgrades:     []upgradeTrackerItem{},
		stateFileDir: tmpdir,
	}
	_, _, err = allUpgrades.getDevelopUpgrades()
	r.Error(err)
}
