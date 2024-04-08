package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestVerifyInTxBody(t *testing.T) {
	sampleTo := sample.EthAddress()
	sampleEthTx, sampleEthTxBytes := sample.EthTx(t, chains.EthChain().ChainId, sampleTo, 42)

	// NOTE: errContains == "" means no error
	for _, tc := range []struct {
		desc        string
		msg         types.MsgAddToInTxTracker
		txBytes     []byte
		chainParams observertypes.ChainParams
		tss         observertypes.QueryGetTssAddressResponse
		errContains string
	}{
		{
			desc: "can't verify btc tx tx body",
			msg: types.MsgAddToInTxTracker{
				ChainId: chains.BtcMainnetChain().ChainId,
			},
			txBytes:     sample.Bytes(),
			errContains: "cannot verify inTx body for chain",
		},
		{
			desc: "txBytes can't be unmarshaled",
			msg: types.MsgAddToInTxTracker{
				ChainId: chains.EthChain().ChainId,
			},
			txBytes:     []byte("invalid"),
			errContains: "failed to unmarshal transaction",
		},
		{
			desc: "txHash doesn't correspond",
			msg: types.MsgAddToInTxTracker{
				ChainId: chains.EthChain().ChainId,
				TxHash:  sample.Hash().Hex(),
			},
			txBytes:     sampleEthTxBytes,
			errContains: "invalid hash",
		},
		{
			desc: "chain id doesn't correspond",
			msg: types.MsgAddToInTxTracker{
				ChainId: chains.SepoliaChain().ChainId,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			txBytes:     sampleEthTxBytes,
			errContains: "invalid chain id",
		},
		{
			desc: "invalid coin type",
			msg: types.MsgAddToInTxTracker{
				ChainId:  chains.EthChain().ChainId,
				TxHash:   sampleEthTx.Hash().Hex(),
				CoinType: coin.CoinType(1000),
			},
			txBytes:     sampleEthTxBytes,
			errContains: "coin type not supported",
		},
		{
			desc: "coin types is zeta, but connector contract address is wrong",
			msg: types.MsgAddToInTxTracker{
				ChainId:  chains.EthChain().ChainId,
				TxHash:   sampleEthTx.Hash().Hex(),
				CoinType: coin.CoinType_Zeta,
			},
			txBytes:     sampleEthTxBytes,
			chainParams: observertypes.ChainParams{ConnectorContractAddress: sample.EthAddress().Hex()},
			errContains: "receiver is not connector contract for coin type",
		},
		{
			desc: "coin types is zeta, connector contract address is correct",
			msg: types.MsgAddToInTxTracker{
				ChainId:  chains.EthChain().ChainId,
				TxHash:   sampleEthTx.Hash().Hex(),
				CoinType: coin.CoinType_Zeta,
			},
			txBytes:     sampleEthTxBytes,
			chainParams: observertypes.ChainParams{ConnectorContractAddress: sampleTo.Hex()},
		},
		{
			desc: "coin types is erc20, but erc20 custody contract address is wrong",
			msg: types.MsgAddToInTxTracker{
				ChainId:  chains.EthChain().ChainId,
				TxHash:   sampleEthTx.Hash().Hex(),
				CoinType: coin.CoinType_ERC20,
			},
			txBytes:     sampleEthTxBytes,
			chainParams: observertypes.ChainParams{Erc20CustodyContractAddress: sample.EthAddress().Hex()},
			errContains: "receiver is not erc20Custory contract for coin type",
		},
		{
			desc: "coin types is erc20, erc20 custody contract address is correct",
			msg: types.MsgAddToInTxTracker{
				ChainId:  chains.EthChain().ChainId,
				TxHash:   sampleEthTx.Hash().Hex(),
				CoinType: coin.CoinType_ERC20,
			},
			txBytes:     sampleEthTxBytes,
			chainParams: observertypes.ChainParams{Erc20CustodyContractAddress: sampleTo.Hex()},
		},
		{
			desc: "coin types is gas, but tss address is not found",
			msg: types.MsgAddToInTxTracker{
				ChainId:  chains.EthChain().ChainId,
				TxHash:   sampleEthTx.Hash().Hex(),
				CoinType: coin.CoinType_Gas,
			},
			txBytes:     sampleEthTxBytes,
			tss:         observertypes.QueryGetTssAddressResponse{},
			errContains: "tss address not found",
		},
		{
			desc: "coin types is gas, but tss address is wrong",
			msg: types.MsgAddToInTxTracker{
				ChainId:  chains.EthChain().ChainId,
				TxHash:   sampleEthTx.Hash().Hex(),
				CoinType: coin.CoinType_Gas,
			},
			txBytes:     sampleEthTxBytes,
			tss:         observertypes.QueryGetTssAddressResponse{Eth: sample.EthAddress().Hex()},
			errContains: "receiver is not tssAddress contract for coin type",
		},
		{
			desc: "coin types is gas, tss address is correct",
			msg: types.MsgAddToInTxTracker{
				ChainId:  chains.EthChain().ChainId,
				TxHash:   sampleEthTx.Hash().Hex(),
				CoinType: coin.CoinType_Gas,
			},
			txBytes: sampleEthTxBytes,
			tss:     observertypes.QueryGetTssAddressResponse{Eth: sampleTo.Hex()},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := types.VerifyInTxBody(tc.msg, tc.txBytes, tc.chainParams, tc.tss)
			if tc.errContains == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.errContains)
			}
		})
	}
}

