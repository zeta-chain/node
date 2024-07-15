package config

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/yaml.v3"
)

// DoubleQuotedString forces a string to be double quoted when marshaling to yaml.
// This is required because of https://github.com/go-yaml/yaml/issues/847
type DoubleQuotedString string

func (s DoubleQuotedString) MarshalYAML() (interface{}, error) {
	return yaml.Node{
		Value: string(s),
		Kind:  yaml.ScalarNode,
		Style: yaml.DoubleQuotedStyle,
	}, nil
}

func (s DoubleQuotedString) String() string {
	return string(s)
}

func (s DoubleQuotedString) AsEVMAddress() (ethcommon.Address, error) {
	if !ethcommon.IsHexAddress(string(s)) {
		return ethcommon.Address{}, fmt.Errorf("invalid evm address: %s", string(s))
	}
	return ethcommon.HexToAddress(string(s)), nil
}

// Config contains the configuration for the e2e test
type Config struct {
	// Default account to use when running tests and running setup
	DefaultAccount     Account            `yaml:"default_account"`
	AdditionalAccounts AdditionalAccounts `yaml:"additional_accounts"`
	PolicyAccounts     PolicyAccounts     `yaml:"policy_accounts"`
	RPCs               RPCs               `yaml:"rpcs"`
	Contracts          Contracts          `yaml:"contracts"`
	ZetaChainID        string             `yaml:"zeta_chain_id"`
}

// Account contains configuration for an account
type Account struct {
	RawBech32Address DoubleQuotedString `yaml:"bech32_address"`
	RawEVMAddress    DoubleQuotedString `yaml:"evm_address"`
	RawPrivateKey    DoubleQuotedString `yaml:"private_key"`
}

// AdditionalAccounts are extra accounts required to run specific tests
type AdditionalAccounts struct {
	UserERC20      Account `yaml:"user_erc20"`
	UserZetaTest   Account `yaml:"user_zeta_test"`
	UserZEVMMPTest Account `yaml:"user_zevm_mp_test"`
	UserBitcoin    Account `yaml:"user_bitcoin"`
	UserEther      Account `yaml:"user_ether"`
	UserMisc       Account `yaml:"user_misc"`
	UserAdmin      Account `yaml:"user_admin"`
}

type PolicyAccounts struct {
	EmergencyPolicyAccount   Account `yaml:"emergency_policy_account"`
	OperationalPolicyAccount Account `yaml:"operational_policy_account"`
	AdminPolicyAccount       Account `yaml:"admin_policy_account"`
}

// RPCs contains the configuration for the RPC endpoints
type RPCs struct {
	Zevm         string     `yaml:"zevm"`
	EVM          string     `yaml:"evm"`
	Bitcoin      BitcoinRPC `yaml:"bitcoin"`
	ZetaCoreGRPC string     `yaml:"zetacore_grpc"`
	ZetaCoreRPC  string     `yaml:"zetacore_rpc"`
}

// BitcoinRPC contains the configuration for the Bitcoin RPC endpoint
type BitcoinRPC struct {
	User         string             `yaml:"user"`
	Pass         string             `yaml:"pass"`
	Host         string             `yaml:"host"`
	HTTPPostMode bool               `yaml:"http_post_mode"`
	DisableTLS   bool               `yaml:"disable_tls"`
	Params       BitcoinNetworkType `yaml:"params"`
}

// Contracts contains the addresses of predeployed contracts
type Contracts struct {
	EVM  EVM  `yaml:"evm"`
	ZEVM ZEVM `yaml:"zevm"`
}

// EVM contains the addresses of predeployed contracts on the EVM chain
type EVM struct {
	ZetaEthAddr      DoubleQuotedString `yaml:"zeta_eth"`
	ConnectorEthAddr DoubleQuotedString `yaml:"connector_eth"`
	CustodyAddr      DoubleQuotedString `yaml:"custody"`
	ERC20            DoubleQuotedString `yaml:"erc20"`
	TestDappAddr     DoubleQuotedString `yaml:"test_dapp"`
}

