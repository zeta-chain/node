package observer

import (
	"context"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/mode"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func createObserver(t *testing.T,
	zetacoreClient zrepo.ZetacoreClient,
	solanaRepo SolanaRepo,
) *Observer {
	chain := chains.SolanaDevnet
	chainParams := *sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = sample.SolanaAddress(t)
	zetaRepo := zrepo.New(zetacoreClient, chain, mode.StandardMode)

	db, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	logger := base.DefaultLogger()

	base, err := base.NewObserver(chain, chainParams, zetaRepo, nil, 1000, nil, db, logger)
	require.NoError(t, err)

	observer, err := New(base, nil, chainParams.GatewayAddress)
	require.NoError(t, err)
	observer.solanaRepo = solanaRepo

	return observer
}

func TestPostGasPrice(t *testing.T) {
	anything := mock.Anything

	priorityFee := uint64(5)
	slot := uint64(100)

	t.Run("Ok", func(t *testing.T) {
		solanaRepo := mocks.NewSolanaRepo(t)
		solanaRepo.On("GetPriorityFee", anything).Return(priorityFee, nil)
		solanaRepo.On("GetSlot", anything, rpc.CommitmentConfirmed).Return(slot, nil)

		zetacoreClient := mocks.NewZetacoreClient(t)
		zetacoreClient.
			On("PostVoteGasPrice", anything, anything, anything, priorityFee, slot).
			Return(anything, nil)

		ob := createObserver(t, zetacoreClient, solanaRepo)
		err := ob.PostGasPrice(context.Background())
		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		errTest := errors.New("test error")

		t.Run("GetPriorityFee", func(t *testing.T) {
			solanaRepo := mocks.NewSolanaRepo(t)
			solanaRepo.On("GetPriorityFee", anything).Return(uint64(0), errTest)

			zetacoreClient := mocks.NewZetacoreClient(t)

			ob := createObserver(t, zetacoreClient, solanaRepo)
			err := ob.PostGasPrice(context.Background())

			require.Error(t, err)
			require.ErrorIs(t, err, errTest)
		})

		t.Run("GetSlot", func(t *testing.T) {
			solanaRepo := mocks.NewSolanaRepo(t)
			solanaRepo.On("GetPriorityFee", anything).Return(priorityFee, nil)
			solanaRepo.On("GetSlot", anything, rpc.CommitmentConfirmed).Return(uint64(0), errTest)

			zetacoreClient := mocks.NewZetacoreClient(t)

			ob := createObserver(t, zetacoreClient, solanaRepo)
			err := ob.PostGasPrice(context.Background())

			require.Error(t, err)
			require.ErrorIs(t, err, errTest)
		})

		t.Run("VoteGasPrice", func(t *testing.T) {
			solanaRepo := mocks.NewSolanaRepo(t)
			solanaRepo.On("GetPriorityFee", anything).Return(priorityFee, nil)
			solanaRepo.On("GetSlot", anything, rpc.CommitmentConfirmed).Return(slot, nil)

			zetacoreClient := mocks.NewZetacoreClient(t)
			zetacoreClient.
				On("PostVoteGasPrice", anything, anything, anything, priorityFee, slot).
				Return(anything, errTest)

			ob := createObserver(t, zetacoreClient, solanaRepo)
			err := ob.PostGasPrice(context.Background())

			require.Error(t, err)
			require.ErrorIs(t, err, errTest)
		})
	})
}
