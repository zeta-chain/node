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

func TestCctxBanned(t *testing.T) {
	// load archived cctx
	var cctx crosschaintypes.CrossChainTx
	err := testutils.LoadObjectFromJSONFile(&cctx, path.Join("../", testutils.TestDataPathCctx, "cctx_1_6270.json"))
	require.NoError(t, err)

	// create config
	cfg := &config.Config{
		ComplianceConfig: &config.ComplianceConfig{},
	}

	t.Run("should return true if sender is banned", func(t *testing.T) {
		cfg.ComplianceConfig.BannedAddresses = []string{cctx.InboundTxParams.Sender}
		config.LoadComplianceConfig(cfg)
		require.True(t, IsCctxBanned(&cctx))
	})
	t.Run("should return true if receiver is banned", func(t *testing.T) {
		cfg.ComplianceConfig.BannedAddresses = []string{cctx.GetCurrentOutTxParam().Receiver}
		config.LoadComplianceConfig(cfg)
		require.True(t, IsCctxBanned(&cctx))
	})
	t.Run("should return false if sender and receiver are not banned", func(t *testing.T) {
		// ban other address
		cfg.ComplianceConfig.BannedAddresses = []string{"0x27104b8dB4aEdDb054fCed87c346C0758Ff5dFB1"}
		config.LoadComplianceConfig(cfg)
		require.False(t, IsCctxBanned(&cctx))
	})
	t.Run("should be able to ban coinbase address", func(t *testing.T) {
		cfg.ComplianceConfig.BannedAddresses = []string{ethcommon.Address{}.String()}
		config.LoadComplianceConfig(cfg)
		cctx.InboundTxParams.Sender = ethcommon.Address{}.String()
		require.True(t, IsCctxBanned(&cctx))
	})
}
