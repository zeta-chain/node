package keeper

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/server/config"

	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
)

func TestKeeper_CallEVMWithData(t *testing.T) {

	t.Run("apply new message", func(t *testing.T) {
		k, ctx := testkeeper.FungibleKeeper(t)

		mockAuthKeeper := testkeeper.GetFungibleAccountKeeper(t, k.authKeeper)
		mockEVMKeeper := testkeeper.GetFungibleEVMKeeper(t, k.evmKeeper)

		// Set up expectations
		fromAddr := common.HexToAddress("0xSomeAddress")
		nonce := uint64(1)
		mockAuthKeeper.On("GetSequence", mock.Anything, fromAddr.Bytes()).Return(nonce, nil)

		data := []byte("some data")
		args, _ := json.Marshal(evmtypes.TransactionArgs{
			From: &fromAddr,
			Data: (*hexutil.Bytes)(&data),
		})
		gasRes := &evmtypes.EstimateGasResponse{Gas: 1000}
		mockEVMKeeper.On("EstimateGas", mock.Anything, &evmtypes.EthCallRequest{Args: args, GasCap: config.DefaultGasCap}).Return(gasRes, nil)

		msgRes := &evmtypes.MsgEthereumTxResponse{}
		mockEVMKeeper.On("ApplyMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(msgRes, nil)

		// Call the method
		res, err := k.CallEVMWithData(
			ctx,
			fromAddr,
			nil,
			data,
			true,
			nil,
			nil,
		)

		// Assertions
		require.NoError(t, err)
		require.Equal(t, msgRes, res)

		// Assert that the expected methods were called
		mockAuthKeeper.AssertExpectations(t)
		mockEVMKeeper.AssertExpectations(t)
	})

	// k.authKeeper.GetSequence(ctx, from.Bytes())
	// k.evmKeeper.EstimateGas(sdk.WrapSDKContext(ctx), &evmtypes.EthCallRequest{})
	// k.evmKeeper.WithChainID(ctx)
	// k.evmKeeper.ApplyMessage(ctx, msg, evmtypes.NewNoOpTracer(), commit)
	// k.evmKeeper.GetBlockBloomTransient(ctx)
	// k.evmKeeper.SetBlockBloomTransient(ctx, bloomReceipt.Big())
	// k.evmKeeper.SetLogSizeTransient(ctx, (k.evmKeeper.GetLogSizeTransient(ctx))+uint64(len(logs)))
}
