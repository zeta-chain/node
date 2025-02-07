package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	. "gopkg.in/check.v1"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
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
		require.ErrorContains(t, err, "duplicated chain id")
	})
}

type UpdateChainParamsSuite struct {
	suite.Suite
	zetaParams *types.ChainParams
	evmParams  *types.ChainParams
	btcParams  *types.ChainParams
}

var _ = Suite(&UpdateChainParamsSuite{})

func TestUpdateChainParamsSuiteSuite(t *testing.T) {
	suite.Run(t, new(UpdateChainParamsSuite))
}

func TestChainParamsEqual(t *testing.T) {
	params := sample.ChainParams(1)

	require.True(t, types.ChainParamsEqual(*params, *params))

	// ChainId matters
	copy := copyParams(params)
	copy.ChainId = 2
	require.False(t, types.ChainParamsEqual(*params, *copy))

	// ConfirmationCount matters
	copy = copyParams(params)
	copy.ConfirmationCount = params.ConfirmationCount + 1
	require.False(t, types.ChainParamsEqual(*params, *copy))

	// GasPriceTicker matters
	copy = copyParams(params)
	copy.GasPriceTicker = params.GasPriceTicker + 1
	require.False(t, types.ChainParamsEqual(*params, *copy))

	// InboundTicker matters
	copy = copyParams(params)
	copy.InboundTicker = params.InboundTicker + 1
	require.False(t, types.ChainParamsEqual(*params, *copy))

	// OutboundTicker matters
	copy = copyParams(params)
	copy.OutboundTicker = params.OutboundTicker + 1
	require.False(t, types.ChainParamsEqual(*params, *copy))

	// WatchUtxoTicker matters
	copy = copyParams(params)
	copy.WatchUtxoTicker = params.WatchUtxoTicker + 1
	require.False(t, types.ChainParamsEqual(*params, *copy))

	// ZetaTokenContractAddress matters
	copy = copyParams(params)
	copy.ZetaTokenContractAddress = "0x_something_else"
	require.False(t, types.ChainParamsEqual(*params, *copy))

	// ConnectorContractAddress matters
	copy = copyParams(params)
	copy.ConnectorContractAddress = "0x_something_else"
	require.False(t, types.ChainParamsEqual(*params, *copy))

	// Erc20CustodyContractAddress matters
	copy = copyParams(params)
	copy.Erc20CustodyContractAddress = "0x_something_else"
	require.False(t, types.ChainParamsEqual(*params, *copy))

	// OutboundScheduleInterval matters
	copy = copyParams(params)
	copy.OutboundScheduleInterval = params.OutboundScheduleInterval + 1
	require.False(t, types.ChainParamsEqual(*params, *copy))

	// OutboundScheduleLookahead matters
	copy = copyParams(params)
	copy.OutboundScheduleLookahead = params.OutboundScheduleLookahead + 1
	require.False(t, types.ChainParamsEqual(*params, *copy))

	// BallotThreshold matters
	copy = copyParams(params)
	copy.BallotThreshold = params.BallotThreshold.Add(sdkmath.LegacySmallestDec())
	require.False(t, types.ChainParamsEqual(*params, *copy))

	// MinObserverDelegation matters
	copy = copyParams(params)
	copy.MinObserverDelegation = params.MinObserverDelegation.Add(sdkmath.LegacySmallestDec())
	require.False(t, types.ChainParamsEqual(*params, *copy))

	// IsSupported matters
	copy = copyParams(params)
	copy.IsSupported = !params.IsSupported
	require.False(t, types.ChainParamsEqual(*params, *copy))

	// GatewayAddress matters
	copy = copyParams(params)
	copy.GatewayAddress = "0x_something_else"
	require.False(t, types.ChainParamsEqual(*params, *copy))

	// ConfirmationParams matters
	copy = copyParams(params)
	copy.ConfirmationParams = nil
	require.False(t, types.ChainParamsEqual(*params, *copy))

	copy = copyParams(params)
	copy.ConfirmationParams.SafeInboundCount = params.ConfirmationParams.SafeInboundCount + 1
	require.False(t, types.ChainParamsEqual(*params, *copy))

	copy = copyParams(params)
	copy.ConfirmationParams.FastInboundCount = params.ConfirmationParams.FastInboundCount + 1
	require.False(t, types.ChainParamsEqual(*params, *copy))

	copy = copyParams(params)
	copy.ConfirmationParams.SafeOutboundCount = params.ConfirmationParams.SafeOutboundCount + 1
	require.False(t, types.ChainParamsEqual(*params, *copy))

	copy = copyParams(params)
	copy.ConfirmationParams.FastOutboundCount = params.ConfirmationParams.FastOutboundCount + 1
	require.False(t, types.ChainParamsEqual(*params, *copy))
}

