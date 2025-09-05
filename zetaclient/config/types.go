package config

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/showa-93/go-mask"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// KeyringBackend is the type of keyring backend to use for the hotkey
type KeyringBackend string

const (
	// KeyringBackendUndefined is undefined keyring backend
	KeyringBackendUndefined KeyringBackend = ""

	// KeyringBackendTest is the test Cosmos keyring backend
	KeyringBackendTest KeyringBackend = "test"

	// KeyringBackendFile is the file Cosmos keyring backend
	KeyringBackendFile KeyringBackend = "file"

	DefaultRelayerDir = "relayer-keys"

	// DefaultRelayerKeyPath is the default path that relayer keys are stored
	DefaultRelayerKeyPath = "~/.zetacored/" + DefaultRelayerDir
)

var (
	// CredsInsecureGRPC is a grpc.DialOption that uses insecure transport credentials
	// this is used when establishing gRPC connection to zetacore node via IP address
	CredsInsecureGRPC = grpc.WithTransportCredentials(insecure.NewCredentials())

	// CredsTLSGRPC is a grpc.DialOption that uses TLS transport credentials
	// this is used when establishing gRPC connection to zetacore node via hostname
	CredsTLSGRPC = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		// #nosec G402 - InsecureSkipVerify required for non-standard certificates (e.g., CloudFlare Origin)
		InsecureSkipVerify: true,
		NextProtos:         []string{"h2"},
	}))
)

// ZetacoreClientConfig is a subset of zetaclient config that is used by zetacore client
type ZetacoreClientConfig struct {
	GRPCURL     string          `json:"grpc_url"`
	WSRemote    string          `json:"ws_remote"`
	SignerName  string          `json:"signer_name"`
	GRPCDialOpt grpc.DialOption `json:"grpc_dial_opt"`
}

// EVMConfig is the config for EVM chain
type EVMConfig struct {
	Endpoint string `mask:"filled"`
}

// BTCConfig is the config for Bitcoin chain
type BTCConfig struct {
	// the following are rpcclient ConnConfig fields
	RPCUsername string `mask:"filled"`
	RPCPassword string `mask:"filled"`
	RPCHost     string `mask:"filled"`
	RPCParams   string // "regtest", "mainnet", "testnet3" , "signet", "testnet4"
}

// SolanaConfig is the config for Solana chain
type SolanaConfig struct {
	Endpoint string `mask:"filled"`
}

// SuiConfig is the config for Sui chain
type SuiConfig struct {
	Endpoint string `mask:"filled"`
}

// TONConfig is the config for TON chain
type TONConfig struct {
	// Endpoint url (toncenter V2 api) e.g. https://toncenter.com/api/v2/
	Endpoint string `mask:"filled"`
}

// ComplianceConfig is the config for compliance
type ComplianceConfig struct {
	LogPath string `json:"LogPath"`
	// Deprecated: use the separate restricted addresses config
	RestrictedAddresses []string `json:"RestrictedAddresses" mask:"zero"`
}

// Config is the config for ZetaClient
// TODO: use snake case for json fields
// https://github.com/zeta-chain/node/issues/1020
type Config struct {
	Peer          string `json:"Peer"`
	PublicIP      string `json:"PublicIP"`
	LogFormat     string `json:"LogFormat"`
	LogLevel      int8   `json:"LogLevel"`
	LogSampler    bool   `json:"LogSampler"`
	PreParamsPath string `json:"PreParamsPath"`
	ZetaCoreHome  string `json:"ZetaCoreHome"`
	ChainID       string `json:"ChainID"`
	// The old name tag 'ZetaCoreURL' is still used for backward compatibility
	ZetacoreIP              string         `json:"ZetaCoreURL"`
	ZetacoreURLGRPC         string         `json:"ZetaCoreURLGRPC"`
	ZetacoreURLWSS          string         `json:"ZetaCoreURLWSS"`
	AuthzGranter            string         `json:"AuthzGranter"`
	AuthzHotkey             string         `json:"AuthzHotkey"`
	P2PDiagnostic           bool           `json:"P2PDiagnostic"`
	ConfigUpdateTicker      uint64         `json:"ConfigUpdateTicker"`
	P2PDiagnosticTicker     uint64         `json:"P2PDiagnosticTicker"`
	TssPath                 string         `json:"TssPath"`
	TSSMaxPendingSignatures uint64         `json:"TSSMaxPendingSignatures"`
	TestTssKeysign          bool           `json:"TestTssKeysign"`
	KeyringBackend          KeyringBackend `json:"KeyringBackend"`
	RelayerKeyPath          string         `json:"RelayerKeyPath"`

	// chain configs
	EVMChainConfigs map[int64]EVMConfig `json:"EVMChainConfigs"`
	BTCChainConfigs map[int64]BTCConfig `json:"BTCChainConfigs"`
	SolanaConfig    SolanaConfig        `json:"SolanaConfig"`
	SuiConfig       SuiConfig           `json:"SuiConfig"`
	TONConfig       TONConfig           `json:"TONConfig"`

	// compliance config
	ComplianceConfig ComplianceConfig `json:"ComplianceConfig"`

	mu *sync.RWMutex
}

// GetZetacoreClientConfig returns the zetacore client config
func (c Config) GetZetacoreClientConfig() ZetacoreClientConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var (
		gRPCURL    string
		wsRemote   string
		dialOption grpc.DialOption
	)

	// zetaclient accepts both zetacore node IP address and hostnames.
	// To be compatible with IP address users, try IP address if set, otherwise use hostnames.
	// Note: leave the IP address field empty to use hostnames.
	if c.ZetacoreIP != "" {
		gRPCURL = cosmosGRPCFromIP(c.ZetacoreIP)
		wsRemote = cosmosWSSRemoteFromIP(c.ZetacoreIP)
		dialOption = CredsInsecureGRPC
	} else {
		gRPCURL = cosmosGRPCFromHost(c.ZetacoreURLGRPC)
		wsRemote = cosmosWSSRemoteFromHost(c.ZetacoreURLWSS)
		dialOption = CredsTLSGRPC
	}

	return ZetacoreClientConfig{
		GRPCURL:     gRPCURL,
		WSRemote:    wsRemote,
		SignerName:  c.AuthzHotkey,
		GRPCDialOpt: dialOption,
	}
}

