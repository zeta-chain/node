package v5_test

import (
	"fmt"
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	v5 "github.com/zeta-chain/zetacore/x/crosschain/migrations/v5"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("TestMigrateStore", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		cctxList := CrossChainTxList(100)
		v5ZetaAccountingAmount := math.ZeroUint()
		v4ZetaAccountingAmount := math.ZeroUint()
		for _, cctx := range cctxList {
			k.SetCrossChainTx(ctx, cctx)
			if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_Aborted || cctx.GetCurrentOutTxParam().CoinType != common.CoinType_Zeta {
				continue
			}
			v5ZetaAccountingAmount = v5ZetaAccountingAmount.Add(crosschainkeeper.GetAbortedAmount(cctx))
			v4ZetaAccountingAmount = v4ZetaAccountingAmount.Add(cctx.GetCurrentOutTxParam().Amount)
		}

		require.True(t, v5ZetaAccountingAmount.GT(v4ZetaAccountingAmount))
		// Previously set the zeta accounting
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{
			AbortedZetaAmount: v4ZetaAccountingAmount,
		})
		err := v5.MigrateStore(ctx, k, k.GetObserverKeeper())
		require.NoError(t, err)
		zetaAccounting, found := k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.True(t, v5ZetaAccountingAmount.Equal(zetaAccounting.AbortedZetaAmount))
		cctxListUpdated := k.GetAllCrossChainTx(ctx)
		// Check refund status of the cctx
		for _, cctx := range cctxListUpdated {
			switch cctx.InboundTxParams.CoinType {
			case common.CoinType_ERC20:
				receiverChain := zk.ObserverKeeper.GetSupportedChainFromChainID(ctx, cctx.GetCurrentOutTxParam().ReceiverChainId)
				require.NotNil(t, receiverChain)
				if receiverChain.IsZetaChain() {
					require.True(t, cctx.CctxStatus.IsAbortRefunded)
				} else {
					require.False(t, cctx.CctxStatus.IsAbortRefunded)
				}
			case common.CoinType_Zeta:
				require.False(t, cctx.CctxStatus.IsAbortRefunded)
			case common.CoinType_Gas:
				require.False(t, cctx.CctxStatus.IsAbortRefunded)
			}
		}
	})

}

func TestResetTestnetNonce(t *testing.T) {
	t.Run("reset only testnet nonce without changing mainnet chains", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		testnetChains := []common.Chain{common.GoerliChain(), common.MumbaiChain(), common.BscTestnetChain(), common.BtcTestNetChain()}
		mainnetChains := []common.Chain{common.EthChain(), common.BscMainnetChain(), common.BtcMainnetChain()}
		nonceLow := int64(1)
		nonceHigh := int64(10)
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)
		for _, chain := range mainnetChains {
			zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{
				Index:   chain.ChainName.String(),
				ChainId: chain.ChainId,
				Nonce:   uint64(nonceHigh),
			})
			zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{
				Tss:       tss.TssPubkey,
				ChainId:   chain.ChainId,
				NonceLow:  nonceLow,
				NonceHigh: nonceHigh,
			})
		}
		for _, chain := range testnetChains {
			zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{
				Tss:       tss.TssPubkey,
				ChainId:   chain.ChainId,
				NonceLow:  nonceLow,
				NonceHigh: nonceHigh,
			})
			zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{
				Index:   chain.ChainName.String(),
				ChainId: chain.ChainId,
				Nonce:   uint64(nonceHigh),
			})
		}
		err := v5.MigrateStore(ctx, k, zk.ObserverKeeper)
		require.NoError(t, err)
		assertValues := map[common.Chain]int64{
			common.GoerliChain():     226841,
			common.MumbaiChain():     200599,
			common.BscTestnetChain(): 110454,
			common.BtcTestNetChain(): 4881,
		}

		for _, chain := range testnetChains {
			pn, found := zk.ObserverKeeper.GetPendingNonces(ctx, tss.TssPubkey, chain.ChainId)
			require.True(t, found)
			require.Equal(t, assertValues[chain], pn.NonceHigh)
			require.Equal(t, assertValues[chain], pn.NonceLow)
			cn, found := zk.ObserverKeeper.GetChainNonces(ctx, chain.ChainName.String())
			require.True(t, found)
			require.Equal(t, uint64(assertValues[chain]), cn.Nonce)
		}
		for _, chain := range mainnetChains {
			pn, found := zk.ObserverKeeper.GetPendingNonces(ctx, tss.TssPubkey, chain.ChainId)
			require.True(t, found)
			require.Equal(t, nonceHigh, pn.NonceHigh)
			require.Equal(t, nonceLow, pn.NonceLow)
			cn, found := zk.ObserverKeeper.GetChainNonces(ctx, chain.ChainName.String())
			require.True(t, found)
			require.Equal(t, uint64(nonceHigh), cn.Nonce)
		}
	})

	t.Run("reset nonce even if some chain values are missing", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		testnetChains := []common.Chain{common.GoerliChain()}
		nonceLow := int64(1)
		nonceHigh := int64(10)
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)
		for _, chain := range testnetChains {
			zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{
				Tss:       tss.TssPubkey,
				ChainId:   chain.ChainId,
				NonceLow:  nonceLow,
				NonceHigh: nonceHigh,
			})
			zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{
				Index:   chain.ChainName.String(),
				ChainId: chain.ChainId,
				Nonce:   uint64(nonceHigh),
			})
		}
		err := v5.MigrateStore(ctx, k, zk.ObserverKeeper)
		require.NoError(t, err)
		assertValuesSet := map[common.Chain]int64{
			common.GoerliChain(): 226841,
		}
		assertValuesNotSet := []common.Chain{common.MumbaiChain(), common.BscTestnetChain(), common.BtcTestNetChain()}

		for _, chain := range testnetChains {
			pn, found := zk.ObserverKeeper.GetPendingNonces(ctx, tss.TssPubkey, chain.ChainId)
			require.True(t, found)
			require.Equal(t, assertValuesSet[chain], pn.NonceHigh)
			require.Equal(t, assertValuesSet[chain], pn.NonceLow)
			cn, found := zk.ObserverKeeper.GetChainNonces(ctx, chain.ChainName.String())
			require.True(t, found)
			require.Equal(t, uint64(assertValuesSet[chain]), cn.Nonce)
		}
		for _, chain := range assertValuesNotSet {
			_, found := zk.ObserverKeeper.GetPendingNonces(ctx, tss.TssPubkey, chain.ChainId)
			require.False(t, found)
			_, found = zk.ObserverKeeper.GetChainNonces(ctx, chain.ChainName.String())
			require.False(t, found)
		}
	})
}

