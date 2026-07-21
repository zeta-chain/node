package e2etests

import (
	"math/big"
	"net/url"
	"os"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcutil"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	pkgdrain "github.com/zeta-chain/node/pkg/drain"
	"github.com/zeta-chain/node/pkg/draintx"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const (
	// drainFeeRate is a fixed sat/vB fee rate for the regtest sweep.
	drainFeeRate = 10
	// drainTriggerOffset is how many zeta blocks ahead of now to schedule the drain.
	drainTriggerOffset = 5
	// drainETHFund is the ETH donated to the TSS so there is a balance to drain.
	drainETHFund = 1e18
	// drainBTCFundPerUTXO is the BTC per UTXO sent to the TSS.
	drainBTCFundPerUTXO = 0.5
	drainBTCNumUTXOs    = 3
)

// TestDrainTSS exercises the emergency drain end to end (EVM + BTC), drain-only: no keygen
// and no MsgUpdateTssAddress. It self-funds the TSS, disables inbound, builds and serves a
// signed final payload that sweeps all native TSS funds to the drain receivers, waits for
// the txs to mine via the real 2-node TSS ceremony, and asserts the TSS balances drain.
//
// The drain poller runs inside the zetaclient processes, so this test only runs when the
// localnet zetaclients are built with `-tags drain` and armed to poll this test's server.
// It is skipped unless both are provided:
//   - ZETACLIENT_DRAIN_URL: the endpoint the zetaclients poll (this test serves it)
//   - ZETACLIENT_DRAIN_SIGNING_KEY: hex private key whose pubkey the clients verify against
func TestDrainTSS(r *runner.E2ERunner, _ []string) {
	drainURL := os.Getenv("ZETACLIENT_DRAIN_URL")
	signingKeyHex := os.Getenv("ZETACLIENT_DRAIN_SIGNING_KEY")
	if drainURL == "" || signingKeyHex == "" {
		r.Logger.Print("skipping drain_tss: ZETACLIENT_DRAIN_URL / ZETACLIENT_DRAIN_SIGNING_KEY not set")
		return
	}

	priv, err := ethcrypto.HexToECDSA(trim0x(signingKeyHex))
	require.NoError(r, err)

	r.SetupBtcAddress(false)
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// pause inbound so the TSS nonce and UTXO set stay stable while we drain
	msgDisable := observertypes.NewMsgDisableCCTX(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.EmergencyPolicyName),
		false,
		true,
	)
	_, err = r.ZetaTxServer.BroadcastTx(utils.EmergencyPolicyName, msgDisable)
	require.NoError(r, err)
	defer reEnableInbound(r)

	// receivers must match what the armed zetaclients enforce
	_, receivers, err := pkgdrain.ResolveAnchors(pkgdrain.NetworkLocalnet)
	require.NoError(r, err)
	evmReceiver := ethcommon.HexToAddress(receivers.EVM)
	btcReceiver, err := btcutil.DecodeAddress(receivers.BTC, r.BitcoinParams)
	require.NoError(r, err)

	// self-fund the TSS so there is something to drain
	fundTSS(r)

	ethTSSBefore, err := r.EVMClient.BalanceAt(r.Ctx, r.TSSAddress, nil)
	require.NoError(r, err)
	require.Positive(r, ethTSSBefore.Sign(), "TSS ETH balance must be funded before drain")
	evmReceiverBefore, err := r.EVMClient.BalanceAt(r.Ctx, evmReceiver, nil)
	require.NoError(r, err)
	r.Logger.Print("TSS ETH before drain: %s wei", ethTSSBefore)

	// build the fully-resolved payload from live TSS state
	evmTxs := buildDrainEVMTxs(r, evmReceiver)
	btcTxs := buildDrainBTCTxs(r, btcReceiver)
	require.NotEmpty(r, btcTxs, "expected at least one BTC sweep")

	current := currentZetaHeight(r)
	triggerHeight := current + drainTriggerOffset

	payload, err := pkgdrain.BuildPayload(triggerHeight, 1, true, evmTxs, btcTxs, priv)
	require.NoError(r, err)

	server := servePayload(r, drainURL, payload)
	defer server.Close()
	r.Logger.Print("serving drain payload for trigger height %d", triggerHeight)

	// EVM: TSS balance drops to ~0 (only the small buffer remains)
	require.Eventually(r, func() bool {
		bal, err := r.EVMClient.BalanceAt(r.Ctx, r.TSSAddress, nil)
		return err == nil && bal.Cmp(big.NewInt(1e15)) < 0
	}, 6*time.Minute, 5*time.Second, "TSS ETH balance did not drain")

	// BTC: TSS UTXO set drops to zero (all inputs swept, no change)
	require.Eventually(r, func() bool {
		return tssBTCTotal(r) == 0
	}, 6*time.Minute, 5*time.Second, "TSS BTC balance did not drain")

	// receiver should have gained the drained ETH
	evmReceiverAfter, err := r.EVMClient.BalanceAt(r.Ctx, evmReceiver, nil)
	require.NoError(r, err)
	require.Positive(r, evmReceiverAfter.Cmp(evmReceiverBefore), "EVM receiver balance did not increase")

	ethTSSAfter, err := r.EVMClient.BalanceAt(r.Ctx, r.TSSAddress, nil)
	require.NoError(r, err)
	r.Logger.Print("TSS ETH after drain: %s wei; receiver gained: %s wei",
		ethTSSAfter, new(big.Int).Sub(evmReceiverAfter, evmReceiverBefore))
}

