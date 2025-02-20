package config

import (
	"crypto/ecdsa"
	"encoding/hex"
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

	sui_utils "github.com/zeta-chain/node/e2e/utils/sui"
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
	DefaultAccount          Account                 `yaml:"default_account"`
	AdditionalAccounts      AdditionalAccounts      `yaml:"additional_accounts"`
	PolicyAccounts          PolicyAccounts          `yaml:"policy_accounts"`
	ObserverRelayerAccounts ObserverRelayerAccounts `yaml:"observer_relayer_accounts"`
	RPCs                    RPCs                    `yaml:"rpcs"`
	Contracts               Contracts               `yaml:"contracts"`
	ZetaChainID             string                  `yaml:"zeta_chain_id"`
}

// Account contains configuration for an account
type Account struct {
	RawBech32Address DoubleQuotedString `yaml:"bech32_address"`
	RawEVMAddress    DoubleQuotedString `yaml:"evm_address"`
	RawPrivateKey    DoubleQuotedString `yaml:"private_key"`
	SolanaAddress    DoubleQuotedString `yaml:"solana_address"`
	SolanaPrivateKey DoubleQuotedString `yaml:"solana_private_key"`
}

// AdditionalAccounts are extra accounts required to run specific tests
type AdditionalAccounts struct {
	UserLegacyERC20       Account `yaml:"user_legacy_erc20"`
	UserLegacyZeta        Account `yaml:"user_legacy_zeta"`
	UserLegacyZEVMMP      Account `yaml:"user_legacy_zevm_mp"`
	UserLegacyEther       Account `yaml:"user_legacy_ether"`
	UserBitcoinDeposit    Account `yaml:"user_bitcoin_deposit"`
	UserBitcoinWithdraw   Account `yaml:"user_bitcoin_withdraw"`
	UserSolana            Account `yaml:"user_solana"`
	UserSPL               Account `yaml:"user_spl"`
	UserMisc              Account `yaml:"user_misc"`
	UserAdmin             Account `yaml:"user_admin"`
	UserMigration         Account `yaml:"user_migration"` // used for TSS migration, TODO: rename (https://github.com/zeta-chain/node/issues/2780)
	UserPrecompile        Account `yaml:"user_precompile"`
	UserEther             Account `yaml:"user_ether"`
	UserERC20             Account `yaml:"user_erc20"`
	UserEtherRevert       Account `yaml:"user_ether_revert"`
	UserERC20Revert       Account `yaml:"user_erc20_revert"`
	UserEmissionsWithdraw Account `yaml:"user_emissions_withdraw"`
}

type PolicyAccounts struct {
	EmergencyPolicyAccount   Account `yaml:"emergency_policy_account"`
	OperationalPolicyAccount Account `yaml:"operational_policy_account"`
	AdminPolicyAccount       Account `yaml:"admin_policy_account"`
}

// ObserverRelayerAccounts are the accounts used by the observers to interact with gateway contracts in non-EVM chains (e.g. Solana)
type ObserverRelayerAccounts struct {
	// RelayerAccounts contains two relayer accounts used by zetaclient0 and zetaclient1
	RelayerAccounts [2]Account `yaml:"relayer_accounts"`
}

// RPCs contains the configuration for the RPC endpoints
type RPCs struct {
	Zevm              string     `yaml:"zevm"`
	EVM               string     `yaml:"evm"`
	Bitcoin           BitcoinRPC `yaml:"bitcoin"`
	Solana            string     `yaml:"solana"`
	TONSidecarURL     string     `yaml:"ton_sidecar_url"`
	Sui               string     `yaml:"sui"`
	SuiFaucet         string     `yaml:"sui_faucet"`
	ZetaCoreGRPC      string     `yaml:"zetacore_grpc"`
	ZetaCoreRPC       string     `yaml:"zetacore_rpc"`
	ZetaclientMetrics string     `yaml:"zetaclient_metrics"`
}

// BitcoinRPC contains the configuration for the Bitcoin RPC endpoint
type BitcoinRPC struct {
	User   string             `yaml:"user"`
	Pass   string             `yaml:"pass"`
	Host   string             `yaml:"host"`
	Params BitcoinNetworkType `yaml:"params"`
}

// Contracts contains the addresses of predeployed contracts
type Contracts struct {
	EVM    EVM    `yaml:"evm"`
	ZEVM   ZEVM   `yaml:"zevm"`
	Solana Solana `yaml:"solana"`
	Sui    Sui    `yaml:"sui"`
}

// Solana contains the addresses of predeployed contracts and accounts on the Solana chain
type Solana struct {
	GatewayProgramID DoubleQuotedString `yaml:"gateway_program_id"`
	SPLAddr          DoubleQuotedString `yaml:"spl"`
}

// Sui contains the addresses of predeployed contracts on the Sui chain
type Sui struct {
	GatewayPackageID DoubleQuotedString `yaml:"gateway_package_id"`
	GatewayObjectID  DoubleQuotedString `yaml:"gateway_object_id"`
	CoinType         DoubleQuotedString `yaml:"coin_type"`
}

