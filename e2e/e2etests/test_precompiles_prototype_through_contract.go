package e2etests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/testprototype"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

func TestPrecompilesPrototypeThroughContract(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	_, testPrototypeTx, testPrototype, err := testprototype.DeployTestPrototype(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, testPrototypeTx, r.Logger, r.ReceiptTimeout)

	res, err := testPrototype.Bech32ify(nil, "zeta", common.HexToAddress("0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE"))
	require.NoError(r, err, "Error calling Bech32ify")
	require.Equal(r, "zeta1h8duy2dltz9xz0qqhm5wvcnj02upy887fyn43u", res, "Failed to validate Bech32ify result")

	addr, err := testPrototype.Bech32ToHexAddr(nil, "zeta1h8duy2dltz9xz0qqhm5wvcnj02upy887fyn43u")
	require.NoError(r, err, "Error calling Bech32ToHexAddr")
	require.Equal(
		r,
		"0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE",
		addr.String(),
		"Failed to validate Bech32ToHexAddr result",
	)

	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err, "Error retrieving ChainID")

	balance, err := testPrototype.GetGasStabilityPoolBalance(nil, chainID.Int64())
	require.NoError(r, err, "Error calling GetGasStabilityPoolBalance")
	require.NotNil(r, balance, "GetGasStabilityPoolBalance returned balance is nil")
}