// ZEVM contains the addresses of predeployed contracts on the zEVM chain
type ZEVM struct {
	SystemContractAddr DoubleQuotedString `yaml:"system_contract"`
	ETHZRC20Addr       DoubleQuotedString `yaml:"eth_zrc20"`
	ERC20ZRC20Addr     DoubleQuotedString `yaml:"erc20_zrc20"`
	BTCZRC20Addr       DoubleQuotedString `yaml:"btc_zrc20"`
	UniswapFactoryAddr DoubleQuotedString `yaml:"uniswap_factory"`
	UniswapRouterAddr  DoubleQuotedString `yaml:"uniswap_router"`
	ConnectorZEVMAddr  DoubleQuotedString `yaml:"connector_zevm"`
	WZetaAddr          DoubleQuotedString `yaml:"wzeta"`
	ZEVMSwapAppAddr    DoubleQuotedString `yaml:"zevm_swap_app"`
	ContextAppAddr     DoubleQuotedString `yaml:"context_app"`
	TestDappAddr       DoubleQuotedString `yaml:"test_dapp"`
}

// DefaultConfig returns the default config using values for localnet testing
func DefaultConfig() Config {
	return Config{
		RPCs: RPCs{
			Zevm: "http://zetacore0:8545",
			EVM:  "http://eth:8545",
			Bitcoin: BitcoinRPC{
				Host:         "bitcoin:18443",
				User:         "smoketest",
				Pass:         "123",
				HTTPPostMode: true,
				DisableTLS:   true,
				Params:       Regnet,
			},
			ZetaCoreGRPC: "zetacore0:9090",
			ZetaCoreRPC:  "http://zetacore0:26657",
		},
		ZetaChainID: "athens_101-1",
		Contracts: Contracts{
			EVM: EVM{
				ERC20: "0xff3135df4F2775f4091b81f4c7B6359CfA07862a",
			},
		},
	}
}

// ReadConfig reads the config file
func ReadConfig(file string) (config Config, err error) {
	if file == "" {
		return Config{}, errors.New("file name cannot be empty")
	}

	// #nosec G304 -- this is a config file
	b, err := os.ReadFile(file)
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return Config{}, err
	}
	if err := config.Validate(); err != nil {
		return Config{}, err
	}

	return
}

// WriteConfig writes the config file
func WriteConfig(file string, config Config) error {
	if file == "" {
		return errors.New("file name cannot be empty")
	}

	b, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, b, 0600)
	if err != nil {
		return err
	}
	return nil
}

// AsSlice gets all accounts as a slice rather than a
// struct
func (a AdditionalAccounts) AsSlice() []Account {
	return []Account{
		a.UserERC20,
		a.UserZetaTest,
		a.UserZEVMMPTest,
		a.UserBitcoin,
		a.UserEther,
		a.UserMisc,
		a.UserAdmin,
	}
}

func (a PolicyAccounts) AsSlice() []Account {
	return []Account{
		a.EmergencyPolicyAccount,
		a.OperationalPolicyAccount,
		a.AdminPolicyAccount,
	}
}

// Validate validates the config
func (c Config) Validate() error {
	if c.RPCs.Bitcoin.Params != Mainnet &&
		c.RPCs.Bitcoin.Params != Testnet3 &&
		c.RPCs.Bitcoin.Params != Regnet {
		return errors.New("invalid bitcoin params")
	}

	err := c.DefaultAccount.Validate()
	if err != nil {
		return fmt.Errorf("validating deployer account: %w", err)
	}

	accounts := c.AdditionalAccounts.AsSlice()
	for i, account := range accounts {
		if account.RawEVMAddress == "" {
			continue
		}
		err := account.Validate()
		if err != nil {
			return fmt.Errorf("validating additional account %d: %w", i, err)
		}
	}

	policyAccounts := c.PolicyAccounts.AsSlice()
	for i, account := range policyAccounts {
		if account.RawEVMAddress == "" {
			continue
		}
		err := account.Validate()
		if err != nil {
			return fmt.Errorf("validating policy account %d: %w", i, err)
		}
	}

	return nil
}

