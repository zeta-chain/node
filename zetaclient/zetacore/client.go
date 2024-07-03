// Package zetacore provides the client to interact with zetacore node via GRPC.
package zetacore

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"cosmossdk.io/simapp/params"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/pkg/authz"
	"github.com/zeta-chain/zetacore/pkg/chains"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	keyinterfaces "github.com/zeta-chain/zetacore/zetaclient/keys/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
)

var _ interfaces.ZetacoreClient = &Client{}

// Client is the client to send tx to zetacore
type Client struct {
	logger zerolog.Logger
	config config.ClientConfiguration

	client              clients
	cosmosClientContext cosmosclient.Context

	blockHeight   int64
	accountNumber map[authz.KeyType]uint64
	seqNumber     map[authz.KeyType]uint64

	encodingCfg params.EncodingConfig
	keys        keyinterfaces.ObserverKeys
	chainID     string
	chain       chains.Chain
	stop        chan struct{}
	pause       chan struct{}
	Telemetry   *metrics.TelemetryServer

	mu sync.RWMutex

	// enableMockSDKClient is a flag that determines whether the mock cosmos sdk client should be used, primarily for
	// unit testing
	enableMockSDKClient bool
	mockSDKClient       rpcclient.Client
}

type clients struct {
	observer   observertypes.QueryClient
	light      lightclienttypes.QueryClient
	crosschain crosschaintypes.QueryClient
	bank       banktypes.QueryClient
	upgrade    upgradetypes.QueryClient
	fees       feemarkettypes.QueryClient
	tendermint tmservice.ServiceClient
}

// NewClient create a new instance of Client
func NewClient(
	keys keyinterfaces.ObserverKeys,
	chainIP string,
	signerName string,
	chainID string,
	hsmMode bool,
	telemetry *metrics.TelemetryServer,
	logger zerolog.Logger,
) (*Client, error) {
	zetaChain, err := chains.ZetaChainFromChainID(chainID)
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

	grpcConn, err := grpc.Dial(
		cosmosGRPC(chainIP),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "grpc dial fail")
	}

	accountsMap := make(map[authz.KeyType]uint64)
	seqMap := make(map[authz.KeyType]uint64)
	for _, keyType := range authz.GetAllKeyTypes() {
		accountsMap[keyType] = 0
		seqMap[keyType] = 0
	}

	c := &Client{
		logger: log,
		config: cfg,

		cosmosClientContext: cosmosclient.Context{},

		client: clients{
			observer:   observertypes.NewQueryClient(grpcConn),
			light:      lightclienttypes.NewQueryClient(grpcConn),
			crosschain: crosschaintypes.NewQueryClient(grpcConn),
			bank:       banktypes.NewQueryClient(grpcConn),
			upgrade:    upgradetypes.NewQueryClient(grpcConn),
			fees:       feemarkettypes.NewQueryClient(grpcConn),
			tendermint: tmservice.NewServiceClient(grpcConn),
		},

		accountNumber: accountsMap,
		seqNumber:     seqMap,

		encodingCfg: app.MakeEncodingConfig(),
		keys:        keys,
		stop:        make(chan struct{}),
		chainID:     chainID,
		chain:       zetaChain,
		pause:       make(chan struct{}),
		Telemetry:   telemetry,

		mu:                  sync.RWMutex{},
		enableMockSDKClient: false,
		mockSDKClient:       nil,
	}

	cosmosClientContext, err := c.buildCosmosClientContext()
	if err != nil {
		return nil, errors.Wrap(err, "fail to resolve cosmos client context")
	}

	c.cosmosClientContext = cosmosClientContext

	return c, nil
}

// buildCosmosClientContext constructs a valid context with all relevant values set
func (c *Client) buildCosmosClientContext() (cosmosclient.Context, error) {
	if c.keys == nil {
		return cosmosclient.Context{}, errors.New("client key are not set")
	}

	addr, err := c.keys.GetAddress()
	if err != nil {
		return cosmosclient.Context{}, errors.Wrap(err, "fail to get address from key")
	}

	var (
		input   = strings.NewReader("")
		client  cosmosclient.TendermintRPC
		nodeURI string
	)

	// if password is needed, set it as input
	password := c.keys.GetHotkeyPassword()
	if password != "" {
		input = strings.NewReader(fmt.Sprintf("%[1]s\n%[1]s\n", password))
	}

	if c.enableMockSDKClient {
		client = c.mockSDKClient
	} else {
		remote := c.config.ChainRPC
		if !strings.HasPrefix(c.config.ChainHost, "http") {
			remote = fmt.Sprintf("tcp://%s", remote)
		}

		wsClient, err := rpchttp.New(remote, "/websocket")
		if err != nil {
			return cosmosclient.Context{}, err
		}

		client = wsClient
		nodeURI = remote
	}

	return cosmosclient.Context{
		Client:            client,
		NodeURI:           nodeURI,
		FromAddress:       addr,
		ChainID:           c.chainID,
		Codec:             c.encodingCfg.Codec,
		InterfaceRegistry: c.encodingCfg.InterfaceRegistry,
		Keyring:           c.keys.GetKeybase(),
		HomeDir:           c.config.ChainHomeFolder,
		BroadcastMode:     "sync",
		FromName:          c.config.SignerName,
		TxConfig:          c.encodingCfg.TxConfig,
		AccountRetriever:  authtypes.AccountRetriever{},
		LegacyAmino:       c.encodingCfg.Amino,
		Input:             input,
	}, nil
}

