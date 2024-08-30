package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/txserver"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestPauseERC20Custody tests the pausing and unpausing of ERC20 custody contracts on the EVM chain
func TestPauseERC20Custody(r *runner.E2ERunner, _ []string) {
	// get EVM chain ID
	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	// check ERC20 custody contract is not paused
	paused, err := r.ERC20Custody.Paused(&bind.CallOpts{})
	require.NoError(r, err)
	require.False(r, paused)

	// Part 1: Pause ERC20 custody contract

	// send command for pausing ERC20 custody contract
	msg := crosschaintypes.NewMsgUpdateERC20CustodyPauseStatus(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		chainID.Int64(),
		true,
	)
	res, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err)

	// fetch cctx index from tx response
	cctxIndex, err := txserver.FetchAttributeFromTxResponse(res, "cctx_index")
	require.NoError(r, err)

	cctxRes, err := r.CctxClient.Cctx(r.Ctx, &crosschaintypes.QueryGetCctxRequest{Index: cctxIndex})
	require.NoError(r, err)

	cctx := cctxRes.CrossChainTx
	r.Logger.CCTX(*cctx, "pausing")

	// wait for the cctx to be mined
	r.WaitForMinedCCTXFromIndex(cctxIndex)

	// check ERC20 custody contract is paused
	paused, err = r.ERC20Custody.Paused(&bind.CallOpts{})
	require.NoError(r, err)
	require.True(r, paused)

	// Part 2: Unpause ERC20 custody contract

	// send command for unpausing ERC20 custody contract
	msg = crosschaintypes.NewMsgUpdateERC20CustodyPauseStatus(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		chainID.Int64(),
		false,
	)
	res, err = r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err)

	// fetch cctx index from tx response
	cctxIndex, err = txserver.FetchAttributeFromTxResponse(res, "cctx_index")
	require.NoError(r, err)

	cctxRes, err = r.CctxClient.Cctx(r.Ctx, &crosschaintypes.QueryGetCctxRequest{Index: cctxIndex})
	require.NoError(r, err)

	cctx = cctxRes.CrossChainTx
	r.Logger.CCTX(*cctx, "unpausing")

	// wait for the cctx to be mined
	r.WaitForMinedCCTXFromIndex(cctxIndex)

	// check ERC20 custody contract is unpaused
	paused, err = r.ERC20Custody.Paused(&bind.CallOpts{})
	require.NoError(r, err)
	require.False(r, paused)
}
