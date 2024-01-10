package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/zeta-chain/zetacore/x/observer/types"
	. "gopkg.in/check.v1"
)

func TestChainParamsList_Validate(t *testing.T) {
	t.Run("should return no error for default list", func(t *testing.T) {
		list := types.GetDefaultChainParams()
		err := list.Validate()
		require.NoError(t, err)
	})

	t.Run("should return error for invalid chain id", func(t *testing.T) {
		list := types.GetDefaultChainParams()
		list.ChainParams[0].ChainId = 999
		err := list.Validate()
		require.Error(t, err)
	})

	t.Run("should return error for duplicated chain ID", func(t *testing.T) {
		list := types.GetDefaultChainParams()
		list.ChainParams = append(list.ChainParams, list.ChainParams[0])
		err := list.Validate()
		require.Error(t, err)
		require.Contains(t, err.Error(), "duplicated chain id")
	})
}

type UpdateChainParamsSuite struct {
	suite.Suite
	evmParams *types.ChainParams
	btcParams *types.ChainParams
}

var _ = Suite(&UpdateChainParamsSuite{})

func TestUpdateChainParamsSuiteSuite(t *testing.T) {
	suite.Run(t, new(UpdateChainParamsSuite))
}

func (s *UpdateChainParamsSuite) SetupTest() {
	s.evmParams = &types.ChainParams{
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
		BallotThreshold:             types.DefaultBallotThreshold,
		MinObserverDelegation:       types.DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
	s.btcParams = &types.ChainParams{
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
		BallotThreshold:             types.DefaultBallotThreshold,
		MinObserverDelegation:       types.DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
}

func (s *UpdateChainParamsSuite) TestValidParams() {
	err := types.ValidateChainParams(s.evmParams)
	require.Nil(s.T(), err)
	err = types.ValidateChainParams(s.btcParams)
	require.Nil(s.T(), err)
}

func (s *UpdateChainParamsSuite) TestCommonParams() {
	s.Validate(s.evmParams)
	s.Validate(s.btcParams)
}

func (s *UpdateChainParamsSuite) TestBTCParams() {
	copy := *s.btcParams
	copy.WatchUtxoTicker = 0
	err := types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)
}

func (s *UpdateChainParamsSuite) TestCoreContractAddresses() {
	copy := *s.evmParams
	copy.ZetaTokenContractAddress = "0x123"
	err := types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)

	copy = *s.evmParams
	copy.ZetaTokenContractAddress = "733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9"
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)

	copy = *s.evmParams
	copy.ConnectorContractAddress = "0x123"
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)

	copy = *s.evmParams
	copy.ConnectorContractAddress = "733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9"
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)

	copy = *s.evmParams
	copy.Erc20CustodyContractAddress = "0x123"
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)

	copy = *s.evmParams
	copy.Erc20CustodyContractAddress = "733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9"
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)
}

func (s *UpdateChainParamsSuite) Validate(params *types.ChainParams) {
	copy := *params
	copy.ConfirmationCount = 0
	err := types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)

	copy = *params
	copy.GasPriceTicker = 0
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)
	copy.GasPriceTicker = 300
	err = types.ValidateChainParams(&copy)
	require.Nil(s.T(), err)
	copy.GasPriceTicker = 301
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)

	copy = *params
	copy.InTxTicker = 0
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)
	copy.InTxTicker = 300
	err = types.ValidateChainParams(&copy)
	require.Nil(s.T(), err)
	copy.InTxTicker = 301
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)

	copy = *params
	copy.OutTxTicker = 0
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)
	copy.OutTxTicker = 300
	err = types.ValidateChainParams(&copy)
	require.Nil(s.T(), err)
	copy.OutTxTicker = 301
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)

	copy = *params
	copy.OutboundTxScheduleInterval = 0
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)
	copy.OutboundTxScheduleInterval = 100
	err = types.ValidateChainParams(&copy)
	require.Nil(s.T(), err)
	copy.OutboundTxScheduleInterval = 101
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)

	copy = *params
	copy.OutboundTxScheduleLookahead = 0
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)
	copy.OutboundTxScheduleLookahead = 500
	err = types.ValidateChainParams(&copy)
	require.Nil(s.T(), err)
	copy.OutboundTxScheduleLookahead = 501
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)
}
