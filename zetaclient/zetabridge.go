package zetaclient

import (
	"fmt"
	//sdk "github.com/cosmos/cosmos-sdk/types"
	"net/url"
	"sync"

	"github.com/hashicorp/go-retryablehttp"
	"google.golang.org/grpc"

	//"fmt"
	"github.com/zeta-chain/zetacore/common/cosmos"
	//"github.com/armon/go-metrics"
	//"github.com/cosmos/cosmos-sdk/Client"
	"github.com/cosmos/cosmos-sdk/codec"

	//"github.com/cosmos/cosmos-sdk/std"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	//"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	//"golang.org/x/tools/go/cfg"
	//"io/ioutil"
	//"net/http"
	//"net/url"
	//"strconv"
	//"strings"

	stypes "github.com/zeta-chain/zetacore/x/zetacore/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

// ZetaCoreBridge will be used to send tx to ZetaCore.
type ZetaCoreBridge struct {
	logger        zerolog.Logger
	blockHeight   int64
	accountNumber uint64
	seqNumber     uint64
	grpcConn      *grpc.ClientConn
	httpClient    *retryablehttp.Client
	cfg           config.ClientConfiguration
	keys          *Keys
	broadcastLock *sync.RWMutex
	ChainNonces   map[string]uint64
}

// NewZetaCoreBridge create a new instance of ZetaCoreBridge
func NewZetaCoreBridge(k *Keys, chainIP string, signerName string) (*ZetaCoreBridge, error) {
	// main module logger
	logger := log.With().Str("module", "zetacore_client").Logger()

	cfg := config.ClientConfiguration{
		ChainHost:    fmt.Sprintf("%s:1317", chainIP),
		SignerName:   signerName,
		SignerPasswd: "password",
		ChainRPC:     fmt.Sprintf("%s:26657", chainIP),
	}

	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil

	grpcConn, err := grpc.Dial(
		fmt.Sprintf("%s:9090", chainIP),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Error().Err(err).Msg("grpc dial fail")
		return nil, err
	}

	return &ZetaCoreBridge{
		logger:        logger,
		grpcConn:      grpcConn,
		httpClient:    httpClient,
		cfg:           cfg,
		keys:          k,
		broadcastLock: &sync.RWMutex{},
		ChainNonces:   map[string]uint64{},
	}, nil
}

// MakeLegacyCodec creates codec
func MakeLegacyCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()
	banktypes.RegisterLegacyAminoCodec(cdc)
	authtypes.RegisterLegacyAminoCodec(cdc)
	cosmos.RegisterCodec(cdc)
	stypes.RegisterCodec(cdc)
	return cdc
}

// getZetaCoreURL with the given path
func (b *ZetaCoreBridge) getZetaCoreURL(path string) string {
	uri := url.URL{
		Scheme: "http",
		Host:   b.cfg.ChainHost,
		Path:   path,
	}
	return uri.String()
}

func (b *ZetaCoreBridge) GetAccountNumberAndSequenceNumber() (uint64, uint64, error) {
	ctx := b.GetContext()
	return ctx.AccountRetriever.GetAccountNumberSequence(ctx, b.keys.GetSignerInfo().GetAddress())
}

func (b *ZetaCoreBridge) GetKeys() *Keys {
	return b.keys
}
