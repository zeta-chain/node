package config

import (
	"errors"
	"fmt"
	"net/netip"
	"path"
	"time"

	"github.com/spf13/viper"

	"github.com/cometbft/cometbft/libs/strings"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/server/config"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// ServerStartTime defines the time duration that the server need to stay running after startup
	// for the startup be considered successful
	ServerStartTime = 5 * time.Second

	// DefaultAPIEnable is the default value for the parameter that defines if the cosmos REST API server is enabled
	DefaultAPIEnable = false

	// DefaultGRPCEnable is the default value for the parameter that defines if the gRPC server is enabled
	DefaultGRPCEnable = false

	// DefaultGRPCWebEnable is the default value for the parameter that defines if the gRPC web server is enabled
	DefaultGRPCWebEnable = false

	// DefaultJSONRPCEnable is the default value for the parameter that defines if the JSON-RPC server is enabled
	DefaultJSONRPCEnable = false

	// DefaultRosettaEnable is the default value for the parameter that defines if the Rosetta API server is enabled
	DefaultRosettaEnable = false

	// DefaultTelemetryEnable is the default value for the parameter that defines if the telemetry is enabled
	DefaultTelemetryEnable = false

	// DefaultGRPCAddress is the default address the gRPC server binds to.
	DefaultGRPCAddress = "0.0.0.0:9900"

	// DefaultJSONRPCAddress is the default address the JSON-RPC server binds to.
	DefaultJSONRPCAddress = "127.0.0.1:8545"

	// DefaultJSONRPCWsAddress is the default address the JSON-RPC WebSocket server binds to.
	DefaultJSONRPCWsAddress = "127.0.0.1:8546"

	// DefaultJsonRPCMetricsAddress is the default address the JSON-RPC Metrics server binds to.
	DefaultJSONRPCMetricsAddress = "127.0.0.1:6065"

	// DefaultEVMTracer is the default vm.Tracer type
	DefaultEVMTracer = ""

	// DefaultEnablePreimageRecording is the default value for EnablePreimageRecording
	DefaultEnablePreimageRecording = false

	// DefaultMaxTxGasWanted is the default gas wanted for each eth tx returned in ante handler in check tx mode
	DefaultMaxTxGasWanted = 0

	// DefaultEVMChainID is the default EVM Chain ID if one is not provided
	DefaultEVMChainID = 262144

	// DefaultEVMMinTip is the default minimum priority fee for the mempool
	DefaultEVMMinTip = 0

	// DefaultGethMetricsAddress is the default port for the geth metrics server.
	DefaultGethMetricsAddress = "127.0.0.1:8100"

	// DefaultGasCap is the default cap on gas that can be used in eth_call/estimateGas
	DefaultGasCap uint64 = 25_000_000

	// DefaultJSONRPCAllowInsecureUnlock is true
	DefaultJSONRPCAllowInsecureUnlock bool = true

	// DefaultFilterCap is the default cap for total number of filters that can be created
	DefaultFilterCap int32 = 200

	// DefaultFeeHistoryCap is the default cap for total number of blocks that can be fetched
	DefaultFeeHistoryCap int32 = 100

	// DefaultLogsCap is the default cap of results returned from single 'eth_getLogs' query
	DefaultLogsCap int32 = 10000

	// DefaultBlockRangeCap is the default cap of block range allowed for 'eth_getLogs' query
	DefaultBlockRangeCap int32 = 10000

	// DefaultEVMTimeout is the default timeout for eth_call
	DefaultEVMTimeout = 5 * time.Second

	// DefaultTxFeeCap is the default tx-fee cap for sending a transaction
	DefaultTxFeeCap float64 = 1.0

	// DefaultHTTPTimeout is the default read/write timeout of the http json-rpc server
	DefaultHTTPTimeout = 30 * time.Second

	// DefaultHTTPIdleTimeout is the default idle timeout of the http json-rpc server
	DefaultHTTPIdleTimeout = 120 * time.Second

	// DefaultAllowUnprotectedTxs value is false
	DefaultAllowUnprotectedTxs = false

	// DefaultBatchRequestLimit is the default maximum batch request limit.
	// https://github.com/ethereum/go-ethereum/blob/v1.15.11/node/defaults.go#L67
	DefaultBatchRequestLimit = 1000

	// DefaultBatchResponseMaxSize is the default maximum batch response size.
	// https://github.com/ethereum/go-ethereum/blob/v1.15.11/node/defaults.go#L68
	DefaultBatchResponseMaxSize = 25 * 1000 * 1000

	// DefaultMaxOpenConnections represents the amount of open connections (unlimited = 0)
	DefaultMaxOpenConnections = 0

	// DefaultGasAdjustment value to use as default in gas-adjustment flag
	DefaultGasAdjustment = 1.2

	// DefaultWSOrigins is the default origin for WebSocket connections
	DefaultWSOrigins = "127.0.0.1"

	// DefaultEnableProfiling toggles whether profiling is enabled in the `debug` namespace
	DefaultEnableProfiling = false
)

