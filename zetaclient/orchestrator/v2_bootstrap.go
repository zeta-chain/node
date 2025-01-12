package orchestrator

import (
	"context"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/chains/bitcoin"
	btcobserver "github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/rpc"
	btcsigner "github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/db"
)

func (oc *V2) bootstrapBitcoin(ctx context.Context, chain zctx.Chain) (*bitcoin.Bitcoin, error) {
	// should not happen
	if !chain.IsBitcoin() {
		return nil, errors.New("chain is not bitcoin")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	cfg, found := app.Config().GetBTCConfig(chain.ID())
	if !found {
		return nil, errors.Wrap(errSkipChain, "unable to find btc config")
	}

	rpcClient, err := rpc.NewRPCClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create rpc client")
	}

	var (
		rawChain       = chain.RawChain()
		rawChainParams = chain.Params()
	)

	dbName := btcDatabaseFileName(*rawChain)

	database, err := db.NewFromSqlite(oc.deps.DBPath, dbName, true)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open database %s", dbName)
	}

	// TODO extract base observer
	// TODO extract base signer
	// https://github.com/zeta-chain/node/issues/3331

	observer, err := btcobserver.NewObserver(
		*rawChain,
		rpcClient,
		*rawChainParams,
		oc.deps.Zetacore,
		oc.deps.TSS,
		database,
		oc.logger.base,
		oc.deps.Telemetry,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create observer")
	}

	signer, err := btcsigner.NewSigner(*rawChain, oc.deps.TSS, oc.logger.base, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create signer")
	}

	return bitcoin.New(oc.scheduler, observer, signer), nil
}
