// Package zetacore provides the client to interact with zetacore node via GRPC.
package zetacore

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	etherminttypes "github.com/zeta-chain/ethermint/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/pkg/authz"
	"github.com/zeta-chain/node/pkg/chains"
	zetacore_rpc "github.com/zeta-chain/node/pkg/rpc"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/config"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	keyinterfaces "github.com/zeta-chain/node/zetaclient/keys/interfaces"
)

var _ interfaces.ZetacoreClient = &Client{}

// Client is the client to send tx to zetacore
type Client struct {
	zetacore_rpc.Clients

	logger zerolog.Logger
	config config.ClientConfiguration

	cosmosClientContext cosmosclient.Context

	blockHeight   int64
	accountNumber map[authz.KeyType]uint64
	seqNumber     map[authz.KeyType]uint64

	encodingCfg          etherminttypes.EncodingConfig
	keys                 keyinterfaces.ObserverKeys
	chainID              string
	chain                chains.Chain
	stop                 chan struct{}
	onBeforeStopCallback []func()

	mu sync.RWMutex
}

var unsecureGRPC = grpc.WithTransportCredentials(insecure.NewCredentials())

type constructOpts struct {
	customTendermint bool
	tendermintClient cosmosclient.TendermintRPC

	customAccountRetriever bool
	accountRetriever       cosmosclient.AccountRetriever
}

type Opt func(cfg *constructOpts)

// WithTendermintClient sets custom tendermint client
func WithTendermintClient(client cosmosclient.TendermintRPC) Opt {
	return func(c *constructOpts) {
		c.customTendermint = true
		c.tendermintClient = client
	}
}

// WithCustomAccountRetriever sets custom tendermint client
func WithCustomAccountRetriever(ac cosmosclient.AccountRetriever) Opt {
	return func(c *constructOpts) {
		c.customAccountRetriever = true
		c.accountRetriever = ac
	}
}

// NewClient create a new instance of Client
func NewClient(
	keys keyinterfaces.ObserverKeys,
	chainIP string,
	signerName string,
	chainID string,
	hsmMode bool,
	logger zerolog.Logger,
	opts ...Opt,
) (*Client, error) {
	var constructOptions constructOpts
	for _, opt := range opts {
		opt(&constructOptions)
	}

	zetaChain, err := chains.ZetaChainFromCosmosChainID(chainID)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid chain id %q", chainID)
	}

	log := logger.With().Str("module", "zetacoreClient").Logger()

	cfg := config.ClientConfiguration{
		ChainHost:    cosmosREST(chainIP),
		SignerName:   signerName,
		SignerPasswd: "password",
		ChainRPC:     tendermintRPC(chainIP),
		HsmMode:      hsmMode,
	}

	encodingCfg := app.MakeEncodingConfig()

	zetacoreClients, err := zetacore_rpc.NewGRPCClients(cosmosGRPC(chainIP), unsecureGRPC)
	if err != nil {
		return nil, errors.Wrap(err, "grpc dial fail")
	}

	accountsMap := make(map[authz.KeyType]uint64)
	seqMap := make(map[authz.KeyType]uint64)
	for _, keyType := range authz.GetAllKeyTypes() {
		accountsMap[keyType] = 0
		seqMap[keyType] = 0
	}

	cosmosContext, err := buildCosmosClientContext(chainID, keys, cfg, encodingCfg, constructOptions)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build cosmos client context")
	}

	return &Client{
		Clients: zetacoreClients,
		logger:  log,
		config:  cfg,

		cosmosClientContext: cosmosContext,

		accountNumber: accountsMap,
		seqNumber:     seqMap,

		encodingCfg: encodingCfg,
		keys:        keys,
		stop:        make(chan struct{}),
		chainID:     chainID,
		chain:       zetaChain,
	}, nil
}

// buildCosmosClientContext constructs a valid context with all relevant values set
func buildCosmosClientContext(
	chainID string,
	keys keyinterfaces.ObserverKeys,
	config config.ClientConfiguration,
	encodingConfig etherminttypes.EncodingConfig,
	opts constructOpts,
) (cosmosclient.Context, error) {
	if keys == nil {
		return cosmosclient.Context{}, errors.New("client key are not set")
	}

	addr, err := keys.GetAddress()
	if err != nil {
		return cosmosclient.Context{}, errors.Wrap(err, "fail to get address from key")
	}

	var (
		input   = strings.NewReader("")
		client  cosmosclient.TendermintRPC
		nodeURI string
	)

	// if password is needed, set it as input
	password := keys.GetHotkeyPassword()
	if password != "" {
		input = strings.NewReader(fmt.Sprintf("%[1]s\n%[1]s\n", password))
	}

	// note that in rare cases, this might give FALSE positive
	// (google "golang nil interface comparison")
	client = opts.tendermintClient
	if !opts.customTendermint {
		remote := config.ChainRPC
		if !strings.HasPrefix(config.ChainHost, "http") {
			remote = fmt.Sprintf("tcp://%s", remote)
		}

		wsClient, err := rpchttp.New(remote, "/websocket")
		if err != nil {
			return cosmosclient.Context{}, err
		}

		client = wsClient
		nodeURI = remote
	}

	var accountRetriever cosmosclient.AccountRetriever
	if opts.customAccountRetriever {
		accountRetriever = opts.accountRetriever
	} else {
		accountRetriever = authtypes.AccountRetriever{}
	}

	return cosmosclient.Context{
		Client:        client,
		NodeURI:       nodeURI,
		FromAddress:   addr,
		ChainID:       chainID,
		Keyring:       keys.GetKeybase(),
		BroadcastMode: "sync",
		HomeDir:       config.ChainHomeFolder,
		FromName:      config.SignerName,

		AccountRetriever: accountRetriever,

		Codec:             encodingConfig.Codec,
		InterfaceRegistry: encodingConfig.InterfaceRegistry,
		TxConfig:          encodingConfig.TxConfig,
		LegacyAmino:       encodingConfig.Amino,

		Input: input,
	}, nil
}

