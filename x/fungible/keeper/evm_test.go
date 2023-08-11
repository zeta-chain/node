package keeper_test

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/server/config"
	"github.com/zeta-chain/zetacore/testutil/sample"

	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
)

func TestKeeper_CallEVMWithData(t *testing.T) {

	t.Run("apply new message", func(t *testing.T) {
		k, ctx := testkeeper.FungibleKeeper(t)

		mockAuthKeeper := testkeeper.GetFungibleAccountKeeper(t, k)
		mockEVMKeeper := testkeeper.GetFungibleEVMKeeper(t, k)

		// Set up expectations
		fromAddr := sample.EthAddress()
		contractAddress := sample.EthAddress()
		nonce := uint64(1)
		chainID := big.NewInt(1)
		mockAuthKeeper.On("GetSequence", mock.Anything, sdk.AccAddress(fromAddr.Bytes())).Return(nonce, nil)

		data := []byte("some data")
		args, _ := json.Marshal(evmtypes.TransactionArgs{
			From: &fromAddr,
			To:   &contractAddress,
			Data: (*hexutil.Bytes)(&data),
		})
		gasRes := &evmtypes.EstimateGasResponse{Gas: 1000}
		mockEVMKeeper.On(
			"EstimateGas",
			mock.Anything,
			&evmtypes.EthCallRequest{Args: args, GasCap: config.DefaultGasCap},
		).Return(gasRes, nil)

		value := big.NewInt(100)
		msgRes := &evmtypes.MsgEthereumTxResponse{}
		mockEVMKeeper.On(
			"ApplyMessage",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(msgRes, nil)

		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(chainID)
		mockEVMKeeper.On("SetBlockBloomTransient", mock.Anything).Maybe()
		mockEVMKeeper.On("SetLogSizeTransient", mock.Anything).Maybe()
		mockEVMKeeper.On("GetLogSizeTransient", mock.Anything, mock.Anything).Maybe()

		// Call the method
		res, err := k.CallEVMWithData(
			ctx,
			fromAddr,
			&contractAddress,
			data,
			true,
			value,
			nil,
		)
		require.NoError(t, err)
		require.Equal(t, msgRes, res)

		// Assert that the expected methods were called
		mockAuthKeeper.AssertExpectations(t)
		mockEVMKeeper.AssertExpectations(t)
	})
}
