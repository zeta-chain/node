package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	. "gopkg.in/check.v1"
)

type UpdateCoreParamsSuite struct {
	suite.Suite
	evmParams *CoreParams
	btcParams *CoreParams
}

var _ = Suite(&UpdateCoreParamsSuite{})

func TestUpdateCoreParamsSuiteSuite(t *testing.T) {
	suite.Run(t, new(UpdateCoreParamsSuite))
}

func (s *UpdateCoreParamsSuite) SetupTest() {
	s.evmParams = &CoreParams{
		ConfirmationCount:           1,
		GasPriceTicker:              1,
		InTxTicker:                  1,
		OutTxTicker:                 1,
		WatchUtxoTicker:             0,
		ZetaTokenContractAddress:    "0xA8D5060feb6B456e886F023709A2795373691E63",
		ConnectorContractAddress:    "0x733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9",
		Erc20CustodyContractAddress: "0xD28D6A0b8189305551a0A8bd247a6ECa9CE781Ca",
		ChainId:                     5,
		OutboundTxScheduleInterval:  1,
		OutboundTxScheduleLookahead: 1,
	}
	s.btcParams = &CoreParams{
		ConfirmationCount:           1,
		GasPriceTicker:              1,
		InTxTicker:                  1,
		OutTxTicker:                 1,
		WatchUtxoTicker:             1,
		ZetaTokenContractAddress:    "",
		ConnectorContractAddress:    "",
		Erc20CustodyContractAddress: "",
		ChainId:                     18332,
		OutboundTxScheduleInterval:  1,
		OutboundTxScheduleLookahead: 1,
	}
}

func (s *UpdateCoreParamsSuite) TestValidParams() {
	err := ValidateCoreParams(s.evmParams)
	require.Nil(s.T(), err)
	err = ValidateCoreParams(s.btcParams)
	require.Nil(s.T(), err)
}

func (s *UpdateCoreParamsSuite) TestCommonParams() {
	s.Validate(s.evmParams)
	s.Validate(s.btcParams)
}

func (s *UpdateCoreParamsSuite) TestBTCParams() {
	copy := *s.btcParams
	copy.WatchUtxoTicker = 0
	err := ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)
}

func (s *UpdateCoreParamsSuite) TestCoreContractAddresses() {
	copy := *s.evmParams
	copy.ZetaTokenContractAddress = "0x123"
	err := ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)

	copy = *s.evmParams
	copy.ZetaTokenContractAddress = "733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9"
	err = ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)

	copy = *s.evmParams
	copy.ConnectorContractAddress = "0x123"
	err = ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)

	copy = *s.evmParams
	copy.ConnectorContractAddress = "733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9"
	err = ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)

	copy = *s.evmParams
	copy.Erc20CustodyContractAddress = "0x123"
	err = ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)

	copy = *s.evmParams
	copy.Erc20CustodyContractAddress = "733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9"
	err = ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)
}

func (s *UpdateCoreParamsSuite) Validate(params *CoreParams) {
	copy := *params
	copy.ConfirmationCount = 0
	err := ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)

	copy = *params
	copy.GasPriceTicker = 0
	err = ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)
	copy.GasPriceTicker = 300
	err = ValidateCoreParams(&copy)
	require.Nil(s.T(), err)
	copy.GasPriceTicker = 301
	err = ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)

	copy = *params
	copy.InTxTicker = 0
	err = ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)
	copy.InTxTicker = 300
	err = ValidateCoreParams(&copy)
	require.Nil(s.T(), err)
	copy.InTxTicker = 301
	err = ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)

	copy = *params
	copy.OutTxTicker = 0
	err = ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)
	copy.OutTxTicker = 300
	err = ValidateCoreParams(&copy)
	require.Nil(s.T(), err)
	copy.OutTxTicker = 301
	err = ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)

	copy = *params
	copy.OutboundTxScheduleInterval = 0
	err = ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)
	copy.OutboundTxScheduleInterval = 100
	err = ValidateCoreParams(&copy)
	require.Nil(s.T(), err)
	copy.OutboundTxScheduleInterval = 101
	err = ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)

	copy = *params
	copy.OutboundTxScheduleLookahead = 0
	err = ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)
	copy.OutboundTxScheduleLookahead = 500
	err = ValidateCoreParams(&copy)
	require.Nil(s.T(), err)
	copy.OutboundTxScheduleLookahead = 501
	err = ValidateCoreParams(&copy)
	require.NotNil(s.T(), err)
}
