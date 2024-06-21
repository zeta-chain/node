package e2etests

import (
	"context"

	sdkmath "cosmossdk.io/math"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// TestEtherWithdraw tests the withdraw of ether
func TestMigrateTssEth(r *runner.E2ERunner, args []string) {

	r.Logger.Info("Pause inbound and outbound processing")
	msg := observertypes.NewMsgDisableCCTX(
		r.ZetaTxServer.GetAccountAddress(0),
		true,
		true)
	_, err := r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}

	// Fetch balance of TSS address
	tssBalance, err := r.EVMClient.BalanceAt(context.Background(), r.TSSAddress, nil)
	if err != nil {
		panic(err)
	}
	r.Logger.Print("TSS Balance: ", tssBalance.String())
	tssBalanceUint := sdkmath.NewUintFromString(tssBalance.String())
	evmChainID, err := r.EVMClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	r.Logger.Print("EVM Chain ID: ", evmChainID.String())
	// Migrate TSS funds for the chain
	msgMigrateFunds := crosschaintypes.NewMsgMigrateTssFunds(
		r.ZetaTxServer.GetAccountAddress(0),
		evmChainID.Int64(),
		tssBalanceUint,
	)
	_, err = r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msgMigrateFunds)
	if err != nil {
		panic(err)
	}

}
