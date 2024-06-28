package e2etests

import (
	"context"
	"fmt"
	"sort"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMigrateTss(r *runner.E2ERunner, args []string) {

	r.Logger.Info("Pause inbound and outbound processing")
	msg := observertypes.NewMsgDisableCCTX(
		r.ZetaTxServer.GetAccountAddress(0),
		false,
		true)
	_, err := r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msg)
	require.NoError(r, err)

	// Migrate btc
	// Fetch balance of BTC address
	utxos, err := r.GetTop20UTXOsForTssAddress()
	require.NoError(r, err)

	var btcBalance float64
	for _, utxo := range utxos {
		r.Logger.Print(fmt.Sprintf("UTXO Amount : %f, Spendable : %t", utxo.Amount, utxo.Spendable))
		r.Logger.Print(fmt.Sprintf("UTXO Amount : %d, Spendable : %t", int64(utxo.Amount*1e8), utxo.Spendable))
		btcBalance += utxo.Amount
	}
	r.Logger.Print("BTC TSS Balance Before fee deduction: %f", btcBalance)

	fees := 0.01
	btcBalance -= fees

	r.Logger.Print("BTC TSS Balance After fee deduction: %f", btcBalance)
	r.Logger.Print("BTC TSS migration amount: %d", int64(btcBalance*1e8))

	btcChain := int64(18444)
	migrationAmount := sdkmath.NewUint(uint64(btcBalance * 1e8))
	msgMigrateFunds := crosschaintypes.NewMsgMigrateTssFunds(
		r.ZetaTxServer.GetAccountAddress(0),
		btcChain,
		migrationAmount,
	)
	tx, err := r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msgMigrateFunds)
	require.NoError(r, err)

	// Fetch migrator cctx for btc migration
	migrator, err := r.ObserverClient.TssFundsMigratorInfo(r.Ctx, &observertypes.QueryTssFundsMigratorInfoRequest{
		ChainId: btcChain})
	require.NoError(r, err)

	r.Logger.Print(fmt.Sprintf("Migrator BTC: %s", migrator.TssFundsMigrator.MigrationCctxIndex))
	cctxBTCMigration := migrator.TssFundsMigrator.MigrationCctxIndex

	tssBalance, err := r.EVMClient.BalanceAt(context.Background(), r.TSSAddress, nil)
	require.NoError(r, err)

	tssBalanceUint := sdkmath.NewUintFromString(tssBalance.String())
	evmChainID, err := r.EVMClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	// Migrate TSS funds for the chain
	msgMigrateFunds = crosschaintypes.NewMsgMigrateTssFunds(
		r.ZetaTxServer.GetAccountAddress(0),
		evmChainID.Int64(),
		tssBalanceUint,
	)
	tx, err = r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msgMigrateFunds)
	require.NoError(r, err)

	r.Logger.Print(fmt.Sprintf("Migrate ETH TSS funds tx: %s", tx.TxHash))
	// Fetch migrator cctx for eth migration
	migrator, err = r.ObserverClient.TssFundsMigratorInfo(r.Ctx, &observertypes.QueryTssFundsMigratorInfoRequest{ChainId: evmChainID.Int64()})
	require.NoError(r, err)
	r.Logger.Print(fmt.Sprintf("Migrator ETH: %s", migrator.TssFundsMigrator.MigrationCctxIndex))
	cctxETHMigration := migrator.TssFundsMigrator.MigrationCctxIndex

	msgEnable := observertypes.NewMsgEnableCCTX(
		r.ZetaTxServer.GetAccountAddress(0),
		true,
		true)
	_, err = r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msgEnable)
	require.NoError(r, err)

	cctxbtc := utils.WaitCCTXMinedByIndex(r.Ctx, cctxBTCMigration, r.CctxClient, r.Logger, r.CctxTimeout)
	if cctxbtc.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic(fmt.Sprintf("expected cctx status to be mined; got %s, message: %s",
			cctxbtc.CctxStatus.Status.String(),
			cctxbtc.CctxStatus.StatusMessage),
		)
	}

	cctxETH := utils.WaitCCTXMinedByIndex(r.Ctx, cctxETHMigration, r.CctxClient, r.Logger, r.CctxTimeout)
	if cctxETH.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic(fmt.Sprintf("expected cctx status to be mined; got %s, message: %s",
			cctxETH.CctxStatus.Status.String(),
			cctxETH.CctxStatus.StatusMessage),
		)
	}

	// TODO Checks for these values
	tssBalance, err = r.EVMClient.BalanceAt(context.Background(), r.TSSAddress, nil)
	require.NoError(r, err)
	r.Logger.Print(fmt.Sprintf("TSS Balance After Old: %s", tssBalance.String()))

	tssBalanceNew, err := r.EVMClient.BalanceAt(context.Background(), common.HexToAddress(cctxETH.GetCurrentOutboundParam().Receiver), nil)
	require.NoError(r, err)
	r.Logger.Print(fmt.Sprintf("TSS Balance After New: %s", tssBalanceNew.String()))

	// Update TSS to new address
	allTss, err := r.ObserverClient.TssHistory(r.Ctx, &observertypes.QueryTssHistoryRequest{})
	require.NoError(r, err)

	require.Len(r, allTss.TssList, 2)

	sort.Slice(allTss.TssList, func(i, j int) bool {
		return allTss.TssList[i].FinalizedZetaHeight < allTss.TssList[j].FinalizedZetaHeight
	})
	msgUpdateTss := crosschaintypes.NewMsgUpdateTssAddress(
		r.ZetaTxServer.GetAccountAddress(0),
		allTss.TssList[1].TssPubkey,
	)
	tx, err = r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msgUpdateTss)
	require.NoError(r, err)

	r.Logger.Print(fmt.Sprintf("Update TSS tx: %s", tx.TxHash))

	time.Sleep(8 * time.Second)

	currentTss, err := r.ObserverClient.TSS(r.Ctx, &observertypes.QueryGetTSSRequest{})
	require.NoError(r, err)

	r.Logger.Print(fmt.Sprintf("Current TSS: %s", currentTss.TSS.TssPubkey))

	//if currentTss.TSS.TssPubkey != allTss.TssList[1].TssPubkey {
	//	panic(fmt.Sprintf("expected tss pubkey to be %s; got %s", allTss.TssList[1].TssPubkey, currentTss.TSS.TssPubkey))
	//}

}
