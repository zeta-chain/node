// Package zetacore provides the client to interact with zetacore node via GRPC.
package zetacore

import (
	"fmt"
	"strings"
	"sync"

	cometbftrpc "github.com/cometbft/cometbft/rpc/client"
	cometbfthttp "github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/types"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/pkg/authz"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/fanout"
	zetacorerpc "github.com/zeta-chain/node/pkg/rpc"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/config"
	keyinterfaces "github.com/zeta-chain/node/zetaclient/keys/interfaces"
	"github.com/zeta-chain/node/zetaclient/logs"
)

var _ interfaces.ZetacoreClient = &Client{}

// Client is the client to send tx to zetacore
type Client struct {
	zetacorerpc.Clients

	logger zerolog.Logger
	config config.ZetacoreClientConfig

	cosmosClientContext cosmosclient.Context
	cometBFTClient      cometbftrpc.Client

	blockHeight   int64
	accountNumber map[authz.KeyType]uint64
	seqNumber     map[authz.KeyType]uint64

	encodingCfg testutil.TestEncodingConfig
	keys        keyinterfaces.ObserverKeys
	chainID     string
	chain       chains.Chain

	// blocksFanout that receives new block events from Zetacore via websockets
	blocksFanout *fanout.FanOut[ctypes.EventDataNewBlock]

	mu sync.RWMutex
}

type constructOpts struct {
	customCometBFT bool
	cometBFTClient cometbftrpc.Client

	customAccountRetriever bool
	accountRetriever       cosmosclient.AccountRetriever
}

type Opt func(cfg *constructOpts)

// WithCometBFTClient sets custom CometBFT client
func WithCometBFTClient(client cometbftrpc.Client) Opt {
	return func(c *constructOpts) {
		c.customCometBFT = true
		c.cometBFTClient = client
	}
}

// WithCustomAccountRetriever sets custom CometBFT client
func WithCustomAccountRetriever(ac cosmosclient.AccountRetriever) Opt {
	return func(c *constructOpts) {
		c.customAccountRetriever = true
		c.accountRetriever = ac
	}
}

// NewClient create a new instance of Client
func NewClient(
	keys keyinterfaces.ObserverKeys,
	cfg config.Config,
	logger zerolog.Logger,
	opts ...Opt,
) (*Client, error) {
	var (
		chainID          = cfg.ChainID
		zetacoreCfg      = cfg.GetZetacoreClientConfig()
		constructOptions constructOpts
	)

	for _, opt := range opts {
		opt(&constructOptions)
	}

	zetaChain, err := chains.ZetaChainFromCosmosChainID(chainID)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid chain id %q", chainID)
	}

	zetacoreClients, err := zetacorerpc.NewGRPCClients(zetacoreCfg.GRPCURL, zetacoreCfg.GRPCDialOpt)
	if err != nil {
		return nil, errors.Wrap(err, "grpc dial fail")
	}

	accountsMap := make(map[authz.KeyType]uint64)
	seqMap := make(map[authz.KeyType]uint64)
	for _, keyType := range authz.GetAllKeyTypes() {
		accountsMap[keyType] = 0
		seqMap[keyType] = 0
	}

	encodingCfg := app.MakeEncodingConfig(uint64(zetaChain.ChainId)) //#nosec G115 won't exceed uint64
	cosmosContext, err := buildCosmosClientContext(chainID, keys, zetacoreCfg, encodingCfg, constructOptions)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build cosmos client context")
	}

	// create a cometbft client if one was not provided in the constructOptions
	cometBFTClient := constructOptions.cometBFTClient
	if !constructOptions.customCometBFT {
		client, err := createCometBFTClient(zetacoreCfg.WSRemote, true)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create cometbft client")
		}
		cometBFTClient = client
	}

	// set account number and sequence number for the zeta client grantee key
	address, err := keys.GetAddress()
	if err != nil {
		return nil, errors.Wrap(err, "fail to get address")
	}

	accN, seq, err := cosmosContext.AccountRetriever.GetAccountNumberSequence(cosmosContext, address)
	if err != nil {
		return nil, errors.Wrap(err, "fail to get account number and sequence number")
	}

	accountsMap[authz.ZetaClientGranteeKey] = accN
	seqMap[authz.ZetaClientGranteeKey] = seq

	return &Client{
		Clients: zetacoreClients,
		logger:  logger.With().Str(logs.FieldModule, "zetacoreClient").Logger(),
		config:  zetacoreCfg,

		cosmosClientContext: cosmosContext,
		cometBFTClient:      cometBFTClient,

		accountNumber: accountsMap,
		seqNumber:     seqMap,

		encodingCfg: encodingCfg,
		keys:        keys,
		chainID:     chainID,
		chain:       zetaChain,
	}, nil
}

// createCometBFTClient creates a cometbft client and optionally starts websocket
func createCometBFTClient(remote string, startWS bool) (cometbftrpc.Client, error) {
	client, err := cometbfthttp.New(remote, "/websocket")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create cometbft client from remote %s", remote)
	}

	// start websocket if needed
	if startWS {
		if err = client.WSEvents.Start(); err != nil {
			_ = client.Stop()
			return nil, errors.Wrap(err, "failed to start cometbft websocket")
		}
	}

	return client, nil
}

// buildCosmosClientContext constructs a valid context with all relevant values set
func buildCosmosClientContext(
	chainID string,
	keys keyinterfaces.ObserverKeys,
	config config.ZetacoreClientConfig,
	encodingConfig testutil.TestEncodingConfig,
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
		input  = strings.NewReader("")
		client cosmosclient.CometRPC
		remote = config.WSRemote
	)

	// if password is needed, set it as input
	password := keys.GetHotkeyPassword()
	if password != "" {
		input = strings.NewReader(fmt.Sprintf("%[1]s\n%[1]s\n", password))
	}

	// note that in rare cases, this might give FALSE positive
	// (google "golang nil interface comparison")
	client = opts.cometBFTClient
	if !opts.customCometBFT {
		client, err = createCometBFTClient(remote, false)
		if err != nil {
			return cosmosclient.Context{}, errors.Wrap(err, "failed to create cometbft client")
		}
	}

	var accountRetriever cosmosclient.AccountRetriever
	if opts.customAccountRetriever {
		accountRetriever = opts.accountRetriever
	} else {
		accountRetriever = authtypes.AccountRetriever{}
	}

	return cosmosclient.Context{
		Client:        client,
		NodeURI:       remote,
		FromAddress:   addr,
		ChainID:       chainID,
		Keyring:       keys.GetKeybase(),
		BroadcastMode: "sync",
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

func (c *Client) GetKeys() keyinterfaces.ObserverKeys {
	return c.keys
}

// GetAccountNumberAndSequenceNumber We do not use multiple KeyType for now , but this can be optionally used in the future to seprate TSS signer from Zetaclient GRantee
func (c *Client) GetAccountNumberAndSequenceNumber(_ authz.KeyType) (uint64, uint64, error) {
	address, err := c.keys.GetAddress()
	if err != nil {
		return 0, 0, err
	}
	return c.cosmosClientContext.AccountRetriever.GetAccountNumberSequence(c.cosmosClientContext, address)
}
