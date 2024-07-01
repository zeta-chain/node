package orchestrator_test

import (
	"encoding/json"
	"path"
	"testing"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/orchestrator"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
)

// the relative path to the testdata directory
var TestDataDir = "../"

func Test_UpdateAppContext(t *testing.T) {
	// define test chains and chain params
	evmChain := chains.Ethereum
	btcChain := chains.BitcoinMainnet
	evmChainParams := sample.ChainParams(evmChain.ChainId)

	// define test config
	evmCfg := config.EVMConfig{
		Chain:    evmChain,
		Endpoint: "http://localhost:8545",
	}
	btcCfg := config.BTCConfig{
		RPCUsername: "user",
	}

	// create app context
	appCtx := createTestAppContext(evmCfg, btcCfg, evmChain, btcChain, evmChainParams, nil)

	t.Run("should update app context", func(t *testing.T) {
		// create orchestrator
		ztacoreClient := mocks.NewMockZetacoreClient()
		oc := orchestrator.NewOrchestrator(appCtx, ztacoreClient, nil, base.Logger{}, testutils.SQLiteMemory, nil)

		// zetacoreHome directory
		zetacoreHome := path.Join(TestDataDir, "testdata")

		// read archived config file and setup test zetacore home
		cfg := testutils.LoadZetaclientConfig(t, TestDataDir)
		cfg.ZetaCoreHome = zetacoreHome

		// set absolute path is needed to make the test pass
		// config reloading will overwrite the relative path in config file
		cfg.TssPath = config.GetPath(cfg.TssPath)

		// set zetacore home directory to the archived testdata directory
		appCtx.Config().ZetaCoreHome = zetacoreHome

		// update app context
		err := oc.UpdateAppContext()
		require.NoError(t, err)

		// get new config
		newCfg := appCtx.Config()

		// serialize old and new config
		oldCfgData, err := json.Marshal(cfg)
		require.NoError(t, err)

		newCfgData, err := json.Marshal(newCfg)
		require.NoError(t, err)

		// compare old and new config
		require.JSONEq(t, string(oldCfgData), string(newCfgData))
	})
	t.Run("should return error if zetacore client fails to update app context", func(t *testing.T) {
		// create orchestrator
		ztacoreClient := mocks.NewMockZetacoreClient()
		oc := orchestrator.NewOrchestrator(appCtx, ztacoreClient, nil, base.Logger{}, testutils.SQLiteMemory, nil)

		// pause zetacore client to simulate error
		ztacoreClient.Pause()

		// update app context
		err := oc.UpdateAppContext()
		require.ErrorContains(t, err, "error updating app context")
	})
	t.Run("should return error if reading config file fails", func(t *testing.T) {
		// create orchestrator
		ztacoreClient := mocks.NewMockZetacoreClient()
		oc := orchestrator.NewOrchestrator(appCtx, ztacoreClient, nil, base.Logger{}, testutils.SQLiteMemory, nil)

		// set invalid zetacore home directory
		appCtx.Config().ZetaCoreHome = "/invalid/path"

		// update app context
		err := oc.UpdateAppContext()
		require.ErrorContains(t, err, "error loading config from path")
	})
}

func Test_UpgradeHeightReached(t *testing.T) {
	t.Run("should return true if upgrade height is reached", func(t *testing.T) {
		// create orchestrator
		zetacoreClient := mocks.NewMockZetacoreClient()
		oc := orchestrator.NewOrchestrator(nil, zetacoreClient, nil, base.Logger{}, testutils.SQLiteMemory, nil)

		// set upgrade plan and current height
		zetacoreClient.WithBlockHeight(99)
		zetacoreClient.WithUpgradedPlan(&upgradetypes.Plan{
			Height: 100,
		})

		// check if upgrade height is reached
		reached, err := oc.UpgradeHeightReached()
		require.NoError(t, err)
		require.True(t, reached)
	})
	t.Run("should return error if failed to get upgrade plan", func(t *testing.T) {
		// create orchestrator
		zetacoreClient := mocks.NewMockZetacoreClient()
		oc := orchestrator.NewOrchestrator(nil, zetacoreClient, nil, base.Logger{}, testutils.SQLiteMemory, nil)

		// pause zetacore client to simulate error
		zetacoreClient.Pause()

		// check if upgrade height is reached
		reached, err := oc.UpgradeHeightReached()
		require.ErrorContains(t, err, "failed to get upgrade plan")
		require.False(t, reached)
	})
	t.Run("should return false if there is no active upgrade plan", func(t *testing.T) {
		// create orchestrator
		zetacoreClient := mocks.NewMockZetacoreClient()
		oc := orchestrator.NewOrchestrator(nil, zetacoreClient, nil, base.Logger{}, testutils.SQLiteMemory, nil)

		// check if upgrade height is reached
		reached, err := oc.UpgradeHeightReached()
		require.NoError(t, err)
		require.False(t, reached)
	})
	t.Run("should return false if upgrade height is not reached", func(t *testing.T) {
		// create orchestrator
		zetacoreClient := mocks.NewMockZetacoreClient()
		oc := orchestrator.NewOrchestrator(nil, zetacoreClient, nil, base.Logger{}, testutils.SQLiteMemory, nil)

		// set upgrade plan and current height
		zetacoreClient.WithBlockHeight(98)
		zetacoreClient.WithUpgradedPlan(&upgradetypes.Plan{
			Height: 100,
		})

		// check if upgrade height is reached
		reached, err := oc.UpgradeHeightReached()
		require.NoError(t, err)
		require.False(t, reached)
	})
}