// EVM contains the addresses of predeployed contracts on the EVM chain
type EVM struct {
	ZetaEthAddr      DoubleQuotedString `yaml:"zeta_eth"`
	ConnectorEthAddr DoubleQuotedString `yaml:"connector_eth"`
	CustodyAddr      DoubleQuotedString `yaml:"custody"`
	ERC20            DoubleQuotedString `yaml:"erc20"`
	TestDappAddr     DoubleQuotedString `yaml:"test_dapp"`
	Gateway          DoubleQuotedString `yaml:"gateway"`
	ERC20CustodyNew  DoubleQuotedString `yaml:"erc20_custody_new"`
	TestDAppV2Addr   DoubleQuotedString `yaml:"test_dapp_v2"`
}

// ZEVM contains the addresses of predeployed contracts on the zEVM chain
type ZEVM struct {
	SystemContractAddr DoubleQuotedString `yaml:"system_contract"`
	ETHZRC20Addr       DoubleQuotedString `yaml:"eth_zrc20"`
	ERC20ZRC20Addr     DoubleQuotedString `yaml:"erc20_zrc20"`
	BTCZRC20Addr       DoubleQuotedString `yaml:"btc_zrc20"`
	SOLZRC20Addr       DoubleQuotedString `yaml:"sol_zrc20"`
	SPLZRC20Addr       DoubleQuotedString `yaml:"spl_zrc20"`
	TONZRC20Addr       DoubleQuotedString `yaml:"ton_zrc20"`
	SUIZRC20Addr       DoubleQuotedString `yaml:"sui_zrc20"`
	UniswapFactoryAddr DoubleQuotedString `yaml:"uniswap_factory"`
	UniswapRouterAddr  DoubleQuotedString `yaml:"uniswap_router"`
	ConnectorZEVMAddr  DoubleQuotedString `yaml:"connector_zevm"`
	WZetaAddr          DoubleQuotedString `yaml:"wzeta"`
	ZEVMSwapAppAddr    DoubleQuotedString `yaml:"zevm_swap_app"`
	ContextAppAddr     DoubleQuotedString `yaml:"context_app"`
	TestDappAddr       DoubleQuotedString `yaml:"test_dapp"`
	Gateway            DoubleQuotedString `yaml:"gateway"`
	TestDAppV2Addr     DoubleQuotedString `yaml:"test_dapp_v2"`
}

// DefaultConfig returns the default config using values for localnet testing
func DefaultConfig() Config {
	return Config{
		RPCs: RPCs{
			Zevm: "http://zetacore0:8545",
			EVM:  "http://eth:8545",
			Bitcoin: BitcoinRPC{
				Host:   "bitcoin:18443",
				User:   "smoketest",
				Pass:   "123",
				Params: Regnet,
			},
			ZetaCoreGRPC: "zetacore0:9090",
			ZetaCoreRPC:  "http://zetacore0:26657",
			Solana:       "http://solana:8899",
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

	// #nosec G304 -- the variable is expected to be controlled by the user
	fHandle, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer fHandle.Close()

	// use a custom encoder so we can set the indentation level
	encoder := yaml.NewEncoder(fHandle)
	defer encoder.Close()
	encoder.SetIndent(2)
	err = encoder.Encode(config)
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	return nil
}

// AsSlice gets all accounts as a slice rather than a
// struct
func (a AdditionalAccounts) AsSlice() []Account {
	return []Account{
		a.UserLegacyERC20,
		a.UserLegacyZeta,
		a.UserLegacyZEVMMP,
		a.UserBitcoinDeposit,
		a.UserBitcoinWithdraw,
		a.UserSolana,
		a.UserSPL,
		a.UserLegacyEther,
		a.UserMisc,
		a.UserAdmin,
		a.UserMigration,
		a.UserPrecompile,
		a.UserEther,
		a.UserERC20,
		a.UserEtherRevert,
		a.UserERC20Revert,
		a.UserEmissionsWithdraw,
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
	c.AdditionalAccounts.UserLegacyERC20, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserLegacyZeta, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserLegacyZEVMMP, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserBitcoinDeposit, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserBitcoinWithdraw, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserSolana, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserSPL, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserLegacyEther, err = generateAccount()
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
	c.AdditionalAccounts.UserMigration, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserPrecompile, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserEther, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserERC20, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserEtherRevert, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserERC20Revert, err = generateAccount()
	if err != nil {
		return err
	}
	c.AdditionalAccounts.UserEmissionsWithdraw, err = generateAccount()
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

// SuiSigner derives the blake2b hash from the private key
// TODO: add a dedicated private key / address pair for Sui so it can be used with a regular account from a wallet
// for a localnet test this works
func (a Account) SuiSigner() (*sui_utils.SignerSecp256k1, error) {
	privateKeyBytes, err := hex.DecodeString(a.RawPrivateKey.String())
	if err != nil {
		return nil, fmt.Errorf("decode private key: %w", err)
	}
	signer := sui_utils.NewSignerSecp256k1FromSecretKey(privateKeyBytes)
	return signer, nil
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