func (c *Client) UpdateChainID(chainID string) error {
	if c.chainID != chainID {
		c.chainID = chainID

		zetaChain, err := chains.ZetaChainFromChainID(chainID)
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

func (c *Client) Stop() {
	c.logger.Info().Msgf("zetacore client is stopping")
	close(c.stop) // this notifies all configupdater to stop
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

// UpdateZetacoreContext updates zetacore context
// zetacore stores zetacore context for all clients
func (c *Client) UpdateZetacoreContext(
	ctx context.Context,
	appContext *zctx.AppContext,
	init bool,
	sampledLogger zerolog.Logger,
) error {
	bn, err := c.GetBlockHeight(ctx)
	if err != nil {
		return fmt.Errorf("failed to get zetablock height: %w", err)
	}

	plan, err := c.GetUpgradePlan(ctx)
	if err != nil {
		// if there is no active upgrade plan, plan will be nil, err will be nil as well.
		return fmt.Errorf("failed to get upgrade plan: %w", err)
	}

	if plan != nil && bn == plan.Height-1 { // stop zetaclients; notify operator to upgrade and restart
		c.logger.Warn().
			Msgf("Active upgrade plan detected and upgrade height reached: %s at height %d; ZetaClient is stopped;"+
				"please kill this process, replace zetaclientd binary with upgraded version, and restart zetaclientd", plan.Name, plan.Height)
		c.pause <- struct{}{} // notify Orchestrator to stop Observers, Signers, and Orchestrator itself
	}

	chainParams, err := c.GetChainParams(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain params: %w", err)
	}

	newEVMParams := make(map[int64]*observertypes.ChainParams)
	var newBTCParams *observertypes.ChainParams

	// check and update chain params for each chain
	for _, chainParam := range chainParams {
		err := observertypes.ValidateChainParams(chainParam)
		if err != nil {
			sampledLogger.Warn().Err(err).Msgf("Invalid chain params for chain %d", chainParam.ChainId)
			continue
		}
		if chains.IsBitcoinChain(chainParam.ChainId) {
			newBTCParams = chainParam
		} else if chains.IsEVMChain(chainParam.ChainId) {
			newEVMParams[chainParam.ChainId] = chainParam
		}
	}

	supportedChains, err := c.GetSupportedChains(ctx)
	if err != nil {
		return fmt.Errorf("failed to get supported chains: %w", err)
	}

	newChains := make([]chains.Chain, len(supportedChains))
	for i, chain := range supportedChains {
		newChains[i] = *chain
	}

	keyGen, err := c.GetKeyGen(ctx)
	if err != nil {
		c.logger.Info().Msg("Unable to fetch keygen from zetacore")
		return fmt.Errorf("failed to get keygen: %w", err)
	}

	tss, err := c.GetCurrentTSS(ctx)
	if err != nil {
		c.logger.Info().Err(err).Msg("Unable to fetch TSS from zetacore")
		return fmt.Errorf("failed to get current tss: %w", err)
	}
	tssPubKey := tss.GetTssPubkey()

	crosschainFlags, err := c.GetCrosschainFlags(ctx)
	if err != nil {
		c.logger.Info().Msg("Unable to fetch cross-chain flags from zetacore")
		return fmt.Errorf("failed to get crosschain flags: %w", err)
	}

	blockHeaderEnabledChains, err := c.GetBlockHeaderEnabledChains(ctx)
	if err != nil {
		c.logger.Info().Msg("Unable to fetch block header enabled chains from zetacore")
		return err
	}

	appContext.Update(
		keyGen,
		newChains,
		newEVMParams,
		newBTCParams,
		tssPubKey,
		crosschainFlags,
		blockHeaderEnabledChains,
		init,
	)

	return nil
}

// Pause pauses the client
func (c *Client) Pause() {
	<-c.pause
}

// Unpause unpauses the client
func (c *Client) Unpause() {
	c.pause <- struct{}{}
}

// EnableMockSDKClient enables the mock cosmos sdk client
// TODO(revamp): move this to a test package
func (c *Client) EnableMockSDKClient(client rpcclient.Client) {
	c.mockSDKClient = client
	c.enableMockSDKClient = true
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
