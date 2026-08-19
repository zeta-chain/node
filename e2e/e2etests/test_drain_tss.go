package e2etests

import (
	"crypto/ecdsa"
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

	// drainRehearsalETHCap is the --evm-max-amount equivalent for the rehearsal phase: well below
	// drainETHFund, so the TSS keeps most of its balance for the real drain that follows.
	drainRehearsalETHCap = 5e16 // 0.05 ETH
	// drainRehearsalBTCCapSats is the --btc-max-sats equivalent. Each funding UTXO is 0.5 BTC
	// (50,000,000 sats), so a cap between one and two of them selects exactly one — which is the
	// point: BTC is spent in whole UTXOs, so a rehearsal moves one chunk and leaves the rest.
	drainRehearsalBTCCapSats = 60_000_000
)

// TestDrainTSS exercises the emergency drain end to end (EVM + BTC), drain-only: no keygen
// and no MsgUpdateTssAddress. It self-funds the TSS, disables inbound, then runs the two-phase
// sequence an operator actually performs, each phase signed by the real 2-node TSS ceremony:
//
//  1. A capped REHEARSAL (--evm-max-amount / --btc-max-sats equivalents), asserting it moves only
//     the capped value and leaves the rest at the TSS. Without this phase the caps have no
//     end-to-end coverage at all, and a capped payload would first be signed on a live network.
//  2. The uncapped REAL drain at a higher trigger height, asserting the TSS balances go to ~0.
//
// The second phase also covers two things the rehearse-then-drain flow depends on and that nothing
// else exercises: the poller re-arming on a newer payload (lastFiredHeight), and the requirement
// that the rehearsal tx confirms before the real payload pins a nonce.
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
	btcUTXOsBefore := tssBTCUTXOCount(r)
	require.GreaterOrEqual(r, btcUTXOsBefore, drainBTCNumUTXOs, "TSS must hold the funding UTXOs")
	r.Logger.Print("TSS ETH before drain: %s wei", ethTSSBefore)

	server := pkgdrain.NewPayloadServer()
	startPayloadServer(r, server, drainURL)
	defer server.Close()

	// ------------------------------------------------------------------
	// Phase 1: the rehearsal. Capped values, so this must move a small amount and leave the rest
	// at the TSS. This is the sequence an operator runs before the real drain, and it is the only
	// coverage that a capped payload survives a real TSS ceremony rather than just unit tests.
	// ------------------------------------------------------------------
	evmCap := sdkmath.NewUint(drainRehearsalETHCap)
	rehearsalEVM := buildDrainEVMTxs(r, evmReceiver, evmCap)
	rehearsalBTC := buildDrainBTCTxs(r, btcReceiver, drainRehearsalBTCCapSats)
	require.Equal(r, evmCap.String(), rehearsalEVM[0].Amount, "the pinned amount must be the cap")
	require.Len(r, rehearsalBTC, 1, "a capped run selects one tx worth of inputs, so one sweep")

	// Assert partiality against what the payload actually pins rather than a hardcoded count, so a
	// stray UTXO left by earlier setup cannot turn this into a false pass or a false failure.
	var rehearsalSweptSats int64
	for _, in := range rehearsalBTC[0].Inputs {
		rehearsalSweptSats += in.AmountSats
	}
	require.LessOrEqual(r, rehearsalSweptSats, int64(drainRehearsalBTCCapSats), "sweep exceeds the cap")
	require.Less(r, len(rehearsalBTC[0].Inputs), btcUTXOsBefore, "the cap must leave UTXOs behind")

	rehearsalNonce := rehearsalEVM[0].Nonce
	rehearsalHeight := currentZetaHeight(r) + drainTriggerOffset
	publishDrainPayload(r, server, priv, rehearsalHeight, 1, rehearsalEVM, rehearsalBTC)
	r.Logger.Print("serving REHEARSAL payload for trigger height %d", rehearsalHeight)

	// the capped EVM transfer lands: the receiver gains exactly the cap, since the amount is pinned
	require.Eventually(r, func() bool {
		bal, err := r.EVMClient.BalanceAt(r.Ctx, evmReceiver, nil)
		return err == nil && new(big.Int).Sub(bal, evmReceiverBefore).Cmp(evmCap.BigInt()) == 0
	}, 6*time.Minute, 5*time.Second, "capped EVM drain did not credit the receiver with the cap")

	// only the pinned UTXOs are swept and the rest stay put — a cap cannot split a UTXO, so this
	// is the assertion that a rehearsal really is partial
	require.Eventually(r, func() bool {
		return tssBTCUTXOCount(r) == btcUTXOsBefore-len(rehearsalBTC[0].Inputs)
	}, 6*time.Minute, 5*time.Second, "capped BTC sweep did not leave the remaining UTXOs at the TSS")

	ethTSSMid, err := r.EVMClient.BalanceAt(r.Ctx, r.TSSAddress, nil)
	require.NoError(r, err)
	require.Positive(
		r,
		ethTSSMid.Cmp(new(big.Int).Div(ethTSSBefore, big.NewInt(2))),
		"rehearsal drained most of the balance instead of the cap",
	)
	require.Positive(r, tssBTCTotal(r), "rehearsal swept all BTC instead of one UTXO")
	r.Logger.Print("after REHEARSAL: TSS ETH %s wei, TSS BTC %f in %d UTXOs",
		ethTSSMid, tssBTCTotal(r), tssBTCUTXOCount(r))

	// ------------------------------------------------------------------
	// Phase 2: the real drain, uncapped, at a higher trigger height. This also covers the poller
	// re-arming on a newer payload (lastFiredHeight), which the rehearse-then-drain flow depends on.
	// ------------------------------------------------------------------
	//
	// The rehearsal tx must CONFIRM first. The payload pins the confirmed nonce while the poller
	// checks the pending one and treats an already-consumed nonce as a hard stop, so building the
	// real payload while the rehearsal is unmined would cost this chain the whole firing window.
	// That is the gate operators are told to respect; waiting on it here is what tests it.
	require.Eventually(r, func() bool {
		nonce, err := r.EVMClient.NonceAt(r.Ctx, r.TSSAddress, nil)
		return err == nil && nonce > rehearsalNonce
	}, 3*time.Minute, 5*time.Second, "rehearsal EVM tx did not confirm, so the real payload would pin a stale nonce")

	fullEVM := buildDrainEVMTxs(r, evmReceiver, sdkmath.ZeroUint())
	fullBTC := buildDrainBTCTxs(r, btcReceiver, 0)
	require.NotEmpty(r, fullBTC, "expected at least one BTC sweep")
	require.Greater(r, fullEVM[0].Nonce, rehearsalNonce, "the real drain must pin a fresh nonce")

	fullHeight := currentZetaHeight(r) + drainTriggerOffset
	require.Greater(r, fullHeight, rehearsalHeight, "the real drain needs a higher trigger height")
	publishDrainPayload(r, server, priv, fullHeight, 2, fullEVM, fullBTC)
	r.Logger.Print("serving FULL payload for trigger height %d", fullHeight)

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

	// Wait for every funding UTXO, not just the first. A UTXO confirming late would appear at the
	// TSS mid-test and break the "drained to zero" assertion in the second phase.
	require.Eventually(r, func() bool {
		return tssBTCUTXOCount(r) >= drainBTCNumUTXOs
	}, 3*time.Minute, 5*time.Second, "TSS BTC funding did not confirm")
}

