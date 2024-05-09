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
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"google.golang.org/grpc"
)

var _ interfaces.ZetaCoreClient = &Client{}

// Client is the client to send tx to ZetaCore
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

func (b *Client) UpdateChainID(chainID string) error {
	if b.chainID != chainID {
		b.chainID = chainID

		zetaChain, err := chains.ZetaChainFromChainID(chainID)
		if err != nil {
			return fmt.Errorf("invalid chain id %s, %w", chainID, err)
		}
		b.chain = zetaChain
	}

	return nil
}

// Chain returns the Chain chain object
func (b *Client) Chain() chains.Chain {
	return b.chain
}

func (b *Client) GetLogger() *zerolog.Logger {
	return &b.logger
}

func (b *Client) GetKeys() *keys.Keys {
	return b.keys
}

func (b *Client) Stop() {
	b.logger.Info().Msgf("zetacore client is stopping")
	close(b.stop) // this notifies all configupdater to stop
}

// GetAccountNumberAndSequenceNumber We do not use multiple KeyType for now , but this can be optionally used in the future to seprate TSS signer from Zetaclient GRantee
func (b *Client) GetAccountNumberAndSequenceNumber(_ authz.KeyType) (uint64, uint64, error) {
	ctx, err := b.GetContext()
	if err != nil {
		return 0, 0, err
	}
	address := b.keys.GetAddress()
	return ctx.AccountRetriever.GetAccountNumberSequence(ctx, address)
}

func (b *Client) SetAccountNumber(keyType authz.KeyType) {
	ctx, err := b.GetContext()
	if err != nil {
		b.logger.Error().Err(err).Msg("fail to get context")
		return
	}
	address := b.keys.GetAddress()
	accN, seq, err := ctx.AccountRetriever.GetAccountNumberSequence(ctx, address)
	if err != nil {
		b.logger.Error().Err(err).Msg("fail to get account number and sequence number")
		return
	}
	b.accountNumber[keyType] = accN
	b.seqNumber[keyType] = seq
}

func (b *Client) WaitForCoreToCreateBlocks() {
	retryCount := 0
	for {
		block, err := b.GetLatestZetaBlock()
		if err == nil && block.Header.Height > 1 {
			b.logger.Info().Msgf("Zeta-core height: %d", block.Header.Height)
			break
		}
		retryCount++
		b.logger.Debug().Msgf("Failed to get latest Block , Retry : %d/%d", retryCount, DefaultRetryCount)
		if retryCount > ExtendedRetryCount {
			panic(fmt.Sprintf("ZetaCore is not ready , Waited for %d seconds", DefaultRetryCount*DefaultRetryInterval))
		}
		time.Sleep(DefaultRetryInterval * time.Second)
	}
}

// UpdateZetaCoreContext updates core context
// zetacore stores core context for all clients
func (b *Client) UpdateZetaCoreContext(coreContext *context.ZetaCoreContext, init bool, sampledLogger zerolog.Logger) error {
	bn, err := b.GetBlockHeight()
	if err != nil {
		return fmt.Errorf("failed to get zetablock height: %w", err)
	}
	plan, err := b.GetUpgradePlan()
	if err != nil {
		// if there is no active upgrade plan, plan will be nil, err will be nil as well.
		return fmt.Errorf("failed to get upgrade plan: %w", err)
	}
	if plan != nil && bn == plan.Height-1 { // stop zetaclients; notify operator to upgrade and restart
		b.logger.Warn().Msgf("Active upgrade plan detected and upgrade height reached: %s at height %d; ZetaClient is stopped;"+
			"please kill this process, replace zetaclientd binary with upgraded version, and restart zetaclientd", plan.Name, plan.Height)
		b.pause <- struct{}{} // notify CoreObserver to stop ChainClients, Signers, and CoreObserver itself
	}

	chainParams, err := b.GetChainParams()
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

	supportedChains, err := b.GetSupportedChains()
	if err != nil {
		return fmt.Errorf("failed to get supported chains: %w", err)
	}
	newChains := make([]chains.Chain, len(supportedChains))
	for i, chain := range supportedChains {
		newChains[i] = *chain
	}
	keyGen, err := b.GetKeyGen()
	if err != nil {
		b.logger.Info().Msg("Unable to fetch keygen from zetacore")
		return fmt.Errorf("failed to get keygen: %w", err)
	}

	tss, err := b.GetCurrentTss()
	if err != nil {
		b.logger.Info().Err(err).Msg("Unable to fetch TSS from zetacore")
		return fmt.Errorf("failed to get current tss: %w", err)
	}
	tssPubKey := tss.GetTssPubkey()

	crosschainFlags, err := b.GetCrosschainFlags()
	if err != nil {
		b.logger.Info().Msg("Unable to fetch cross-chain flags from zetacore")
		return fmt.Errorf("failed to get crosschain flags: %w", err)
	}

	verificationFlags, err := b.GetVerificationFlags()
	if err != nil {
		b.logger.Info().Msg("Unable to fetch verification flags from zetacore")

		// The block header functionality is currently disabled on the ZetaCore side
		// The verification flags might not exist and we should not return an error here to prevent the ZetaClient from starting
		// TODO: Uncomment this line when the block header functionality is enabled and we need to get the verification flags
		// https://github.com/zeta-chain/node/issues/1717
		// return fmt.Errorf("failed to get verification flags: %w", err)

		verificationFlags = lightclienttypes.VerificationFlags{}
	}

	coreContext.Update(
		keyGen,
		newChains,
		newEVMParams,
		newBTCParams,
		tssPubKey,
		crosschainFlags,
		verificationFlags,
		init,
		b.logger,
	)

	return nil
}

func (b *Client) Pause() {
	<-b.pause
}

func (b *Client) Unpause() {
	b.pause <- struct{}{}
}

func (b *Client) EnableMockSDKClient(client rpcclient.Client) {
	b.mockSDKClient = client
	b.enableMockSDKClient = true
}
