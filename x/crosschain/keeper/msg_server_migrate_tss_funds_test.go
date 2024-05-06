package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/gas"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func setupTssMigrationParams(
	zk keepertest.ZetaKeepers,
	k *keeper.Keeper,
	ctx sdk.Context,
	chain chains.Chain,
	amount sdkmath.Uint,
	setNewTss bool,
	setCurrentTSS bool,
) (string, string) {
	zk.ObserverKeeper.SetCrosschainFlags(ctx, observertypes.CrosschainFlags{
		IsInboundEnabled:  false,
		IsOutboundEnabled: true,
	})

	zk.ObserverKeeper.SetChainParamsList(ctx, observertypes.ChainParamsList{
		ChainParams: []*observertypes.ChainParams{
			{
				ChainId:               chain.ChainId,
				BallotThreshold:       sdk.NewDec(0),
				MinObserverDelegation: sdk.OneDec(),
				IsSupported:           true,
			},
		},
	})

	currentTss := sample.Tss()
	newTss := sample.Tss()
	newTss.FinalizedZetaHeight = currentTss.FinalizedZetaHeight + 1
	newTss.KeyGenZetaHeight = currentTss.KeyGenZetaHeight + 1
	k.GetObserverKeeper().SetTSS(ctx, currentTss)
	if setCurrentTSS {
		k.GetObserverKeeper().SetTSSHistory(ctx, currentTss)
	}
	if setNewTss {
		k.GetObserverKeeper().SetTSSHistory(ctx, newTss)
	}
	k.GetObserverKeeper().SetPendingNonces(ctx, observertypes.PendingNonces{
		NonceLow:  1,
		NonceHigh: 1,
		ChainId:   chain.ChainId,
		Tss:       currentTss.TssPubkey,
	})
	k.SetGasPrice(ctx, crosschaintypes.GasPrice{
		Creator:     "",
		Index:       "",
		ChainId:     chain.ChainId,
		Signers:     nil,
		BlockNums:   nil,
		Prices:      []uint64{100000, 100000, 100000},
		MedianIndex: 1,
	})
	k.GetObserverKeeper().SetChainNonces(ctx, observertypes.ChainNonces{
		Index:   chain.ChainName.String(),
		ChainId: chain.ChainId,
		Nonce:   1,
	})
	indexString := keeper.GetIndexStringForTssMigration(currentTss.TssPubkey, newTss.TssPubkey, chain.ChainId, amount, ctx.BlockHeight())
	return indexString, currentTss.TssPubkey
}

func TestKeeper_MigrateTSSFundsForChain(t *testing.T) {
	t.Run("test evm chain", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		indexString, _ := setupTssMigrationParams(zk, k, ctx, *chain, amount, true, true)
		gp, found := k.GetMedianGasPriceInUint(ctx, chain.ChainId)
		require.True(t, found)
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		require.NoError(t, err)
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		cctx, found := k.GetCrossChainTx(ctx, index)
		require.True(t, found)
		multipliedValue, err := gas.MultiplyGasPrice(gp, crosschaintypes.TssMigrationGasMultiplierEVM)
		require.NoError(t, err)
		require.Equal(t, multipliedValue.String(), cctx.GetCurrentOutboundParam().GasPrice)
	})

	t.Run("test btc chain", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidBTCChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		indexString, _ := setupTssMigrationParams(zk, k, ctx, *chain, amount, true, true)
		gp, found := k.GetMedianGasPriceInUint(ctx, chain.ChainId)
		require.True(t, found)
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		require.NoError(t, err)
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		cctx, found := k.GetCrossChainTx(ctx, index)
		require.True(t, found)
		require.Equal(t, gp.MulUint64(2).String(), cctx.GetCurrentOutboundParam().GasPrice)
	})
}

