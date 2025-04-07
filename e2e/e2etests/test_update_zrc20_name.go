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
		"New ETH",
		"ETH.NEW",
	)
	res, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err)

	r.Logger.Info("Update eth zrc20 bytecode tx hash: %s", res.TxHash)

	// Get new info of the ZRC20
	newName, err := r.ETHZRC20.Name(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, "New ETH", newName)

	newSymbol, err := r.ETHZRC20.Symbol(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, "ETH.NEW", newSymbol)

	qRes, err := r.FungibleClient.ForeignCoins(r.Ctx, &fungibletypes.QueryGetForeignCoinsRequest{
		Index: r.ETHZRC20Addr.Hex(),
	})
	require.NoError(r, err)
	require.EqualValues(r, "New ETH", qRes.ForeignCoins.Name)
	require.EqualValues(r, "ETH.NEW", qRes.ForeignCoins.Symbol)

	// try another zrc20
	msg = fungibletypes.NewMsgUpdateZRC20Name(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		r.ERC20ZRC20Addr.Hex(),
		"New USDT",
		"USDT.NEW",
	)
	res, err = r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err)

	r.Logger.Info("Update erc20 zrc20 bytecode tx hash: %s", res.TxHash)

	// Get new info of the ZRC20
	newName, err = r.ERC20ZRC20.Name(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, "New USDT", newName)

	newSymbol, err = r.ERC20ZRC20.Symbol(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, "USDT.NEW", newSymbol)

	qRes, err = r.FungibleClient.ForeignCoins(r.Ctx, &fungibletypes.QueryGetForeignCoinsRequest{
		Index: r.ERC20ZRC20Addr.Hex(),
	})
	require.NoError(r, err)
	require.EqualValues(r, "New USDT", qRes.ForeignCoins.Name)
	require.EqualValues(r, "USDT.NEW", qRes.ForeignCoins.Symbol)
}