func CrossChainTxList(count int) []crosschaintypes.CrossChainTx {
	cctxList := make([]crosschaintypes.CrossChainTx, count+100)
	i := 0
	r := rand.New(rand.NewSource(9))
	for ; i < count/2; i++ {
		amount := math.NewUint(uint64(r.Uint32()))
		cctxList[i] = crosschaintypes.CrossChainTx{
			Index:      fmt.Sprintf("%d", i),
			CctxStatus: &crosschaintypes.Status{Status: crosschaintypes.CctxStatus_Aborted},
			InboundTxParams: &crosschaintypes.InboundTxParams{
				Amount:   amount.Add(math.NewUint(uint64(r.Uint32()))),
				CoinType: common.CoinType_Zeta,
			},
			OutboundTxParams: []*crosschaintypes.OutboundTxParams{{
				Amount:   amount,
				CoinType: common.CoinType_Zeta,
			}},
		}
		for ; i < count; i++ {
			amount := math.NewUint(uint64(r.Uint32()))
			cctxList[i] = crosschaintypes.CrossChainTx{
				Index:      fmt.Sprintf("%d", i),
				CctxStatus: &crosschaintypes.Status{Status: crosschaintypes.CctxStatus_Aborted},
				InboundTxParams: &crosschaintypes.InboundTxParams{
					Amount:   amount,
					CoinType: common.CoinType_Zeta,
				},
				OutboundTxParams: []*crosschaintypes.OutboundTxParams{{
					Amount:   math.ZeroUint(),
					CoinType: common.CoinType_Zeta,
				}},
			}
		}
		for ; i < count+20; i++ {
			amount := math.NewUint(uint64(r.Uint32()))
			cctxList[i] = crosschaintypes.CrossChainTx{
				Index:      fmt.Sprintf("%d", i),
				CctxStatus: &crosschaintypes.Status{Status: crosschaintypes.CctxStatus_Aborted},
				InboundTxParams: &crosschaintypes.InboundTxParams{
					Amount:   amount,
					CoinType: common.CoinType_ERC20,
				},
				OutboundTxParams: []*crosschaintypes.OutboundTxParams{{
					Amount:          math.ZeroUint(),
					CoinType:        common.CoinType_ERC20,
					ReceiverChainId: common.ZetaPrivnetChain().ChainId,
				}},
			}
		}
		for ; i < count+50; i++ {
			amount := math.NewUint(uint64(r.Uint32()))
			cctxList[i] = crosschaintypes.CrossChainTx{
				Index:      fmt.Sprintf("%d", i),
				CctxStatus: &crosschaintypes.Status{Status: crosschaintypes.CctxStatus_Aborted},
				InboundTxParams: &crosschaintypes.InboundTxParams{
					Amount:   amount,
					CoinType: common.CoinType_ERC20,
				},
				OutboundTxParams: []*crosschaintypes.OutboundTxParams{{
					Amount:          math.ZeroUint(),
					CoinType:        common.CoinType_ERC20,
					ReceiverChainId: common.GoerliLocalnetChain().ChainId,
				}},
			}
		}
		for ; i < count+100; i++ {
			amount := math.NewUint(uint64(r.Uint32()))
			cctxList[i] = crosschaintypes.CrossChainTx{
				Index:      fmt.Sprintf("%d", i),
				CctxStatus: &crosschaintypes.Status{Status: crosschaintypes.CctxStatus_Aborted},
				InboundTxParams: &crosschaintypes.InboundTxParams{
					Amount:   amount,
					CoinType: common.CoinType_Gas,
				},
				OutboundTxParams: []*crosschaintypes.OutboundTxParams{{
					Amount:   amount,
					CoinType: common.CoinType_Gas,
				}},
			}
		}
	}
	return cctxList
}
