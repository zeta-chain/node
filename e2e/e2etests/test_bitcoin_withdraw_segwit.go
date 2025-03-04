package e2etests

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestBitcoinWithdrawSegWit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// parse arguments
	defaultReceiver := r.BTCDeployerAddress.EncodeAddress()
	receiver, amount := utils.ParseBitcoinWithdrawArgs(r, args, defaultReceiver, r.GetBitcoinChainID())
	_, ok := receiver.(*btcutil.AddressWitnessPubKeyHash)
	require.True(r, ok, "Invalid receiver address specified for TestBitcoinWithdrawSegWit.")

	// perform the withdraw
	tx := withdrawBTCZRC20(r, receiver, amount)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
}
