package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// AddToInTxTracker adds a new record to the inbound transaction tracker.
func (k msgServer) AddToInTxTracker(goCtx context.Context, msg *types.MsgAddToInTxTracker) (*types.MsgAddToInTxTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chain := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, msg.ChainId)
	if chain == nil {
		return nil, observertypes.ErrSupportedChains
	}

	// emergency or observer group can submit tracker without proof
	isEmergencyGroup := k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, authoritytypes.PolicyType_groupEmergency)
	isObserver := k.GetObserverKeeper().IsNonTombstonedObserver(ctx, msg.Creator)

	if !(isEmergencyGroup || isObserver) {
		// if not directly authorized, check the proof, if not provided, return unauthorized
		if msg.Proof == nil {
			return nil, errorsmod.Wrap(authoritytypes.ErrUnauthorized, fmt.Sprintf("Creator %s", msg.Creator))
		}

		// verify the proof and tx body
		if err := verifyProofAndInTxBody(ctx, k, msg); err != nil {
			return nil, err
		}
	}

	// add the inTx tracker
	k.SetInTxTracker(ctx, types.InTxTracker{
		ChainId:  msg.ChainId,
		TxHash:   msg.TxHash,
		CoinType: msg.CoinType,
	})

	return &types.MsgAddToInTxTrackerResponse{}, nil
}

// verifyProofAndInTxBody verifies the proof and inbound tx body
func verifyProofAndInTxBody(ctx sdk.Context, k msgServer, msg *types.MsgAddToInTxTracker) error {
	txBytes, err := k.GetLightclientKeeper().VerifyProof(ctx, msg.Proof, msg.ChainId, msg.BlockHash, msg.TxIndex)
	if err != nil {
		return types.ErrProofVerificationFail.Wrapf(err.Error())
	}

	// get chain params and tss addresses to verify the inTx body
	chainParams, found := k.GetObserverKeeper().GetChainParamsByChainID(ctx, msg.ChainId)
	if !found || chainParams == nil {
		return types.ErrUnsupportedChain.Wrapf("chain params not found for chain %d", msg.ChainId)
	}
	tss, err := k.GetObserverKeeper().GetTssAddress(ctx, &observertypes.QueryGetTssAddressRequest{
		BitcoinChainId: msg.ChainId,
	})
	if err != nil || tss == nil {
		reason := "tss response is nil"
		if err != nil {
			reason = err.Error()
		}
		return observertypes.ErrTssNotFound.Wrapf("tss address not found %s", reason)
	}

	if err := types.VerifyInTxBody(*msg, txBytes, *chainParams, *tss); err != nil {
		return types.ErrTxBodyVerificationFail.Wrapf(err.Error())
	}

	return nil
}

// TODO: Implement tests for verifyOutTxBodyBTC
// https://github.com/zeta-chain/node/issues/1994

//func TestKeeper_VerifyEVMInTxBody(t *testing.T) {
//to := sample.EthAddress()
//tx := ethtypes.NewTx(&ethtypes.DynamicFeeTx{
//	ChainID:   big.NewInt(5),
//	Nonce:     1,
//	GasTipCap: nil,
//	GasFeeCap: nil,
//	Gas:       21000,
//	To:        &to,
//	Value:     big.NewInt(5),
//	Data:      nil,
//})
//
//t.Run("should error if msg tx hash not correct", func(t *testing.T) {
//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
//		UseObserverMock: true,
//	})
//
//	txBytes, err := tx.MarshalBinary()
//	require.NoError(t, err)
//	msg := &types.MsgAddToInTxTracker{
//		TxHash: "0x0",
//	}
//
//	err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
//	require.Error(t, err)
//})

