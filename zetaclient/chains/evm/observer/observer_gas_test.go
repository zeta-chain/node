package observer

import (
	"context"
	"math/big"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostGasPrice(t *testing.T) {
	const (
		gwei     = 10e9
		anything = mock.Anything
	)

	ctx := context.Background()

	t.Run("Pre EIP-1559 doesn't support priorityFee", func(t *testing.T) {
		// ARRANGE
		// Given an observer
		observer := newTestSuite(t)

		// Given empty baseFee from RPC
		observer.evmMock.On("HeaderByNumber", anything, anything).Maybe().Return(&ethtypes.Header{BaseFee: nil}, nil)

		// Given gasPrice and priorityFee from RPC
		observer.evmMock.On("SuggestGasPrice", anything).Maybe().Return(big.NewInt(3*gwei), nil)
		observer.evmMock.On("SuggestGasTipCap", anything).Maybe().Return(big.NewInt(0), nil)

		// Given mock collector for zetacore call
		// PostVoteGasPrice(ctx, chain, gasPrice, priorityFee, blockNum)
		var gasPrice, priorityFee uint64
		collector := func(args mock.Arguments) {
			gasPrice = args.Get(2).(uint64)
			priorityFee = args.Get(3).(uint64)
		}

		observer.zetacore.
			On("PostVoteGasPrice", anything, anything, anything, anything, anything).
			Run(collector).
			Return("0xABC123...", nil)

		// ACT
		err := observer.PostGasPrice(ctx)

		// ASSERT
		assert.NoError(t, err)

		// Check that gas price is posted with proper gasPrice and priorityFee
		assert.Equal(t, uint64(3*gwei), gasPrice)
		assert.Equal(t, uint64(0), priorityFee)
	})

	// TODO: https://github.com/zeta-chain/node/issues/3221
	// t.Run("Post EIP-1559 supports priorityFee", func(t *testing.T) {
	// 	// ARRANGE
	// 	// Given an observer
	// 	observer := newTestSuite(t)

	// 	// Given 1 gwei baseFee from RPC
	// 	observer.evmMock.On("HeaderByNumber", anything, anything).
	// 		Return(&ethtypes.Header{BaseFee: big.NewInt(gwei)}, nil)

	// 	// Given gasPrice and priorityFee from RPC
	// 	observer.evmMock.On("SuggestGasPrice", anything).Return(big.NewInt(3*gwei), nil)
	// 	observer.evmMock.On("SuggestGasTipCap", anything).Return(big.NewInt(2*gwei), nil)

	// 	// Given mock collector for zetacore call
	// 	// PostVoteGasPrice(ctx, chain, gasPrice, priorityFee, blockNum)
	// 	var gasPrice, priorityFee uint64
	// 	collector := func(args mock.Arguments) {
	// 		gasPrice = args.Get(2).(uint64)
	// 		priorityFee = args.Get(3).(uint64)
	// 	}

	// 	observer.zetacore.
	// 		On("PostVoteGasPrice", anything, anything, anything, anything, anything).
	// 		Run(collector).
	// 		Return("0xABC123...", nil)

	// 	// ACT
	// 	err := observer.PostGasPrice(ctx)

	// 	// ASSERT
	// 	assert.NoError(t, err)

	// 	// Check that gas price is posted with proper gasPrice and priorityFee
	// 	assert.Equal(t, uint64(3*gwei), gasPrice)
	// 	assert.Equal(t, uint64(2*gwei), priorityFee)
	// })
}
