package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	solanacontract "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestZEVMToSolanaCall executes simple call from ZEVM to Solana, calling connected program
func TestZEVMToSolanaCall(r *runner.E2ERunner, _ []string) {
	// approve amount is 1 SOL
	approvedAmount := new(big.Int).SetUint64(solana.LAMPORTS_PER_SOL)

	r.ZEVMAuth.GasLimit = 10000000
	// withdraw and call
	tx := r.CallSOLZRC20(
		r.ConnectedProgram,
		approvedAmount,
		[]byte("simple call"),
		gatewayzevm.RevertOptions{
			OnRevertGasLimit: big.NewInt(0),
		},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// check pda account info of connected program
	connectedPda, err := solanacontract.ComputeConnectedPdaAddress(r.ConnectedProgram)
	require.NoError(r, err)

	pda := r.ParseConnectedPda(connectedPda)
	require.Equal(r, "simple call", pda.LastMessage)
	require.Equal(r, r.ZEVMAuth.From.String(), common.BytesToAddress(pda.LastSender[:]).String())
}
