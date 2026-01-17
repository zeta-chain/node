package observer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/testutil/sample"
	cc "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/ton/encoder"
	"github.com/zeta-chain/node/zetaclient/chains/ton/repo"
	"github.com/zeta-chain/node/zetaclient/testutils"
)

func TestOutbound(t *testing.T) {
	gw := toncontracts.NewGateway(
		ton.MustParseAccountID(testutils.GatewayAddresses[chains.TONMainnet.ChainId]),
	)

	t.Run("observeOutboundTrackers", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		tonRepo := repo.NewTONRepo(ts.rpc, gw, ts.baseObserver.Chain())
		ob, err := New(ts.baseObserver, tonRepo, gw)
		require.NoError(t, err)

		// Given withdrawal
		withdrawal := toncontracts.Withdrawal{
			Recipient: ton.MustParseAccountID("0:552f6db5da0cae7f0b3ab4ab58d85927f6beb962cda426a6a6ee751c82cead1f"),
			Amount:    toncontracts.Coins(2),
			Seqno:     3,
		}
		ts.sign(&withdrawal)

		nonce := uint64(withdrawal.Seqno)

		// Given TON tx
		withdrawalTX := sample.TONWithdrawal(t, gw.AccountID(), withdrawal)

		ts.MockGetTransaction(gw.AccountID(), withdrawalTX)

		// Given outbound tracker
		tracker := cc.OutboundTracker{
			Index:    "index123",
			ChainId:  ts.chain.ChainId,
			Nonce:    nonce,
			HashList: []*cc.TxHash{{TxHash: encoder.EncodeTx(withdrawalTX)}},
		}

		ts.OnGetOutboundTrackers([]cc.OutboundTracker{tracker})

		// Given cctx
		cctx := sample.CrossChainTx(t, "index456")
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		cctx.GetCurrentOutboundParam().TssNonce = nonce

		ts.MockCCTXByNonce(cctx)

		// ACT
		err = ob.ProcessOutboundTrackers(ts.ctx)

		// ASSERT
		require.NoError(t, err)

		// Check that tx exists in outbounds
		res := ob.getOutbound(nonce)
		assert.NotNil(t, res)

		assert.Equal(t, nonce, res.nonce)
		assert.Equal(t, chains.ReceiveStatus_success, res.receiveStatus)
		assert.Equal(t, true, res.tx.IsSuccess())
		assert.Equal(t, int32(0), res.tx.ExitCode)

		w2, err := res.tx.Withdrawal()
		assert.NoError(t, err)
		assert.Equal(t, withdrawal, w2)
	})
}
