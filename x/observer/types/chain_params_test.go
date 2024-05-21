package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

func TestChainParamsEqual(t *testing.T) {
	params := types.GetDefaultChainParams()
	require.True(t, types.ChainParamsEqual(*params.ChainParams[0], *params.ChainParams[0]))
	require.False(t, types.ChainParamsEqual(*params.ChainParams[0], *params.ChainParams[1]))
}

func (s *UpdateChainParamsSuite) SetupTest() {
	s.evmParams = &types.ChainParams{
		ConfirmationCount:           1,
		GasPriceTicker:              1,
		InboundTicker:               1,
		OutboundTicker:              1,
		WatchUtxoTicker:             0,
		ZetaTokenContractAddress:    "0xA8D5060feb6B456e886F023709A2795373691E63",
		ConnectorContractAddress:    "0x733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9",
		Erc20CustodyContractAddress: "0xD28D6A0b8189305551a0A8bd247a6ECa9CE781Ca",
		ChainId:                     5,
		OutboundScheduleInterval:    1,
		OutboundScheduleLookahead:   1,
		BallotThreshold:             types.DefaultBallotThreshold,
		MinObserverDelegation:       types.DefaultMinObserverDelegation,
		IsSupported:                 false,
	}
	s.btcParams = &types.ChainParams{
		ConfirmationCount:           1,
		GasPriceTicker:              1,
		InboundTicker:               1,
		OutboundTicker:              1,
		WatchUtxoTicker:             1,
		ZetaTokenContractAddress:    "",
		ConnectorContractAddress:    "",
		Erc20CustodyContractAddress: "",
		ChainId:                     18332,
		OutboundScheduleInterval:    1,
		OutboundScheduleLookahead:   1,
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
	copy.InboundTicker = 0
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)
	copy.InboundTicker = 300
	err = types.ValidateChainParams(&copy)
	require.Nil(s.T(), err)
	copy.InboundTicker = 301
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)

	copy = *params
	copy.OutboundTicker = 0
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)
	copy.OutboundTicker = 300
	err = types.ValidateChainParams(&copy)
	require.Nil(s.T(), err)
	copy.OutboundTicker = 301
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)

	copy = *params
	copy.OutboundScheduleInterval = 0
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)
	copy.OutboundScheduleInterval = 100
	err = types.ValidateChainParams(&copy)
	require.Nil(s.T(), err)
	copy.OutboundScheduleInterval = 101
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)

	copy = *params
	copy.OutboundScheduleLookahead = 0
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)
	copy.OutboundScheduleLookahead = 500
	err = types.ValidateChainParams(&copy)
	require.Nil(s.T(), err)
	copy.OutboundScheduleLookahead = 501
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)

	copy = *params
	copy.BallotThreshold = sdk.Dec{}
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)
	copy.BallotThreshold = sdk.MustNewDecFromStr("1.2")
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)
	copy.BallotThreshold = sdk.MustNewDecFromStr("0.9")
	err = types.ValidateChainParams(&copy)
	require.Nil(s.T(), err)

	copy = *params
	copy.MinObserverDelegation = sdk.Dec{}
	err = types.ValidateChainParams(&copy)
	require.NotNil(s.T(), err)
	copy.MinObserverDelegation = sdk.MustNewDecFromStr("0.9")
	err = types.ValidateChainParams(&copy)
	require.Nil(s.T(), err)
}
