// Package zetacore provides functionalities for interacting with ZetaChain
package zetacore

import (
	"fmt"
	"sync"
	"time"

	"cosmossdk.io/simapp/params"
	"github.com/btcsuite/btcd/chaincfg"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/pkg/authz"
	"github.com/zeta-chain/zetacore/pkg/chains"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	keyinterfaces "github.com/zeta-chain/zetacore/zetaclient/keys/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
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
	keys          keyinterfaces.ObserverKeys
	broadcastLock *sync.RWMutex
	chainID       string
	chain         chains.Chain
	Telemetry     *metrics.TelemetryServer

	// enableMockSDKClient is a flag that determines whether the mock cosmos sdk client should be used, primarily for
	// unit testing
	enableMockSDKClient bool
	mockSDKClient       rpcclient.Client
}

// CreateClient is a helper function to create a new instance of Client
func CreateClient(
	cfg *config.Config,
	telemetry *metrics.TelemetryServer,
	hotkeyPassword string,
) (*Client, error) {
	hotKey := cfg.AuthzHotkey
	if cfg.HsmMode {
		hotKey = cfg.HsmHotKey
	}

	chainIP := cfg.ZetaCoreURL

	kb, _, err := keys.GetKeyringKeybase(cfg, hotkeyPassword)
	if err != nil {
		return nil, err
	}

	granterAddreess, err := sdk.AccAddressFromBech32(cfg.AuthzGranter)
	if err != nil {
		return nil, err
	}

	keys := keys.NewKeysWithKeybase(kb, granterAddreess, cfg.AuthzHotkey, hotkeyPassword)

	client, err := NewClient(keys, chainIP, hotKey, cfg.ChainID, cfg.HsmMode, telemetry)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewClient create a new instance of Client
func NewClient(
	keys keyinterfaces.ObserverKeys,
	chainIP string,
	signerName string,
	chainID string,
	hsmMode bool,
	telemetry *metrics.TelemetryServer,
) (*Client, error) {
	// main module logger
	logger := log.With().Str("module", "ZetacoreClient").Logger()
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
		keys:                keys,
		broadcastLock:       &sync.RWMutex{},
		chainID:             chainID,
		chain:               zetaChain,
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

func (c *Client) GetKeys() keyinterfaces.ObserverKeys {
	return c.keys
}

// GetAccountNumberAndSequenceNumber We do not use multiple KeyType for now , but this can be optionally used in the future to seprate TSS signer from Zetaclient GRantee
func (c *Client) GetAccountNumberAndSequenceNumber(_ authz.KeyType) (uint64, uint64, error) {
	ctx, err := c.GetContext()
	if err != nil {
		return 0, 0, err
	}
	address, err := c.keys.GetAddress()
	if err != nil {
		return 0, 0, err
	}
	return ctx.AccountRetriever.GetAccountNumberSequence(ctx, address)
}

// SetAccountNumber sets the account number and sequence number for the given keyType
func (c *Client) SetAccountNumber(keyType authz.KeyType) error {
	ctx, err := c.GetContext()
	if err != nil {
		return errors.Wrap(err, "fail to get context")
	}
	address, err := c.keys.GetAddress()
	if err != nil {
		return errors.Wrap(err, "fail to get address")
	}
	accN, seq, err := ctx.AccountRetriever.GetAccountNumberSequence(ctx, address)
	if err != nil {
		return errors.Wrap(err, "fail to get account number and sequence number")
	}
	c.accountNumber[keyType] = accN
	c.seqNumber[keyType] = seq

	return nil
}

// WaitForZetacoreToCreateBlocks waits for zetacore to create blocks
func (c *Client) WaitForZetacoreToCreateBlocks() error {
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
			return fmt.Errorf("zetacore is not ready, waited for %d seconds", DefaultRetryCount*DefaultRetryInterval)
		}
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return nil
}

// UpdateAppContext queries zetacore to update app context fields
func (c *Client) UpdateAppContext(appContext *context.AppContext, logger zerolog.Logger) error {
	// get latest supported chains
	supportedChains, err := c.GetSupportedChains()
	if err != nil {
		return errors.Wrap(err, "GetSupportedChains failed")
	}
	supportedChainsMap := make(map[int64]chains.Chain)
	for _, chain := range supportedChains {
		supportedChainsMap[chain.ChainId] = *chain
	}

	// get latest chain parameters
	chainParams, err := c.GetChainParams()
	if err != nil {
		return errors.Wrap(err, "GetChainParams failed")
	}

	var btcNetParams *chaincfg.Params
	chainsEnabled := make([]chains.Chain, 0)
	chainParamMap := make(map[int64]*observertypes.ChainParams)

	for _, chainParam := range chainParams {
		// skip unsupported chain
		if !chainParam.IsSupported {
			continue
		}

		// chain should exist in chain list
		chain, found := supportedChainsMap[chainParam.ChainId]
		if !found {
			continue
		}

		// skip ZetaChain
		if !chain.IsExternalChain() {
			continue
		}

		// zetaclient detects Bitcoin network (regnet, testnet, mainnet) from chain params in zetacore
		// The network params will be used by TSS to calculate the correct TSS address.
		if chains.IsBitcoinChain(chainParam.ChainId) {
			btcNetParams, err = chains.BitcoinNetParamsFromChainID(chainParam.ChainId)
			if err != nil {
				return errors.Wrapf(err, "BitcoinNetParamsFromChainID failed for chain %d", chainParam.ChainId)
			}
		}

		// zetaclient should observe this chain
		chainsEnabled = append(chainsEnabled, chain)
		chainParamMap[chainParam.ChainId] = chainParam
	}

	// get latest keygen
	keyGen, err := c.GetKeyGen()
	if err != nil {
		return errors.Wrap(err, "GetKeyGen failed")
	}

	// get latest TSS public key
	tss, err := c.GetCurrentTss()
	if err != nil {
		return errors.Wrap(err, "GetCurrentTss failed")
	}
	currentTssPubkey := tss.GetTssPubkey()

	// get latest crosschain flags
	crosschainFlags, err := c.GetCrosschainFlags()
	if err != nil {
		return errors.Wrap(err, "GetCrosschainFlags failed")
	}

	// get latest block header enabled chains
	blockHeaderEnabledChains, err := c.GetBlockHeaderEnabledChains()
	if err != nil {
		return errors.Wrap(err, "GetBlockHeaderEnabledChains failed")
	}

	// update app context fields
	appContext.Update(
		*keyGen,
		currentTssPubkey,
		chainsEnabled,
		chainParamMap,
		btcNetParams,
		crosschainFlags,
		blockHeaderEnabledChains,
		logger,
	)

	return nil
}

// EnableMockSDKClient enables the mock cosmos sdk client
// TODO(revamp): move this to a test package
func (c *Client) EnableMockSDKClient(client rpcclient.Client) {
	c.mockSDKClient = client
	c.enableMockSDKClient = true
}
