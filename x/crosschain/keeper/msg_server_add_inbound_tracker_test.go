package keeper_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/zeta-chain/zetacore/pkg/chains"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
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

		keepertest.MockIsAuthorized(&authorityMock.Mock, nonAdmin, authoritytypes.PolicyType_groupEmergency, false)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)

		txHash := "string"
		chainID := getValidEthChainID()

		_, err := msgServer.AddInboundTracker(ctx, &types.MsgAddInboundTracker{
			Creator:   nonAdmin,
			ChainId:   chainID,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Zeta,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
		})
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

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(nil)

		txHash := "string"
		chainID := getValidEthChainID()

		_, err := msgServer.AddInboundTracker(ctx, &types.MsgAddInboundTracker{
			Creator:   sample.AccAddress(),
			ChainId:   chainID + 1,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Zeta,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
		})
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

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)

		txHash := "string"
		chainID := getValidEthChainID()
		setSupportedChain(ctx, zk, chainID)

		_, err := msgServer.AddInboundTracker(ctx, &types.MsgAddInboundTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Zeta,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
		})
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

		keepertest.MockIsAuthorized(&authorityMock.Mock, mock.Anything, authoritytypes.PolicyType_groupEmergency, false)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(true)

		txHash := "string"
		chainID := getValidEthChainID()

		_, err := msgServer.AddInboundTracker(ctx, &types.MsgAddInboundTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Zeta,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
		})
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

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		lightclientMock := keepertest.GetCrosschainLightclientMock(t, k)

		keepertest.MockIsAuthorized(&authorityMock.Mock, mock.Anything, authoritytypes.PolicyType_groupEmergency, false)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		lightclientMock.On("VerifyProof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))

		txHash := "string"
		chainID := getValidEthChainID()

		_, err := msgServer.AddInboundTracker(ctx, &types.MsgAddInboundTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Zeta,
			Proof:     &proofs.Proof{},
			BlockHash: "",
			TxIndex:   0,
		})
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

		keepertest.MockIsAuthorized(&authorityMock.Mock, mock.Anything, authoritytypes.PolicyType_groupEmergency, false)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		lightclientMock.On("VerifyProof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(sample.Bytes(), nil)
		observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(nil, false)

		txHash := "string"
		chainID := getValidEthChainID()

		_, err := msgServer.AddInboundTracker(ctx, &types.MsgAddInboundTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Zeta,
			Proof:     &proofs.Proof{},
			BlockHash: "",
			TxIndex:   0,
		})
		require.ErrorIs(t, err, types.ErrUnsupportedChain)
	})

	t.Run("fail if proof is provided but can't find tss to verify body", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock:   true,
			UseLightclientMock: true,
			UseObserverMock:    true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		lightclientMock := keepertest.GetCrosschainLightclientMock(t, k)

		keepertest.MockIsAuthorized(&authorityMock.Mock, mock.Anything, authoritytypes.PolicyType_groupEmergency, false)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		lightclientMock.On("VerifyProof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(sample.Bytes(), nil)
		observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(sample.ChainParams(chains.EthChain.ChainId), true)
		observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

		txHash := "string"
		chainID := getValidEthChainID()
		setSupportedChain(ctx, zk, chainID)

		_, err := msgServer.AddInboundTracker(ctx, &types.MsgAddInboundTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Zeta,
			Proof:     &proofs.Proof{},
			BlockHash: "",
			TxIndex:   0,
		})
		require.ErrorIs(t, err, observertypes.ErrTssNotFound)
	})

	t.Run("fail if proof is provided but error while verifying tx body", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock:   true,
			UseLightclientMock: true,
			UseObserverMock:    true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		lightclientMock := keepertest.GetCrosschainLightclientMock(t, k)

		keepertest.MockIsAuthorized(&authorityMock.Mock, mock.Anything, authoritytypes.PolicyType_groupEmergency, false)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(sample.ChainParams(chains.EthChain.ChainId), true)
		observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(&observertypes.QueryGetTssAddressResponse{
			Eth: sample.EthAddress().Hex(),
		}, nil)

		// verifying the body will fail because the bytes are tried to be unmarshaled but they are not valid
		lightclientMock.On("VerifyProof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]byte("invalid"), nil)

		txHash := "string"
		chainID := getValidEthChainID()
		setSupportedChain(ctx, zk, chainID)

		_, err := msgServer.AddInboundTracker(ctx, &types.MsgAddInboundTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Zeta,
			Proof:     &proofs.Proof{},
			BlockHash: "",
			TxIndex:   0,
		})
		require.ErrorIs(t, err, types.ErrTxBodyVerificationFail)
	})

	t.Run("can add a in tx tracker with a proof", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock:   true,
			UseLightclientMock: true,
			UseObserverMock:    true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()

		chainID := chains.EthChain.ChainId
		tssAddress := sample.EthAddress()
		ethTx, ethTxBytes := sample.EthTx(t, chainID, tssAddress, 42)
		txHash := ethTx.Hash().Hex()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		lightclientMock := keepertest.GetCrosschainLightclientMock(t, k)

		keepertest.MockIsAuthorized(&authorityMock.Mock, mock.Anything, authoritytypes.PolicyType_groupEmergency, false)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		observerMock.On("GetChainParamsByChainID", mock.Anything, mock.Anything).Return(sample.ChainParams(chains.EthChain.ChainId), true)
		observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(&observertypes.QueryGetTssAddressResponse{
			Eth: tssAddress.Hex(),
		}, nil)
		lightclientMock.On("VerifyProof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ethTxBytes, nil)

		_, err := msgServer.AddInboundTracker(ctx, &types.MsgAddInboundTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			CoinType:  coin.CoinType_Gas, // use coin types gas: the receiver must be the tss address
			Proof:     &proofs.Proof{},
			BlockHash: "",
			TxIndex:   0,
		})
		require.NoError(t, err)
		_, found := k.GetInboundTracker(ctx, chainID, txHash)
		require.True(t, found)
	})
}
