package zetabridge

import (
	"fmt"
	"sync"
	"time"

	"cosmossdk.io/simapp/params"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/pkg/authz"
	"github.com/zeta-chain/zetacore/pkg/chains"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"google.golang.org/grpc"
)

var _ interfaces.ZetaCoreBridger = &ZetaCoreBridge{}

// ZetaCoreBridge will be used to send tx to ZetaCore.
type ZetaCoreBridge struct {
	logger        zerolog.Logger
	blockHeight   int64
	accountNumber map[authz.KeyType]uint64
	seqNumber     map[authz.KeyType]uint64
	grpcConn      *grpc.ClientConn
	cfg           config.ClientConfiguration
	encodingCfg   params.EncodingConfig
	keys          *keys.Keys
	broadcastLock *sync.RWMutex
	zetaChainID   string
	zetaChain     chains.Chain
	stop          chan struct{}
	pause         chan struct{}
	Telemetry     *metrics.TelemetryServer

	// enableMockSDKClient is a flag that determines whether the mock cosmos sdk client should be used, primarily for
	// unit testing
	enableMockSDKClient bool
	mockSDKClient       rpcclient.Client
}

// NewZetaCoreBridge create a new instance of ZetaCoreBridge
func NewZetaCoreBridge(
	k *keys.Keys,
	chainIP string,
	signerName string,
	chainID string,
	hsmMode bool,
	telemetry *metrics.TelemetryServer,
) (*ZetaCoreBridge, error) {

	// main module logger
	logger := log.With().Str("module", "CoreBridge").Logger()
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

	return &ZetaCoreBridge{
		logger:              logger,
		grpcConn:            grpcConn,
		accountNumber:       accountsMap,
		seqNumber:           seqMap,
		cfg:                 cfg,
		encodingCfg:         app.MakeEncodingConfig(),
		keys:                k,
		broadcastLock:       &sync.RWMutex{},
		stop:                make(chan struct{}),
		zetaChainID:         chainID,
		zetaChain:           zetaChain,
		pause:               make(chan struct{}),
		Telemetry:           telemetry,
		enableMockSDKClient: false,
		mockSDKClient:       nil,
	}, nil
}

// MakeLegacyCodec creates codec
func MakeLegacyCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()
	banktypes.RegisterLegacyAminoCodec(cdc)
	authtypes.RegisterLegacyAminoCodec(cdc)
	sdk.RegisterLegacyAminoCodec(cdc)
	crosschaintypes.RegisterCodec(cdc)
	return cdc
}

func (b *ZetaCoreBridge) GetLogger() *zerolog.Logger {
	return &b.logger
}

func (b *ZetaCoreBridge) UpdateChainID(chainID string) error {
	if b.zetaChainID != chainID {
		b.zetaChainID = chainID

		zetaChain, err := chains.ZetaChainFromChainID(chainID)
		if err != nil {
			return fmt.Errorf("invalid chain id %s, %w", chainID, err)
		}
		b.zetaChain = zetaChain
	}

	return nil
}

// ZetaChain returns the ZetaChain chain object
func (b *ZetaCoreBridge) ZetaChain() chains.Chain {
	return b.zetaChain
}

func (b *ZetaCoreBridge) Stop() {
	b.logger.Info().Msgf("ZetaBridge is stopping")
	close(b.stop) // this notifies all configupdater to stop
}

// GetAccountNumberAndSequenceNumber We do not use multiple KeyType for now , but this can be optionally used in the future to seprate TSS signer from Zetaclient GRantee
func (b *ZetaCoreBridge) GetAccountNumberAndSequenceNumber(_ authz.KeyType) (uint64, uint64, error) {
	ctx, err := b.GetContext()
	if err != nil {
		return 0, 0, err
	}
	address := b.keys.GetAddress()
	return ctx.AccountRetriever.GetAccountNumberSequence(ctx, address)
}

func (b *ZetaCoreBridge) SetAccountNumber(keyType authz.KeyType) {
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

func (b *ZetaCoreBridge) WaitForCoreToCreateBlocks() {
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

func (b *ZetaCoreBridge) GetKeys() *keys.Keys {
	return b.keys
}

// UpdateZetaCoreContext updates core context
// zetacore stores core context for all clients
func (b *ZetaCoreBridge) UpdateZetaCoreContext(coreContext *corecontext.ZetaCoreContext, init bool) error {
	bn, err := b.GetZetaBlockHeight()
	if err != nil {
		return err
	}
	plan, err := b.GetUpgradePlan()
	// if there is no active upgrade plan, plan will be nil, err will be nil as well.
	if err != nil {
		return err
	}
	if plan != nil && bn == plan.Height-1 { // stop zetaclients; notify operator to upgrade and restart
		b.logger.Warn().Msgf("Active upgrade plan detected and upgrade height reached: %s at height %d; ZetaClient is stopped;"+
			"please kill this process, replace zetaclientd binary with upgraded version, and restart zetaclientd", plan.Name, plan.Height)
		b.pause <- struct{}{} // notify CoreObserver to stop ChainClients, Signers, and CoreObserver itself
	}

	chainParams, err := b.GetChainParams()
	if err != nil {
		return err
	}

	newEVMParams := make(map[int64]*observertypes.ChainParams)
	var newBTCParams *observertypes.ChainParams

	// check and update chain params for each chain
	sampledLogger := b.logger.Sample(&zerolog.BasicSampler{N: 10})
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
		return err
	}
	newChains := make([]chains.Chain, len(supportedChains))
	for i, chain := range supportedChains {
		newChains[i] = *chain
	}
	keyGen, err := b.GetKeyGen()
	if err != nil {
		b.logger.Info().Msg("Unable to fetch keygen from zetabridge")
		return err
	}

	tss, err := b.GetCurrentTss()
	if err != nil {
		b.logger.Info().Err(err).Msg("Unable to fetch TSS from zetabridge")
		return err
	}
	tssPubKey := tss.GetTssPubkey()

	crosschainFlags, err := b.GetCrosschainFlags()
	if err != nil {
		b.logger.Info().Msg("Unable to fetch cross-chain flags from zetabridge")
		return err
	}

	coreContext.Update(keyGen, newChains, newEVMParams, newBTCParams, tssPubKey, crosschainFlags, init, b.logger)
	return nil
}

func (b *ZetaCoreBridge) Pause() {
	<-b.pause
}

func (b *ZetaCoreBridge) Unpause() {
	b.pause <- struct{}{}
}

func (b *ZetaCoreBridge) EnableMockSDKClient(client rpcclient.Client) {
	b.mockSDKClient = client
	b.enableMockSDKClient = true
}
