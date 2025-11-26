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
	solobserver "github.com/zeta-chain/node/zetaclient/chains/solana/observer"
	solrepo "github.com/zeta-chain/node/zetaclient/chains/solana/repo"
	solsigner "github.com/zeta-chain/node/zetaclient/chains/solana/signer"
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
	// should not happen
	if !chain.IsBitcoin() {
		return nil, errors.New("chain is not bitcoin")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	config := app.Config()
	clientMode := config.ClientMode

	btcConfig, found := config.GetBTCConfig(chain.ID())
	if !found {
		return nil, errors.Wrap(errSkipChain, "unable to find BTC config")
	}

	standardBitcoinClient, err := btcclient.New(btcConfig, chain.ID(), oc.logger.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create RPC client")
	}
	var bitcoinClient bitcoin.BitcoinClient = standardBitcoinClient
	if clientMode.IsChaosMode() {
		bitcoinClient = oc.chaosSource.WrapBitcoinClient(bitcoinClient)
	}
	if clientMode.IsDryMode() {
		bitcoinClient = dry.WrapBitcoinClient(bitcoinClient)
	}

	var (
		rawChain = chain.RawChain()
		dbName   = btcDatabaseFileName(*rawChain)
	)

	baseObserver, err := oc.newBaseObserver(chain, dbName, clientMode)
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
	// should not happen
	if !chain.IsEVM() {
		return nil, errors.New("chain is not EVM")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	config := app.Config()
	clientMode := config.ClientMode

	evmConfig, found := config.GetEVMConfig(chain.ID())
	if !found || evmConfig.Empty() {
		return nil, errors.Wrap(errSkipChain, "unable to find evm config")
	}

	standardEvmClient, err := evmclient.NewFromEndpoint(ctx, evmConfig.Endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create evm client (%s)", evmConfig.Endpoint)
	}
	var evmClient evm.EVMClient = standardEvmClient
	if clientMode.IsChaosMode() {
		evmClient = oc.chaosSource.WrapEVMClient(evmClient)
	}
	if clientMode.IsDryMode() {
		evmClient = dry.WrapEVMClient(evmClient)
	}

	baseObserver, err := oc.newBaseObserver(chain, chain.Name(), clientMode)
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
	// should not happen
	if !chain.IsSolana() {
		return nil, errors.New("chain is not Solana")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	config := app.Config()
	clientMode := config.ClientMode

	solanaConfig, found := config.GetSolanaConfig()
	if !found {
		return nil, errors.Wrap(errSkipChain, "unable to find solana config")
	}

	baseObserver, err := oc.newBaseObserver(chain, chain.Name(), clientMode)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create base observer")
	}

	gwAddress := chain.Params().GatewayAddress

	standardSolanaClient := solrpc.New(solanaConfig.Endpoint)
	if standardSolanaClient == nil {
		return nil, errors.New("unable to create RPC client")
	}
	var solanaClient solrepo.SolanaClient = standardSolanaClient
	if clientMode.IsChaosMode() {
		solanaClient = oc.chaosSource.WrapSolanaClient(solanaClient)
	}
	if clientMode.IsDryMode() {
		solanaClient = dry.WrapSolanaClient(solanaClient)
	}

	observer, err := solobserver.New(baseObserver, solanaClient, gwAddress)
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
	signer, err := solsigner.New(baseSigner, solanaClient, gwAddress, relayerKey)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create signer")
	}

	return solana.New(oc.scheduler, observer, signer), nil
}

