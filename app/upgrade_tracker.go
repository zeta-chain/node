package app

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const incrementalUpgradeTrackerStateFile = "incrementalupgradetracker"

type upgradeHandlerFn func(ctx context.Context, vm module.VersionMap) (module.VersionMap, error)

type upgradeTrackerItem struct {
	// Monotonically increasing index to order and track migrations. Typically the current unix epoch timestamp.
	index int64
	// Function that will run during the SetUpgradeHandler callback. The VersionMap must always be returned.
	upgradeHandler upgradeHandlerFn
	// StoreUpgrades that will be provided to UpgradeStoreLoader
	storeUpgrade *storetypes.StoreUpgrades
}

// upgradeTracker allows us to track needed upgrades/migrations across both release and develop builds
type upgradeTracker struct {
	upgrades []upgradeTrackerItem
	// directory the incremental state file is stored
	stateFileDir string
}

// getIncrementalUpgrades gets all upgrades that have not been applied. This is typically
// used for developnet upgrades since we need to run migrations as the are committed rather than
// all at once during a release
func (t upgradeTracker) getIncrementalUpgrades() ([]upgradeHandlerFn, *storetypes.StoreUpgrades, error) {
	neededUpgrades := &storetypes.StoreUpgrades{}
	neededUpgradeHandlers := []upgradeHandlerFn{}
	stateFilePath := path.Join(t.stateFileDir, incrementalUpgradeTrackerStateFile)

	currentIndex := int64(0)
	stateFileContents, err := os.ReadFile(stateFilePath) // #nosec G304 -- stateFilePath is not user controllable
	if err == nil {
		currentIndex, err = strconv.ParseInt(string(stateFileContents), 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to decode upgrade tracker: %w", err)
		}
	} else {
		fmt.Printf("unable to load upgrade tracker: %v\n", err)
	}

	maxIndex := currentIndex
	for _, item := range t.upgrades {
		index := item.index
		upgrade := item.storeUpgrade
		upgradeHandler := item.upgradeHandler
		if index <= currentIndex {
			continue
		}
		if upgradeHandler != nil {
			neededUpgradeHandlers = append(neededUpgradeHandlers, upgradeHandler)
		}
		if upgrade != nil {
			neededUpgrades.Added = append(neededUpgrades.Added, upgrade.Added...)
			neededUpgrades.Deleted = append(neededUpgrades.Deleted, upgrade.Deleted...)
			neededUpgrades.Renamed = append(neededUpgrades.Renamed, upgrade.Renamed...)
		}
		maxIndex = index
	}
	err = os.WriteFile(stateFilePath, []byte(strconv.FormatInt(maxIndex, 10)), 0o600)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to write upgrade state file: %w", err)
	}
	return neededUpgradeHandlers, neededUpgrades, nil
}

// mergeAllUpgrades unconditionally merges all upgrades. Typically used to gather the
// migrations used during a release upgrade.
func (t upgradeTracker) mergeAllUpgrades() ([]upgradeHandlerFn, *storetypes.StoreUpgrades) {
	upgrades := &storetypes.StoreUpgrades{}
	upgradeHandlers := []upgradeHandlerFn{}
	for _, item := range t.upgrades {
		upgrade := item.storeUpgrade
		versionModifier := item.upgradeHandler
		if versionModifier != nil {
			upgradeHandlers = append(upgradeHandlers, versionModifier)
		}
		if upgrade != nil {
			upgrades.Added = append(upgrades.Added, upgrade.Added...)
			upgrades.Deleted = append(upgrades.Deleted, upgrade.Deleted...)
			upgrades.Renamed = append(upgrades.Renamed, upgrade.Renamed...)
		}
	}
	return upgradeHandlers, upgrades
}
