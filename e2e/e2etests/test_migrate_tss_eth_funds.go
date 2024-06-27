package e2etests

import (
	"context"
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
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
		false,
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
	r.Logger.Print(fmt.Sprintf("TSS Balance: %s", tssBalance.String()))
	tssBalanceUint := sdkmath.NewUintFromString(tssBalance.String())
	evmChainID, err := r.EVMClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	// Migrate TSS funds for the chain
	msgMigrateFunds := crosschaintypes.NewMsgMigrateTssFunds(
		r.ZetaTxServer.GetAccountAddress(0),
		evmChainID.Int64(),
		tssBalanceUint,
	)
	tx, err := r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msgMigrateFunds)
	if err != nil {
		panic(err)
	}
	r.Logger.Print(fmt.Sprintf("Migrate TSS funds tx: %s", tx.TxHash))
	// Fetch migrator cctx
	migrator, err := r.ObserverClient.TssFundsMigratorInfo(r.Ctx, &observertypes.QueryTssFundsMigratorInfoRequest{ChainId: evmChainID.Int64()})
	if err != nil {
		r.Logger.Print("Error fetching migrator: ", err)
		return
	}

	r.Logger.Print(fmt.Sprintf("Migrator: %s", migrator.TssFundsMigrator.MigrationCctxIndex))

	cctx := utils.WaitCCTXMinedByIndex(r.Ctx, migrator.TssFundsMigrator.MigrationCctxIndex, r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic(fmt.Sprintf("expected cctx status to be mined; got %s, message: %s",
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage),
		)
	}

	// TODO Checks for these values
	tssBalance, err = r.EVMClient.BalanceAt(context.Background(), r.TSSAddress, nil)
	if err != nil {
		panic(err)
	}
	r.Logger.Print(fmt.Sprintf("TSS Balance After Old: %s", tssBalance.String()))

	tssBalanceNew, err := r.EVMClient.BalanceAt(context.Background(), common.HexToAddress(cctx.GetCurrentOutboundParam().Receiver), nil)
	if err != nil {
		panic(err)
	}
	r.Logger.Print(fmt.Sprintf("TSS Balance After New: %s", tssBalanceNew.String()))

	allTss, err := r.ObserverClient.TssHistory(r.Ctx, &observertypes.QueryTssHistoryRequest{})
	if err != nil {
		panic(err)
	}

	if len(allTss.TssList) != 2 {
		panic(fmt.Sprintf("expected 2 tss addresses; got %d", len(allTss.TssList)))
	}

	msgUpdateTss := crosschaintypes.NewMsgUpdateTssAddress(
		r.ZetaTxServer.GetAccountAddress(0),
		allTss.TssList[0].TssPubkey,
	)
	tx, err = r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msgUpdateTss)
	if err != nil {
		panic(err)
	}

	r.Logger.Print(fmt.Sprintf("Update TSS tx: %s", tx.TxHash))

	time.Sleep(8 * time.Second)

	currentTss, err := r.ObserverClient.TSS(r.Ctx, &observertypes.QueryGetTSSRequest{})
	if err != nil {
		panic(err)
	}

	r.Logger.Print(fmt.Sprintf("Current TSS: %s", currentTss.TSS.TssPubkey))

	//if currentTss.TSS.TssPubkey != allTss.TssList[1].TssPubkey {
	//	panic(fmt.Sprintf("expected tss pubkey to be %s; got %s", allTss.TssList[1].TssPubkey, currentTss.TSS.TssPubkey))
	//}

}
