package keeper_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/proofs"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestMsgServer_AddToInboundTracker(t *testing.T) {
	t.Run("fail normal user submit without proof", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		nonAdmin := sample.AccAddress()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		txHash := "string"
		chainID := getValidEthChainID()

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)

		msg := types.MsgAddInboundTracker{
			Creator:   nonAdmin,
			ChainId:   chainID,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Zeta,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.AddInboundTracker(ctx, &msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
		_, found := k.GetInboundTracker(ctx, chainID, txHash)
		require.False(t, found)
	})

	t.Run("fail for unsupported chain id", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		txHash := "string"
		chainID := getValidEthChainID()

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, false)

		msg := types.MsgAddInboundTracker{
			Creator:   sample.AccAddress(),
			ChainId:   chainID + 1,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Zeta,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
		}
		_, err := msgServer.AddInboundTracker(ctx, &msg)
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
		_, found := k.GetInboundTracker(ctx, chainID, txHash)
		require.False(t, found)
	})

	t.Run("admin add tx tracker", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		txHash := "string"
		chainID := getValidEthChainID()

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)

		setSupportedChain(ctx, zk, chainID)

		msg := types.MsgAddInboundTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Zeta,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.AddInboundTracker(ctx, &msg)

		require.NoError(t, err)
		_, found := k.GetInboundTracker(ctx, chainID, txHash)
		require.True(t, found)
	})

	t.Run("observer add tx tracker", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		txHash := "string"
		chainID := getValidEthChainID()

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(true)

		msg := types.MsgAddInboundTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Zeta,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.AddInboundTracker(ctx, &msg)
		require.NoError(t, err)
		_, found := k.GetInboundTracker(ctx, chainID, txHash)
		require.True(t, found)
	})

	t.Run("fail if proof is provided but not verified", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock:   true,
			UseLightclientMock: true,
			UseObserverMock:    true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		lightclientMock := keepertest.GetCrosschainLightclientMock(t, k)

		admin := sample.AccAddress()
		txHash := "string"
		chainID := getValidEthChainID()

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		lightclientMock.On("VerifyProof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("error"))

		msg := types.MsgAddInboundTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Zeta,
			Proof:     &proofs.Proof{},
			BlockHash: "",
			TxIndex:   0,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.AddInboundTracker(ctx, &msg)
		require.ErrorIs(t, err, types.ErrProofVerificationFail)
	})

	t.Run("fail if proof is provided but can't find chain params to verify body", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock:   true,
			UseLightclientMock: true,
			UseObserverMock:    true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		lightclientMock := keepertest.GetCrosschainLightclientMock(t, k)
		txHash := "string"
		chainID := getValidEthChainID()

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		lightclientMock.On("VerifyProof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(sample.Bytes(), nil)
		observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(nil, false)

		msg := types.MsgAddInboundTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Zeta,
			Proof:     &proofs.Proof{},
			BlockHash: "",
			TxIndex:   0,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.AddInboundTracker(ctx, &msg)
		require.ErrorIs(t, err, types.ErrUnsupportedChain)
	})

	t.Run("fail if proof is provided but can't find tss to verify body", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock:   true,
			UseLightclientMock: true,
			UseObserverMock:    true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		lightclientMock := keepertest.GetCrosschainLightclientMock(t, k)

		admin := sample.AccAddress()
		txHash := "string"
		chainID := getValidEthChainID()

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		lightclientMock.On("VerifyProof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(sample.Bytes(), nil)
		observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).
			Return(sample.ChainParams(chains.Ethereum.ChainId), true)
		observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

		setSupportedChain(ctx, zk, chainID)

		msg := types.MsgAddInboundTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Zeta,
			Proof:     &proofs.Proof{},
			BlockHash: "",
			TxIndex:   0,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.AddInboundTracker(ctx, &msg)
		require.ErrorIs(t, err, observertypes.ErrTssNotFound)
	})

	t.Run("fail if proof is provided but error while verifying tx body", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock:   true,
			UseLightclientMock: true,
			UseObserverMock:    true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		lightclientMock := keepertest.GetCrosschainLightclientMock(t, k)

		admin := sample.AccAddress()
		txHash := "string"
		chainID := getValidEthChainID()

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).
			Return(sample.ChainParams(chains.Ethereum.ChainId), true)
		observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(&observertypes.QueryGetTssAddressResponse{
			Eth: sample.EthAddress().Hex(),
		}, nil)

		// verifying the body will fail because the bytes are tried to be unmarshaled but they are not valid
		lightclientMock.On("VerifyProof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]byte("invalid"), nil)

		setSupportedChain(ctx, zk, chainID)

		msg := types.MsgAddInboundTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Zeta,
			Proof:     &proofs.Proof{},
			BlockHash: "",
			TxIndex:   0,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.AddInboundTracker(ctx, &msg)
		require.ErrorIs(t, err, types.ErrTxBodyVerificationFail)
	})

	t.Run("can add a in tx tracker with a proof", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock:   true,
			UseLightclientMock: true,
			UseObserverMock:    true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		chainID := chains.Ethereum.ChainId
		tssAddress := sample.EthAddress()
		ethTx, ethTxBytes := sample.EthTx(t, chainID, tssAddress, 42)
		admin := sample.AccAddress()
		txHash := ethTx.Hash().Hex()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		lightclientMock := keepertest.GetCrosschainLightclientMock(t, k)

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).
			Return(sample.ChainParams(chains.Ethereum.ChainId), true)
		observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(&observertypes.QueryGetTssAddressResponse{
			Eth: tssAddress.Hex(),
		}, nil)
		lightclientMock.On("VerifyProof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(ethTxBytes, nil)

		msg := types.MsgAddInboundTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Gas, // use coin types gas: the receiver must be the tss address
			Proof:     &proofs.Proof{},
			BlockHash: "",
			TxIndex:   0,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.AddInboundTracker(ctx, &msg)
		require.NoError(t, err)
		_, found := k.GetInboundTracker(ctx, chainID, txHash)
		require.True(t, found)
	})
}
