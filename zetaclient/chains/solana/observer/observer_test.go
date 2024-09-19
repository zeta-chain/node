package observer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/keys"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/chains/solana/observer"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

// MockSolanaObserver creates a mock Solana observer with custom chain, TSS, params etc
func MockSolanaObserver(
	t *testing.T,
	chain chains.Chain,
	solClient interfaces.SolanaRPCClient,
	chainParams observertypes.ChainParams,
	zetacoreClient interfaces.ZetacoreClient,
	tss interfaces.TSSSigner,
) *observer.Observer {
	// use mock zetacore client if not provided
	if zetacoreClient == nil {
		zetacoreClient = mocks.NewZetacoreClient(t).WithKeys(&keys.Keys{})
	}
	// use mock tss if not provided
	if tss == nil {
		tss = mocks.NewTSSMainnet()
	}

	database, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	// create observer
	ob, err := observer.NewObserver(
		chain,
		solClient,
		chainParams,
		zetacoreClient,
		tss,
		60,
		database,
		base.DefaultLogger(),
		nil,
	)
	require.NoError(t, err)

	return ob
}

func Test_LoadLastTxScanned(t *testing.T) {
	// parepare params
	chain := chains.SolanaDevnet
	params := sample.ChainParams(chain.ChainId)
	params.GatewayAddress = sample.SolanaAddress(t)

	// create observer
	ob := MockSolanaObserver(t, chain, nil, *params, nil, nil)

	t.Run("should load last block scanned", func(t *testing.T) {
		// write sample last tx to db
		lastTx := sample.SolanaSignature(t).String()
		ob.WriteLastTxScannedToDB(lastTx)

		// load last tx scanned
		err := ob.LoadLastTxScanned()
		require.NoError(t, err)
		require.Equal(t, lastTx, ob.LastTxScanned())
	})
}
