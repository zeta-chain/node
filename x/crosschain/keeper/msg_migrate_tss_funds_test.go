package keeper_test

import (
	"fmt"
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
		chain := common.Chain{
			ChainId:   1,
			ChainName: common.ChainName_eth_mainnet,
		}
		amount := sdkmath.NewUint(100)
		indexString := setupTssMigrationParams(zk, k, ctx, chain, amount)
		_, err := msgServer.MigrateTssFunds(ctx, &crosschaintypes.MsgMigrateTssFunds{
			Creator: admin,
			ChainId: 1,
			Amount:  amount,
		})
		assert.NoError(t, err)
		hash := crypto.Keccak256Hash([]byte(indexString))
		index := hash.Hex()
		_, found := k.GetCrossChainTx(ctx, index)
		assert.True(t, found)

	})
}
func setupTssMigrationParams(zk keepertest.ZetaKeepers, k *keeper.Keeper, ctx sdk.Context, chain common.Chain, amount sdkmath.Uint) string {
	zk.ObserverKeeper.SetCrosschainFlags(ctx, observerTypes.CrosschainFlags{
		IsInboundEnabled:  false,
		IsOutboundEnabled: true,
	})
	currentTss := sample.Tss()
	newTss := sample.Tss()
	newTss.FinalizedZetaHeight = currentTss.FinalizedZetaHeight + 1
	newTss.KeyGenZetaHeight = currentTss.KeyGenZetaHeight + 1
	k.SetTSS(ctx, *currentTss)
	k.SetTSSHistory(ctx, *currentTss)
	k.SetTSSHistory(ctx, *newTss)
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
	indexString := fmt.Sprintf("%s-%s-%d-%s-%d", currentTss.TssPubkey, newTss.TssPubkey, chain.ChainId, amount.String(), ctx.BlockHeight())
	return indexString
}