var evmTracers = []string{"json", "markdown", "struct", "access_list"}

// Config defines the server's top level configuration. It includes the default app config
// from the SDK as well as the EVM configuration to enable the JSON-RPC APIs.
type Config struct {
	config.Config `mapstructure:",squash"`

	EVM     EVMConfig     `mapstructure:"evm"`
	JSONRPC JSONRPCConfig `mapstructure:"json-rpc"`
	TLS     TLSConfig     `mapstructure:"tls"`
}

// EVMConfig defines the application configuration values for the EVM.
type EVMConfig struct {
	// Tracer defines vm.Tracer type that the EVM will use if the node is run in
	// trace mode. Default: 'json'.
	Tracer string `mapstructure:"tracer"`
	// MaxTxGasWanted defines the gas wanted for each eth tx returned in ante handler in check tx mode.
	MaxTxGasWanted uint64 `mapstructure:"max-tx-gas-wanted"`
	// Enables tracking of SHA3 preimages in the VM
	EnablePreimageRecording bool `mapstructure:"cache-preimage"`
	// EVMChainID defines the EIP-155 replay-protection chain ID.
	EVMChainID uint64 `mapstructure:"evm-chain-id"`
	// MinTip defines the minimum priority fee for the mempool
	MinTip uint64 `mapstructure:"min-tip"`
	// GethMetricsAddress is the address the geth metrics server will bind to. Default 127.0.0.1:8100
	GethMetricsAddress string `mapstructure:"geth-metrics-address"`
	// Mempool defines the EVM mempool configuration
	Mempool MempoolConfig `mapstructure:"mempool"`
}

// MempoolConfig defines the configuration for the EVM mempool transaction pool.
type MempoolConfig struct {
	// PriceLimit is the minimum gas price to enforce for acceptance into the pool
	PriceLimit uint64 `mapstructure:"price-limit"`
	// PriceBump is the minimum price bump percentage to replace an already existing transaction (nonce)
	PriceBump uint64 `mapstructure:"price-bump"`
	// AccountSlots is the number of executable transaction slots guaranteed per account
	AccountSlots uint64 `mapstructure:"account-slots"`
	// GlobalSlots is the maximum number of executable transaction slots for all accounts
	GlobalSlots uint64 `mapstructure:"global-slots"`
	// AccountQueue is the maximum number of non-executable transaction slots permitted per account
	AccountQueue uint64 `mapstructure:"account-queue"`
	// GlobalQueue is the maximum number of non-executable transaction slots for all accounts
	GlobalQueue uint64 `mapstructure:"global-queue"`
	// Lifetime is the maximum amount of time non-executable transaction are queued
	Lifetime time.Duration `mapstructure:"lifetime"`
}

// DefaultMempoolConfig returns the default mempool configuration
func DefaultMempoolConfig() MempoolConfig {
	return MempoolConfig{
		PriceLimit:   1,             // Minimum gas price of 1 wei
		PriceBump:    10,            // 10% price bump to replace transaction
		AccountSlots: 16,            // 16 executable transaction slots per account
		GlobalSlots:  5120,          // 4096 + 1024 = 5120 global executable slots
		AccountQueue: 64,            // 64 non-executable transaction slots per account
		GlobalQueue:  1024,          // 1024 global non-executable slots
		Lifetime:     3 * time.Hour, // 3 hour lifetime for queued transactions
	}
}

