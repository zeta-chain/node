package types_test

import (
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"math/big"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/coin"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_VerifyEVMInTxBody(t *testing.T) {
	to := sample.EthAddress()
	tx := ethtypes.NewTx(&ethtypes.DynamicFeeTx{
		ChainID:   big.NewInt(5),
		Nonce:     1,
		GasTipCap: nil,
		GasFeeCap: nil,
		Gas:       21000,
		To:        &to,
		Value:     big.NewInt(5),
		Data:      nil,
	})
	t.Run("should error if msg tx hash not correct", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		txBytes, err := tx.MarshalBinary()
		require.NoError(t, err)
		msg := &types.MsgAddToInTxTracker{
			TxHash: "0x0",
		}

		err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
		require.Error(t, err)
	})

	t.Run("should error if msg chain id not correct", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		txBytes, err := tx.MarshalBinary()
		require.NoError(t, err)
		msg := &types.MsgAddToInTxTracker{
			TxHash:  tx.Hash().Hex(),
			ChainId: 1,
		}

		err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
		require.Error(t, err)
	})

	t.Run("should error if not supported coin type", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		txBytes, err := tx.MarshalBinary()
		require.NoError(t, err)
		msg := &types.MsgAddToInTxTracker{
			TxHash:   tx.Hash().Hex(),
			ChainId:  tx.ChainId().Int64(),
			CoinType: coin.CoinType_Cmd,
		}

		err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
		require.Error(t, err)
	})

	t.Run("should error for cointype_zeta if chain params not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(&observertypes.ChainParams{}, false)

		txBytes, err := tx.MarshalBinary()
		require.NoError(t, err)
		msg := &types.MsgAddToInTxTracker{
			TxHash:   tx.Hash().Hex(),
			ChainId:  tx.ChainId().Int64(),
			CoinType: coin.CoinType_Zeta,
		}

		err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
		require.Error(t, err)
	})

	t.Run("should error for cointype_zeta if tx.to wrong", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(&observertypes.ChainParams{
			ConnectorContractAddress: sample.EthAddress().Hex(),
		}, true)

		txBytes, err := tx.MarshalBinary()
		require.NoError(t, err)
		msg := &types.MsgAddToInTxTracker{
			TxHash:   tx.Hash().Hex(),
			ChainId:  tx.ChainId().Int64(),
			CoinType: coin.CoinType_Zeta,
		}

		err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
		require.Error(t, err)
	})

	t.Run("should not error for cointype_zeta", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(&observertypes.ChainParams{
			ConnectorContractAddress: to.Hex(),
		}, true)

		txBytes, err := tx.MarshalBinary()
		require.NoError(t, err)
		msg := &types.MsgAddToInTxTracker{
			TxHash:   tx.Hash().Hex(),
			ChainId:  tx.ChainId().Int64(),
			CoinType: coin.CoinType_Zeta,
		}

		err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
		require.NoError(t, err)
	})

	t.Run("should error for cointype_erc20 if chain params not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(&observertypes.ChainParams{}, false)

		txBytes, err := tx.MarshalBinary()
		require.NoError(t, err)
		msg := &types.MsgAddToInTxTracker{
			TxHash:   tx.Hash().Hex(),
			ChainId:  tx.ChainId().Int64(),
			CoinType: coin.CoinType_ERC20,
		}

		err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
		require.Error(t, err)
	})

	t.Run("should error for cointype_erc20 if tx.to wrong", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(&observertypes.ChainParams{
			Erc20CustodyContractAddress: sample.EthAddress().Hex(),
		}, true)

		txBytes, err := tx.MarshalBinary()
		require.NoError(t, err)
		msg := &types.MsgAddToInTxTracker{
			TxHash:   tx.Hash().Hex(),
			ChainId:  tx.ChainId().Int64(),
			CoinType: coin.CoinType_ERC20,
		}

		err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
		require.Error(t, err)
	})

	t.Run("should not error for cointype_erc20", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(&observertypes.ChainParams{
			Erc20CustodyContractAddress: to.Hex(),
		}, true)

		txBytes, err := tx.MarshalBinary()
		require.NoError(t, err)
		msg := &types.MsgAddToInTxTracker{
			TxHash:   tx.Hash().Hex(),
			ChainId:  tx.ChainId().Int64(),
			CoinType: coin.CoinType_ERC20,
		}

		err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
		require.NoError(t, err)
	})

	t.Run("should error for cointype_gas if tss address not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(&observertypes.QueryGetTssAddressResponse{}, errors.New("err"))

		txBytes, err := tx.MarshalBinary()
		require.NoError(t, err)
		msg := &types.MsgAddToInTxTracker{
			TxHash:   tx.Hash().Hex(),
			ChainId:  tx.ChainId().Int64(),
			CoinType: coin.CoinType_Gas,
		}

		err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
		require.Error(t, err)
	})

	t.Run("should error for cointype_gas if tss eth address is empty", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(&observertypes.QueryGetTssAddressResponse{
			Eth: "0x",
		}, nil)

		txBytes, err := tx.MarshalBinary()
		require.NoError(t, err)
		msg := &types.MsgAddToInTxTracker{
			TxHash:   tx.Hash().Hex(),
			ChainId:  tx.ChainId().Int64(),
			CoinType: coin.CoinType_Gas,
		}

		err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
		require.Error(t, err)
	})

	t.Run("should error for cointype_gas if tss eth address is wrong", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(&observertypes.QueryGetTssAddressResponse{
			Eth: sample.EthAddress().Hex(),
		}, nil)

		txBytes, err := tx.MarshalBinary()
		require.NoError(t, err)
		msg := &types.MsgAddToInTxTracker{
			TxHash:   tx.Hash().Hex(),
			ChainId:  tx.ChainId().Int64(),
			CoinType: coin.CoinType_Gas,
		}

		err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
		require.Error(t, err)
	})

	t.Run("should not error for cointype_gas", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(&observertypes.QueryGetTssAddressResponse{
			Eth: to.Hex(),
		}, nil)

		txBytes, err := tx.MarshalBinary()
		require.NoError(t, err)
		msg := &types.MsgAddToInTxTracker{
			TxHash:   tx.Hash().Hex(),
			ChainId:  tx.ChainId().Int64(),
			CoinType: coin.CoinType_Gas,
		}

		err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
		require.NoError(t, err)
	})
}
