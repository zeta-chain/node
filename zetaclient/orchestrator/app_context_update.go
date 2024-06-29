package orchestrator

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

// WatchAppContext watches for app context changes and updates app context
func (oc *Orchestrator) WatchAppContext() {
	oc.logger.Std.Info().Msg("UpdateAppContext started")

	ticker := time.NewTicker(time.Duration(oc.appContext.Config().ConfigUpdateTicker) * time.Second)
	for {
		select {
		case <-ticker.C:
			err := oc.UpdateAppContext()
			if err != nil {
				oc.logger.Std.Err(err).Msg("error updating zetaclient app context")
			}
		case <-oc.stop:
			oc.logger.Std.Info().Msg("UpdateAppContext stopped")
			return
		}
	}
}

// WatchUpgradePlan watches for upgrade plan and stops orchestrator if upgrade height is reached
func (oc *Orchestrator) WatchUpgradePlan() {
	oc.logger.Std.Info().Msg("WatchUpgradePlan started")

	// detect upgrade plan every half Zeta block in order to hit every height
	ticker := time.NewTicker(common.ZetaBlockTime / 2)
	for range ticker.C {
		reached, err := oc.UpgradeHeightReached()
		if err != nil {
			oc.logger.Sampled.Error().Err(err).Msg("error detecting upgrade plan")
		} else if reached {
			oc.Stop()
			oc.logger.Std.Info().Msg("WatchUpgradePlan stopped")
			return
		}
	}
}

// UpdateAppContext updates zetaclient app context
func (oc *Orchestrator) UpdateAppContext() error {
	// fetch latest app context from zetacore
	err := oc.zetacoreClient.UpdateAppContext(oc.appContext, oc.logger.Std)
	if err != nil {
		return errors.Wrap(err, "UpdateAppContext: error updating app context from zetacore")
	}

	// reload config from file to allow for runtime config changes
	// this allows operator to update zetaclient config without restarting zetaclient
	zetazoreHome := oc.appContext.Config().ZetaCoreHome
	newConfig, err := config.Load(zetazoreHome)
	if err != nil {
		return errors.Wrapf(err, "UpdateAppContext: error loading config from path %s", zetazoreHome)
	}

	// set new config to app context
	// we keep the old ZetaCoreHome because newConfig.ZetaCoreHome is empty after reload
	newConfig.ZetaCoreHome = zetazoreHome
	oc.appContext.SetConfig(newConfig)

	return nil
}

// UpgradeHeightReached returns true if upgrade height is reached
func (oc *Orchestrator) UpgradeHeightReached() (bool, error) {
	// query for active upgrade plan
	plan, err := oc.zetacoreClient.GetUpgradePlan()
	if err != nil {
		return false, fmt.Errorf("failed to get upgrade plan: %w", err)
	}

	// if there is no active upgrade plan, plan will be nil.
	if plan == nil {
		return false, nil
	}

	// get ZetaChain block height
	height, err := oc.zetacoreClient.GetBlockHeight()
	if err != nil {
		return false, fmt.Errorf("failed to get block height: %w", err)
	}

	// if upgrade height is not reached, do nothing
	if height != plan.Height-1 {
		return false, nil
	}

	// stop zetaclients if upgrade height is reached; notify operator to upgrade and restart
	oc.logger.Std.Warn().
		Msgf("Active upgrade plan detected and upgrade height reached: %s at height %d; ZetaClient is stopped;"+
			"please kill this process, replace zetaclientd binary with upgraded version, and restart zetaclientd", plan.Name, plan.Height)

	return true, nil
}
