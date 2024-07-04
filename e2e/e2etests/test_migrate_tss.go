package e2etests

import (
	"context"
	"fmt"
	"sort"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	zetacrypto "github.com/zeta-chain/zetacore/pkg/crypto"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMigrateTss(r *runner.E2ERunner, args []string) {

	r.SetBtcAddress(r.Name, false)
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// Pause inbound procoessing for tss migration
	r.Logger.Info("Pause inbound  processing")
	msg := observertypes.NewMsgDisableCCTX(
		r.ZetaTxServer.GetAccountAddress(0),
		false,
		true)
	_, err := r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msg)
	require.NoError(r, err)

	// Migrate btc

	// Fetch balance of BTC TSS address
	utxos, err := r.GetTop20UTXOsForTssAddress()
	require.NoError(r, err)

	var btcBalance float64
	for _, utxo := range utxos {
		r.Logger.Print(fmt.Sprintf("UTXO Amount old : %d, Spendable : %t", int64(utxo.Amount*1e8), utxo.Spendable))
		btcBalance += utxo.Amount
	}

	// Use fixed fee for migration
	btcTSSBalanceOld := btcBalance
	fees := 0.01
	btcBalance -= fees
	btcChain := int64(18444)

	r.Logger.Info("BTC TSS migration amount: %d", int64(btcBalance*1e8))

	//migrate btc funds
	migrationAmountBTC := sdkmath.NewUint(uint64(btcBalance * 1e8))
	msgMigrateFunds := crosschaintypes.NewMsgMigrateTssFunds(
		r.ZetaTxServer.GetAccountAddress(0),
		btcChain,
		migrationAmountBTC,
	)
	_, err = r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msgMigrateFunds)
	require.NoError(r, err)

	// Fetch migrator cctx for btc migration
	migrator, err := r.ObserverClient.TssFundsMigratorInfo(r.Ctx, &observertypes.QueryTssFundsMigratorInfoRequest{
		ChainId: btcChain})
	require.NoError(r, err)
	cctxBTCMigration := migrator.TssFundsMigrator.MigrationCctxIndex

	// ETH migration
	// Fetch balance of ETH TSS address
	tssBalance, err := r.EVMClient.BalanceAt(context.Background(), r.TSSAddress, nil)
	require.NoError(r, err)
	ethTSSBalanceOld := tssBalance

	tssBalanceUint := sdkmath.NewUintFromString(tssBalance.String())
	evmChainID, err := r.EVMClient.ChainID(context.Background())
	require.NoError(r, err)

	// Migrate TSS funds for the eth chain
	msgMigrateFunds = crosschaintypes.NewMsgMigrateTssFunds(
		r.ZetaTxServer.GetAccountAddress(0),
		evmChainID.Int64(),
		tssBalanceUint,
	)
	_, err = r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msgMigrateFunds)
	require.NoError(r, err)

	// Fetch migrator cctx for eth migration
	migrator, err = r.ObserverClient.TssFundsMigratorInfo(r.Ctx, &observertypes.QueryTssFundsMigratorInfoRequest{ChainId: evmChainID.Int64()})
	require.NoError(r, err)
	cctxETHMigration := migrator.TssFundsMigrator.MigrationCctxIndex

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
	_, err = r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msgUpdateTss)
	require.NoError(r, err)

	time.Sleep(8 * time.Second)

	currentTss, err := r.ObserverClient.TSS(r.Ctx, &observertypes.QueryGetTSSRequest{})
	require.NoError(r, err)
	require.Equal(r, allTss.TssList[1].TssPubkey, currentTss.TSS.TssPubkey)

	newTss, err := r.ObserverClient.GetTssAddress(r.Ctx, &observertypes.QueryGetTssAddressRequest{})
	require.NoError(r, err)

	// BTC

	btcTssAddress, err := zetacrypto.GetTssAddrBTC(currentTss.TSS.TssPubkey, r.BitcoinParams)
	require.NoError(r, err)

	btcTssAddressNew, err := btcutil.DecodeAddress(btcTssAddress, r.BitcoinParams)
	require.NoError(r, err)

	r.Logger.Print(fmt.Sprintf("Pubkey New : %s", currentTss.TSS.TssPubkey))
	r.Logger.Print(fmt.Sprintf("BTC AddressFromPubkey : %s   Decoded : %s,Receiver %s", btcTssAddress, btcTssAddressNew, cctxbtc.GetCurrentOutboundParam().Receiver))
	r.Logger.Print(fmt.Sprintf("BTC TSS address from zetacore : %s ", newTss.Btc))

	r.BTCTSSAddress = btcTssAddressNew
	r.AddTssToNode()

	utxos, err = r.GetTop20UTXOsForTssAddress()
	require.NoError(r, err)

	var btcTSSBalanceNew float64
	for _, utxo := range utxos {
		r.Logger.Print(fmt.Sprintf("UTXO Amount new : %d, Spendable : %t", int64(utxo.Amount*1e8), utxo.Spendable))
		btcTSSBalanceNew += utxo.Amount
	}

	pubkeyOld := allTss.TssList[0].TssPubkey
	bO, err := zetacrypto.GetTssAddrBTC(pubkeyOld, r.BitcoinParams)
	require.NoError(r, err)

	bAO, err := btcutil.DecodeAddress(bO, r.BitcoinParams)
	require.NoError(r, err)

	utxos, err = r.BtcRPCClient.ListUnspentMinMaxAddresses(
		1,
		9999999,
		[]btcutil.Address{bAO},
	)
	require.NoError(r, err)

	sort.SliceStable(utxos, func(i, j int) bool {
		return utxos[i].Amount < utxos[j].Amount
	})

	if len(utxos) > 20 {
		utxos = utxos[:20]
	}
	var bOB float64
	for _, utxo := range utxos {
		r.Logger.Print(fmt.Sprintf("UTXO Amount old recalculate: %d, Spendable : %t", int64(utxo.Amount*1e8), utxo.Spendable))
		bOB += utxo.Amount
	}

	r.Logger.Print(fmt.Sprintf("BTC Balance Old: %f", btcTSSBalanceOld*1e8))
	r.Logger.Print(fmt.Sprintf("BTC Balance New: %f", btcTSSBalanceNew*1e8))
	r.Logger.Print(fmt.Sprintf("Migrator amount : %s", cctxbtc.GetCurrentOutboundParam().Amount))
	r.Logger.Print(fmt.Sprintf("BTC Balance Old Recalcuated: %f", bOB*1e8))

	// ETH

	r.TSSAddress = common.HexToAddress(newTss.Eth)

	ethTSSBalanceNew, err := r.EVMClient.BalanceAt(context.Background(), r.TSSAddress, nil)
	require.NoError(r, err)

	r.Logger.Print(fmt.Sprintf("TSS Balance Old: %s", ethTSSBalanceOld.String()))
	r.Logger.Print(fmt.Sprintf("TSS Balance New: %s", ethTSSBalanceNew.String()))
	r.Logger.Print(fmt.Sprintf("Migrator amount : %s", cctxETH.GetCurrentOutboundParam().Amount.String()))

	msgEnable := observertypes.NewMsgEnableCCTX(
		r.ZetaTxServer.GetAccountAddress(0),
		true,
		true)
	_, err = r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msgEnable)
	require.NoError(r, err)
}