func (s *UpdateChainParamsSuite) SetupTest() {
	s.zetaParams = &types.ChainParams{
		ConfirmationCount:           0,
		GasPriceTicker:              0,
		InboundTicker:               0,
		OutboundTicker:              0,
		WatchUtxoTicker:             0,
		ZetaTokenContractAddress:    "0x0000000000000000000000000000000000000000",
		ConnectorContractAddress:    "0x0000000000000000000000000000000000000000",
		Erc20CustodyContractAddress: "0x0000000000000000000000000000000000000000",
		ChainId:                     7000,
		OutboundScheduleInterval:    0,
		OutboundScheduleLookahead:   0,
		BallotThreshold:             types.DefaultBallotThreshold,
		MinObserverDelegation:       types.DefaultMinObserverDelegation,
		IsSupported:                 true,
		GatewayAddress:              "",
	}
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
		GatewayAddress:              "0xF0deebCB0E9C829519C4baa794c5445171973826",
		ConfirmationParams: &types.ConfirmationParams{
			SafeInboundCount:  2,
			FastInboundCount:  0, // zero means fast observation is disabled
			SafeOutboundCount: 2,
			FastOutboundCount: 0, // zero means fast observation is disabled
		},
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
		ConfirmationParams: &types.ConfirmationParams{
			SafeInboundCount:  2,
			FastInboundCount:  1,
			SafeOutboundCount: 2,
			FastOutboundCount: 1,
		},
	}
}

func (s *UpdateChainParamsSuite) TestValidParams() {
	err := types.ValidateChainParams(s.zetaParams)
	require.Nil(s.T(), err)
	err = types.ValidateChainParams(s.evmParams)
	require.Nil(s.T(), err)
	err = types.ValidateChainParams(s.btcParams)
	require.Nil(s.T(), err)
}

func (s *UpdateChainParamsSuite) TestCommonParams() {
	s.Validate(s.evmParams)
	s.Validate(s.btcParams)
}

