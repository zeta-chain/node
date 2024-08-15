package e2etests

import (
	"fmt"
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/x/crosschain/types"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

func TestZetaDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestZetaDeposit.")

	hash := r.DepositZetaWithAmount(r.EVMAddress().Bytes(), amount)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, hash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
}

func TestZetaDepositToInvalidAddress(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestZetaDeposit.")

	hash := r.DepositZetaWithAmount([]byte("invalid"), amount)

	r.Logger.Print(fmt.Sprintf("Deposit to invalid address tx hass: %s", hash.Hex()))

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, hash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	r.Logger.Print(fmt.Sprintf("CCTX: %v", cctx.Index))
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_Reverted, "deposit with invalid address")
	utils.ContainsStringInCCTXStatusMessage(r, cctx, types.ErrInvalidReceiverAddress.Error(), "deposit with invalid address")
}
