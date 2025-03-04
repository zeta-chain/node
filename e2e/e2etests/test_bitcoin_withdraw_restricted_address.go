package e2etests

import (
	"math/big"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestBitcoinWithdrawRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	withdrawalAmount := utils.ParseFloat(r, args[0])
	amount := utils.BTCAmountFromFloat64(r, withdrawalAmount)

	withdrawBitcoinRestricted(r, amount)
}

func withdrawBitcoinRestricted(r *runner.E2ERunner, amount *big.Int) {
	// use restricted BTC P2WPKH address
	addressRestricted, err := chains.DecodeBtcAddress(
		sample.RestrictedBtcAddressTest,
		chains.BitcoinRegtest.ChainId,
	)
	require.NoError(r, err)

	// perform the withdraw
	tx := withdrawBTCZRC20(r, addressRestricted, amount)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// get bitcoin tx hash from cctx struct
	outTxHash := cctx.GetCurrentOutboundParam().Hash
	hash, err := chainhash.NewHashFromStr(outTxHash)
	require.NoError(r, err)

	// the cctx should be cancelled
	rawTx, err := r.BtcRPCClient.GetRawTransactionVerbose(r.Ctx, hash)
	require.NoError(r, err)
	require.Len(r, rawTx.Vout, 2, "BTC cancelled outtx rawTx.Vout should have 2 outputs")
}