// Validate returns an error if the mempool configuration is invalid
func (c MempoolConfig) Validate() error {
	if c.PriceLimit < 1 {
		return fmt.Errorf("price limit must be at least 1, got %d", c.PriceLimit)
	}
	if c.PriceBump < 1 {
		return fmt.Errorf("price bump must be at least 1, got %d", c.PriceBump)
	}
	if c.AccountSlots < 1 {
		return fmt.Errorf("account slots must be at least 1, got %d", c.AccountSlots)
	}
	if c.GlobalSlots < 1 {
		return fmt.Errorf("global slots must be at least 1, got %d", c.GlobalSlots)
	}
	if c.AccountQueue < 1 {
		return fmt.Errorf("account queue must be at least 1, got %d", c.AccountQueue)
	}
	if c.GlobalQueue < 1 {
		return fmt.Errorf("global queue must be at least 1, got %d", c.GlobalQueue)
	}
	if c.Lifetime < 1 {
		return fmt.Errorf("lifetime must be at least 1 nanosecond, got %s", c.Lifetime)
	}
	return nil
}

// JSONRPCConfig defines configuration for the EVM RPC server.
type JSONRPCConfig struct {
	// API defines a list of JSON-RPC namespaces that should be enabled
	API []string `mapstructure:"api"`
	// Address defines the HTTP server to listen on
	Address string `mapstructure:"address"`
	// WsAddress defines the WebSocket server to listen on
	WsAddress string `mapstructure:"ws-address"`
	// GasCap is the global gas cap for eth-call variants.
	GasCap uint64 `mapstructure:"gas-cap"`
	// AllowInsecureUnlock toggles if account unlocking is enabled when account-related RPCs are exposed by http.
	AllowInsecureUnlock bool `mapstructure:"allow-insecure-unlock"`
	// EVMTimeout is the global timeout for eth-call.
	EVMTimeout time.Duration `mapstructure:"evm-timeout"`
	// TxFeeCap is the global tx-fee cap for send transaction
	TxFeeCap float64 `mapstructure:"txfee-cap"`
	// FilterCap is the global cap for total number of filters that can be created.
	FilterCap int32 `mapstructure:"filter-cap"`
	// FeeHistoryCap is the global cap for total number of blocks that can be fetched
	FeeHistoryCap int32 `mapstructure:"feehistory-cap"`
	// Enable defines if the EVM RPC server should be enabled.
	Enable bool `mapstructure:"enable"`
	// LogsCap defines the max number of results can be returned from single `eth_getLogs` query.
	LogsCap int32 `mapstructure:"logs-cap"`
	// BlockRangeCap defines the max block range allowed for `eth_getLogs` query.
	BlockRangeCap int32 `mapstructure:"block-range-cap"`
	// HTTPTimeout is the read/write timeout of http json-rpc server.
	HTTPTimeout time.Duration `mapstructure:"http-timeout"`
	// HTTPIdleTimeout is the idle timeout of http json-rpc server.
	HTTPIdleTimeout time.Duration `mapstructure:"http-idle-timeout"`
	// AllowUnprotectedTxs restricts unprotected (non EIP155 signed) transactions to be submitted via
	// the node's RPC when global parameter is disabled.
	AllowUnprotectedTxs bool `mapstructure:"allow-unprotected-txs"`
	// BatchRequestLimit is the maximum number of requests in a batch.
	BatchRequestLimit int `mapstructure:"batch-request-limit"`
	// BatchResponseMaxSize is the maximum number of bytes returned from a batched rpc call.
	BatchResponseMaxSize int `mapstructure:"batch-response-max-size"`
	// MaxOpenConnections sets the maximum number of simultaneous connections
	// for the server listener.
	MaxOpenConnections int `mapstructure:"max-open-connections"`
	// EnableIndexer defines if enable the custom indexer service.
	EnableIndexer bool `mapstructure:"enable-indexer"`
	// MetricsAddress defines the metrics server to listen on
	MetricsAddress string `mapstructure:"metrics-address"`
	// WSOrigins defines the allowed origins for WebSocket connections
	WSOrigins []string `mapstructure:"ws-origins"`
	// EnableProfiling enables the profiling in the `debug` namespace. SHOULD NOT be used on public tracing nodes
	EnableProfiling bool `mapstructure:"enable-profiling"`
}

