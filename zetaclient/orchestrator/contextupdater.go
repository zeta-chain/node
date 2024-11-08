package orchestrator

import (
	"context"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/ticker"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	zctx "github.com/zeta-chain/node/zetaclient/context"
)

type Zetacore interface {
	GetBlockHeight(ctx context.Context) (int64, error)
	GetUpgradePlan(ctx context.Context) (*upgradetypes.Plan, error)
	GetSupportedChains(ctx context.Context) ([]chains.Chain, error)
	GetAdditionalChains(ctx context.Context) ([]chains.Chain, error)
	GetCrosschainFlags(ctx context.Context) (observertypes.CrosschainFlags, error)
	GetChainParams(ctx context.Context) ([]*observertypes.ChainParams, error)
	GetTSS(ctx context.Context) (observertypes.TSS, error)
	GetKeyGen(ctx context.Context) (observertypes.Keygen, error)
}

var ErrUpgradeRequired = errors.New("upgrade required")

func (oc *Orchestrator) runAppContextUpdater(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	interval := ticker.DurationFromUint64Seconds(app.Config().ConfigUpdateTicker)

	oc.logger.Info().Msg("UpdateAppContext worker started")

	task := func(ctx context.Context, t *ticker.Ticker) error {
		err := UpdateAppContext(ctx, app, oc.zetacoreClient, oc.logger.Sampled)
		switch {
		case errors.Is(err, ErrUpgradeRequired):
			oc.onUpgradeDetected(err)
			t.Stop()
			return nil
		case err != nil:
			oc.logger.Err(err).Msg("UpdateAppContext failed")
		}

		return nil
	}

	return ticker.Run(
		ctx,
		interval,
		task,
		ticker.WithLogger(oc.logger.Logger, "UpdateAppContext"),
		ticker.WithStopChan(oc.stop),
	)
}

// UpdateAppContext fetches latest data from Zetacore and updates the AppContext.
// Also detects if an upgrade is required. If an upgrade is required, it returns ErrUpgradeRequired.
func UpdateAppContext(ctx context.Context, app *zctx.AppContext, zc Zetacore, logger zerolog.Logger) error {
	bn, err := zc.GetBlockHeight(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get zeta block height")
	}

	if err := checkForZetacoreUpgrade(ctx, bn, zc); err != nil {
		return err
	}

	supportedChains, err := zc.GetSupportedChains(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch supported chains")
	}

	additionalChains, err := zc.GetAdditionalChains(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch additional chains")
	}

	chainParams, err := zc.GetChainParams(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch chain params")
	}

	keyGen, err := zc.GetKeyGen(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch keygen from zetacore")
	}

	crosschainFlags, err := zc.GetCrosschainFlags(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch crosschain flags from zetacore")
	}

	tss, err := zc.GetTSS(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch current TSS")
	}

	freshParams := make(map[int64]*observertypes.ChainParams, len(chainParams))

	// check and update chain params for each chain
	// Note that we are EXCLUDING ZetaChain from the chainParams if it's present
	for i := range chainParams {
		cp := chainParams[i]

		if !cp.IsSupported {
			logger.Warn().Int64("chain.id", cp.ChainId).Msg("Skipping unsupported chain")
			continue
		}

		if chains.IsZetaChain(cp.ChainId, nil) {
			continue
		}

		if err := observertypes.ValidateChainParams(cp); err != nil {
			logger.Warn().Err(err).Int64("chain.id", cp.ChainId).Msg("Skipping invalid chain params")
			continue
		}

		freshParams[cp.ChainId] = cp
	}

	return app.Update(
		keyGen,
		supportedChains,
		additionalChains,
		freshParams,
		tss.GetTssPubkey(),
		crosschainFlags,
	)
}

// returns an error if an upgrade is required
func checkForZetacoreUpgrade(ctx context.Context, zetaHeight int64, zc Zetacore) error {
	plan, err := zc.GetUpgradePlan(ctx)
	switch {
	case err != nil:
		return errors.Wrap(err, "unable to get upgrade plan")
	case plan == nil:
		// no upgrade planned
		return nil
	}

	upgradeHeight := plan.Height

	// We can return an error in a few blocks ahead.
	// It's okay because the ticker might have a long interval.
	const upgradeRange = 2

	// Note that after plan.Height's block `x/upgrade` module deletes the plan
	if (upgradeHeight - zetaHeight) <= upgradeRange {
		return errors.Wrapf(ErrUpgradeRequired, "current height: %d, upgrade height: %d", zetaHeight, upgradeHeight)
	}

	return nil
}

// onUpgradeDetected is called when an upgrade is detected.
func (oc *Orchestrator) onUpgradeDetected(errDetected error) {
	const msg = "Upgrade detected." +
		" Kill the process, replace the binary with upgraded version, and restart zetaclientd"

	oc.logger.Warn().Str("upgrade", errDetected.Error()).Msg(msg)
	oc.Stop()
}
