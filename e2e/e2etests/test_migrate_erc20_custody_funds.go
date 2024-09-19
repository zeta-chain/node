package e2etests

import (
	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/txserver"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestMigrateERC20CustodyFunds tests the migration of ERC20 custody funds
func TestMigrateERC20CustodyFunds(r *runner.E2ERunner, _ []string) {
	// get erc20 balance on ERC20 custody contract
	balance, err := r.ERC20.BalanceOf(&bind.CallOpts{}, r.ERC20CustodyAddr)
	require.NoError(r, err)

	// get EVM chain ID
	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	newAddr := sample.EthAddress()

	// send MigrateERC20CustodyFunds command
	msg := crosschaintypes.NewMsgMigrateERC20CustodyFunds(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		chainID.Int64(),
		newAddr.Hex(),
		r.ERC20Addr.Hex(),
		sdkmath.NewUintFromBigInt(balance),
	)
	res, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err)

	// fetch cctx index from tx response
	cctxIndex, err := txserver.FetchAttributeFromTxResponse(res, "cctx_index")
	require.NoError(r, err)

	cctxRes, err := r.CctxClient.Cctx(r.Ctx, &crosschaintypes.QueryGetCctxRequest{Index: cctxIndex})
	require.NoError(r, err)

	cctx := cctxRes.CrossChainTx
	r.Logger.CCTX(*cctx, "migration")

	// wait for the cctx to be mined
	r.WaitForMinedCCTXFromIndex(cctxIndex)

	// check ERC20 balance on new address
	newAddrBalance, err := r.ERC20.BalanceOf(&bind.CallOpts{}, newAddr)
	require.NoError(r, err)
	require.Equal(r, balance, newAddrBalance)

	// artificially set the ERC20 Custody address to the new address to prevent accounting check from failing
	r.ERC20CustodyAddr = newAddr
}