// TLSConfig defines the certificate and matching private key for the server.
type TLSConfig struct {
	// CertificatePath the file path for the certificate .pem file
	CertificatePath string `mapstructure:"certificate-path"`
	// KeyPath the file path for the key .pem file
	KeyPath string `mapstructure:"key-path"`
}

// DefaultEVMConfig returns the default EVM configuration
func DefaultEVMConfig() *EVMConfig {
	return &EVMConfig{
		Tracer:                  DefaultEVMTracer,
		MaxTxGasWanted:          DefaultMaxTxGasWanted,
		EVMChainID:              DefaultEVMChainID,
		EnablePreimageRecording: DefaultEnablePreimageRecording,
		MinTip:                  DefaultEVMMinTip,
		GethMetricsAddress:      DefaultGethMetricsAddress,
		Mempool:                 DefaultMempoolConfig(),
	}
}

// Validate returns an error if the tracer type is invalid.
func (c EVMConfig) Validate() error {
	if c.Tracer != "" && !strings.StringInSlice(c.Tracer, evmTracers) {
		return fmt.Errorf("invalid tracer type %s, available types: %v", c.Tracer, evmTracers)
	}

	if _, err := netip.ParseAddrPort(c.GethMetricsAddress); err != nil {
		return fmt.Errorf("invalid geth metrics address %q: %w", c.GethMetricsAddress, err)
	}

	if err := c.Mempool.Validate(); err != nil {
		return fmt.Errorf("invalid mempool config: %w", err)
	}

	return nil
}

// GetDefaultAPINamespaces returns the default list of JSON-RPC namespaces that should be enabled
func GetDefaultAPINamespaces() []string {
	return []string{"eth", "net", "web3"}
}

// GetAPINamespaces returns the all the available JSON-RPC API namespaces.
func GetAPINamespaces() []string {
	return []string{"web3", "eth", "personal", "net", "txpool", "debug", "miner"}
}

// GetDefaultWSOrigins returns the default WebSocket origins.
func GetDefaultWSOrigins() []string {
	return []string{DefaultWSOrigins, "localhost"}
}

// DefaultJSONRPCConfig returns an EVM config with the JSON-RPC API enabled by default
func DefaultJSONRPCConfig() *JSONRPCConfig {
	return &JSONRPCConfig{
		Enable:               false,
		API:                  GetDefaultAPINamespaces(),
		Address:              DefaultJSONRPCAddress,
		WsAddress:            DefaultJSONRPCWsAddress,
		GasCap:               DefaultGasCap,
		AllowInsecureUnlock:  DefaultJSONRPCAllowInsecureUnlock,
		EVMTimeout:           DefaultEVMTimeout,
		TxFeeCap:             DefaultTxFeeCap,
		FilterCap:            DefaultFilterCap,
		FeeHistoryCap:        DefaultFeeHistoryCap,
		BlockRangeCap:        DefaultBlockRangeCap,
		LogsCap:              DefaultLogsCap,
		HTTPTimeout:          DefaultHTTPTimeout,
		HTTPIdleTimeout:      DefaultHTTPIdleTimeout,
		AllowUnprotectedTxs:  DefaultAllowUnprotectedTxs,
		BatchRequestLimit:    DefaultBatchRequestLimit,
		BatchResponseMaxSize: DefaultBatchResponseMaxSize,
		MaxOpenConnections:   DefaultMaxOpenConnections,
		EnableIndexer:        false,
		MetricsAddress:       DefaultJSONRPCMetricsAddress,
		WSOrigins:            GetDefaultWSOrigins(),
		EnableProfiling:      DefaultEnableProfiling,
	}
}

