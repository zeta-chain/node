package e2etests

import (
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

func TestV2Migration(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// Part 1: add new admin authorization
	err := r.ZetaTxServer.AddAuthorization("/zetachain.zetacore.crosschain.MsgUpdateERC20CustodyPauseStatus")
	require.NoError(r, err)

	err = r.ZetaTxServer.AddAuthorization("/zetachain.zetacore.crosschain.MsgMigrateERC20CustodyFunds")
	require.NoError(r, err)

	err = r.ZetaTxServer.AddAuthorization("/zetachain.zetacore.fungible.MsgUpdateGatewayContract")
	require.NoError(r, err)

	// Part 2: deploy v2 contracts on EVM chain
	r.SetupEVMV2()

	// Part 3: deploy gateway on ZetaChain
	r.SetZEVMContractsV2()

	// Part 4: upgrade all ZRC20s
	upgradeZRC20s(r)

	// Part 5: migrate ERC20 custody funds
	migrateERC20CustodyFunds(r)
}

func upgradeZRC20s(r *runner.E2ERunner) {

}

func migrateERC20CustodyFunds(r *runner.E2ERunner) {

}
