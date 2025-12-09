package types

import (
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	ethchains "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
)

var (
	DefaultMinObserverDelegation = sdkmath.LegacyMustNewDecFromStr("1000000000000000000000")
	DefaultBallotThreshold       = sdkmath.LegacyMustNewDecFromStr("0.66")

	// DefaultGasPriceMultiplier is the default gas price multiplier for all chains
	DefaultGasPriceMultiplier = sdkmath.LegacyMustNewDecFromStr("1.0")

	// DefaultEVMOutboundGasPriceMultiplier is the default gas price multiplier for EVM-chain outbond txs
	DefaultEVMOutboundGasPriceMultiplier = sdkmath.LegacyMustNewDecFromStr("1.2")

	// DefaultBTCOutboundGasPriceMultiplier is the default gas price multiplier for BTC outbond txs
	DefaultBTCOutboundGasPriceMultiplier = sdkmath.LegacyMustNewDecFromStr("2.0")
)

// Validate checks that the ConfirmationParams is valid
func (cp ConfirmationParams) Validate() error {
	switch {
	case cp.SafeInboundCount == 0:
		return errors.New("SafeInboundCount must be greater than 0")
	case cp.FastInboundCount > cp.SafeInboundCount:
		return errors.New("FastInboundCount must be less than or equal to SafeInboundCount")
	case cp.SafeOutboundCount == 0:
		return errors.New("SafeOutboundCount must be greater than 0")
	case cp.FastOutboundCount > cp.SafeOutboundCount:
		return errors.New("FastOutboundCount must be less than or equal to SafeOutboundCount")
	default:
		return nil
	}
}

// Validate checks all chain params correspond to a chain and there is no duplicate chain id
func (cpl ChainParamsList) Validate() error {
	// check all chain params correspond to a chain
	chainMap := make(map[int64]struct{})
	existingChainMap := make(map[int64]struct{})

	externalChainList := chains.DefaultChainsList()
	for _, chain := range externalChainList {
		chainMap[chain.ChainId] = struct{}{}
	}

	// validate the chain params and check for duplicates
	for _, chainParam := range cpl.ChainParams {
		if err := chainParam.Validate(); err != nil {
			return err
		}

		if _, ok := chainMap[chainParam.ChainId]; !ok {
			return fmt.Errorf("chain id %d not found in chain list", chainParam.ChainId)
		}
		if _, ok := existingChainMap[chainParam.ChainId]; ok {
			return fmt.Errorf("duplicated chain id %d found", chainParam.ChainId)
		}
		existingChainMap[chainParam.ChainId] = struct{}{}
	}
	return nil
}

// Validate performs basic checks on chain params
func (cp ChainParams) Validate() error {
	// don't validate ZetaChain params, because the validation will fail on existing params in the store
	// we might remove the ZetaChain params in the future, this is TBD
	if chains.IsZetaChain(cp.ChainId, nil) {
		return nil
	}

	// validate confirmation counts
	// TODO: ensure 'ConfirmationCount === 0' in the chain params update msg after migration (no other values allowed)
	// zetaclient will unconditionally use 'ConfirmationParams' fields in the future
	// https://github.com/zeta-chain/node/issues/3466
	if cp.ConfirmationCount == 0 {
		return errors.New("ConfirmationCount must be greater than 0")
	}
	if cp.ConfirmationParams == nil {
		return errors.New("confirmation params cannot be nil")
	}
	if err := cp.ConfirmationParams.Validate(); err != nil {
		return errors.Wrap(err, "invalid confirmation params")
	}

	// validate tickers and intervals
	if cp.GasPriceTicker <= 0 || cp.GasPriceTicker > 300 {
		return fmt.Errorf("GasPriceTicker %d out of range", cp.GasPriceTicker)
	}
	if cp.InboundTicker <= 0 || cp.InboundTicker > 300 {
		return fmt.Errorf("InboundTicker %d out of range", cp.InboundTicker)
	}
	if cp.OutboundTicker <= 0 || cp.OutboundTicker > 300 {
		return fmt.Errorf("OutboundTicker %d out of range", cp.OutboundTicker)
	}
	if cp.OutboundScheduleInterval == 0 || cp.OutboundScheduleInterval > 100 { // 600 secs
		return fmt.Errorf(
			"OutboundScheduleInterval %d out of range",
			cp.OutboundScheduleInterval,
		)
	}
	if cp.OutboundScheduleLookahead == 0 || cp.OutboundScheduleLookahead > 500 { // 500 cctxs
		return fmt.Errorf(
			"OutboundScheduleLookahead %d out of range",
			cp.OutboundScheduleLookahead,
		)
	}

	// if WatchUtxoTicker defined, check validity
	if cp.WatchUtxoTicker > 300 {
		return fmt.Errorf(
			"WatchUtxoTicker %d out of range",
			cp.WatchUtxoTicker,
		)
	}

	// if contract addresses are defined, check validity
	if cp.ZetaTokenContractAddress != "" && !validChainContractAddress(cp.ZetaTokenContractAddress) {
		return fmt.Errorf(
			"invalid ZetaTokenContractAddress %s",
			cp.ZetaTokenContractAddress,
		)
	}
	if cp.ConnectorContractAddress != "" && !validChainContractAddress(cp.ConnectorContractAddress) {
		return fmt.Errorf(
			"invalid ConnectorContractAddress %s",
			cp.ConnectorContractAddress,
		)
	}
	if cp.Erc20CustodyContractAddress != "" && !validChainContractAddress(cp.Erc20CustodyContractAddress) {
		return fmt.Errorf(
			"invalid Erc20CustodyContractAddress %s",
			cp.Erc20CustodyContractAddress,
		)
	}

	if cp.BallotThreshold.IsNil() || cp.BallotThreshold.GT(sdkmath.LegacyOneDec()) {
		return ErrParamsThreshold
	}

	if cp.MinObserverDelegation.IsNil() {
		return ErrParamsMinObserverDelegation
	}

	// validate gas price multiplier
	if cp.GasPriceMultiplier.IsNil() {
		return fmt.Errorf("GasPriceMultiplier cannot be nil")
	}

	if !cp.GasPriceMultiplier.IsPositive() ||
		cp.GasPriceMultiplier.GT(sdkmath.LegacyNewDec(10)) ||
		cp.GasPriceMultiplier.LT(sdkmath.LegacyMustNewDecFromStr("0.1")) {
		return fmt.Errorf("GasPriceMultiplier %s out of range", cp.GasPriceMultiplier.String())
	}

	return nil
}