// Validate returns an error if the JSON-RPC configuration fields are invalid.
func (c JSONRPCConfig) Validate() error {
	if c.Enable && len(c.API) == 0 {
		return errors.New("cannot enable JSON-RPC without defining any API namespace")
	}

	if c.FilterCap < 0 {
		return errors.New("JSON-RPC filter-cap cannot be negative")
	}

	if c.FeeHistoryCap <= 0 {
		return errors.New("JSON-RPC feehistory-cap cannot be negative or 0")
	}

	if c.TxFeeCap < 0 {
		return errors.New("JSON-RPC tx fee cap cannot be negative")
	}

	if c.EVMTimeout < 0 {
		return errors.New("JSON-RPC EVM timeout duration cannot be negative")
	}

	if c.LogsCap < 0 {
		return errors.New("JSON-RPC logs cap cannot be negative")
	}

	if c.BlockRangeCap < 0 {
		return errors.New("JSON-RPC block range cap cannot be negative")
	}

	if c.HTTPTimeout < 0 {
		return errors.New("JSON-RPC HTTP timeout duration cannot be negative")
	}

	if c.HTTPIdleTimeout < 0 {
		return errors.New("JSON-RPC HTTP idle timeout duration cannot be negative")
	}

	if c.BatchRequestLimit < 0 {
		return errors.New("JSON-RPC batch request limit cannot be negative")
	}

	if c.BatchResponseMaxSize < 0 {
		return errors.New("JSON-RPC batch response max size cannot be negative")
	}

	// check for duplicates
	seenAPIs := make(map[string]bool)
	for _, api := range c.API {
		if seenAPIs[api] {
			return fmt.Errorf("repeated API namespace '%s'", api)
		}

		seenAPIs[api] = true
	}

	return nil
}

// DefaultTLSConfig returns the default TLS configuration
func DefaultTLSConfig() *TLSConfig {
	return &TLSConfig{
		CertificatePath: "",
		KeyPath:         "",
	}
}

// Validate returns an error if the TLS certificate and key file extensions are invalid.
func (c TLSConfig) Validate() error {
	certExt := path.Ext(c.CertificatePath)

	if c.CertificatePath != "" && certExt != ".pem" {
		return fmt.Errorf("invalid extension %s for certificate path %s, expected '.pem'", certExt, c.CertificatePath)
	}

	keyExt := path.Ext(c.KeyPath)

	if c.KeyPath != "" && keyExt != ".pem" {
		return fmt.Errorf("invalid extension %s for key path %s, expected '.pem'", keyExt, c.KeyPath)
	}

	return nil
}

// DefaultConfig returns server's default configuration.
func DefaultConfig() *Config {
	defaultSDKConfig := config.DefaultConfig()
	defaultSDKConfig.API.Enable = DefaultAPIEnable
	defaultSDKConfig.GRPC.Enable = DefaultGRPCEnable
	defaultSDKConfig.GRPCWeb.Enable = DefaultGRPCWebEnable
	defaultSDKConfig.Telemetry.Enabled = DefaultTelemetryEnable

	return &Config{
		Config:  *defaultSDKConfig,
		EVM:     *DefaultEVMConfig(),
		JSONRPC: *DefaultJSONRPCConfig(),
		TLS:     *DefaultTLSConfig(),
	}
}

// GetConfig returns a fully parsed Config object.
func GetConfig(v *viper.Viper) (Config, error) {
	conf := DefaultConfig()
	if err := v.Unmarshal(conf); err != nil {
		return Config{}, fmt.Errorf("error extracting app config: %w", err)
	}
	return *conf, nil
}

// ValidateBasic returns an error any of the application configuration fields are invalid
func (c Config) ValidateBasic() error {
	if err := c.EVM.Validate(); err != nil {
		return errorsmod.Wrapf(errortypes.ErrAppConfig, "invalid evm config value: %s", err.Error())
	}

	if err := c.JSONRPC.Validate(); err != nil {
		return errorsmod.Wrapf(errortypes.ErrAppConfig, "invalid json-rpc config value: %s", err.Error())
	}

	if err := c.TLS.Validate(); err != nil {
		return errorsmod.Wrapf(errortypes.ErrAppConfig, "invalid tls config value: %s", err.Error())
	}

	return c.Config.ValidateBasic()
}