// GetEVMConfig returns the EVM config for the given chain ID
func (c Config) GetEVMConfig(chainID int64) (EVMConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	evmCfg := c.EVMChainConfigs[chainID]
	return evmCfg, !evmCfg.Empty()
}

// GetAllEVMConfigs returns a map of all EVM configs
func (c Config) GetAllEVMConfigs() map[int64]EVMConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// deep copy evm configs
	copied := make(map[int64]EVMConfig, len(c.EVMChainConfigs))
	for chainID, evmConfig := range c.EVMChainConfigs {
		copied[chainID] = evmConfig
	}
	return copied
}

// GetBTCConfig returns the BTC config for the given chain ID
func (c Config) GetBTCConfig(chainID int64) (BTCConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	btcCfg := c.BTCChainConfigs[chainID]

	return btcCfg, !btcCfg.Empty()
}

// GetSolanaConfig returns the Solana config
func (c Config) GetSolanaConfig() (SolanaConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.SolanaConfig, c.SolanaConfig != (SolanaConfig{})
}

// GetSuiConfig returns the Sui config
func (c Config) GetSuiConfig() (SuiConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.SuiConfig, c.SuiConfig != (SuiConfig{})
}

// GetTONConfig returns the TONConfig and a bool indicating if it's present.
func (c Config) GetTONConfig() (TONConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.TONConfig, c.TONConfig != TONConfig{}
}

// StringMasked returns the string representation of the config with sensitive fields masked.
// Currently only the endpoints and bitcoin credentials are masked.
func (c Config) StringMasked() string {
	// create a masker
	masker := mask.NewMasker()
	masker.RegisterMaskStringFunc(mask.MaskTypeFilled, masker.MaskFilledString)
	masker.RegisterMaskAnyFunc(mask.MaskTypeFilled, masker.MaskZero)

	// mask the config
	masked, err := masker.Mask(c)
	if err != nil {
		return ""
	}

	s, err := json.MarshalIndent(masked, "", "\t")
	if err != nil {
		return ""
	}
	return string(s)
}

// GetRestrictedAddressBook returns a map of restricted addresses
// Note: the restricted address book contains both ETH and BTC addresses
func (c Config) GetRestrictedAddressBook() map[string]bool {
	restrictedAddresses := make(map[string]bool)
	for _, address := range c.ComplianceConfig.RestrictedAddresses {
		if address != "" {
			restrictedAddresses[strings.ToLower(address)] = true
		}
	}
	return restrictedAddresses
}

// GetKeyringBackend returns the keyring backend
func (c Config) GetKeyringBackend() KeyringBackend {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.KeyringBackend
}

// GetRelayerKeyPath returns the relayer key path
func (c Config) GetRelayerKeyPath() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// use default path if not configured
	if c.RelayerKeyPath == "" {
		return DefaultRelayerKeyPath
	}
	return c.RelayerKeyPath
}

func (c EVMConfig) Empty() bool {
	return c.Endpoint == ""
}

func (c BTCConfig) Empty() bool {
	return c.RPCHost == ""
}

// cosmosGRPCFromIP returns the gRPC URL for the given IP address
// Note: this function does not strictly enforce the IP address format.
// In E2E tests, the IP addresses passed to zetaclientd are ['zetacore0' ~ 'zetacore3']
// and these are not IP addresses, but they should still work without issues.
// Any wrong format of 'ipAddress' will trigger gRPC connection error, no worries.
func cosmosGRPCFromIP(ipAddress string) string {
	return fmt.Sprintf("%s:9090", ipAddress)
}

// cosmosWSSRemoteFromIP returns the websocket remote URI for the given IP address
func cosmosWSSRemoteFromIP(ipAddress string) string {
	remote := cometBFTRPC(ipAddress)

	// given an IP address, both remote URI formats will work:
	// 1. http://zetacore_ip_address:26657
	// 2. tcp://zetacore_ip_address:26657
	// append http:// prefix if not present
	if !strings.HasPrefix(remote, "http://") &&
		!strings.HasPrefix(remote, "tcp://") {
		remote = fmt.Sprintf("http://%s", remote)
	}

	return remote
}

// cometBFTRPC returns the CometBFT RPC endpoint for the given IP address
func cometBFTRPC(ipAddress string) string {
	return fmt.Sprintf("%s:26657", ipAddress)
}

// cosmosGRPCFromHost returns the gRPC URL for the given host
// Note: there is no assumption on node provider's gRPC URL format.
// The given host is expected to be a valid gRPC URL.
func cosmosGRPCFromHost(host string) string {
	return host
}

// cosmosWSSRemoteFromHost returns the websocket remote URI for the given host.
func cosmosWSSRemoteFromHost(host string) string {
	// A typical WSS URLs may look like below, and we need to convert them to the remote URI.
	// wss://rpc.provider.com/zetachain/websocket
	// wss://zetachain-mainnet.provider.com/websocket

	// remove "wss://" prefix and replace with "https://"
	remote := "https://" + strings.TrimPrefix(host, "wss://")

	// remove "/websocket" endpoint suffix if present
	// the suffix will be passed to http.New() as a separate argument
	remote = strings.TrimSuffix(remote, "/websocket")

	return remote
}
