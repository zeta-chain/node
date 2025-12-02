package e2etests

import (
	"context"
	"sort"
	"strconv"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	zetacrypto "github.com/zeta-chain/node/pkg/crypto"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestMigrateTSS(r *runner.E2ERunner, _ []string) {
	r.SetupBtcAddress(false)
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// Pause inbound processing for tss migration
	r.Logger.Info("Pause inbound  processing")
	msg := observertypes.NewMsgDisableCCTX(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.EmergencyPolicyName),
		false,
		true)
	_, err := r.ZetaTxServer.BroadcastTx(utils.EmergencyPolicyName, msg)
	require.NoError(r, err)

	// Migrate btc
	// Fetch balance of BTC TSS address
	utxos, err := r.GetTop20UTXOsForTssAddress()
	require.NoError(r, err)

	var btcBalance float64
	for _, utxo := range utxos {
		btcBalance += utxo.Amount
	}

	btcTSSBalanceOld := btcBalance
	halfBtcBalance := btcBalance / 2
	// Use fixed fee of 0.02 for migration
	// Code links for reference to calculate fee:
	// https://github.com/zeta-chain/node/blob/5c73a45b6096ed7f7387ae8f7ccc7d5e2c668e0d/zetaclient/chains/bitcoin/signer/sign.go#L54
	// https://github.com/zeta-chain/node/blob/5c73a45b6096ed7f7387ae8f7ccc7d5e2c668e0d/zetaclient/chains/bitcoin/signer/sign.go#L32
	halfBtcBalance = halfBtcBalance - 0.02
	btcChain := chains.BitcoinRegtest.ChainId

	r.WaitForTSSGeneration(2)

	//migrate btc funds
	// #nosec G701 e2eTest - always in range
	cctxBTCAmount := sdkmath.ZeroUint()
	for i := 0; i < 2; i++ {
		migrationAmountBTC := sdkmath.NewUint(uint64(halfBtcBalance * 1e8))
		msgMigrateFunds := crosschaintypes.NewMsgMigrateTssFunds(
			r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
			btcChain,
			migrationAmountBTC,
		)
		_, err = r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msgMigrateFunds)
		require.NoError(r, err)

		// Fetch migrator cctx for btc migration
		migrator, err := r.ObserverClient.TssFundsMigratorInfo(r.Ctx, &observertypes.QueryTssFundsMigratorInfoRequest{
			ChainId: btcChain})
		require.NoError(r, err)
		cctxBTCMigration := migrator.TssFundsMigrator.MigrationCctxIndex
		cctxBTC := utils.WaitCCTXMinedByIndex(r.Ctx, cctxBTCMigration, r.CctxClient, r.Logger, r.CctxTimeout)
		require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctxBTC.CctxStatus.Status)

		cctxBTCAmount = cctxBTCAmount.Add(cctxBTC.GetCurrentOutboundParam().Amount)

	}

	// ETH migration
	// Fetch balance of ETH TSS address
	tssBalance, err := r.EVMClient.BalanceAt(context.Background(), r.TSSAddress, nil)
	require.NoError(r, err)
	ethTSSBalanceOld := tssBalance

	tssBalanceUint := sdkmath.NewUintFromString(tssBalance.String())
	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	// Migrate TSS funds for the eth chain
	msgMigrateFunds := crosschaintypes.NewMsgMigrateTssFunds(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		evmChainID.Int64(),
		tssBalanceUint,
	)
	_, err = r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msgMigrateFunds)
	require.NoError(r, err)

	// Fetch migrator cctx for eth migration
	migrator, err := r.ObserverClient.TssFundsMigratorInfo(
		r.Ctx,
		&observertypes.QueryTssFundsMigratorInfoRequest{ChainId: evmChainID.Int64()},
	)
	require.NoError(r, err)
	cctxETHMigration := migrator.TssFundsMigrator.MigrationCctxIndex

	cctxETH := utils.WaitCCTXMinedByIndex(r.Ctx, cctxETHMigration, r.CctxClient, r.Logger, r.CctxTimeout)
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctxETH.CctxStatus.Status)

	// Check if new TSS is added to list
	allTss, err := r.ObserverClient.TssHistory(r.Ctx, &observertypes.QueryTssHistoryRequest{})
	require.NoError(r, err)

	require.Len(r, allTss.TssList, 2)

	// Update TSS to new address
	sort.Slice(allTss.TssList, func(i, j int) bool {
		return allTss.TssList[i].FinalizedZetaHeight < allTss.TssList[j].FinalizedZetaHeight
	})
	msgUpdateTss := crosschaintypes.NewMsgUpdateTssAddress(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		allTss.TssList[1].TssPubkey,
	)
	_, err = r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msgUpdateTss)
	require.NoError(r, err)

	// Wait for at least one block for the TSS to be updated
	time.Sleep(8 * time.Second)

	currentTss, err := r.ObserverClient.TSS(r.Ctx, &observertypes.QueryGetTSSRequest{})
	require.NoError(r, err)
	require.Equal(r, allTss.TssList[1].TssPubkey, currentTss.TSS.TssPubkey)

	newTss, err := r.ObserverClient.GetTssAddress(r.Ctx, &observertypes.QueryGetTssAddressRequest{})
	require.NoError(r, err)

	// Check balance of new TSS address to make sure all funds have been transferred
	// BTC
	btcTssAddress, err := zetacrypto.GetTSSAddrBTC(currentTss.TSS.TssPubkey, r.BitcoinParams)
	require.NoError(r, err)

	btcTssAddressNew, err := btcutil.DecodeAddress(btcTssAddress, r.BitcoinParams)
	require.NoError(r, err)

	r.BTCTSSAddress = btcTssAddressNew
	r.AddTSSToNode()

	utxos, err = r.GetTop20UTXOsForTssAddress()
	require.NoError(r, err)

	var btcTSSBalanceNew float64
	// #nosec G701 e2eTest - always in range
	for _, utxo := range utxos {
		btcTSSBalanceNew += utxo.Amount
	}
	require.Equal(
		r,
		strconv.FormatInt(int64(btcTSSBalanceNew*1e8), 10),
		cctxBTCAmount.String(),
	)
	r.Logger.Info("BTC TSS Balance Old: %f", btcTSSBalanceOld*1e8)
	r.Logger.Info("BTC TSS Balance New: %f", btcTSSBalanceNew*1e8)
	r.Logger.Info("Total Migration amount : %s", cctxBTCAmount)

	// btcTSSBalanceNew should be less than btcTSSBalanceOld as there is some loss of funds during migration
	// #nosec G701 e2eTest - always in range
	require.LessOrEqual(r, btcTSSBalanceNew*1e8, btcTSSBalanceOld*1e8)

	// ETH
	r.TSSAddress = common.HexToAddress(newTss.Eth)
	ethTSSBalanceNew, err := r.EVMClient.BalanceAt(context.Background(), r.TSSAddress, nil)
	require.NoError(r, err)

	r.Logger.Info("ETH TSS Balance Old: %s", ethTSSBalanceOld.String())
	r.Logger.Info("ETH TSS Balance New: %s", ethTSSBalanceNew.String())
	r.Logger.Info("Total Migration amount : %s", cctxETH.GetCurrentOutboundParam().Amount.String())

	// ethTSSBalanceNew should be less than ethTSSBalanceOld as there is some loss of funds during migration
	require.Equal(r, ethTSSBalanceNew.String(), cctxETH.GetCurrentOutboundParam().Amount.String())
	require.True(r, ethTSSBalanceNew.Cmp(ethTSSBalanceOld) < 0)

	msgEnable := observertypes.NewMsgEnableCCTX(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		true,
		true)
	_, err = r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msgEnable)
	require.NoError(r, err)
}
