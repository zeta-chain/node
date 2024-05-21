package compliance

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func TestCctxRestricted(t *testing.T) {
	// load archived cctx
	chain := chains.EthChain
	cctx := testutils.LoadCctxByNonce(t, chain.ChainId, 6270)

	// create config
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	t.Run("should return true if sender is restricted", func(t *testing.T) {
		cfg.ComplianceConfig.RestrictedAddresses = []string{cctx.InboundParams.Sender}
		config.LoadComplianceConfig(cfg)
		require.True(t, IsCctxRestricted(cctx))
	})
	t.Run("should return true if receiver is restricted", func(t *testing.T) {
		cfg.ComplianceConfig.RestrictedAddresses = []string{cctx.GetCurrentOutboundParam().Receiver}
		config.LoadComplianceConfig(cfg)
		require.True(t, IsCctxRestricted(cctx))
	})
	t.Run("should return false if sender and receiver are not restricted", func(t *testing.T) {
		// restrict other address
		cfg.ComplianceConfig.RestrictedAddresses = []string{"0x27104b8dB4aEdDb054fCed87c346C0758Ff5dFB1"}
		config.LoadComplianceConfig(cfg)
		require.False(t, IsCctxRestricted(cctx))
	})
	t.Run("should be able to restrict coinbase address", func(t *testing.T) {
		cfg.ComplianceConfig.RestrictedAddresses = []string{ethcommon.Address{}.String()}
		config.LoadComplianceConfig(cfg)
		cctx.InboundParams.Sender = ethcommon.Address{}.String()
		require.True(t, IsCctxRestricted(cctx))
	})
	t.Run("should ignore empty address", func(t *testing.T) {
		cfg.ComplianceConfig.RestrictedAddresses = []string{""}
		config.LoadComplianceConfig(cfg)
		cctx.InboundParams.Sender = ""
		require.False(t, IsCctxRestricted(cctx))
	})
}