func TestVerifyOutTxBody(t *testing.T) {

	sampleTo := sample.EthAddress()
	sampleEthTx, sampleEthTxBytes, sampleFrom := sample.EthTxSigned(t, chains.EthChain().ChainId, sampleTo, 42)
	_, sampleEthTxBytesNonSigned := sample.EthTx(t, chains.EthChain().ChainId, sampleTo, 42)

	// NOTE: errContains == "" means no error
	for _, tc := range []struct {
		desc        string
		msg         types.MsgAddToOutTxTracker
		txBytes     []byte
		tss         observertypes.QueryGetTssAddressResponse
		errContains string
	}{
		{
			desc: "invalid chain id",
			msg: types.MsgAddToOutTxTracker{
				ChainId: int64(1000),
				Nonce:   42,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			txBytes:     sample.Bytes(),
			errContains: "cannot verify outTx body for chain",
		},
		{
			desc: "txBytes can't be unmarshaled",
			msg: types.MsgAddToOutTxTracker{
				ChainId: chains.EthChain().ChainId,
				Nonce:   42,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			txBytes:     []byte("invalid"),
			errContains: "failed to unmarshal transaction",
		},
		{
			desc: "can't recover sender address",
			msg: types.MsgAddToOutTxTracker{
				ChainId: chains.EthChain().ChainId,
				Nonce:   42,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			txBytes:     sampleEthTxBytesNonSigned,
			errContains: "failed to recover sender",
		},
		{
			desc: "tss address not found",
			msg: types.MsgAddToOutTxTracker{
				ChainId: chains.EthChain().ChainId,
				Nonce:   42,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			tss:         observertypes.QueryGetTssAddressResponse{},
			txBytes:     sampleEthTxBytes,
			errContains: "tss address not found",
		},
		{
			desc: "tss address is wrong",
			msg: types.MsgAddToOutTxTracker{
				ChainId: chains.EthChain().ChainId,
				Nonce:   42,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			tss:         observertypes.QueryGetTssAddressResponse{Eth: sample.EthAddress().Hex()},
			txBytes:     sampleEthTxBytes,
			errContains: "sender is not tss address",
		},
		{
			desc: "chain id doesn't correspond",
			msg: types.MsgAddToOutTxTracker{
				ChainId: chains.SepoliaChain().ChainId,
				Nonce:   42,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			tss:         observertypes.QueryGetTssAddressResponse{Eth: sampleFrom.Hex()},
			txBytes:     sampleEthTxBytes,
			errContains: "invalid chain id",
		},
		{
			desc: "nonce doesn't correspond",
			msg: types.MsgAddToOutTxTracker{
				ChainId: chains.EthChain().ChainId,
				Nonce:   100,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			tss:         observertypes.QueryGetTssAddressResponse{Eth: sampleFrom.Hex()},
			txBytes:     sampleEthTxBytes,
			errContains: "invalid nonce",
		},
		{
			desc: "tx hash doesn't correspond",
			msg: types.MsgAddToOutTxTracker{
				ChainId: chains.EthChain().ChainId,
				Nonce:   42,
				TxHash:  sample.Hash().Hex(),
			},
			tss:         observertypes.QueryGetTssAddressResponse{Eth: sampleFrom.Hex()},
			txBytes:     sampleEthTxBytes,
			errContains: "invalid tx hash",
		},
		{
			desc: "valid out tx body",
			msg: types.MsgAddToOutTxTracker{
				ChainId: chains.EthChain().ChainId,
				Nonce:   42,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			tss:     observertypes.QueryGetTssAddressResponse{Eth: sampleFrom.Hex()},
			txBytes: sampleEthTxBytes,
		},
		// TODO: Implement tests for verifyOutTxBodyBTC
		// https://github.com/zeta-chain/node/issues/1994
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := types.VerifyOutTxBody(tc.msg, tc.txBytes, tc.tss)
			if tc.errContains == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.errContains)
			}
		})
	}
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