// GenerateKeys generates new key pairs for all accounts
func (c *Config) GenerateKeys() error {
	var err error
	c.DefaultAccount, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserERC20, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserZetaTest, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserZEVMMPTest, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserBitcoin, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserEther, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserMisc, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserAdmin, err = generateAccount()
	if err != nil {
		return err
	}
	c.PolicyAccounts.EmergencyPolicyAccount, err = generateAccount()
	if err != nil {
		return err
	}
	c.PolicyAccounts.OperationalPolicyAccount, err = generateAccount()
	if err != nil {
		return err
	}
	c.PolicyAccounts.AdminPolicyAccount, err = generateAccount()
	if err != nil {
		return err
	}
	return nil
}

func (a Account) EVMAddress() ethcommon.Address {
	return ethcommon.HexToAddress(a.RawEVMAddress.String())
}

func (a Account) PrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(a.RawPrivateKey.String())
}

// Validate that the address and the private key specified in the
// config actually match
func (a Account) Validate() error {
	privateKey, err := a.PrivateKey()
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}
	privateKeyAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	if a.RawEVMAddress.String() != privateKeyAddress.Hex() {
		return fmt.Errorf(
			"address derived from private key (%s) does not match configured address (%s)",
			privateKeyAddress,
			a.RawEVMAddress,
		)
	}
	_, _, err = bech32.DecodeAndConvert(a.RawBech32Address.String())
	if err != nil {
		return fmt.Errorf("decoding bech32 address: %w", err)
	}
	bech32PrivateKeyAddress, err := bech32.ConvertAndEncode("zeta", privateKeyAddress.Bytes())
	if err != nil {
		return fmt.Errorf("encoding private key to bech32: %w", err)
	}
	if a.RawBech32Address.String() != bech32PrivateKeyAddress {
		return fmt.Errorf(
			"address derived from private key (%s) does not match configured address (%s)",
			bech32PrivateKeyAddress,
			a.RawBech32Address,
		)
	}
	return nil
}

// BitcoinNetworkType is a custom type to represent allowed network types
type BitcoinNetworkType string

// Enum values for BitcoinNetworkType
const (
	Mainnet  BitcoinNetworkType = "mainnet"
	Testnet3 BitcoinNetworkType = "testnet3"
	Regnet   BitcoinNetworkType = "regnet"
)

// GetParams returns the chaincfg.Params for the BitcoinNetworkType
func (bnt BitcoinNetworkType) GetParams() (chaincfg.Params, error) {
	switch bnt {
	case Mainnet:
		return chaincfg.MainNetParams, nil
	case Testnet3:
		return chaincfg.TestNet3Params, nil
	case Regnet:
		return chaincfg.RegressionNetParams, nil
	default:
		return chaincfg.Params{}, fmt.Errorf("invalid bitcoin params %s", bnt)
	}
}

func generateAccount() (Account, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return Account{}, fmt.Errorf("generating private key: %w", err)
	}
	// encode private key and strip 0x prefix
	encodedPrivateKey := hexutil.Encode(crypto.FromECDSA(privateKey))
	encodedPrivateKey = strings.TrimPrefix(encodedPrivateKey, "0x")

	evmAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	bech32Address, err := bech32.ConvertAndEncode("zeta", evmAddress.Bytes())
	if err != nil {
		return Account{}, fmt.Errorf("bech32 convert and encode: %w", err)
	}

	return Account{
		RawEVMAddress:    DoubleQuotedString(evmAddress.String()),
		RawPrivateKey:    DoubleQuotedString(encodedPrivateKey),
		RawBech32Address: DoubleQuotedString(bech32Address),
	}, nil
}
