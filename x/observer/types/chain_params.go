package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethchains "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/pkg/chains"
)

const (
	zeroAddress = "0x0000000000000000000000000000000000000000"
)

var (
	DefaultMinObserverDelegation = sdk.MustNewDecFromStr("1000000000000000000000")
	DefaultBallotThreshold       = sdk.MustNewDecFromStr("0.66")
)

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
		if err := ValidateChainParams(chainParam); err != nil {
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

// ValidateChainParams performs some basic checks on chain params
func ValidateChainParams(params *ChainParams) error {
	if params == nil {
		return fmt.Errorf("chain params cannot be nil")
	}
	chain := chains.GetChainFromChainID(params.ChainId)
	if chain == nil {
		return fmt.Errorf("ChainId %d not supported", params.ChainId)
	}
	// zeta chain skips the rest of the checks for now
	if chain.IsZetaChain() {
		return nil
	}

	if params.ConfirmationCount == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "ConfirmationCount must be greater than 0")
	}
	if params.GasPriceTicker <= 0 || params.GasPriceTicker > 300 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "GasPriceTicker %d out of range", params.GasPriceTicker)
	}
	if params.InboundTicker <= 0 || params.InboundTicker > 300 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "InboundTicker %d out of range", params.InboundTicker)
	}
	if params.OutboundTicker <= 0 || params.OutboundTicker > 300 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "OutboundTicker %d out of range", params.OutboundTicker)
	}
	if params.OutboundScheduleInterval == 0 || params.OutboundScheduleInterval > 100 { // 600 secs
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "OutboundScheduleInterval %d out of range", params.OutboundScheduleInterval)
	}
	if params.OutboundScheduleLookahead == 0 || params.OutboundScheduleLookahead > 500 { // 500 cctxs
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "OutboundScheduleLookahead %d out of range", params.OutboundScheduleLookahead)
	}

	// chain type specific checks
	if chains.IsBitcoinChain(params.ChainId) {
		if params.WatchUtxoTicker == 0 || params.WatchUtxoTicker > 300 {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "WatchUtxoTicker %d out of range", params.WatchUtxoTicker)
		}
	}
	if chains.IsEVMChain(params.ChainId) {
		if !validChainContractAddress(params.ZetaTokenContractAddress) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ZetaTokenContractAddress %s", params.ZetaTokenContractAddress)
		}
		if !validChainContractAddress(params.ConnectorContractAddress) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ConnectorContractAddress %s", params.ConnectorContractAddress)
		}
		if !validChainContractAddress(params.Erc20CustodyContractAddress) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid Erc20CustodyContractAddress %s", params.Erc20CustodyContractAddress)
		}
	}

	if params.BallotThreshold.IsNil() || params.BallotThreshold.GT(sdk.OneDec()) {
		return ErrParamsThreshold
	}

	if params.MinObserverDelegation.IsNil() {
		return ErrParamsMinObserverDelegation
	}

	return nil
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
			GetDefaultZetaPrivnetChainParams(),
		},
	}
}

func GetDefaultEthMainnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.EthChain.ChainId,
		ConfirmationCount:           14,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		InboundTicker:               12,
		OutboundTicker:              15,
		WatchUtxoTicker:             0,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultBscMainnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.BscMainnetChain.ChainId,
		ConfirmationCount:           14,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		InboundTicker:               5,
		OutboundTicker:              15,
		WatchUtxoTicker:             0,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultBtcMainnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.BtcMainnetChain.ChainId,
		ConfirmationCount:           2,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		WatchUtxoTicker:             30,
		InboundTicker:               120,
		OutboundTicker:              60,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultGoerliTestnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:           chains.GoerliChain.ChainId,
		ConfirmationCount: 6,
		// This is the actual Zeta token Goerli testnet, we need to specify this address for the integration tests to pass
		ZetaTokenContractAddress:    "0x0000c304d2934c00db1d51995b9f6996affd17c0",
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		InboundTicker:               12,
		OutboundTicker:              15,
		WatchUtxoTicker:             0,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultBscTestnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.BscTestnetChain.ChainId,
		ConfirmationCount:           6,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		InboundTicker:               5,
		OutboundTicker:              15,
		WatchUtxoTicker:             0,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultMumbaiTestnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.MumbaiChain.ChainId,
		ConfirmationCount:           12,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		InboundTicker:               2,
		OutboundTicker:              15,
		WatchUtxoTicker:             0,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultBtcTestnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.BtcTestNetChain.ChainId,
		ConfirmationCount:           2,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		WatchUtxoTicker:             30,
		InboundTicker:               120,
		OutboundTicker:              12,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   100,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultBtcRegtestChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.BtcRegtestChain.ChainId,
		ConfirmationCount:           1,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		GasPriceTicker:              5,
		WatchUtxoTicker:             1,
		InboundTicker:               1,
		OutboundTicker:              2,
		OutboundScheduleInterval:    2,
		OutboundScheduleLookahead:   5,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultGoerliLocalnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.GoerliLocalnetChain.ChainId,
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
	}
}
func GetDefaultZetaPrivnetChainParams() *ChainParams {
	return &ChainParams{
		ChainId:                     chains.ZetaPrivnetChain.ChainId,
		ConfirmationCount:           1,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		InboundTicker:               2,
		OutboundTicker:              2,
		WatchUtxoTicker:             0,
		GasPriceTicker:              5,
		OutboundScheduleInterval:    0,
		OutboundScheduleLookahead:   0,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
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
		params1.IsSupported == params2.IsSupported
}
