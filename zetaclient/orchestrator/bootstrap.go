package orchestrator

import (
	"context"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	tontools "github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/pkg/chains"
	suigateway "github.com/zeta-chain/node/pkg/contracts/sui"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin"
	btcclient "github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	btcobserver "github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	btcsigner "github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
	"github.com/zeta-chain/node/zetaclient/chains/evm"
	evmclient "github.com/zeta-chain/node/zetaclient/chains/evm/client"
	evmobserver "github.com/zeta-chain/node/zetaclient/chains/evm/observer"
	evmsigner "github.com/zeta-chain/node/zetaclient/chains/evm/signer"
	"github.com/zeta-chain/node/zetaclient/chains/solana"
	solbserver "github.com/zeta-chain/node/zetaclient/chains/solana/observer"
	solanasigner "github.com/zeta-chain/node/zetaclient/chains/solana/signer"
	"github.com/zeta-chain/node/zetaclient/chains/sui"
	suiclient "github.com/zeta-chain/node/zetaclient/chains/sui/client"
	suiobserver "github.com/zeta-chain/node/zetaclient/chains/sui/observer"
	suisigner "github.com/zeta-chain/node/zetaclient/chains/sui/signer"
	"github.com/zeta-chain/node/zetaclient/chains/ton"
	tonobserver "github.com/zeta-chain/node/zetaclient/chains/ton/observer"
	tonrepo "github.com/zeta-chain/node/zetaclient/chains/ton/repo"
	tonclient "github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
	tonsigner "github.com/zeta-chain/node/zetaclient/chains/ton/signer"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/dry"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/mode"
)

const btcBlocksPerDay = 144

func (oc *Orchestrator) bootstrapBitcoin(ctx context.Context, chain zctx.Chain) (*bitcoin.Bitcoin, error) {
	// TODO: hardcoded for now
	// See: https://github.com/zeta-chain/node/issues/2865
	clientMode := mode.StandardMode

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

	standardBitcoinClient, err := btcclient.New(cfg, chain.ID(), oc.logger.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create rpc client")
	}
	var bitcoinClient bitcoin.Client = standardBitcoinClient

	var (
		rawChain = chain.RawChain()
		dbName   = btcDatabaseFileName(*rawChain)
	)

	baseObserver, err := oc.newBaseObserver(clientMode, chain, dbName)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create base observer")
	}

	observer, err := btcobserver.New(baseObserver, bitcoinClient, *rawChain)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create observer")
	}

	baseSigner := oc.newBaseSigner(chain, clientMode)
	signer := btcsigner.New(baseSigner, bitcoinClient)

	return bitcoin.New(oc.scheduler, observer, signer), nil
}

func (oc *Orchestrator) bootstrapEVM(ctx context.Context, chain zctx.Chain) (*evm.EVM, error) {
	// TODO: hardcoded for now
	// See: https://github.com/zeta-chain/node/issues/2865
	clientMode := mode.StandardMode

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

	standardEvmClient, err := evmclient.NewFromEndpoint(ctx, cfg.Endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create evm client (%s)", cfg.Endpoint)
	}
	var evmClient evm.Client = standardEvmClient

	baseObserver, err := oc.newBaseObserver(clientMode, chain, chain.Name())
	if err != nil {
		return nil, errors.Wrap(err, "unable to create base observer")
	}

	observer, err := evmobserver.New(baseObserver, evmClient)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create observer")
	}

	var (
		zetaConnectorAddress = ethcommon.HexToAddress(chain.Params().ConnectorContractAddress)
		erc20CustodyAddress  = ethcommon.HexToAddress(chain.Params().Erc20CustodyContractAddress)
		gatewayAddress       = ethcommon.HexToAddress(chain.Params().GatewayAddress)
	)

	baseSigner := oc.newBaseSigner(chain, clientMode)
	signer, err := evmsigner.New(
		baseSigner,
		evmClient,
		zetaConnectorAddress,
		erc20CustodyAddress,
		gatewayAddress,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create signer")
	}

	return evm.New(oc.scheduler, observer, signer), nil
}

func (oc *Orchestrator) bootstrapSolana(ctx context.Context, chain zctx.Chain) (*solana.Solana, error) {
	// TODO: hardcoded for now
	// See: https://github.com/zeta-chain/node/issues/2865
	clientMode := mode.StandardMode

	// should not happen
	if !chain.IsSolana() {
		return nil, errors.New("chain is not Solana")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	cfg, found := app.Config().GetSolanaConfig()
	if !found {
		return nil, errors.Wrap(errSkipChain, "unable to find solana config")
	}

	baseObserver, err := oc.newBaseObserver(clientMode, chain, chain.Name())
	if err != nil {
		return nil, errors.Wrap(err, "unable to create base observer")
	}

	gwAddress := chain.Params().GatewayAddress

	standardSolanaClient := solrpc.New(cfg.Endpoint)
	if standardSolanaClient == nil {
		return nil, errors.New("unable to create rpc client")
	}
	var solanaClient solana.Client = standardSolanaClient

	observer, err := solbserver.New(baseObserver, solanaClient, gwAddress)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create observer")
	}

	// Try loading Solana relayer key if present
	// Note that relayerKey might be nil if the key is not present. It's okay.
	password := chain.RelayerKeyPassword()
	relayerKey, err := keys.LoadRelayerKey(app.Config().GetRelayerKeyPath(), chain.RawChain().Network, password)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load relayer key")
	}

	baseSigner := oc.newBaseSigner(chain, clientMode)

	// create Solana signer
	signer, err := solanasigner.New(baseSigner, solanaClient, gwAddress, relayerKey)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create signer")
	}

	return solana.New(oc.scheduler, observer, signer), nil
}