// IsInboundFastConfirmationEnabled returns true if fast inbound confirmation is enabled.
func (cp ChainParams) IsInboundFastConfirmationEnabled() bool {
	return cp.ConfirmationParams.FastInboundCount > 0 &&
		cp.ConfirmationParams.FastInboundCount < cp.ConfirmationParams.SafeInboundCount
}

// InboundConfirmationSafe returns the safe number of confirmation for inbound observation.
func (cp ChainParams) InboundConfirmationSafe() uint64 {
	return cp.ConfirmationParams.SafeInboundCount
}

// InboundConfirmationFast returns the fast number of confirmation for inbound observation.
// It falls back to safe confirmation count if fast mode is disabled.
func (cp ChainParams) InboundConfirmationFast() uint64 {
	if cp.ConfirmationParams.FastInboundCount > 0 {
		return cp.ConfirmationParams.FastInboundCount
	}
	return cp.ConfirmationParams.SafeInboundCount
}

// OutboundConfirmationSafe returns the safe number of confirmation for outbound observation.
func (cp ChainParams) OutboundConfirmationSafe() uint64 {
	return cp.ConfirmationParams.SafeOutboundCount
}

// OutboundConfirmationFast returns the fast number of confirmation for outbound observation.
// It falls back to safe confirmation count if fast mode is disabled.
func (cp ChainParams) OutboundConfirmationFast() uint64 {
	if cp.ConfirmationParams.FastOutboundCount > 0 {
		return cp.ConfirmationParams.FastOutboundCount
	}
	return cp.ConfirmationParams.SafeOutboundCount
}

func validChainContractAddress(address string) bool {
	if !strings.HasPrefix(address, "0x") {
		return false
	}
	return ethchains.IsHexAddress(address)
}

// GetDefaultChainParams returns a list of default chain params
// TODO: remove this function
// https://github.com/zeta-chain/node-private/issues/100
func GetDefaultChainParams() ChainParamsList {
	return ChainParamsList{
		ChainParams: []*ChainParams{
			GetDefaultEthMainnetChainParams(),
			GetDefaultBscMainnetChainParams(),
			GetDefaultBtcMainnetChainParams(),
			GetDefaultGoerliTestnetChainParams(),
			GetDefaultBscTestnetChainParams(),
			GetDefaultMumbaiTestnetChainParams(),
			GetDefaultBtcTestnetChainParams(),
			GetDefaultBtcRegtestChainParams(),
			GetDefaultGoerliLocalnetChainParams(),
		},
	}
}

func GetDefaultEthMainnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.Ethereum.ChainId,
		ConfirmationCount:           14,
		ZetaTokenContractAddress:    constant.EVMZeroAddress,
		ConnectorContractAddress:    constant.EVMZeroAddress,
		Erc20CustodyContractAddress: constant.EVMZeroAddress,
		InboundTicker:               12,
		OutboundTicker:              15,
		WatchUtxoTicker:             0,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
		ConfirmationParams: &ConfirmationParams{
			SafeInboundCount:  14,
			FastInboundCount:  14,
			SafeOutboundCount: 14,
			FastOutboundCount: 14,
		},
		GasPriceMultiplier: DefaultEVMOutboundGasPriceMultiplier,
	}
}

func GetDefaultBscMainnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.BscMainnet.ChainId,
		ConfirmationCount:           14,
		ZetaTokenContractAddress:    constant.EVMZeroAddress,
		ConnectorContractAddress:    constant.EVMZeroAddress,
		Erc20CustodyContractAddress: constant.EVMZeroAddress,
		InboundTicker:               5,
		OutboundTicker:              15,
		WatchUtxoTicker:             0,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
		ConfirmationParams: &ConfirmationParams{
			SafeInboundCount:  14,
			FastInboundCount:  14,
			SafeOutboundCount: 14,
			FastOutboundCount: 14,
		},
		GasPriceMultiplier: DefaultEVMOutboundGasPriceMultiplier,
	}
}

func GetDefaultBtcMainnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.BitcoinMainnet.ChainId,
		ConfirmationCount:           2,
		ZetaTokenContractAddress:    constant.EVMZeroAddress,
		ConnectorContractAddress:    constant.EVMZeroAddress,
		Erc20CustodyContractAddress: constant.EVMZeroAddress,
		WatchUtxoTicker:             30,
		InboundTicker:               120,
		OutboundTicker:              60,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
		ConfirmationParams: &ConfirmationParams{
			SafeInboundCount:  2,
			FastInboundCount:  2,
			SafeOutboundCount: 2,
			FastOutboundCount: 2,
		},
		GasPriceMultiplier: DefaultBTCOutboundGasPriceMultiplier,
	}
}

func GetDefaultGoerliTestnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:           chains.Goerli.ChainId,
		ConfirmationCount: 6,
		// This is the actual Zeta token Goerli testnet, we need to specify this address for the integration tests to pass
		ZetaTokenContractAddress:    "0x0000c304d2934c00db1d51995b9f6996affd17c0",
		ConnectorContractAddress:    constant.EVMZeroAddress,
		Erc20CustodyContractAddress: constant.EVMZeroAddress,
		InboundTicker:               12,
		OutboundTicker:              15,
		WatchUtxoTicker:             0,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
		ConfirmationParams: &ConfirmationParams{
			SafeInboundCount:  6,
			FastInboundCount:  6,
			SafeOutboundCount: 6,
			FastOutboundCount: 6,
		},
		GasPriceMultiplier: DefaultEVMOutboundGasPriceMultiplier,
	}
}

func GetDefaultBscTestnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.BscTestnet.ChainId,
		ConfirmationCount:           6,
		ZetaTokenContractAddress:    constant.EVMZeroAddress,
		ConnectorContractAddress:    constant.EVMZeroAddress,
		Erc20CustodyContractAddress: constant.EVMZeroAddress,
		InboundTicker:               5,
		OutboundTicker:              15,
		WatchUtxoTicker:             0,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
		ConfirmationParams: &ConfirmationParams{
			SafeInboundCount:  6,
			FastInboundCount:  6,
			SafeOutboundCount: 6,
			FastOutboundCount: 6,
		},
		GasPriceMultiplier: DefaultEVMOutboundGasPriceMultiplier,
	}
}

func GetDefaultMumbaiTestnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.Mumbai.ChainId,
		ConfirmationCount:           12,
		ZetaTokenContractAddress:    constant.EVMZeroAddress,
		ConnectorContractAddress:    constant.EVMZeroAddress,
		Erc20CustodyContractAddress: constant.EVMZeroAddress,
		InboundTicker:               2,
		OutboundTicker:              15,
		WatchUtxoTicker:             0,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
		ConfirmationParams: &ConfirmationParams{
			SafeInboundCount:  12,
			FastInboundCount:  12,
			SafeOutboundCount: 12,
			FastOutboundCount: 12,
		},
		GasPriceMultiplier: DefaultEVMOutboundGasPriceMultiplier,
	}
}

func GetDefaultBtcTestnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.BitcoinTestnet.ChainId,
		ConfirmationCount:           2,
		ZetaTokenContractAddress:    constant.EVMZeroAddress,
		ConnectorContractAddress:    constant.EVMZeroAddress,
		Erc20CustodyContractAddress: constant.EVMZeroAddress,
		WatchUtxoTicker:             30,
		InboundTicker:               120,
		OutboundTicker:              12,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   100,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
		ConfirmationParams: &ConfirmationParams{
			SafeInboundCount:  2,
			FastInboundCount:  2,
			SafeOutboundCount: 2,
			FastOutboundCount: 2,
		},
		GasPriceMultiplier: DefaultBTCOutboundGasPriceMultiplier,
	}
}

func GetDefaultBtcRegtestChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.BitcoinRegtest.ChainId,
		ConfirmationCount:           1,
		ZetaTokenContractAddress:    constant.EVMZeroAddress,
		ConnectorContractAddress:    constant.EVMZeroAddress,
		Erc20CustodyContractAddress: constant.EVMZeroAddress,
		GasPriceTicker:              5,
		WatchUtxoTicker:             1,
		InboundTicker:               1,
		OutboundTicker:              2,
		OutboundScheduleInterval:    1,
		OutboundScheduleLookahead:   5,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
		ConfirmationParams: &ConfirmationParams{
			SafeInboundCount:  1,
			FastInboundCount:  1,
			SafeOutboundCount: 1,
			FastOutboundCount: 1,
		},
		GasPriceMultiplier: DefaultBTCOutboundGasPriceMultiplier,
	}
}

func GetDefaultGoerliLocalnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.GoerliLocalnet.ChainId,
		ConfirmationCount:           1,
		ZetaTokenContractAddress:    "0x733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9",
		ConnectorContractAddress:    "0xD28D6A0b8189305551a0A8bd247a6ECa9CE781Ca",
		Erc20CustodyContractAddress: "0xff3135df4F2775f4091b81f4c7B6359CfA07862a",
		InboundTicker:               2,
		OutboundTicker:              1,
		WatchUtxoTicker:             0,
		GasPriceTicker:              5,
		OutboundScheduleInterval:    1,
		OutboundScheduleLookahead:   50,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
		GatewayAddress:              "0xF0deebCB0E9C829519C4baa794c5445171973826",
		ConfirmationParams: &ConfirmationParams{
			SafeInboundCount:  1,
			FastInboundCount:  1,
			SafeOutboundCount: 1,
			FastOutboundCount: 1,
		},
		DisableTssBlockScan: true,
		GasPriceMultiplier:  DefaultEVMOutboundGasPriceMultiplier,
	}
}

func GetDefaultZetaPrivnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.ZetaChainPrivnet.ChainId,
		ConfirmationCount:           1,
		ZetaTokenContractAddress:    constant.EVMZeroAddress,
		ConnectorContractAddress:    constant.EVMZeroAddress,
		Erc20CustodyContractAddress: constant.EVMZeroAddress,
		InboundTicker:               2,
		OutboundTicker:              2,
		WatchUtxoTicker:             0,
		GasPriceTicker:              5,
		OutboundScheduleInterval:    0,
		OutboundScheduleLookahead:   0,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
		GasPriceMultiplier:          DefaultGasPriceMultiplier,
	}
}

// ChainParamsEqual returns true if two chain params are equal
func ChainParamsEqual(params1, params2 ChainParams) bool {
	return params1.ChainId == params2.ChainId &&
		params1.ConfirmationCount == params2.ConfirmationCount &&
		params1.ZetaTokenContractAddress == params2.ZetaTokenContractAddress &&
		params1.ConnectorContractAddress == params2.ConnectorContractAddress &&
		params1.Erc20CustodyContractAddress == params2.Erc20CustodyContractAddress &&
		params1.InboundTicker == params2.InboundTicker &&
		params1.OutboundTicker == params2.OutboundTicker &&
		params1.WatchUtxoTicker == params2.WatchUtxoTicker &&
		params1.GasPriceTicker == params2.GasPriceTicker &&
		params1.OutboundScheduleInterval == params2.OutboundScheduleInterval &&
		params1.OutboundScheduleLookahead == params2.OutboundScheduleLookahead &&
		params1.BallotThreshold.Equal(params2.BallotThreshold) &&
		params1.MinObserverDelegation.Equal(params2.MinObserverDelegation) &&
		params1.IsSupported == params2.IsSupported &&
		params1.GatewayAddress == params2.GatewayAddress &&
		confirmationParamsEqual(params1.ConfirmationParams, params2.ConfirmationParams) &&
		params1.DisableTssBlockScan == params2.DisableTssBlockScan &&
		params1.GasPriceMultiplier.Equal(params2.GasPriceMultiplier)
}

// confirmationParamsEqual returns true if two confirmation params are equal
func confirmationParamsEqual(a, b *ConfirmationParams) bool {
	if a == b {
		return true
	}
	if (a == nil) || (b == nil) {
		return false
	}

	return *a == *b
}
