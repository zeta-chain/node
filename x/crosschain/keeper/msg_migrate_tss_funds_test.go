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
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_MigrateTssFunds(t *testing.T) {
	t.Run("successfully create tss migration cctx", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain(t)
		amount := sdkmath.NewUint(100)
		indexString, _ := setupTssMigrationParams(zk, k, ctx, *chain, amount, true)
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		assert.NoError(t, err)
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		_, found := k.GetCrossChainTx(ctx, index)
		assert.True(t, found)
	})
	t.Run("unable to migrate funds if new TSS is not created ", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain(t)
		amount := sdkmath.NewUint(100)
		indexString, _ := setupTssMigrationParams(zk, k, ctx, *chain, amount, false)
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: chain.ChainId,
			Amount:  amount,
		})
		assert.ErrorIs(t, err, crosschaintypes.ErrUnableToUpdateTss)
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		_, found := k.GetCrossChainTx(ctx, index)
		assert.False(t, found)
	})
	t.Run("unable to migrate funds pending cctx is present", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		chain := getValidEthChain(t)
		amount := sdkmath.NewUint(100)
		indexString, tssPubkey := setupTssMigrationParams(zk, k, ctx, *chain, amount, false)
		k.SetPendingNonces(ctx, crosschaintypes.PendingNonces{
			NonceLow:  1,
			NonceHigh: 10,
			ChainId:   chain.ChainId,
			Tss:       tssPubkey,
		})
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: 1,
			Amount:  amount,
		})
		assert.ErrorIs(t, err, crosschaintypes.ErrUnableToUpdateTss)
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		_, found := k.GetCrossChainTx(ctx, index)
		assert.False(t, found)
	})
}
func setupTssMigrationParams(zk keepertest.ZetaKeepers, k *keeper.Keeper, ctx sdk.Context, chain common.Chain, amount sdkmath.Uint, setNewTss bool) (string, string) {
	zk.ObserverKeeper.SetCrosschainFlags(ctx, observerTypes.CrosschainFlags{
		IsInboundEnabled:  false,
		IsOutboundEnabled: true,
	})
	params := zk.ObserverKeeper.GetParamsIfExists(ctx)
	params.ObserverParams = append(params.ObserverParams, &observerTypes.ObserverParams{
		Chain:                 &chain,
		BallotThreshold:       sdk.NewDec(0),
		MinObserverDelegation: sdk.OneDec(),
		IsSupported:           true,
	})
	zk.ObserverKeeper.SetParams(ctx, params)
	currentTss := sample.Tss()
	newTss := sample.Tss()
	newTss.FinalizedZetaHeight = currentTss.FinalizedZetaHeight + 1
	newTss.KeyGenZetaHeight = currentTss.KeyGenZetaHeight + 1
	k.GetObserverKeeper().SetTSS(ctx, currentTss)
	k.GetObserverKeeper().SetTSSHistory(ctx, currentTss)
	if setNewTss {
		k.GetObserverKeeper().SetTSSHistory(ctx, newTss)
	}
	k.SetPendingNonces(ctx, crosschaintypes.PendingNonces{
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
		Prices:      []uint64{1, 1, 1},
		MedianIndex: 1,
	})
	k.SetChainNonces(ctx, crosschaintypes.ChainNonces{
		Index:   chain.ChainName.String(),
		ChainId: chain.ChainId,
		Nonce:   1,
	})
	indexString := keeper.GetIndexStringForTssMigration(currentTss.TssPubkey, newTss.TssPubkey, chain.ChainId, amount, ctx.BlockHeight())
	return indexString, currentTss.TssPubkey
}
