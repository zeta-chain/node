package metaclientd

import (
	"encoding/json"
	//"context"
	"errors"
	"fmt"
	"github.com/Meta-Protocol/metacore/cmd/metaclientd/types"

	//sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashicorp/go-retryablehttp"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"

	"net/url"
	"sync"

	//"fmt"
	"github.com/Meta-Protocol/metacore/common/cosmos"
	//"github.com/armon/go-metrics"
	//"github.com/cosmos/cosmos-sdk/client"
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

	"github.com/Meta-Protocol/metacore/cmd/metaclientd/config"
	stypes "github.com/Meta-Protocol/metacore/x/metacore/types"
)

const (
	AuthAccountEndpoint = "/auth/accounts"
)

type TxStatus int64

const (
	Pending TxStatus = iota
	Processed
	Confirmed
)

// MetachainBridge will be used to send tx to MetaChain.
type MetachainBridge struct {
	logger                zerolog.Logger
	blockHeight           int64
	accountNumber         uint64
	seqNumber             uint64
	grpcConn              *grpc.ClientConn
	httpClient            *retryablehttp.Client
	cfg                   config.ClientConfiguration
	keys                  *Keys
	broadcastLock         *sync.RWMutex
	ProcessedTransactions map[string]TxStatus
	ChainNonces           map[string]uint64
}

// NewMetachainBridge create a new instance of MetachainBridge
func NewMetachainBridge(k *Keys, chainIP string) (*MetachainBridge, error) {
	// main module logger
	logger := log.With().Str("module", "metachain_client").Logger()

	cfg := config.ClientConfiguration{
		ChainHost:    fmt.Sprintf("%s:1317", chainIP),
		SignerName:   "val",
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

	return &MetachainBridge{
		logger:                logger,
		grpcConn:              grpcConn,
		httpClient:            httpClient,
		cfg:                   cfg,
		keys:                  k,
		broadcastLock:         &sync.RWMutex{},
		ProcessedTransactions: map[string]TxStatus{},
		ChainNonces:           map[string]uint64{},
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

// getMetachainURL with the given path
func (b *MetachainBridge) getMetachainURL(path string) string {
	uri := url.URL{
		Scheme: "http",
		Host:   b.cfg.ChainHost,
		Path:   path,
	}
	return uri.String()
}

//
// getAccountNumberAndSequenceNumber returns account and Sequence number required to post into thorchain
func (b *MetachainBridge) GetAccountNumberAndSequenceNumber() (uint64, uint64, error) {

	path := fmt.Sprintf("%s/%s", AuthAccountEndpoint, b.keys.GetSignerInfo().GetAddress())

	body, _, err := b.GetWithPath(path)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get auth accounts: %w", err)
	}

	var resp types.AccountResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, 0, fmt.Errorf("failed to unmarshal account resp: %w", err)
	}

	fmt.Printf("acct # %d, seq # %d\n", resp.Result.Value.AccountNumber, resp.Result.Value.Sequence)

	return resp.Result.Value.AccountNumber, resp.Result.Value.Sequence, nil
}

// get handle all the low level http GET calls using retryablehttp.ThorchainBridge
func (b *MetachainBridge) get(url string) ([]byte, int, error) {
	resp, err := b.httpClient.Get(url)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("failed to GET from thorchain: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	buf, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return buf, resp.StatusCode, errors.New("Status code: " + resp.Status + " returned")
	}
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}
	return buf, resp.StatusCode, nil
}

func (b *MetachainBridge) GetWithPath(path string) ([]byte, int, error) {
	return b.get(b.getMetachainURL(path))
}

func (b *MetachainBridge) GetKeys() *Keys {
	return b.keys
}