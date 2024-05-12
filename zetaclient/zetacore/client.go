package zetacore

import (
	"fmt"
	"sync"
	"time"

	"cosmossdk.io/simapp/params"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/pkg/authz"
	"github.com/zeta-chain/zetacore/pkg/chains"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"google.golang.org/grpc"
)

var _ interfaces.ZetacoreClient = &Client{}

// Client is the client to send tx to zetacore
type Client struct {
	logger        zerolog.Logger
	blockHeight   int64
	accountNumber map[authz.KeyType]uint64
	seqNumber     map[authz.KeyType]uint64
	grpcConn      *grpc.ClientConn
	cfg           config.ClientConfiguration
	encodingCfg   params.EncodingConfig
	keys          *keys.Keys
	broadcastLock *sync.RWMutex
	chainID       string
	chain         chains.Chain
	stop          chan struct{}
	pause         chan struct{}
	Telemetry     *metrics.TelemetryServer

	// enableMockSDKClient is a flag that determines whether the mock cosmos sdk client should be used, primarily for
	// unit testing
	enableMockSDKClient bool
	mockSDKClient       rpcclient.Client
}

// NewClient create a new instance of Client
func NewClient(
	k *keys.Keys,
	chainIP string,
	signerName string,
	chainID string,
	hsmMode bool,
	telemetry *metrics.TelemetryServer,
) (*Client, error) {

	// main module logger
	logger := log.With().Str("module", "ZetaCoreClient").Logger()
	cfg := config.ClientConfiguration{
		ChainHost:    fmt.Sprintf("%s:1317", chainIP),
		SignerName:   signerName,
		SignerPasswd: "password",
		ChainRPC:     fmt.Sprintf("%s:26657", chainIP),
		HsmMode:      hsmMode,
	}

	grpcConn, err := grpc.Dial(
		fmt.Sprintf("%s:9090", chainIP),
		grpc.WithInsecure(),
	)
	if err != nil {
		logger.Error().Err(err).Msg("grpc dial fail")
		return nil, err
	}
	accountsMap := make(map[authz.KeyType]uint64)
	seqMap := make(map[authz.KeyType]uint64)
	for _, keyType := range authz.GetAllKeyTypes() {
		accountsMap[keyType] = 0
		seqMap[keyType] = 0
	}

	zetaChain, err := chains.ZetaChainFromChainID(chainID)
	if err != nil {
		return nil, fmt.Errorf("invalid chain id %s, %w", chainID, err)
	}

	return &Client{
		logger:              logger,
		grpcConn:            grpcConn,
		accountNumber:       accountsMap,
		seqNumber:           seqMap,
		cfg:                 cfg,
		encodingCfg:         app.MakeEncodingConfig(),
		keys:                k,
		broadcastLock:       &sync.RWMutex{},
		stop:                make(chan struct{}),
		chainID:             chainID,
		chain:               zetaChain,
		pause:               make(chan struct{}),
		Telemetry:           telemetry,
		enableMockSDKClient: false,
		mockSDKClient:       nil,
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

func (c *Client) GetKeys() *keys.Keys {
	return c.keys
}

func (c *Client) Stop() {
	c.logger.Info().Msgf("zetacore client is stopping")
	close(c.stop) // this notifies all configupdater to stop
}

// GetAccountNumberAndSequenceNumber We do not use multiple KeyType for now , but this can be optionally used in the future to seprate TSS signer from Zetaclient GRantee
func (c *Client) GetAccountNumberAndSequenceNumber(_ authz.KeyType) (uint64, uint64, error) {
	ctx, err := c.GetContext()
	if err != nil {
		return 0, 0, err
	}
	address := c.keys.GetAddress()
	return ctx.AccountRetriever.GetAccountNumberSequence(ctx, address)
}

func (c *Client) SetAccountNumber(keyType authz.KeyType) {
	ctx, err := c.GetContext()
	if err != nil {
		c.logger.Error().Err(err).Msg("fail to get context")
		return
	}
	address := c.keys.GetAddress()
	accN, seq, err := ctx.AccountRetriever.GetAccountNumberSequence(ctx, address)
	if err != nil {
		c.logger.Error().Err(err).Msg("fail to get account number and sequence number")
		return
	}
	c.accountNumber[keyType] = accN
	c.seqNumber[keyType] = seq
}

func (c *Client) WaitForCoreToCreateBlocks() {
	retryCount := 0
	for {
		block, err := c.GetLatestZetaBlock()
		if err == nil && block.Header.Height > 1 {
			c.logger.Info().Msgf("Zetacore height: %d", block.Header.Height)
			break
		}
		retryCount++
		c.logger.Debug().Msgf("Failed to get latest Block , Retry : %d/%d", retryCount, DefaultRetryCount)
		if retryCount > ExtendedRetryCount {
			panic(fmt.Sprintf("Zetacore is not ready, waited for %d seconds", DefaultRetryCount*DefaultRetryInterval))
		}
		time.Sleep(DefaultRetryInterval * time.Second)
	}
}

// UpdateZetaCoreContext updates core context
// zetacore stores core context for all clients
func (c *Client) UpdateZetaCoreContext(coreContext *context.ZetaCoreContext, init bool, sampledLogger zerolog.Logger) error {
	bn, err := c.GetBlockHeight()
	if err != nil {
		return fmt.Errorf("failed to get zetablock height: %w", err)
	}
	plan, err := c.GetUpgradePlan()
	if err != nil {
		// if there is no active upgrade plan, plan will be nil, err will be nil as well.
		return fmt.Errorf("failed to get upgrade plan: %w", err)
	}
	if plan != nil && bn == plan.Height-1 { // stop zetaclients; notify operator to upgrade and restart
		c.logger.Warn().Msgf("Active upgrade plan detected and upgrade height reached: %s at height %d; ZetaClient is stopped;"+
			"please kill this process, replace zetaclientd binary with upgraded version, and restart zetaclientd", plan.Name, plan.Height)
		c.pause <- struct{}{} // notify Orchestrator to stop Observers, Signers, and Orchestrator itself
	}

	chainParams, err := c.GetChainParams()
	if err != nil {
		return fmt.Errorf("failed to get chain params: %w", err)
	}

	newEVMParams := make(map[int64]*observertypes.ChainParams)
	var newBTCParams *observertypes.ChainParams

	// check and update chain params for each chain
	for _, chainParam := range chainParams {
		if !chainParam.GetIsSupported() {
			sampledLogger.Info().Msgf("Chain %d is not supported yet", chainParam.ChainId)
			continue
		}
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

	supportedChains, err := c.GetSupportedChains()
	if err != nil {
		return fmt.Errorf("failed to get supported chains: %w", err)
	}
	newChains := make([]chains.Chain, len(supportedChains))
	for i, chain := range supportedChains {
		newChains[i] = *chain
	}
	keyGen, err := c.GetKeyGen()
	if err != nil {
		c.logger.Info().Msg("Unable to fetch keygen from zetacore")
		return fmt.Errorf("failed to get keygen: %w", err)
	}

	tss, err := c.GetCurrentTss()
	if err != nil {
		c.logger.Info().Err(err).Msg("Unable to fetch TSS from zetacore")
		return fmt.Errorf("failed to get current tss: %w", err)
	}
	tssPubKey := tss.GetTssPubkey()

	crosschainFlags, err := c.GetCrosschainFlags()
	if err != nil {
		c.logger.Info().Msg("Unable to fetch cross-chain flags from zetacore")
		return fmt.Errorf("failed to get crosschain flags: %w", err)
	}

	blockHeaderEnabledChains, err := c.GetBlockHeaderEnabledChains()
	if err != nil {
		c.logger.Info().Msg("Unable to fetch block header enabled chains from zetacore")
		return err
	}

	coreContext.Update(
		keyGen,
		newChains,
		newEVMParams,
		newBTCParams,
		tssPubKey,
		crosschainFlags,
		blockHeaderEnabledChains,
		init,
		c.logger,
	)

	return nil
}

func (c *Client) Pause() {
	<-c.pause
}

func (c *Client) Unpause() {
	c.pause <- struct{}{}
}

func (c *Client) EnableMockSDKClient(client rpcclient.Client) {
	c.mockSDKClient = client
	c.enableMockSDKClient = true
}