func (oc *Orchestrator) bootstrapSui(ctx context.Context, chain zctx.Chain) (*sui.Sui, error) {
	// TODO: hardcoded for now
	// See: https://github.com/zeta-chain/node/issues/2865
	clientMode := mode.StandardMode

	// should not happen
	if !chain.IsSui() {
		return nil, errors.New("chain is not sui")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	cfg, found := app.Config().GetSuiConfig()
	if !found {
		return nil, errors.Wrap(errSkipChain, "unable to find sui config")
	}

	// note that gw address should be in format of `$packageID,$gatewayObjectID`
	gateway, err := suigateway.NewGatewayFromPairID(chain.Params().GatewayAddress)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create gateway")
	}

	standardSuiClient := suiclient.New(cfg.Endpoint)
	var suiClient sui.Client = standardSuiClient
	if clientMode.IsDryMode() {
		suiClient = dry.WrapSuiClient(suiClient)
	}

	baseObserver, err := oc.newBaseObserver(clientMode, chain, chain.Name())
	if err != nil {
		return nil, errors.Wrap(err, "unable to create base observer")
	}

	observer := suiobserver.New(baseObserver, suiClient, gateway)

	baseSigner := oc.newBaseSigner(chain, clientMode)
	signer := suisigner.New(baseSigner, oc.deps.Zetacore, suiClient, gateway)

	return sui.New(oc.scheduler, observer, signer), nil
}

func (oc *Orchestrator) bootstrapTON(ctx context.Context, chain zctx.Chain) (*ton.TON, error) {
	// TODO: hardcoded for now
	// See: https://github.com/zeta-chain/node/issues/2865
	clientMode := mode.StandardMode

	// should not happen
	if !chain.IsTON() {
		return nil, errors.New("chain is not TON")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	cfg, found := app.Config().GetTONConfig()
	if !found {
		return nil, errors.Wrap(errSkipChain, "unable to find TON config")
	}

	gwAddress := chain.Params().GatewayAddress
	if gwAddress == "" {
		return nil, errors.New("gateway address is empty")
	}

	gatewayID, err := tontools.ParseAccountID(gwAddress)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse gateway address %q", gwAddress)
	}

	gw := toncontracts.NewGateway(gatewayID)

	if cfg.Endpoint == "" {
		return nil, errors.New("rpc url is empty")
	}

	rpcClient, err := metrics.GetInstrumentedHTTPClient(cfg.Endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create instrumented rpc client")
	}

	standardTONClient := tonclient.New(cfg.Endpoint, chain.ID(), tonclient.WithHTTPClient(rpcClient))
	var tonClient ton.Client = standardTONClient
	if clientMode.IsDryMode() {
		tonClient = dry.WrapTONClient(tonClient)
	}

	baseObserver, err := oc.newBaseObserver(clientMode, chain, chain.Name())
	if err != nil {
		return nil, errors.Wrap(err, "unable to create base observer")
	}

	tonRepo := tonrepo.NewTONRepo(tonClient, gw, baseObserver.Chain())
	observer, err := tonobserver.New(baseObserver, tonRepo, gw)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create observer")
	}

	baseSigner := oc.newBaseSigner(chain, clientMode)
	signer := tonsigner.New(baseSigner, tonClient, gw)

	return ton.New(oc.scheduler, observer, signer), nil
}

func (oc *Orchestrator) newBaseObserver(
	clientMode mode.ClientMode,
	chain zctx.Chain,
	dbName string,
) (*base.Observer, error) {
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

	zetaRepo := zrepo.New(oc.deps.Zetacore, *rawChain, clientMode)

	return base.NewObserver(
		*rawChain,
		*rawChainParams,
		zetaRepo,
		oc.deps.TSS,
		blocksCacheSize,
		oc.deps.Telemetry,
		database,
		oc.logger.base,
	)
}

func (oc *Orchestrator) newBaseSigner(chain zctx.Chain, clientMode mode.ClientMode) *base.Signer {
	return base.NewSigner(*chain.RawChain(), oc.deps.TSS, oc.logger.base, clientMode)
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
