package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/gas"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func setupTssMigrationParams(
	zk keepertest.ZetaKeepers,
	k *keeper.Keeper,
	ctx sdk.Context,
	chain chains.Chain,
	amount sdkmath.Uint,
	setNewTss bool,
	setCurrentTSS bool,
	setChainNonces bool,
) (string, string) {
	zk.ObserverKeeper.SetCrosschainFlags(ctx, observertypes.CrosschainFlags{
		IsInboundEnabled:  false,
		IsOutboundEnabled: true,
	})

	zk.ObserverKeeper.SetChainParamsList(ctx, observertypes.ChainParamsList{
		ChainParams: []*observertypes.ChainParams{
			{
				ChainId:               chain.ChainId,
				BallotThreshold:       sdkmath.LegacyNewDec(0),
				MinObserverDelegation: sdkmath.LegacyOneDec(),
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
		Creator:      "",
		Index:        "",
		ChainId:      chain.ChainId,
		Signers:      nil,
		BlockNums:    nil,
		Prices:       []uint64{100000, 100000, 100000},
		PriorityFees: []uint64{100, 300, 200},
		MedianIndex:  1,
	})
	if setChainNonces {
		k.GetObserverKeeper().SetChainNonces(ctx, observertypes.ChainNonces{
			ChainId: chain.ChainId,
			Nonce:   1,
		})
	}
	indexString := crosschaintypes.GetTssMigrationCCTXIndexString(
		currentTss.TssPubkey,
		newTss.TssPubkey,
		chain.ChainId,
		amount,
		ctx.BlockHeight(),
	)
	return indexString, currentTss.TssPubkey
}

func TestKeeper_MigrateTSSFundsForChain(t *testing.T) {
	t.Run("test evm chain", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)

		indexString, _ := setupTssMigrationParams(zk, k, ctx, chain, amount, true, true, true)
		gp, priorityFee, found := k.GetMedianGasValues(ctx, chain.ChainId)
		require.True(t, found)
		msg := crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)
		_, err := msgServer.MigrateTssFunds(ctx, &msg)
		require.NoError(t, err)

		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		cctx, found := k.GetCrossChainTx(ctx, index)
		require.True(t, found)

		multipliedValue, err := gas.MultiplyGasPrice(gp, crosschaintypes.TssMigrationGasMultiplierEVM)
		require.NoError(t, err)
		require.Equal(t, multipliedValue.String(), cctx.GetCurrentOutboundParam().GasPrice)
		require.Equal(t, priorityFee.MulUint64(2).String(), cctx.GetCurrentOutboundParam().GasPriorityFee)
	})

	t.Run("test btc chain", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chain := getValidBTCChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)
		indexString, _ := setupTssMigrationParams(zk, k, ctx, chain, amount, true, true, true)
		gp, priorityFee, found := k.GetMedianGasValues(ctx, chain.ChainId)
		require.True(t, found)

		msg := crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)
		_, err := msgServer.MigrateTssFunds(ctx, &msg)
		require.NoError(t, err)
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		cctx, found := k.GetCrossChainTx(ctx, index)
		require.True(t, found)
		require.Equal(t, gp.MulUint64(2).String(), cctx.GetCurrentOutboundParam().GasPrice)
		require.Equal(t, priorityFee.MulUint64(2).String(), cctx.GetCurrentOutboundParam().GasPriorityFee)
	})
}