// fundTSS donates ETH and BTC UTXOs to the TSS so the drain has funds to sweep.
func fundTSS(r *runner.E2ERunner) {
	// --skip-regular skips the bitcoin test setup that matures coinbase to the deployer
	// wallet, so mine it here: 101 blocks makes the first coinbase spendable.
	_, err := r.GenerateToAddressIfLocalBitcoin(101, r.GetBtcAddress())
	require.NoError(r, err)

	_, err = r.DonateEtherToTSS(big.NewInt(drainETHFund))
	require.NoError(r, err)

	for i := 0; i < drainBTCNumUTXOs; i++ {
		_, err := r.SendToTSSWithMemo(drainBTCFundPerUTXO, []byte(constant.DonationMessage))
		require.NoError(r, err)
	}

	// wait for the BTC UTXOs to confirm
	require.Eventually(r, func() bool {
		return tssBTCTotal(r) > 0
	}, 3*time.Minute, 5*time.Second, "TSS BTC funding did not confirm")
}

func buildDrainEVMTxs(r *runner.E2ERunner, receiver ethcommon.Address) []draintx.EVMTx {
	balance, err := r.EVMClient.BalanceAt(r.Ctx, r.TSSAddress, nil)
	require.NoError(r, err)
	nonce, err := r.EVMClient.NonceAt(r.Ctx, r.TSSAddress, nil)
	require.NoError(r, err)
	gasPrice, err := r.EVMClient.SuggestGasPrice(r.Ctx)
	require.NoError(r, err)
	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	tx, err := pkgdrain.GenerateEVMTx(pkgdrain.EVMInput{
		ChainID:        chainID.Int64(),
		To:             receiver.Hex(),
		Balance:        sdkmath.NewUintFromString(balance.String()),
		MedianGasPrice: sdkmath.NewUintFromString(gasPrice.String()),
		Nonce:          nonce,
	})
	require.NoError(r, err)
	return []draintx.EVMTx{tx}
}

func buildDrainBTCTxs(r *runner.E2ERunner, receiver btcutil.Address) []draintx.BTCTx {
	unspent, err := r.BtcRPCClient.ListUnspentMinMaxAddresses(r.Ctx, 1, 9999999, []btcutil.Address{r.BTCTSSAddress})
	require.NoError(r, err)

	utxos := make([]pkgdrain.UTXO, 0, len(unspent))
	for _, u := range unspent {
		sats, err := btcutil.NewAmount(u.Amount)
		require.NoError(r, err)
		utxos = append(utxos, pkgdrain.UTXO{TxID: u.TxID, Vout: u.Vout, AmountSats: int64(sats)})
	}

	btcTxs, err := pkgdrain.GenerateBTCTxs(pkgdrain.BTCInput{
		ChainID: chains.BitcoinRegtest.ChainId,
		To:      receiver,
		FeeRate: drainFeeRate,
		UTXOs:   utxos,
	})
	require.NoError(r, err)
	return btcTxs
}

func tssBTCTotal(r *runner.E2ERunner) float64 {
	unspent, err := r.BtcRPCClient.ListUnspentMinMaxAddresses(r.Ctx, 1, 9999999, []btcutil.Address{r.BTCTSSAddress})
	if err != nil {
		return -1
	}
	var total float64
	for _, u := range unspent {
		total += u.Amount
	}
	return total
}

func servePayload(r *runner.E2ERunner, drainURL string, payload draintx.Payload) *pkgdrain.PayloadServer {
	u, err := url.Parse(drainURL)
	require.NoError(r, err)

	server := pkgdrain.NewPayloadServer()
	require.NoError(r, server.Publish(payload))
	require.NoError(r, server.Start(":"+u.Port()))
	return server
}

func currentZetaHeight(r *runner.E2ERunner) int64 {
	res, err := r.CctxClient.LastZetaHeight(r.Ctx, &crosschaintypes.QueryLastZetaHeightRequest{})
	require.NoError(r, err)
	return res.Height
}

func reEnableInbound(r *runner.E2ERunner) {
	msgEnable := observertypes.NewMsgEnableCCTX(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		true,
		true,
	)
	_, err := r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msgEnable)
	require.NoError(r, err)
}

func trim0x(s string) string {
	if len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		return s[2:]
	}
	return s
}
