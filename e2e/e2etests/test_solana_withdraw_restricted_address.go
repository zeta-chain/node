package e2etests

import (
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestSolanaWithdrawRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// parse restricted address
	receiverRestricted, err := chains.DecodeSolanaWalletAddress(args[0])
	require.NoError(r, err, fmt.Sprintf("unable to decode solana wallet address: %s", args[0]))

	// parse withdraw amount (in lamports), approve amount is 1 SOL
	approvedAmount := new(big.Int).SetUint64(solana.LAMPORTS_PER_SOL)
	withdrawAmount := utils.ParseBigInt(r, args[1])
	require.Equal(
		r,
		-1,
		withdrawAmount.Cmp(approvedAmount),
		"Withdrawal amount must be less than the approved amount (1e9).",
	)

	// withdraw
	tx := r.WithdrawSOLZRC20(receiverRestricted, withdrawAmount, approvedAmount)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// the cctx should be cancelled with zero value
	r.VerifySolanaWithdrawalAmountFromCCTX(cctx, 0)
}
