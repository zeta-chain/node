package e2etests

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	testcontract "github.com/zeta-chain/node/testutil/contracts"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestSolanaDepositSPL(r *runner.E2ERunner, _ []string) {
	// require.Len(r, args, 1)

	// deploy an example contract in ZEVM
	contractAddr, _, contract, err := testcontract.DeployExample(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Example contract deployed at: %s", contractAddr.String())

	// load deployer private key
	privKey, err := solana.PrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
	require.NoError(r, err)

	wall, err := solana.WalletFromPrivateKeyBase58(r.SPLPrivateKey.String())
	require.NoError(r, err)

	// execute the deposit transaction
	data := []byte("hello lamports")
	sig := r.DepositSPL(&privKey, *wall, contractAddr, data)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "solana_deposit")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// check if example contract has been called, bar value should be set to amount
	utils.MustHaveCalledExampleContract(r, contract, big.NewInt(500_000))
}