func (oc *Orchestrator) bootstrapSui(ctx context.Context, chain zctx.Chain) (*sui.Sui, error) {
	// should not happen
	if !chain.IsSui() {
		return nil, errors.New("chain is not sui")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	config := app.Config()
	clientMode := config.ClientMode

	suiConfig, found := config.GetSuiConfig()
	if !found {
		return nil, errors.Wrap(errSkipChain, "unable to find sui config")
	}

	// note that gateway address should be in either of the following formats:
	//   - `$packageID,$gatewayObjectID`
	//   - `$packageID,$gatewayObjectID,$withdrawCapID,$previousPackageID,$originalPackageID`
	gateway, err := suigateway.NewGatewayFromPairID(chain.Params().GatewayAddress)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create gateway")
	}

	standardSuiClient := suiclient.New(suiConfig.Endpoint)
	var suiClient sui.SuiClient = standardSuiClient
	if clientMode.IsChaosMode() {
		suiClient = oc.chaosSource.WrapSuiClient(suiClient)
	}
	if clientMode.IsDryMode() {
		suiClient = dry.WrapSuiClient(suiClient)
	}

	baseObserver, err := oc.newBaseObserver(chain, chain.Name(), clientMode)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create base observer")
	}

	observer := suiobserver.New(baseObserver, suiClient, gateway)

	// migrate inbound cursor to adopt authenticated call upgrade.
	// after upgrade, we might have to deal with multiple packages,
	// so any auxiliary data should be managed under different package IDs.
	// TODO: https://github.com/zeta-chain/node/issues/4164
	if err = observer.MigrateCursorForAuthenticatedCallUpgrade(); err != nil {
		return nil, errors.Wrap(err, "unable to migrate inbound cursor")
	}

	baseSigner := oc.newBaseSigner(chain, clientMode)
	signer := suisigner.New(baseSigner, baseObserver.ZetaRepo(), suiClient, gateway)

	return sui.New(oc.scheduler, observer, signer), nil
}

func (oc *Orchestrator) bootstrapTON(ctx context.Context, chain zctx.Chain) (*ton.TON, error) {
	// should not happen
	if !chain.IsTON() {
		return nil, errors.New("chain is not TON")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	config := app.Config()
	clientMode := config.ClientMode

	tonConfig, found := config.GetTONConfig()
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

	if tonConfig.Endpoint == "" {
		return nil, errors.New("rpc url is empty")
	}

	rpcClient, err := metrics.GetInstrumentedHTTPClient(tonConfig.Endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create instrumented rpc client")
	}

	standardTONClient := tonclient.New(tonConfig.Endpoint, chain.ID(), tonclient.WithHTTPClient(rpcClient))
	var tonClient ton.TONClient = standardTONClient
	if clientMode.IsChaosMode() {
		tonClient = oc.chaosSource.WrapTONClient(tonClient)
	}
	if clientMode.IsDryMode() {
		tonClient = dry.WrapTONClient(tonClient)
	}

	baseObserver, err := oc.newBaseObserver(chain, chain.Name(), clientMode)
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
	chain zctx.Chain,
	dbName string,
	clientMode mode.ClientMode,
) (*base.Observer, error) {
	var (
		rawChain       = chain.RawChain()
		rawChainParams = chain.Params()
	)

	database, err := db.NewFromSqlite(oc.dbPath, dbName, true)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open database %s", dbName)
	}

	blocksCacheSize := base.DefaultBlockCacheSize
	if chain.IsBitcoin() {
		blocksCacheSize = btcBlocksPerDay
	}

	zetacoreClient := oc.zetacoreClient.(zrepo.ZetacoreClient)
	tssClient := oc.tssClient
	if clientMode.IsChaosMode() {
		zetacoreClient = oc.chaosSource.WrapZetacoreClient(zetacoreClient)
		tssClient = oc.chaosSource.WrapTSSClient(tssClient)
	}

	return base.NewObserver(
		*rawChain,
		*rawChainParams,
		zrepo.New(zetacoreClient, *rawChain, clientMode),
		tssClient,
		blocksCacheSize,
		oc.telemetry,
		database,
		oc.logger.base,
	)
}

func (oc *Orchestrator) newBaseSigner(chain zctx.Chain, clientMode mode.ClientMode) *base.Signer {
	tssClient := oc.tssClient
	if clientMode.IsChaosMode() {
		tssClient = oc.chaosSource.WrapTSSClient(tssClient)
	}
	return base.NewSigner(*chain.RawChain(), tssClient, oc.logger.base, clientMode)
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
