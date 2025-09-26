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

func TestConfirmationParams_Validate(t *testing.T) {
	tests := []struct {
		name  string
		cp    types.ConfirmationParams
		isErr bool
	}{
		{
			name: "valid confirmation params",
			cp: types.ConfirmationParams{
				SafeInboundCount:  1,
				FastInboundCount:  1,
				SafeOutboundCount: 1,
				FastOutboundCount: 1,
			},
		},
		{
			name: "invalid SafeInboundCount",
			cp: types.ConfirmationParams{
				SafeInboundCount:  0,
				FastInboundCount:  1,
				SafeOutboundCount: 1,
				FastOutboundCount: 1,
			},
			isErr: true,
		},
		{
			name: "invalid FastInboundCount",
			cp: types.ConfirmationParams{
				SafeInboundCount:  1,
				FastInboundCount:  2,
				SafeOutboundCount: 1,
				FastOutboundCount: 1,
			},
			isErr: true,
		},
		{
			name: "invalid SafeOutboundCount",
			cp: types.ConfirmationParams{
				SafeInboundCount:  1,
				FastInboundCount:  1,
				SafeOutboundCount: 0,
				FastOutboundCount: 1,
			},
			isErr: true,
		},
		{
			name: "invalid FastOutboundCount",
			cp: types.ConfirmationParams{
				SafeInboundCount:  1,
				FastInboundCount:  1,
				SafeOutboundCount: 1,
				FastOutboundCount: 2,
			},
			isErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cp.Validate()
			if tt.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

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

	t.Run("should return error if stability pool percentage is greater than 100", func(t *testing.T) {
		list := types.GetDefaultChainParams()
		list.ChainParams[0].StabilityPoolPercentage = 101
		err := list.Validate()
		require.ErrorIs(t, err, types.ErrParamsStabilityPoolPercentage)
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
	cp := copyParams(params)
	cp.ChainId = 2
	require.False(t, types.ChainParamsEqual(*params, *cp))

	// GasPriceTicker matters
	cp = copyParams(params)
	cp.GasPriceTicker = params.GasPriceTicker + 1
	require.False(t, types.ChainParamsEqual(*params, *cp))

	// InboundTicker matters
	cp = copyParams(params)
	cp.InboundTicker = params.InboundTicker + 1
	require.False(t, types.ChainParamsEqual(*params, *cp))

	// OutboundTicker matters
	cp = copyParams(params)
	cp.OutboundTicker = params.OutboundTicker + 1
	require.False(t, types.ChainParamsEqual(*params, *cp))

	// WatchUtxoTicker matters
	cp = copyParams(params)
	cp.WatchUtxoTicker = params.WatchUtxoTicker + 1
	require.False(t, types.ChainParamsEqual(*params, *cp))

	// ZetaTokenContractAddress matters
	cp = copyParams(params)
	cp.ZetaTokenContractAddress = "0x_something_else"
	require.False(t, types.ChainParamsEqual(*params, *cp))

	// ConnectorContractAddress matters
	cp = copyParams(params)
	cp.ConnectorContractAddress = "0x_something_else"
	require.False(t, types.ChainParamsEqual(*params, *cp))

	// Erc20CustodyContractAddress matters
	cp = copyParams(params)
	cp.Erc20CustodyContractAddress = "0x_something_else"
	require.False(t, types.ChainParamsEqual(*params, *cp))

	// OutboundScheduleInterval matters
	cp = copyParams(params)
	cp.OutboundScheduleInterval = params.OutboundScheduleInterval + 1
	require.False(t, types.ChainParamsEqual(*params, *cp))

	// OutboundScheduleLookahead matters
	cp = copyParams(params)
	cp.OutboundScheduleLookahead = params.OutboundScheduleLookahead + 1
	require.False(t, types.ChainParamsEqual(*params, *cp))

	// BallotThreshold matters
	cp = copyParams(params)
	cp.BallotThreshold = params.BallotThreshold.Add(sdkmath.LegacySmallestDec())
	require.False(t, types.ChainParamsEqual(*params, *cp))

	// MinObserverDelegation matters
	cp = copyParams(params)
	cp.MinObserverDelegation = params.MinObserverDelegation.Add(sdkmath.LegacySmallestDec())
	require.False(t, types.ChainParamsEqual(*params, *cp))

	// IsSupported matters
	cp = copyParams(params)
	cp.IsSupported = !params.IsSupported
	require.False(t, types.ChainParamsEqual(*params, *cp))

	// GatewayAddress matters
	cp = copyParams(params)
	cp.GatewayAddress = "0x_something_else"
	require.False(t, types.ChainParamsEqual(*params, *cp))

	// ConfirmationParams matters
	cp = copyParams(params)
	cp.ConfirmationParams = nil
	require.False(t, types.ChainParamsEqual(*params, *cp))

	cp = copyParams(params)
	cp.ConfirmationParams.SafeInboundCount = params.ConfirmationParams.SafeInboundCount + 1
	require.False(t, types.ChainParamsEqual(*params, *cp))

	cp = copyParams(params)
	cp.ConfirmationParams.FastInboundCount = params.ConfirmationParams.FastInboundCount + 1
	require.False(t, types.ChainParamsEqual(*params, *cp))

	cp = copyParams(params)
	cp.ConfirmationParams.SafeOutboundCount = params.ConfirmationParams.SafeOutboundCount + 1
	require.False(t, types.ChainParamsEqual(*params, *cp))

	cp = copyParams(params)
	cp.ConfirmationParams.FastOutboundCount = params.ConfirmationParams.FastOutboundCount + 1
	require.False(t, types.ChainParamsEqual(*params, *cp))

	// DisableTSSBlockScan matters
	cp = copyParams(params)
	cp.DisableTssBlockScan = !params.DisableTssBlockScan
	require.False(t, types.ChainParamsEqual(*params, *cp))
}

func (s *UpdateChainParamsSuite) SetupTest() {
	s.zetaParams = &types.ChainParams{
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
	require.NoError(s.T(), s.zetaParams.Validate())
	require.NoError(s.T(), s.evmParams.Validate())
	require.NoError(s.T(), s.btcParams.Validate())
}

func (s *UpdateChainParamsSuite) TestCommonParams() {
	s.Validate(s.evmParams)
	s.Validate(s.btcParams)
}

func (s *UpdateChainParamsSuite) TestBTCParamsInvalid() {
	cp := *s.btcParams
	cp.WatchUtxoTicker = 301
	require.Error(s.T(), cp.Validate())
}

func (s *UpdateChainParamsSuite) TestCoreContractAddresses() {
	cp := *s.evmParams
	cp.ZetaTokenContractAddress = "0x123"
	require.Error(s.T(), cp.Validate())

	cp = *s.evmParams
	cp.ZetaTokenContractAddress = "733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9"
	require.Error(s.T(), cp.Validate())

	cp = *s.evmParams
	cp.ConnectorContractAddress = "0x123"
	require.Error(s.T(), cp.Validate())

	cp = *s.evmParams
	cp.ConnectorContractAddress = "733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9"
	require.Error(s.T(), cp.Validate())

	cp = *s.evmParams
	cp.Erc20CustodyContractAddress = "0x123"
	require.Error(s.T(), cp.Validate())

	cp = *s.evmParams
	cp.Erc20CustodyContractAddress = "733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9"
	require.Error(s.T(), cp.Validate())
}

func Test_IsInboundFastConfirmationEnabled(t *testing.T) {
	cp := sample.ChainParams(1)

	// fast confirmation is enabled
	cp.ConfirmationParams.SafeInboundCount = 2
	cp.ConfirmationParams.FastInboundCount = 1
	require.True(t, cp.IsInboundFastConfirmationEnabled())

	// fast confirmation is disabled if fast count = 0
	cp.ConfirmationParams.FastInboundCount = 0
	require.False(t, cp.IsInboundFastConfirmationEnabled())

	// fast confirmation is disabled if fast count == safe count
	cp.ConfirmationParams.FastInboundCount = 2
	require.False(t, cp.IsInboundFastConfirmationEnabled())
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
	cp := copyParams(params)
	cp.ConfirmationParams = nil
	require.Error(s.T(), cp.Validate())

	cp = copyParams(params)
	cp.ConfirmationParams.SafeInboundCount = 0
	require.Error(s.T(), cp.Validate())
	cp = copyParams(params)
	cp.ConfirmationParams.FastInboundCount = cp.ConfirmationParams.SafeInboundCount + 1
	require.Error(s.T(), cp.Validate())

	cp = copyParams(params)
	cp.ConfirmationParams.SafeOutboundCount = 0
	require.Error(s.T(), cp.Validate())
	cp = copyParams(params)
	cp.ConfirmationParams.FastOutboundCount = cp.ConfirmationParams.SafeOutboundCount + 1
	require.Error(s.T(), cp.Validate())

	cp = copyParams(params)
	cp.GasPriceTicker = 0
	require.Error(s.T(), cp.Validate())
	cp.GasPriceTicker = 300
	require.NoError(s.T(), cp.Validate())
	cp.GasPriceTicker = 301
	require.Error(s.T(), cp.Validate())

	cp = copyParams(params)
	cp.InboundTicker = 0
	require.Error(s.T(), cp.Validate())
	cp.InboundTicker = 300
	require.NoError(s.T(), cp.Validate())
	cp.InboundTicker = 301
	require.Error(s.T(), cp.Validate())

	cp = copyParams(params)
	cp.OutboundTicker = 0
	require.Error(s.T(), cp.Validate())
	cp.OutboundTicker = 300
	require.NoError(s.T(), cp.Validate())
	cp.OutboundTicker = 301
	require.Error(s.T(), cp.Validate())

	cp = copyParams(params)
	cp.OutboundScheduleInterval = 0
	require.Error(s.T(), cp.Validate())
	cp.OutboundScheduleInterval = 100
	require.NoError(s.T(), cp.Validate())
	cp.OutboundScheduleInterval = 101
	require.Error(s.T(), cp.Validate())

	cp = copyParams(params)
	cp.OutboundScheduleLookahead = 0
	require.Error(s.T(), cp.Validate())
	cp.OutboundScheduleLookahead = 500
	require.NoError(s.T(), cp.Validate())
	cp.OutboundScheduleLookahead = 501
	require.Error(s.T(), cp.Validate())

	cp = copyParams(params)
	cp.BallotThreshold = sdkmath.LegacyDec{}
	require.Error(s.T(), cp.Validate())
	cp.BallotThreshold = sdkmath.LegacyMustNewDecFromStr("1.2")
	require.Error(s.T(), cp.Validate())
	cp.BallotThreshold = sdkmath.LegacyMustNewDecFromStr("0.9")
	require.NoError(s.T(), cp.Validate())

	cp = copyParams(params)
	cp.MinObserverDelegation = sdkmath.LegacyDec{}
	require.Error(s.T(), cp.Validate())
	cp.MinObserverDelegation = sdkmath.LegacyMustNewDecFromStr("0.9")
	require.NoError(s.T(), cp.Validate())
}

// copyParams creates a deep copy of the given ChainParams.
func copyParams(src *types.ChainParams) *types.ChainParams {
	if src == nil {
		return nil
	}

	cp := *src
	if src.ConfirmationParams == nil {
		return &cp
	}

	cp.ConfirmationParams = &types.ConfirmationParams{
		SafeInboundCount:  src.ConfirmationParams.SafeInboundCount,
		FastInboundCount:  src.ConfirmationParams.FastInboundCount,
		SafeOutboundCount: src.ConfirmationParams.SafeOutboundCount,
		FastOutboundCount: src.ConfirmationParams.FastOutboundCount,
	}

	return &cp
}