// buildDrainEVMTxs builds the pinned EVM drain tx from live TSS state. maxAmount caps the transfer
// for a rehearsal; a nil/zero Uint means the full balance, i.e. the real drain.
func buildDrainEVMTxs(
	r *runner.E2ERunner,
	receiver ethcommon.Address,
	maxAmount sdkmath.Uint,
) []draintx.EVMTx {
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
		MaxAmount:      maxAmount,
	})
	require.NoError(r, err)
	return []draintx.EVMTx{tx}
}

// buildDrainBTCTxs builds the pinned BTC sweeps from the live TSS UTXO set. maxSats caps the total
// swept for a rehearsal; zero means sweep everything.
func buildDrainBTCTxs(r *runner.E2ERunner, receiver btcutil.Address, maxSats int64) []draintx.BTCTx {
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
		MaxSats: maxSats,
	})
	require.NoError(r, err)
	return btcTxs
}

// tssBTCUTXOCount is the number of UTXOs still held by the TSS. The rehearsal assertion needs the
// count, not just the total: a capped sweep must leave whole UTXOs behind, since it cannot split one.
func tssBTCUTXOCount(r *runner.E2ERunner) int {
	unspent, err := r.BtcRPCClient.ListUnspentMinMaxAddresses(r.Ctx, 1, 9999999, []btcutil.Address{r.BTCTSSAddress})
	require.NoError(r, err)
	return len(unspent)
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

// startPayloadServer binds the payload server on the port the zetaclients are armed to poll. It
// stays up across both phases so the second payload can replace the first in place — rebinding the
// port between phases would race the pollers' fetches.
func startPayloadServer(r *runner.E2ERunner, server *pkgdrain.PayloadServer, drainURL string) {
	u, err := url.Parse(drainURL)
	require.NoError(r, err)
	require.NoError(r, server.Start(":"+u.Port()))
}

// publishDrainPayload signs and serves one final payload. Replacing the served payload with a
// higher trigger height is how the operator retries or follows a rehearsal with the real drain: the
// poller ignores anything at or below the height it last acted on, and arms on anything above it.
func publishDrainPayload(
	r *runner.E2ERunner,
	server *pkgdrain.PayloadServer,
	priv *ecdsa.PrivateKey,
	triggerHeight int64,
	seq uint64,
	evmTxs []draintx.EVMTx,
	btcTxs []draintx.BTCTx,
) {
	payload, err := pkgdrain.BuildPayload(
		triggerHeight,
		seq,
		true,
		pkgdrain.NetworkLocalnet,
		evmTxs,
		btcTxs,
		priv,
	)
	require.NoError(r, err)
	require.NoError(r, server.Publish(payload))
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
