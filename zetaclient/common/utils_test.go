package common

import (
	"path"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func TestCctxRestricted(t *testing.T) {
	// load archived cctx
	var cctx crosschaintypes.CrossChainTx
	err := testutils.LoadObjectFromJSONFile(&cctx, path.Join("../", testutils.TestDataPathCctx, "cctx_1_6270.json"))
	require.NoError(t, err)

	// create config
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	t.Run("should return true if sender is restricted", func(t *testing.T) {
		cfg.ComplianceConfig.RestrictedAddresses = []string{cctx.InboundTxParams.Sender}
		config.LoadComplianceConfig(cfg)
		require.True(t, IsCctxRestricted(&cctx))
	})
	t.Run("should return true if receiver is restricted", func(t *testing.T) {
		cfg.ComplianceConfig.RestrictedAddresses = []string{cctx.GetCurrentOutTxParam().Receiver}
		config.LoadComplianceConfig(cfg)
		require.True(t, IsCctxRestricted(&cctx))
	})
	t.Run("should return false if sender and receiver are not restricted", func(t *testing.T) {
		// restrict other address
		cfg.ComplianceConfig.RestrictedAddresses = []string{"0x27104b8dB4aEdDb054fCed87c346C0758Ff5dFB1"}
		config.LoadComplianceConfig(cfg)
		require.False(t, IsCctxRestricted(&cctx))
	})
	t.Run("should be able to restrict coinbase address", func(t *testing.T) {
		cfg.ComplianceConfig.RestrictedAddresses = []string{ethcommon.Address{}.String()}
		config.LoadComplianceConfig(cfg)
		cctx.InboundTxParams.Sender = ethcommon.Address{}.String()
		require.True(t, IsCctxRestricted(&cctx))
	})
	t.Run("should ignore empty address", func(t *testing.T) {
		cfg.ComplianceConfig.RestrictedAddresses = []string{""}
		config.LoadComplianceConfig(cfg)
		cctx.InboundTxParams.Sender = ""
		require.False(t, IsCctxRestricted(&cctx))
	})
}

func Test_GasPriceMultiplier(t *testing.T) {
	tt := []struct {
		name       string
		chainID    int64
		multiplier float64
		fail       bool
	}{
		{
			name:       "get Ethereum multiplier",
			chainID:    1,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get Goerli multiplier",
			chainID:    5,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get BSC multiplier",
			chainID:    56,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get BSC Testnet multiplier",
			chainID:    97,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get Polygon multiplier",
			chainID:    137,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get Mumbai Testnet multiplier",
			chainID:    80001,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get Bitcoin multiplier",
			chainID:    8332,
			multiplier: 2.0,
			fail:       false,
		},
		{
			name:       "get Bitcoin Testnet multiplier",
			chainID:    18332,
			multiplier: 2.0,
			fail:       false,
		},
		{
			name:       "get unknown chain gas price multiplier",
			chainID:    1234,
			multiplier: 1.0,
			fail:       true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			multiplier, err := GasPriceMultiplier(tc.chainID)
			if tc.fail {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.multiplier, multiplier)
		})
	}
}
