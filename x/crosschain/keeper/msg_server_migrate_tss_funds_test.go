package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_MigrateTSSFundsForChain(t *testing.T) {
	t.Run("test gas price multiplier", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain(t)
		amount := sdkmath.NewUintFromString("10000000000000000000")
		indexString, _ := setupTssMigrationParams(zk, k, ctx, *chain, amount, true, true)
		gp, found := k.GetMedianGasPriceInUint(ctx, chain.ChainId)
		assert.True(t, found)
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		assert.NoError(t, err)
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		cctx, found := k.GetCrossChainTx(ctx, index)
		assert.True(t, found)
		multipliedValue, err := common.MultiplyGasPrice(gp, crosschaintypes.TssMigrationGasMultiplierEVM)
		assert.NoError(t, err)
		assert.Equal(t, multipliedValue.String(), cctx.GetCurrentOutTxParam().OutboundTxGasPrice)

	})
}
func TestMsgServer_MigrateTssFunds(t *testing.T) {
	t.Run("successfully create tss migration cctx", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain(t)
		amount := sdkmath.NewUintFromString("10000000000000000000")
		indexString, _ := setupTssMigrationParams(zk, k, ctx, *chain, amount, true, true)
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		assert.NoError(t, err)
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		cctx, found := k.GetCrossChainTx(ctx, index)
		assert.True(t, found)
		feeCalculated := sdk.NewUint(cctx.GetCurrentOutTxParam().OutboundTxGasLimit).
			Mul(sdkmath.NewUintFromString(cctx.GetCurrentOutTxParam().OutboundTxGasPrice))
		assert.Equal(t, cctx.GetCurrentOutTxParam().Amount.String(), amount.Sub(feeCalculated).String())
	})
	t.Run("not enough funds in tss address for migration", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain(t)
		amount := sdkmath.NewUintFromString("100")
		indexString, _ := setupTssMigrationParams(zk, k, ctx, *chain, amount, true, true)
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		assert.ErrorContains(t, err, crosschaintypes.ErrCannotMigrateTssFunds.Error())
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		_, found := k.GetCrossChainTx(ctx, index)
		assert.False(t, found)
	})
	t.Run("unable to migrate funds if new TSS is not created ", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain(t)
		amount := sdkmath.NewUintFromString("10000000000000000000")
		indexString, _ := setupTssMigrationParams(zk, k, ctx, *chain, amount, false, true)
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		assert.ErrorContains(t, err, "no new tss address has been generated")
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		_, found := k.GetCrossChainTx(ctx, index)
		assert.False(t, found)
	})
	t.Run("unable to migrate funds when nonce low does not match nonce high", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain(t)
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
		assert.ErrorIs(t, err, crosschaintypes.ErrCannotMigrateTssFunds)
		assert.ErrorContains(t, err, "cannot migrate funds when there are pending nonces")
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		_, found := k.GetCrossChainTx(ctx, index)
		assert.False(t, found)
	})
	t.Run("unable to migrate funds when a pending cctx is presnt in migration info", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain(t)
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
		assert.ErrorIs(t, err, crosschaintypes.ErrCannotMigrateTssFunds)
		assert.ErrorContains(t, err, "cannot migrate funds while there are pending migrations")
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		_, found := k.GetCrossChainTx(ctx, index)
		assert.False(t, found)
		_, found = k.GetCrossChainTx(ctx, existingCctx.Index)
		assert.True(t, found)
	})

	t.Run("unable to migrate funds if current TSS is not present in TSSHistory and no new TSS has been generated", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain(t)
		amount := sdkmath.NewUintFromString("10000000000000000000")
		indexString, _ := setupTssMigrationParams(zk, k, ctx, *chain, amount, false, false)
		currentTss, found := k.GetObserverKeeper().GetTSS(ctx)
		assert.True(t, found)
		newTss := sample.Tss()
		newTss.FinalizedZetaHeight = currentTss.FinalizedZetaHeight - 10
		newTss.KeyGenZetaHeight = currentTss.KeyGenZetaHeight - 10
		k.GetObserverKeeper().SetTSSHistory(ctx, newTss)
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		assert.ErrorIs(t, err, crosschaintypes.ErrCannotMigrateTssFunds)
		assert.ErrorContains(t, err, "current tss is the latest")
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		_, found = k.GetCrossChainTx(ctx, index)
		assert.False(t, found)
	})
}
func setupTssMigrationParams(
	zk keepertest.ZetaKeepers,
	k *keeper.Keeper,
	ctx sdk.Context,
	chain common.Chain,
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
