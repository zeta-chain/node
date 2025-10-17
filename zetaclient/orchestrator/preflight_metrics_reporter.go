package orchestrator

import (
	"context"
	"time"

	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	evmclient "github.com/zeta-chain/node/zetaclient/chains/evm/client"
	solrepo "github.com/zeta-chain/node/zetaclient/chains/solana/repo"
	suiclient "github.com/zeta-chain/node/zetaclient/chains/sui/client"
	tonrpc "github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// ReportPreflightMetrics performs a preflight check on preflight chains (where IsSuported=false) and updates metrics.
// This helps to visualize the readiness of new chains to be enabled and observed by zetaclient.
func ReportPreflightMetrics(ctx context.Context, app *zctx.AppContext, zc ZetacoreClient, logger zerolog.Logger) error {
	additionalChains, err := zc.GetAdditionalChains(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch additional chains")
	}

	chainParams, err := zc.GetChainParams(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch chain params")
	}

	// chains that exist in zetacore but are not enabled yet (e.g. in preflight mode)
	unsupportedChains := make([]chains.Chain, 0)

	// We treat chain params with IsSupported = false as preflight chains, because
	// having a new flag 'IsPreflight' in the chain params incurs additional complexity.
	for i := range chainParams {
		cp := chainParams[i]

		if chains.IsZetaChain(cp.ChainId, nil) {
			continue
		}

		if err := cp.Validate(); err != nil {
			continue
		}

		chain, found := chains.GetChainFromChainID(cp.ChainId, additionalChains)
		if !found {
			continue
		}

		if !cp.IsSupported {
			unsupportedChains = append(unsupportedChains, chain)
		}
	}

	// perform preflight check
	start := time.Now()
	for _, chain := range unsupportedChains {
		switch {
		case chains.IsBitcoinChain(chain.ChainId, additionalChains):
			err = reportPreflightMetricsBitcoin(ctx, app, chain, logger)
		case chains.IsEVMChain(chain.ChainId, additionalChains):
			err = reportPreflightMetricsEVM(ctx, app, chain)
		case chains.IsSolanaChain(chain.ChainId, additionalChains):
			err = reportPreflightMetricsSolana(ctx, app, chain)
		case chains.IsSuiChain(chain.ChainId, additionalChains):
			err = reportPreflightMetricsSui(ctx, app, chain)
		case chains.IsTONChain(chain.ChainId, additionalChains):
			err = reportPreflightMetricsTON(ctx, app, chain)
		default:
			err = errors.New("unable to perform preflight check for unsupported chain")
		}

		if err != nil {
			logger.Error().
				Err(err).
				Int64(logs.FieldChain, chain.ChainId).
				Float64("time_taken", time.Since(start).Seconds()).
				Msg("unable to report preflight metrics")
		}
	}

	return nil
}

func reportPreflightMetricsBitcoin(
	ctx context.Context,
	app *zctx.AppContext,
	chain chains.Chain,
	logger zerolog.Logger,
) error {
	cfg, found := app.Config().GetBTCConfig(chain.ChainId)
	if !found {
		return nil
	}

	rpcClient, err := client.New(cfg, chain.ChainId, logger)
	if err != nil {
		return errors.Wrap(err, "unable to create btc rpc client")
	}

	blockTime, err := rpcClient.Healthcheck(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get bitcoin last block time")
	}
	metrics.ReportBlockLatency(chain.Name, blockTime)

	return nil
}

func reportPreflightMetricsEVM(ctx context.Context, app *zctx.AppContext, chain chains.Chain) error {
	cfg, found := app.Config().GetEVMConfig(chain.ChainId)
	if !found {
		return nil
	}

	evmClient, err := evmclient.NewFromEndpoint(ctx, cfg.Endpoint)
	if err != nil {
		return errors.Wrapf(err, "unable to create evm client (%s)", cfg.Endpoint)
	}

	blockTime, err := evmClient.HealthCheck(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get evm last block time")
	}
	metrics.ReportBlockLatency(chain.Name, blockTime)

	return nil
}

func reportPreflightMetricsSolana(ctx context.Context, app *zctx.AppContext, chain chains.Chain) error {
	cfg, found := app.Config().GetSolanaConfig()
	if !found {
		return nil
	}

	rpcClient := solrpc.New(cfg.Endpoint)
	if rpcClient == nil {
		return errors.New("unable to create solana rpc client")
	}

	// TODO: The Solana repository should be injected as a dependency into this function. We
	// should not have to instantiate the Solana client here.
	blockTime, err := solrepo.New(rpcClient).HealthCheck(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get solana last block time")
	}
	metrics.ReportBlockLatency(chain.Name, blockTime)

	return nil
}

func reportPreflightMetricsTON(ctx context.Context, app *zctx.AppContext, chain chains.Chain) error {
	cfg, found := app.Config().GetTONConfig()
	if !found {
		return nil
	}

	client := tonrpc.New(cfg.Endpoint, chain.ChainId)

	blockTime, err := client.HealthCheck(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get ton last block time")
	}

	metrics.ReportBlockLatency(chain.Name, blockTime)

	return nil
}

func reportPreflightMetricsSui(ctx context.Context, app *zctx.AppContext, chain chains.Chain) error {
	cfg, found := app.Config().GetSuiConfig()
	if !found {
		return nil
	}

	suiClient := suiclient.New(cfg.Endpoint)

	blockTime, err := suiClient.HealthCheck(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get sui last block time")
	}
	metrics.ReportBlockLatency(chain.Name, blockTime)

	return nil
}
