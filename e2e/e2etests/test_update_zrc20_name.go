package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

// TestUpdateZRC20Name tests updating name and symbol of a ZRC20
func TestUpdateZRC20Name(r *runner.E2ERunner, _ []string) {
	msg := fungibletypes.NewMsgUpdateZRC20Name(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		r.ETHZRC20Addr.Hex(),
		"New USDT",
		"USDT.NEW",
	)
	res, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err)

	r.Logger.Info("Update zrc20 bytecode tx hash: %s", res.TxHash)

	// Get new info of the ZRC20
	r.Logger.Info("Checking the new values of the ZRC20")
	newName, err := r.ETHZRC20.Name(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, "New USDT", newName)

	newSymbol, err := r.ETHZRC20.Symbol(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, "USDT.NEW", newSymbol)
}
