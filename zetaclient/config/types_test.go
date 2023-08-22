package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	. "gopkg.in/check.v1"
)

type TypesSuite struct {
	suite.Suite
	evmParams  *observertypes.CoreParams
	btcParams  *observertypes.CoreParams
	zetaParams *observertypes.CoreParams
}

var _ = Suite(&TypesSuite{})

func TestTypesSuite(t *testing.T) {
	suite.Run(t, new(TypesSuite))
}

func (s *TypesSuite) SetupTest() {
	s.evmParams = &observertypes.CoreParams{
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
	s.btcParams = &observertypes.CoreParams{
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
	s.zetaParams = &observertypes.CoreParams{
		ConfirmationCount:           1,
		GasPriceTicker:              0,
		InTxTicker:                  0,
		OutTxTicker:                 0,
		WatchUtxoTicker:             0,
		ZetaTokenContractAddress:    "",
		ConnectorContractAddress:    "",
		Erc20CustodyContractAddress: "",
		ChainId:                     7001,
		OutboundTxScheduleInterval:  0,
		OutboundTxScheduleLookahead: 0,
	}
}

func (s *TypesSuite) TestValidParams() {
	err := ValidateCoreParams(s.evmParams)
	require.Nil(s.T(), err)
	err = ValidateCoreParams(s.btcParams)
	require.Nil(s.T(), err)
	err = ValidateCoreParams(s.zetaParams)
	require.Nil(s.T(), err)
}

func (s *TypesSuite) TestEVMInvalidParams() {
	params := *s.evmParams
	params.ConfirmationCount = 0
	err := ValidateCoreParams(&params)
	require.NotNil(s.T(), err)

	params = *s.evmParams
	params.GasPriceTicker = 0
	err = ValidateCoreParams(&params)
	require.NotNil(s.T(), err)

	params = *s.evmParams
	params.InTxTicker = 0
	err = ValidateCoreParams(&params)
	require.NotNil(s.T(), err)

	params = *s.evmParams
	params.OutTxTicker = 0
	err = ValidateCoreParams(&params)
	require.NotNil(s.T(), err)

	params = *s.evmParams
	params.OutboundTxScheduleInterval = 0
	err = ValidateCoreParams(&params)
	require.NotNil(s.T(), err)

	params = *s.evmParams
	params.OutboundTxScheduleLookahead = 0
	err = ValidateCoreParams(&params)
	require.NotNil(s.T(), err)
}

func (s *TypesSuite) TestBTCInvalidTickers() {
	params := *s.btcParams
	params.GasPriceTicker = 0
	err := ValidateCoreParams(&params)
	require.NotNil(s.T(), err)

	params = *s.btcParams
	params.InTxTicker = 0
	err = ValidateCoreParams(&params)
	require.NotNil(s.T(), err)

	params = *s.btcParams
	params.OutTxTicker = 0
	err = ValidateCoreParams(&params)
	require.NotNil(s.T(), err)

	params = *s.btcParams
	params.WatchUtxoTicker = 0
	err = ValidateCoreParams(&params)
	require.NotNil(s.T(), err)

	params = *s.btcParams
	params.OutboundTxScheduleInterval = 0
	err = ValidateCoreParams(&params)
	require.NotNil(s.T(), err)

	params = *s.btcParams
	params.OutboundTxScheduleLookahead = 0
	err = ValidateCoreParams(&params)
	require.NotNil(s.T(), err)
}

func (s *TypesSuite) TestZetaInvalidParams() {
	params := *s.zetaParams
	params.ConfirmationCount = 0
	err := ValidateCoreParams(&params)
	require.NotNil(s.T(), err)
}

func (s *TypesSuite) TestInvalidCoreContractAddresses() {
	params := *s.evmParams
	params.ZetaTokenContractAddress = "0x123"
	err := ValidateCoreParams(&params)
	require.NotNil(s.T(), err)

	params = *s.evmParams
	params.ZetaTokenContractAddress = "733aB8b06DDDEf27Eaa72294B0d7c9cEF7f12db9"
	err = ValidateCoreParams(&params)
	require.NotNil(s.T(), err)
}