func (s *UpdateChainParamsSuite) TestBTCParamsInvalid() {
	copy := *s.btcParams
	copy.WatchUtxoTicker = 301
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

func Test_InboundConfirmationSafe(t *testing.T) {
	cp := sample.ChainParams(1)

	// set and check safe inbound count
	cp.ConfirmationParams.SafeInboundCount = 10
	require.Equal(t, uint64(10), cp.InboundConfirmationSafe())
}

func Test_OutboundConfirmationSafe(t *testing.T) {
	cp := sample.ChainParams(1)

	// set and check safe outbound count
	cp.ConfirmationParams.SafeOutboundCount = 10
	require.Equal(t, uint64(10), cp.OutboundConfirmationSafe())
}

func Test_InboundConfirmationFast(t *testing.T) {
	t.Run("should return fast inbound confirmation count if enabled", func(t *testing.T) {
		cp := sample.ChainParams(1)
		cp.ConfirmationParams.SafeInboundCount = 2
		cp.ConfirmationParams.FastInboundCount = 1
		confirmation := cp.InboundConfirmationFast()
		require.Equal(t, uint64(1), confirmation)
	})

	t.Run("should fallback to safe inbound confirmation count if fast confirmation is disabled", func(t *testing.T) {
		cp := sample.ChainParams(1)
		cp.ConfirmationParams.SafeInboundCount = 2
		cp.ConfirmationParams.FastInboundCount = 0
		confirmation := cp.InboundConfirmationFast()
		require.Equal(t, uint64(2), confirmation)
	})
}

func Test_OutboundConfirmationFast(t *testing.T) {
	t.Run("should return fast outbound confirmation count if enabled", func(t *testing.T) {
		cp := sample.ChainParams(1)
		cp.ConfirmationParams.SafeOutboundCount = 2
		cp.ConfirmationParams.FastOutboundCount = 1
		confirmation := cp.OutboundConfirmationFast()
		require.Equal(t, uint64(1), confirmation)
	})

	t.Run("should fallback to safe outbound confirmation count if fast confirmation is disabled", func(t *testing.T) {
		cp := sample.ChainParams(1)
		cp.ConfirmationParams.SafeOutboundCount = 2
		cp.ConfirmationParams.FastOutboundCount = 0
		confirmation := cp.OutboundConfirmationFast()
		require.Equal(t, uint64(2), confirmation)
	})
}

func (s *UpdateChainParamsSuite) Validate(params *types.ChainParams) {
	copy := copyParams(params)
	copy.ConfirmationCount = 0
	err := types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)

	copy = copyParams(params)
	copy.ConfirmationParams = nil
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)

	copy = copyParams(params)
	copy.ConfirmationParams.SafeInboundCount = 0
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)
	copy = copyParams(params)
	copy.ConfirmationParams.FastInboundCount = copy.ConfirmationParams.SafeInboundCount + 1
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)

	copy = copyParams(params)
	copy.ConfirmationParams.SafeOutboundCount = 0
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)
	copy = copyParams(params)
	copy.ConfirmationParams.FastOutboundCount = copy.ConfirmationParams.SafeOutboundCount + 1
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)

	copy = copyParams(params)
	copy.GasPriceTicker = 0
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)
	copy.GasPriceTicker = 300
	err = types.ValidateChainParams(copy)
	require.Nil(s.T(), err)
	copy.GasPriceTicker = 301
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)

	copy = copyParams(params)
	copy.InboundTicker = 0
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)
	copy.InboundTicker = 300
	err = types.ValidateChainParams(copy)
	require.Nil(s.T(), err)
	copy.InboundTicker = 301
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)

	copy = copyParams(params)
	copy.OutboundTicker = 0
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)
	copy.OutboundTicker = 300
	err = types.ValidateChainParams(copy)
	require.Nil(s.T(), err)
	copy.OutboundTicker = 301
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)

	copy = copyParams(params)
	copy.OutboundScheduleInterval = 0
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)
	copy.OutboundScheduleInterval = 100
	err = types.ValidateChainParams(copy)
	require.Nil(s.T(), err)
	copy.OutboundScheduleInterval = 101
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)

	copy = copyParams(params)
	copy.OutboundScheduleLookahead = 0
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)
	copy.OutboundScheduleLookahead = 500
	err = types.ValidateChainParams(copy)
	require.Nil(s.T(), err)
	copy.OutboundScheduleLookahead = 501
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)

	copy = copyParams(params)
	copy.BallotThreshold = sdkmath.LegacyDec{}
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)
	copy.BallotThreshold = sdkmath.LegacyMustNewDecFromStr("1.2")
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)
	copy.BallotThreshold = sdkmath.LegacyMustNewDecFromStr("0.9")
	err = types.ValidateChainParams(copy)
	require.Nil(s.T(), err)

	copy = copyParams(params)
	copy.MinObserverDelegation = sdkmath.LegacyDec{}
	err = types.ValidateChainParams(copy)
	require.NotNil(s.T(), err)
	copy.MinObserverDelegation = sdkmath.LegacyMustNewDecFromStr("0.9")
	err = types.ValidateChainParams(copy)
	require.Nil(s.T(), err)
}

// copyParams creates a deep copy of the given ChainParams.
func copyParams(src *types.ChainParams) *types.ChainParams {
	if src == nil {
		return nil
	}

	copy := *src
	if src.ConfirmationParams == nil {
		return &copy
	}

	copy.ConfirmationParams = &types.ConfirmationParams{
		SafeInboundCount:  src.ConfirmationParams.SafeInboundCount,
		FastInboundCount:  src.ConfirmationParams.FastInboundCount,
		SafeOutboundCount: src.ConfirmationParams.SafeOutboundCount,
		FastOutboundCount: src.ConfirmationParams.FastOutboundCount,
	}

	return &copy
}