func TestMsgServer_MigrateTssFunds(t *testing.T) {
	t.Run("should error if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)

		msg := crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.MigrateTssFunds(ctx, &msg)
		require.Error(t, err)
	})

	t.Run("should error if inbound enabled", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})

		admin := sample.AccAddress()
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("IsInboundEnabled", mock.Anything).Return(true)
		msgServer := keeper.NewMsgServerImpl(*k)

		msg := crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.MigrateTssFunds(ctx, &msg)
		require.Error(t, err)
	})

	t.Run("should error if tss not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})

		admin := sample.AccAddress()
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("IsInboundEnabled", mock.Anything).Return(false)
		observerMock.On("GetTSS", mock.Anything).Return(observertypes.TSS{}, false)

		msgServer := keeper.NewMsgServerImpl(*k)

		msg := crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.MigrateTssFunds(ctx, &msg)
		require.Error(t, err)
	})

	t.Run("should error if tss history empty", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})

		admin := sample.AccAddress()
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("IsInboundEnabled", mock.Anything).Return(false)
		observerMock.On("GetTSS", mock.Anything).Return(sample.Tss(), true)
		observerMock.On("GetAllTSS", mock.Anything).Return([]observertypes.TSS{})

		msgServer := keeper.NewMsgServerImpl(*k)

		msg := crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.MigrateTssFunds(ctx, &msg)
		require.Error(t, err)
	})

	t.Run("should error if no new tss generated", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})

		admin := sample.AccAddress()
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("IsInboundEnabled", mock.Anything).Return(false)
		tss := sample.Tss()
		observerMock.On("GetTSS", mock.Anything).Return(tss, true)
		observerMock.On("GetAllTSS", mock.Anything).Return([]observertypes.TSS{tss})

		msgServer := keeper.NewMsgServerImpl(*k)

		msg := crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.MigrateTssFunds(ctx, &msg)
		require.Error(t, err)
	})

	t.Run("should error if current tss is the latest", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})

		admin := sample.AccAddress()
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("IsInboundEnabled", mock.Anything).Return(false)
		tss1 := sample.Tss()
		tss1.FinalizedZetaHeight = 2
		tss2 := sample.Tss()
		tss2.FinalizedZetaHeight = 1
		observerMock.On("GetTSS", mock.Anything).Return(tss1, true)
		observerMock.On("GetAllTSS", mock.Anything).Return([]observertypes.TSS{tss2})

		msgServer := keeper.NewMsgServerImpl(*k)

		msg := crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.MigrateTssFunds(ctx, &msg)
		require.Error(t, err)
	})

	t.Run("should error if pending nonces not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})

		admin := sample.AccAddress()
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("IsInboundEnabled", mock.Anything).Return(false)
		tss1 := sample.Tss()
		tss1.FinalizedZetaHeight = 2
		tss2 := sample.Tss()
		tss2.FinalizedZetaHeight = 3
		observerMock.On("GetTSS", mock.Anything).Return(tss1, true)
		observerMock.On("GetAllTSS", mock.Anything).Return([]observertypes.TSS{tss2})
		observerMock.On("GetPendingNonces", mock.Anything, mock.Anything, mock.Anything).
			Return(observertypes.PendingNonces{}, false)

		msgServer := keeper.NewMsgServerImpl(*k)

		msg := crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.MigrateTssFunds(ctx, &msg)
		require.Error(t, err)
	})

	t.Run("successfully create tss migration cctx", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)

		indexString, _ := setupTssMigrationParams(zk, k, ctx, chain, amount, true, true, true)

		msg := crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)
		_, err := msgServer.MigrateTssFunds(ctx, &msg)
		require.NoError(t, err)
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		cctx, found := k.GetCrossChainTx(ctx, index)
		require.True(t, found)
		feeCalculated := sdkmath.NewUint(cctx.GetCurrentOutboundParam().CallOptions.GasLimit).
			Mul(sdkmath.NewUintFromString(cctx.GetCurrentOutboundParam().GasPrice)).
			Add(sdkmath.NewUintFromString(crosschaintypes.TSSMigrationBufferAmountEVM))
		require.Equal(t, cctx.GetCurrentOutboundParam().Amount.String(), amount.Sub(feeCalculated).String())
	})

	t.Run("not enough funds in tss address for migration", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("100")

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		msgServer := keeper.NewMsgServerImpl(*k)
		indexString, _ := setupTssMigrationParams(zk, k, ctx, chain, amount, true, true, true)

		msg := crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)
		_, err := msgServer.MigrateTssFunds(ctx, &msg)
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
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		msgServer := keeper.NewMsgServerImpl(*k)
		indexString, _ := setupTssMigrationParams(zk, k, ctx, chain, amount, false, true, true)

		msg := crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.MigrateTssFunds(ctx, &msg)
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
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		msgServer := keeper.NewMsgServerImpl(*k)

		indexString, tssPubkey := setupTssMigrationParams(zk, k, ctx, chain, amount, true, true, true)
		k.GetObserverKeeper().SetPendingNonces(ctx, observertypes.PendingNonces{
			NonceLow:  1,
			NonceHigh: 10,
			ChainId:   chain.ChainId,
			Tss:       tssPubkey,
		})

		msg := crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.MigrateTssFunds(ctx, &msg)
		require.ErrorIs(t, err, crosschaintypes.ErrCannotMigrateTssFunds)
		require.ErrorContains(t, err, "cannot migrate funds when there are pending nonces")
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		_, found := k.GetCrossChainTx(ctx, index)
		require.False(t, found)
	})

	t.Run("unable to migrate funds when a pending cctx is present in migration info", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		msgServer := keeper.NewMsgServerImpl(*k)

		indexString, tssPubkey := setupTssMigrationParams(zk, k, ctx, chain, amount, true, true, true)
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

		msg := crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)
		_, err := msgServer.MigrateTssFunds(ctx, &msg)
		require.ErrorIs(t, err, crosschaintypes.ErrCannotMigrateTssFunds)
		require.ErrorContains(t, err, "cannot migrate funds while there are pending migrations")
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		_, found := k.GetCrossChainTx(ctx, index)
		require.False(t, found)
		_, found = k.GetCrossChainTx(ctx, existingCctx.Index)
		require.True(t, found)
	})

	t.Run(
		"unable to migrate funds if current TSS is not present in TSSHistory and no new TSS has been generated",
		func(t *testing.T) {
			k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
				UseAuthorityMock: true,
			})

			admin := sample.AccAddress()
			chain := getValidEthChain()
			amount := sdkmath.NewUintFromString("10000000000000000000")
			authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

			msgServer := keeper.NewMsgServerImpl(*k)

			indexString, _ := setupTssMigrationParams(zk, k, ctx, chain, amount, false, false, true)
			currentTss, found := k.GetObserverKeeper().GetTSS(ctx)
			require.True(t, found)
			newTss := sample.Tss()
			newTss.FinalizedZetaHeight = currentTss.FinalizedZetaHeight - 10
			newTss.KeyGenZetaHeight = currentTss.KeyGenZetaHeight - 10
			k.GetObserverKeeper().SetTSSHistory(ctx, newTss)

			msg := crosschaintypes.MsgMigrateTssFunds{
				Creator: admin,
				ChainId: chain.ChainId,
				Amount:  amount,
			}

			keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
			_, err := msgServer.MigrateTssFunds(ctx, &msg)
			require.ErrorIs(t, err, crosschaintypes.ErrCannotMigrateTssFunds)
			require.ErrorContains(t, err, "current tss is the latest")
			hash := crypto.Keccak256Hash([]byte(indexString))
			index := hash.Hex()
			_, found = k.GetCrossChainTx(ctx, index)
			require.False(t, found)
		},
	)

	t.Run("unable to process migration if SetObserverOutboundInfo fails", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chain := getValidEthChain()
		amount := sdkmath.NewUintFromString("10000000000000000000")
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)

		_, _ = setupTssMigrationParams(zk, k, ctx, chain, amount, true, true, false)
		msg := crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		_, err := msgServer.MigrateTssFunds(ctx, &msg)
		require.ErrorContains(t, err, crosschaintypes.ErrUnableToSetOutboundInfo.Error())
	})
}
