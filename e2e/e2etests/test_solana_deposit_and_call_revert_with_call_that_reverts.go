package e2etests

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"

	testcontract "github.com/zeta-chain/node/e2e/contracts/reverter"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	solanacontracts "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSolanaDepositAndCallRevertWithCallThatReverts tests deposit of lamports
// with revert options when call on revert program reverts, and cctx is aborted
func TestSolanaDepositAndCallRevertWithCallThatReverts(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse deposit amount (in lamports)
	depositAmount := utils.ParseBigInt(r, args[0])

	// deploy a reverter contract in ZEVM
	// TODO: consider removing repeated deployments of reverter contract
	reverterAddr, _, _, err := testcontract.DeployReverter(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Reverter contract deployed at: %s", reverterAddr.String())

	// execute the deposit transaction
	data := []byte("hello deposit and call revert")

	// check balances before deposit
	connectedPda, err := solanacontracts.ComputeConnectedPdaAddress(r.ConnectedProgram)
	require.NoError(r, err)
	connectedPdaInfoBefore, err := r.SolanaClient.GetAccountInfo(r.Ctx, connectedPda)
	require.NoError(r, err)

	// create encoded msg
	msg := solanacontracts.ExecuteMsg{
		Accounts: []solanacontracts.AccountMeta{
			{PublicKey: [32]byte(connectedPda.Bytes()), IsWritable: true},
			{PublicKey: [32]byte(r.ComputePdaAddress().Bytes()), IsWritable: false},
			{PublicKey: [32]byte(solana.SystemProgramID.Bytes()), IsWritable: false},
		},
		Data: data,
	}

	msgEncoded, err := msg.Encode()
	require.NoError(r, err)

	sig := r.SOLDepositAndCall(nil, reverterAddr, depositAmount, data, &solanacontracts.RevertOptions{
		RevertAddress: r.ConnectedProgram,
		CallOnRevert:  true,
		RevertMessage: msgEncoded,
	})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "solana_deposit_and_refund")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Aborted)

	require.Contains(r, cctx.CctxStatus.ErrorMessage, utils.ErrHashRevertFoo)

	// verify that pda of solana connected program balance is not changed
	connectedPdaInfo, err := r.SolanaClient.GetAccountInfo(r.Ctx, connectedPda)
	require.NoError(r, err)
	type ConnectedPdaInfo struct {
		Discriminator     [8]byte
		LastSender        ethcommon.Address
		LastMessage       string
		LastRevertSender  solana.PublicKey
		LastRevertMessage string
	}
	pda := ConnectedPdaInfo{}
	err = borsh.Deserialize(&pda, connectedPdaInfo.Bytes())
	require.NoError(r, err)

	require.Equal(r, connectedPdaInfoBefore.Value.Lamports, connectedPdaInfo.Value.Lamports)
}
