package observer

import (
	"context"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"github.com/zeta-chain/node/zetaclient/testutils/testlog"
)

func TestObserver(t *testing.T) {
	t.Run("PostGasPrice", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		// Given latest checkpoint from RPC
		checkpoint := models.CheckpointResponse{
			// should be used instead of block number
			Epoch:          "333",
			SequenceNumber: "123456",
		}

		ts.suiMock.On("GetLatestCheckpoint", mock.Anything).Return(checkpoint, nil)

		// Given ref price from RPC
		const refGasPrice = uint64(800)
		ts.suiMock.On("SuiXGetReferenceGasPrice", mock.Anything).Return(refGasPrice, nil)

		// Given expected vote for zetacore
		ts.zetaMock.
			On("PostVoteGasPrice", mock.Anything, chains.SuiMainnet, refGasPrice, uint64(0), uint64(333)).
			Return("", nil)

		// ACT
		err := ts.PostGasPrice(ts.ctx)

		// ASSERT
		require.NoError(t, err)
	})
}

type testSuite struct {
	t        *testing.T
	ctx      context.Context
	zetaMock *mocks.ZetacoreClient
	suiMock  *mocks.SUIClient
	*Observer
}

func newTestSuite(t *testing.T) *testSuite {
	ctx := context.Background()

	chain := chains.SuiMainnet
	chainParams := mocks.MockChainParams(chain.ChainId, 10)

	// todo zctx with chain & params (in future PRs)

	zetacore := mocks.NewZetacoreClient(t).
		WithKeys(&keys.Keys{}).
		WithZetaChain().
		WithPostVoteInbound("", "").
		WithPostVoteOutbound("", "")

	tss := mocks.NewTSS(t).FakePubKey(testutils.TSSPubKeyMainnet)

	database, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	log := testlog.New(t)
	logger := base.Logger{
		Std:        log.Logger,
		Compliance: log.Logger,
	}

	baseObserver, err := base.NewObserver(chain, chainParams, zetacore, tss, 1000, nil, database, logger)
	require.NoError(t, err)

	suiMock := mocks.NewSUIClient(t)

	observer := New(baseObserver, suiMock)

	return &testSuite{
		t:        t,
		ctx:      ctx,
		zetaMock: zetacore,
		suiMock:  suiMock,
		Observer: observer,
	}
}
