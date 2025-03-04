package e2etests

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestBitcoinWithdrawP2SH(r *runner.E2ERunner, args []string) {
	// check length of arguments
	require.Len(r, args, 2)

	// parse arguments and withdraw BTC
	defaultReceiver := "2N6AoUj3KPS7wNGZXuCckh8YEWcSYNsGbqd"
	receiver, amount := utils.ParseBitcoinWithdrawArgs(r, args, defaultReceiver, r.GetBitcoinChainID())
	_, ok := receiver.(*btcutil.AddressScriptHash)
	require.True(r, ok, "Invalid receiver address specified for TestBitcoinWithdrawP2SH.")

	// perform the withdraw
	tx := withdrawBTCZRC20(r, receiver, amount)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
}
