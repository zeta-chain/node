package orchestrator

import (
	"context"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/onrik/ethrpc"
	"github.com/pkg/errors"
	"github.com/zeta-chain/node/pkg/chains"

	"github.com/zeta-chain/node/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	btcobserver "github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	btcsigner "github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
	"github.com/zeta-chain/node/zetaclient/chains/evm"
	evmclient "github.com/zeta-chain/node/zetaclient/chains/evm/client"
	evmobserver "github.com/zeta-chain/node/zetaclient/chains/evm/observer"
	evmsigner "github.com/zeta-chain/node/zetaclient/chains/evm/signer"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/metrics"
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

	rpcClient, err := client.New(cfg, chain.ID(), oc.logger.Logger)
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

	signer := btcsigner.New(*rawChain, oc.deps.TSS, rpcClient, oc.logger.base)

	return bitcoin.New(oc.scheduler, observer, signer), nil
}

func (oc *V2) bootstrapEVM(ctx context.Context, chain zctx.Chain) (*evm.EVM, error) {
	// should not happen
	if !chain.IsEVM() {
		return nil, errors.New("chain is not EVM")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	cfg, found := app.Config().GetEVMConfig(chain.ID())
	if !found || cfg.Empty() {
		return nil, errors.Wrap(errSkipChain, "unable to find evm config")
	}

	var (
		rawChain       = chain.RawChain()
		rawChainParams = chain.Params()
	)

	httpClient, err := metrics.GetInstrumentedHTTPClient(cfg.Endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create http client (%s)", cfg.Endpoint)
	}

	evmClient, err := evmclient.NewFromEndpoint(ctx, cfg.Endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create evm client (%s)", cfg.Endpoint)
	}

	database, err := db.NewFromSqlite(oc.deps.DBPath, chain.Name(), true)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open database path")
	}

	evmJSONRPCClient := ethrpc.NewEthRPC(cfg.Endpoint, ethrpc.WithHttpClient(httpClient))

	observer, err := evmobserver.New(
		ctx,
		*rawChain,
		evmClient,
		evmJSONRPCClient,
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

	var (
		zetaConnectorAddress = ethcommon.HexToAddress(chain.Params().ConnectorContractAddress)
		erc20CustodyAddress  = ethcommon.HexToAddress(chain.Params().Erc20CustodyContractAddress)
		gatewayAddress       = ethcommon.HexToAddress(chain.Params().GatewayAddress)
	)

	signer, err := evmsigner.New(
		*rawChain,
		oc.deps.TSS,
		evmClient,
		zetaConnectorAddress,
		erc20CustodyAddress,
		gatewayAddress,
		oc.logger.base,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create signer")
	}

	return evm.New(oc.scheduler, observer, signer), nil
}

func btcDatabaseFileName(chain chains.Chain) string {
	// legacyBTCDatabaseFilename is the Bitcoin database file name now used in mainnet and testnet3
	// so we keep using it here for backward compatibility
	const legacyBTCDatabaseFilename = "btc_chain_client"

	// For additional bitcoin networks, we use the chain name as the database file name
	switch chain.ChainId {
	case chains.BitcoinMainnet.ChainId, chains.BitcoinTestnet.ChainId:
		return legacyBTCDatabaseFilename
	default:
		return fmt.Sprintf("%s_%s", legacyBTCDatabaseFilename, chain.Name)
	}
}