//t.Run("should error if msg chain id not correct", func(t *testing.T) {
//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
//		UseObserverMock: true,
//	})
//
//	txBytes, err := tx.MarshalBinary()
//	require.NoError(t, err)
//	msg := &types.MsgAddToInTxTracker{
//		TxHash:  tx.Hash().Hex(),
//		ChainId: 1,
//	}
//
//	err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
//	require.Error(t, err)
//})
//
//t.Run("should error if not supported coin type", func(t *testing.T) {
//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
//		UseObserverMock: true,
//	})
//
//	txBytes, err := tx.MarshalBinary()
//	require.NoError(t, err)
//	msg := &types.MsgAddToInTxTracker{
//		TxHash:   tx.Hash().Hex(),
//		ChainId:  tx.ChainId().Int64(),
//		CoinType: coin.CoinType_Cmd,
//	}
//
//	err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
//	require.Error(t, err)
//})
//
//t.Run("should error for cointype_zeta if chain params not found", func(t *testing.T) {
//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
//		UseObserverMock: true,
//	})
//	observerMock := keepertest.GetCrosschainObserverMock(t, k)
//	observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(&observertypes.ChainParams{}, false)
//
//	txBytes, err := tx.MarshalBinary()
//	require.NoError(t, err)
//	msg := &types.MsgAddToInTxTracker{
//		TxHash:   tx.Hash().Hex(),
//		ChainId:  tx.ChainId().Int64(),
//		CoinType: coin.CoinType_Zeta,
//	}
//
//	err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
//	require.Error(t, err)
//})
//
//t.Run("should error for cointype_zeta if tx.to wrong", func(t *testing.T) {
//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
//		UseObserverMock: true,
//	})
//	observerMock := keepertest.GetCrosschainObserverMock(t, k)
//	observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(&observertypes.ChainParams{
//		ConnectorContractAddress: sample.EthAddress().Hex(),
//	}, true)
//
//	txBytes, err := tx.MarshalBinary()
//	require.NoError(t, err)
//	msg := &types.MsgAddToInTxTracker{
//		TxHash:   tx.Hash().Hex(),
//		ChainId:  tx.ChainId().Int64(),
//		CoinType: coin.CoinType_Zeta,
//	}
//
//	err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
//	require.Error(t, err)
//})
//
//t.Run("should not error for cointype_zeta", func(t *testing.T) {
//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
//		UseObserverMock: true,
//	})
//	observerMock := keepertest.GetCrosschainObserverMock(t, k)
//	observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(&observertypes.ChainParams{
//		ConnectorContractAddress: to.Hex(),
//	}, true)
//
//	txBytes, err := tx.MarshalBinary()
//	require.NoError(t, err)
//	msg := &types.MsgAddToInTxTracker{
//		TxHash:   tx.Hash().Hex(),
//		ChainId:  tx.ChainId().Int64(),
//		CoinType: coin.CoinType_Zeta,
//	}
//
//	err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
//	require.NoError(t, err)
//})
//
//t.Run("should error for cointype_erc20 if chain params not found", func(t *testing.T) {
//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
//		UseObserverMock: true,
//	})
//	observerMock := keepertest.GetCrosschainObserverMock(t, k)
//	observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(&observertypes.ChainParams{}, false)
//
//	txBytes, err := tx.MarshalBinary()
//	require.NoError(t, err)
//	msg := &types.MsgAddToInTxTracker{
//		TxHash:   tx.Hash().Hex(),
//		ChainId:  tx.ChainId().Int64(),
//		CoinType: coin.CoinType_ERC20,
//	}
//
//	err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
//	require.Error(t, err)
//})
//
//t.Run("should error for cointype_erc20 if tx.to wrong", func(t *testing.T) {
//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
//		UseObserverMock: true,
//	})
//	observerMock := keepertest.GetCrosschainObserverMock(t, k)
//	observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(&observertypes.ChainParams{
//		Erc20CustodyContractAddress: sample.EthAddress().Hex(),
//	}, true)
//
//	txBytes, err := tx.MarshalBinary()
//	require.NoError(t, err)
//	msg := &types.MsgAddToInTxTracker{
//		TxHash:   tx.Hash().Hex(),
//		ChainId:  tx.ChainId().Int64(),
//		CoinType: coin.CoinType_ERC20,
//	}
//
//	err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
//	require.Error(t, err)
//})
//
//t.Run("should not error for cointype_erc20", func(t *testing.T) {
//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
//		UseObserverMock: true,
//	})
//	observerMock := keepertest.GetCrosschainObserverMock(t, k)
//	observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(&observertypes.ChainParams{
//		Erc20CustodyContractAddress: to.Hex(),
//	}, true)
//
//	txBytes, err := tx.MarshalBinary()
//	require.NoError(t, err)
//	msg := &types.MsgAddToInTxTracker{
//		TxHash:   tx.Hash().Hex(),
//		ChainId:  tx.ChainId().Int64(),
//		CoinType: coin.CoinType_ERC20,
//	}
//
//	err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
//	require.NoError(t, err)
//})
//
//t.Run("should error for cointype_gas if tss address not found", func(t *testing.T) {
//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
//		UseObserverMock: true,
//	})
//	observerMock := keepertest.GetCrosschainObserverMock(t, k)
//	observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(&observertypes.QueryGetTssAddressResponse{}, errors.New("err"))
//
//	txBytes, err := tx.MarshalBinary()
//	require.NoError(t, err)
//	msg := &types.MsgAddToInTxTracker{
//		TxHash:   tx.Hash().Hex(),
//		ChainId:  tx.ChainId().Int64(),
//		CoinType: coin.CoinType_Gas,
//	}
//
//	err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
//	require.Error(t, err)
//})
//
//t.Run("should error for cointype_gas if tss eth address is empty", func(t *testing.T) {
//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
//		UseObserverMock: true,
//	})
//	observerMock := keepertest.GetCrosschainObserverMock(t, k)
//	observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(&observertypes.QueryGetTssAddressResponse{
//		Eth: "0x",
//	}, nil)
//
//	txBytes, err := tx.MarshalBinary()
//	require.NoError(t, err)
//	msg := &types.MsgAddToInTxTracker{
//		TxHash:   tx.Hash().Hex(),
//		ChainId:  tx.ChainId().Int64(),
//		CoinType: coin.CoinType_Gas,
//	}
//
//	err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
//	require.Error(t, err)
//})
//
//t.Run("should error for cointype_gas if tss eth address is wrong", func(t *testing.T) {
//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
//		UseObserverMock: true,
//	})
//	observerMock := keepertest.GetCrosschainObserverMock(t, k)
//	observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(&observertypes.QueryGetTssAddressResponse{
//		Eth: sample.EthAddress().Hex(),
//	}, nil)
//
//	txBytes, err := tx.MarshalBinary()
//	require.NoError(t, err)
//	msg := &types.MsgAddToInTxTracker{
//		TxHash:   tx.Hash().Hex(),
//		ChainId:  tx.ChainId().Int64(),
//		CoinType: coin.CoinType_Gas,
//	}
//
//	err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
//	require.Error(t, err)
//})
//
//t.Run("should not error for cointype_gas", func(t *testing.T) {
//	k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
//		UseObserverMock: true,
//	})
//	observerMock := keepertest.GetCrosschainObserverMock(t, k)
//	observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(&observertypes.QueryGetTssAddressResponse{
//		Eth: to.Hex(),
//	}, nil)
//
//	txBytes, err := tx.MarshalBinary()
//	require.NoError(t, err)
//	msg := &types.MsgAddToInTxTracker{
//		TxHash:   tx.Hash().Hex(),
//		ChainId:  tx.ChainId().Int64(),
//		CoinType: coin.CoinType_Gas,
//	}
//
//	err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
//	require.NoError(t, err)
//})
//}
