package orchestrator

import (
	"context"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	btcobserver "github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	btcsigner "github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/db"
)

const btcBlocksPerDay = 144

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

	rpcClient, err := client.New(cfg, chain.ID(), oc.logger.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create rpc client")
	}

	var (
		rawChain = chain.RawChain()
		dbName   = btcDatabaseFileName(*rawChain)
	)

	baseObserver, err := oc.newBaseObserver(chain, dbName)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create base observer")
	}

	observer, err := btcobserver.New(*rawChain, baseObserver, rpcClient)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create observer")
	}

	baseSigner := oc.newBaseSigner(chain)
	signer := btcsigner.New(baseSigner, rpcClient)

	return bitcoin.New(oc.scheduler, observer, signer), nil
}

func (oc *V2) newBaseObserver(chain zctx.Chain, dbName string) (*base.Observer, error) {
	var (
		rawChain       = chain.RawChain()
		rawChainParams = chain.Params()
	)

	database, err := db.NewFromSqlite(oc.deps.DBPath, dbName, true)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open database %s", dbName)
	}

	blocksCacheSize := base.DefaultBlockCacheSize
	if chain.IsBitcoin() {
		blocksCacheSize = btcBlocksPerDay
	}

	return base.NewObserver(
		*rawChain,
		*rawChainParams,
		oc.deps.Zetacore,
		oc.deps.TSS,
		blocksCacheSize,
		oc.deps.Telemetry,
		database,
		oc.logger.base,
	)
}

func (oc *V2) newBaseSigner(chain zctx.Chain) *base.Signer {
	return base.NewSigner(*chain.RawChain(), oc.deps.TSS, oc.logger.base)
}