func TestMsgServer_MigrateTssFunds(t *testing.T) {
	t.Run("should error if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, false)

		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		require.Error(t, err)
	})

	t.Run("should error if inbound enabled", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("IsInboundEnabled", mock.Anything).Return(true)
		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		require.Error(t, err)
	})

	t.Run("should error if tss not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("IsInboundEnabled", mock.Anything).Return(false)
		observerMock.On("GetTSS", mock.Anything).Return(observertypes.TSS{}, false)

		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		require.Error(t, err)
	})

	t.Run("should error if tss history empty", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("IsInboundEnabled", mock.Anything).Return(false)
		observerMock.On("GetTSS", mock.Anything).Return(sample.Tss(), true)
		observerMock.On("GetAllTSS", mock.Anything).Return([]observertypes.TSS{})

		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		require.Error(t, err)
	})

	t.Run("should error if no new tss generated", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("IsInboundEnabled", mock.Anything).Return(false)
		tss := sample.Tss()
		observerMock.On("GetTSS", mock.Anything).Return(tss, true)
		observerMock.On("GetAllTSS", mock.Anything).Return([]observertypes.TSS{tss})

		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		require.Error(t, err)
	})

	t.Run("should error if current tss is the latest", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("IsInboundEnabled", mock.Anything).Return(false)
		tss1 := sample.Tss()
		tss1.FinalizedZetaHeight = 2
		tss2 := sample.Tss()
		tss2.FinalizedZetaHeight = 1
		observerMock.On("GetTSS", mock.Anything).Return(tss1, true)
		observerMock.On("GetAllTSS", mock.Anything).Return([]observertypes.TSS{tss2})

		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		require.Error(t, err)
	})

	t.Run("should error if pending nonces not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("IsInboundEnabled", mock.Anything).Return(false)
		tss1 := sample.Tss()
		tss1.FinalizedZetaHeight = 2
		tss2 := sample.Tss()
		tss2.FinalizedZetaHeight = 3
		observerMock.On("GetTSS", mock.Anything).Return(tss1, true)
		observerMock.On("GetAllTSS", mock.Anything).Return([]observertypes.TSS{tss2})
		observerMock.On("GetPendingNonces", mock.Anything, mock.Anything, mock.Anything).Return(observertypes.PendingNonces{}, false)

		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		require.Error(t, err)
	})

	t.Run("successfully create tss migration cctx", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		indexString, _ := setupTssMigrationParams(zk, k, ctx, *chain, amount, true, true)
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		require.NoError(t, err)
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		cctx, found := k.GetCrossChainTx(ctx, index)
		require.True(t, found)
		feeCalculated := sdk.NewUint(cctx.GetCurrentOutboundParam().GasLimit).
			Mul(sdkmath.NewUintFromString(cctx.GetCurrentOutboundParam().GasPrice))
		require.Equal(t, cctx.GetCurrentOutboundParam().Amount.String(), amount.Sub(feeCalculated).String())
	})

	t.Run("not enough funds in tss address for migration", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("100")
		indexString, _ := setupTssMigrationParams(zk, k, ctx, *chain, amount, true, true)
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		require.ErrorContains(t, err, crosschaintypes.ErrCannotMigrateTssFunds.Error())
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		_, found := k.GetCrossChainTx(ctx, index)
		require.False(t, found)
	})

	t.Run("unable to migrate funds if new TSS is not created ", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		indexString, _ := setupTssMigrationParams(zk, k, ctx, *chain, amount, false, true)
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		require.ErrorContains(t, err, "no new tss address has been generated")
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		_, found := k.GetCrossChainTx(ctx, index)
		require.False(t, found)
	})

	t.Run("unable to migrate funds when nonce low does not match nonce high", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		indexString, tssPubkey := setupTssMigrationParams(zk, k, ctx, *chain, amount, true, true)
		k.GetObserverKeeper().SetPendingNonces(ctx, observertypes.PendingNonces{
			NonceLow:  1,
			NonceHigh: 10,
			ChainId:   chain.ChainId,
			Tss:       tssPubkey,
		})
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		require.ErrorIs(t, err, crosschaintypes.ErrCannotMigrateTssFunds)
		require.ErrorContains(t, err, "cannot migrate funds when there are pending nonces")
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		_, found := k.GetCrossChainTx(ctx, index)
		require.False(t, found)
	})

	t.Run("unable to migrate funds when a pending cctx is presnt in migration info", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		indexString, tssPubkey := setupTssMigrationParams(zk, k, ctx, *chain, amount, true, true)
		k.GetObserverKeeper().SetPendingNonces(ctx, observertypes.PendingNonces{
			NonceLow:  1,
			NonceHigh: 1,
			ChainId:   chain.ChainId,
			Tss:       tssPubkey,
		})
		existingCctx := sample.CrossChainTx(t, "sample_index")
		existingCctx.CctxStatus.Status = crosschaintypes.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, *existingCctx)
		k.GetObserverKeeper().SetFundMigrator(ctx, observertypes.TssFundMigratorInfo{
			ChainId:            chain.ChainId,
			MigrationCctxIndex: existingCctx.Index,
		})
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		require.ErrorIs(t, err, crosschaintypes.ErrCannotMigrateTssFunds)
		require.ErrorContains(t, err, "cannot migrate funds while there are pending migrations")
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		_, found := k.GetCrossChainTx(ctx, index)
		require.False(t, found)
		_, found = k.GetCrossChainTx(ctx, existingCctx.Index)
		require.True(t, found)
	})

	t.Run("unable to migrate funds if current TSS is not present in TSSHistory and no new TSS has been generated", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		indexString, _ := setupTssMigrationParams(zk, k, ctx, *chain, amount, false, false)
		currentTss, found := k.GetObserverKeeper().GetTSS(ctx)
		require.True(t, found)
		newTss := sample.Tss()
		newTss.FinalizedZetaHeight = currentTss.FinalizedZetaHeight - 10
		newTss.KeyGenZetaHeight = currentTss.KeyGenZetaHeight - 10
		k.GetObserverKeeper().SetTSSHistory(ctx, newTss)
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		require.ErrorIs(t, err, crosschaintypes.ErrCannotMigrateTssFunds)
		require.ErrorContains(t, err, "current tss is the latest")
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		_, found = k.GetCrossChainTx(ctx, index)
		require.False(t, found)
	})
}
