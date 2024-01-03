package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
)

const (
	zeroAddress = "0x0000000000000000000000000000000000000000"
)

var (
	DefaultMinObserverDelegation = sdk.MustNewDecFromStr("1000000000000000000000")
	DefaultBallotThreshold       = sdk.MustNewDecFromStr("0.66")
)

// Validate checks all core params correspond to a chain and there is no duplicate chain id
func (cpl CoreParamsList) Validate() error {
	// check all core params correspond to a chain
	chainMap := make(map[int64]struct{})
	existingChainMap := make(map[int64]struct{})

	externalChainList := common.DefaultChainsList()
	for _, chain := range externalChainList {
		chainMap[chain.ChainId] = struct{}{}
	}

	// validate the core params and check for duplicates
	for _, coreParam := range cpl.CoreParams {
		if err := ValidateCoreParams(coreParam); err != nil {
			return err
		}

		if _, ok := chainMap[coreParam.ChainId]; !ok {
			return fmt.Errorf("chain id %d not found in chain list", coreParam.ChainId)
		}
		if _, ok := existingChainMap[coreParam.ChainId]; ok {
			return fmt.Errorf("duplicated chain id %d found", coreParam.ChainId)
		}
		existingChainMap[coreParam.ChainId] = struct{}{}
	}
	return nil
}

