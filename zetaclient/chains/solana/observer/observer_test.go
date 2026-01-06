package observer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/mode"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/solana/observer"
	"github.com/zeta-chain/node/zetaclient/chains/tssrepo"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

// MockSolanaObserver creates a mock Solana observer with custom chain, TSS, params etc
func MockSolanaObserver(
	t *testing.T,
	chain chains.Chain,
	solanaClient observer.SolanaClient,
	chainParams observertypes.ChainParams,
	zetacoreClient zrepo.ZetacoreClient,
	tssSigner tssrepo.TSSClient,
) *observer.Observer {
	// use mock zetacore client if not provided
	if zetacoreClient == nil {
		zetacoreClient = mocks.NewZetacoreClient(t).WithKeys(&keys.Keys{})
	}

	// use mock tss if not provided
	if tssSigner == nil {
		tssSigner = mocks.NewTSS(t)
	}

	database, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	baseObserver, err := base.NewObserver(
		chain,
		chainParams,
		zrepo.New(zetacoreClient, chain, mode.StandardMode),
		tssSigner,
		1000,
		nil,
		database,
		base.DefaultLogger(),
	)
	require.NoError(t, err)

	ob, err := observer.New(baseObserver, solanaClient, chainParams.GatewayAddress)
	require.NoError(t, err)

	return ob
}

func Test_LoadLastTxScanned(t *testing.T) {
	// prepare params
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