func (c *Client) UpdateChainID(chainID string) error {
	if c.chainID != chainID {
		c.chainID = chainID

		zetaChain, err := chains.ZetaChainFromCosmosChainID(chainID)
		if err != nil {
			return fmt.Errorf("invalid chain id %s, %w", chainID, err)
		}
		c.chain = zetaChain
	}

	return nil
}

// Chain returns the Chain chain object
func (c *Client) Chain() chains.Chain {
	return c.chain
}

func (c *Client) GetLogger() *zerolog.Logger {
	return &c.logger
}

func (c *Client) GetKeys() keyinterfaces.ObserverKeys {
	return c.keys
}

// OnBeforeStop adds a callback to be called before the client stops.
func (c *Client) OnBeforeStop(callback func()) {
	c.onBeforeStopCallback = append(c.onBeforeStopCallback, callback)
}

// Stop stops the client and optionally calls the onBeforeStop callbacks.
func (c *Client) Stop() {
	c.logger.Info().Msgf("Stopping zetacore client")

	for i := len(c.onBeforeStopCallback) - 1; i >= 0; i-- {
		c.logger.Info().Int("callback.index", i).Msgf("calling onBeforeStopCallback")
		c.onBeforeStopCallback[i]()
	}

	close(c.stop)
}

// GetAccountNumberAndSequenceNumber We do not use multiple KeyType for now , but this can be optionally used in the future to seprate TSS signer from Zetaclient GRantee
func (c *Client) GetAccountNumberAndSequenceNumber(_ authz.KeyType) (uint64, uint64, error) {
	address, err := c.keys.GetAddress()
	if err != nil {
		return 0, 0, err
	}
	return c.cosmosClientContext.AccountRetriever.GetAccountNumberSequence(c.cosmosClientContext, address)
}

// SetAccountNumber sets the account number and sequence number for the given keyType
// todo remove method and make it part of the client constructor.
func (c *Client) SetAccountNumber(keyType authz.KeyType) error {
	address, err := c.keys.GetAddress()
	if err != nil {
		return errors.Wrap(err, "fail to get address")
	}

	accN, seq, err := c.cosmosClientContext.AccountRetriever.GetAccountNumberSequence(c.cosmosClientContext, address)
	if err != nil {
		return errors.Wrap(err, "fail to get account number and sequence number")
	}

	c.accountNumber[keyType] = accN
	c.seqNumber[keyType] = seq

	return nil
}

// WaitForZetacoreToCreateBlocks waits for zetacore to create blocks
func (c *Client) WaitForZetacoreToCreateBlocks(ctx context.Context) error {
	retryCount := 0
	for {
		block, err := c.GetLatestZetaBlock(ctx)
		if err == nil && block.Header.Height > 1 {
			c.logger.Info().Msgf("Zetacore height: %d", block.Header.Height)
			break
		}
		retryCount++
		c.logger.Debug().Msgf("Failed to get latest Block , Retry : %d/%d", retryCount, DefaultRetryCount)
		if retryCount > ExtendedRetryCount {
			return fmt.Errorf("zetacore is not ready, waited for %d seconds", DefaultRetryCount*DefaultRetryInterval)
		}
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return nil
}

// UpdateAppContext updates zctx.AppContext
// zetacore stores AppContext for all clients
func (c *Client) UpdateAppContext(ctx context.Context, appContext *zctx.AppContext, logger zerolog.Logger) error {
	bn, err := c.GetBlockHeight(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get zetablock height")
	}

	plan, err := c.GetUpgradePlan(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get upgrade plan")
	}

	// Stop client and notify dependant services to stop (Orchestrator, Observers, and Signers)
	if plan != nil && bn == plan.Height-1 {
		c.logger.Warn().Msgf(
			"Active upgrade plan detected and upgrade height reached: %s at height %d; Stopping ZetaClient;"+
				" please kill this process, replace zetaclientd binary with upgraded version, and restart zetaclientd",
			plan.Name,
			plan.Height,
		)

		c.Stop()

		return nil
	}

	supportedChains, err := c.GetSupportedChains(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch supported chains")
	}

	additionalChains, err := c.GetAdditionalChains(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch additional chains")
	}

	chainParams, err := c.GetChainParams(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch chain params")
	}

	keyGen, err := c.GetKeyGen(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch keygen from zetacore")
	}

	crosschainFlags, err := c.GetCrosschainFlags(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch crosschain flags from zetacore")
	}

	tss, err := c.GetTSS(ctx)
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

	return appContext.Update(
		keyGen,
		supportedChains,
		additionalChains,
		freshParams,
		tss.GetTssPubkey(),
		crosschainFlags,
	)
}

func cosmosREST(host string) string {
	return fmt.Sprintf("%s:1317", host)
}

func cosmosGRPC(host string) string {
	return fmt.Sprintf("%s:9090", host)
}

func tendermintRPC(host string) string {
	return fmt.Sprintf("%s:26657", host)
}