// ValidateCoreParams performs some basic checks on core params
func ValidateCoreParams(params *CoreParams) error {
	if params == nil {
		return fmt.Errorf("core params cannot be nil")
	}
	chain := common.GetChainFromChainID(params.ChainId)
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
	if params.InTxTicker <= 0 || params.InTxTicker > 300 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "InTxTicker %d out of range", params.InTxTicker)
	}
	if params.OutTxTicker <= 0 || params.OutTxTicker > 300 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "OutTxTicker %d out of range", params.OutTxTicker)
	}
	if params.OutboundTxScheduleInterval == 0 || params.OutboundTxScheduleInterval > 100 { // 600 secs
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "OutboundTxScheduleInterval %d out of range", params.OutboundTxScheduleInterval)
	}
	if params.OutboundTxScheduleLookahead == 0 || params.OutboundTxScheduleLookahead > 500 { // 500 cctxs
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "OutboundTxScheduleLookahead %d out of range", params.OutboundTxScheduleLookahead)
	}

	// chain type specific checks
	if common.IsBitcoinChain(params.ChainId) {
		if params.WatchUtxoTicker == 0 || params.WatchUtxoTicker > 300 {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "WatchUtxoTicker %d out of range", params.WatchUtxoTicker)
		}
	}
	if common.IsEVMChain(params.ChainId) {
		if !validCoreContractAddress(params.ZetaTokenContractAddress) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ZetaTokenContractAddress %s", params.ZetaTokenContractAddress)
		}
		if !validCoreContractAddress(params.ConnectorContractAddress) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ConnectorContractAddress %s", params.ConnectorContractAddress)
		}
		if !validCoreContractAddress(params.Erc20CustodyContractAddress) {
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

func validCoreContractAddress(address string) bool {
	if !strings.HasPrefix(address, "0x") {
		return false
	}
	return ethcommon.IsHexAddress(address)
}

// GetDefaultCoreParams returns a list of default core params
// TODO: remove this function
// https://github.com/zeta-chain/node-private/issues/100
func GetDefaultCoreParams() CoreParamsList {
	return CoreParamsList{
		CoreParams: []*CoreParams{
			GetDefaultEthMainnetCoreParams(),
			GetDefaultBscMainnetCoreParams(),
			GetDefaultBtcMainnetCoreParams(),
			GetDefaultGoerliTestnetCoreParams(),
			GetDefaultBscTestnetCoreParams(),
			GetDefaultMumbaiTestnetCoreParams(),
			GetDefaultBtcTestnetCoreParams(),
			GetDefaultBtcRegtestCoreParams(),
			GetDefaultGoerliLocalnetCoreParams(),
			GetDefaultZetaPrivnetCoreParams(),
		},
	}
}

func GetDefaultEthMainnetCoreParams() *CoreParams {
	return &CoreParams{
		ChainId:                     common.EthChain().ChainId,
		ConfirmationCount:           14,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		InTxTicker:                  12,
		OutTxTicker:                 15,
		WatchUtxoTicker:             0,
		GasPriceTicker:              30,
		OutboundTxScheduleInterval:  30,
		OutboundTxScheduleLookahead: 60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultBscMainnetCoreParams() *CoreParams {
	return &CoreParams{
		ChainId:                     common.BscMainnetChain().ChainId,
		ConfirmationCount:           14,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		InTxTicker:                  5,
		OutTxTicker:                 15,
		WatchUtxoTicker:             0,
		GasPriceTicker:              30,
		OutboundTxScheduleInterval:  30,
		OutboundTxScheduleLookahead: 60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultBtcMainnetCoreParams() *CoreParams {
	return &CoreParams{
		ChainId:                     common.BtcMainnetChain().ChainId,
		ConfirmationCount:           2,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		WatchUtxoTicker:             30,
		InTxTicker:                  120,
		OutTxTicker:                 60,
		GasPriceTicker:              30,
		OutboundTxScheduleInterval:  30,
		OutboundTxScheduleLookahead: 60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultGoerliTestnetCoreParams() *CoreParams {
	return &CoreParams{
		ChainId:           common.GoerliChain().ChainId,
		ConfirmationCount: 6,
		// This is the actual Zeta token Goerli testnet, we need to specify this address for the integration tests to pass
		ZetaTokenContractAddress:    "0x0000c304d2934c00db1d51995b9f6996affd17c0",
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		InTxTicker:                  12,
		OutTxTicker:                 15,
		WatchUtxoTicker:             0,
		GasPriceTicker:              30,
		OutboundTxScheduleInterval:  30,
		OutboundTxScheduleLookahead: 60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultBscTestnetCoreParams() *CoreParams {
	return &CoreParams{
		ChainId:                     common.BscTestnetChain().ChainId,
		ConfirmationCount:           6,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		InTxTicker:                  5,
		OutTxTicker:                 15,
		WatchUtxoTicker:             0,
		GasPriceTicker:              30,
		OutboundTxScheduleInterval:  30,
		OutboundTxScheduleLookahead: 60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultMumbaiTestnetCoreParams() *CoreParams {
	return &CoreParams{
		ChainId:                     common.MumbaiChain().ChainId,
		ConfirmationCount:           12,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		InTxTicker:                  2,
		OutTxTicker:                 15,
		WatchUtxoTicker:             0,
		GasPriceTicker:              30,
		OutboundTxScheduleInterval:  30,
		OutboundTxScheduleLookahead: 60,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultBtcTestnetCoreParams() *CoreParams {
	return &CoreParams{
		ChainId:                     common.BtcTestNetChain().ChainId,
		ConfirmationCount:           2,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		WatchUtxoTicker:             30,
		InTxTicker:                  120,
		OutTxTicker:                 12,
		GasPriceTicker:              30,
		OutboundTxScheduleInterval:  30,
		OutboundTxScheduleLookahead: 100,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultBtcRegtestCoreParams() *CoreParams {
	return &CoreParams{
		ChainId:                     common.BtcRegtestChain().ChainId,
		ConfirmationCount:           1,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		GasPriceTicker:              5,
		WatchUtxoTicker:             1,
		InTxTicker:                  1,
		OutTxTicker:                 2,
		OutboundTxScheduleInterval:  2,
		OutboundTxScheduleLookahead: 5,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultGoerliLocalnetCoreParams() *CoreParams {
	return &CoreParams{
		ChainId:                     common.GoerliLocalnetChain().ChainId,
		ConfirmationCount:           1,
		ZetaTokenContractAddress:    "0x733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9",
		ConnectorContractAddress:    "0xD28D6A0b8189305551a0A8bd247a6ECa9CE781Ca",
		Erc20CustodyContractAddress: "0xff3135df4F2775f4091b81f4c7B6359CfA07862a",
		InTxTicker:                  2,
		OutTxTicker:                 2,
		WatchUtxoTicker:             0,
		GasPriceTicker:              5,
		OutboundTxScheduleInterval:  2,
		OutboundTxScheduleLookahead: 5,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
func GetDefaultZetaPrivnetCoreParams() *CoreParams {
	return &CoreParams{
		ChainId:                     common.ZetaPrivnetChain().ChainId,
		ConfirmationCount:           1,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    zeroAddress,
		Erc20CustodyContractAddress: zeroAddress,
		InTxTicker:                  2,
		OutTxTicker:                 2,
		WatchUtxoTicker:             0,
		GasPriceTicker:              5,
		OutboundTxScheduleInterval:  0,
		OutboundTxScheduleLookahead: 0,
		BallotThreshold:             DefaultBallotThreshold,
		MinObserverDelegation:       DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}
