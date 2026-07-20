package e2etests

import (
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
	pkgdrain "github.com/zeta-chain/node/pkg/drain"
	"github.com/zeta-chain/node/pkg/draintx"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// drainFeeRate is a fixed sat/vB fee rate for the regtest sweep.
const drainFeeRate = 10

// TestDrainTSS exercises the emergency drain end to end (EVM + BTC), drain-only: no keygen
// and no MsgUpdateTssAddress. It disables inbound, builds and serves a signed final payload
// that sweeps all native TSS funds to the compiled-in localnet safe receivers, waits for the
// txs to mine, and asserts the TSS balances drop to ~0.
//
// The drain poller runs inside the zetaclient processes, so this test only runs when the
// localnet zetaclients are built with `-tags drain` and armed to poll this test's server.
// It is skipped unless both are provided:
//   - ZETACLIENT_DRAIN_URL: the endpoint the zetaclients poll (this test serves it)
//   - ZETACLIENT_DRAIN_SIGNING_KEY: hex private key whose pubkey is compiled into the build
func TestDrainTSS(r *runner.E2ERunner, _ []string) {
	drainURL := os.Getenv("ZETACLIENT_DRAIN_URL")
	signingKeyHex := os.Getenv("ZETACLIENT_DRAIN_SIGNING_KEY")
	if drainURL == "" || signingKeyHex == "" {
		r.Logger.Print("skipping drain_tss: ZETACLIENT_DRAIN_URL / ZETACLIENT_DRAIN_SIGNING_KEY not set")
		return
	}

	priv, err := ethcrypto.HexToECDSA(signingKeyHex)
	require.NoError(r, err)

	r.SetupBtcAddress(false)
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// pause inbound so the TSS nonce stays stable while we drain
	msgDisable := observertypes.NewMsgDisableCCTX(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.EmergencyPolicyName),
		false,
		true,
	)
	_, err = r.ZetaTxServer.BroadcastTx(utils.EmergencyPolicyName, msgDisable)
	require.NoError(r, err)
	defer reEnableInbound(r)

	// compiled-in safe receivers — the poller enforces these
	receivers, err := pkgdrain.ReceiverForNetwork(pkgdrain.NetworkLocalnet)
	require.NoError(r, err)
	evmReceiver := ethcommon.HexToAddress(receivers.EVM)
	btcReceiver, err := btcutil.DecodeAddress(receivers.BTC, r.BitcoinParams)
	require.NoError(r, err)

	// record TSS balances before the drain
	ethTSSBefore, err := r.EVMClient.BalanceAt(r.Ctx, r.TSSAddress, nil)
	require.NoError(r, err)

	// build the fully-resolved payload
	evmTxs := buildDrainEVMTxs(r, evmReceiver)
	btcTxs := buildDrainBTCTxs(r, btcReceiver)

	// trigger a few blocks ahead so the poller sees the final in time
	current := currentZetaHeight(r)
	triggerHeight := current + 10

	payload, err := pkgdrain.BuildPayload(triggerHeight, 1, true, evmTxs, btcTxs, priv)
	require.NoError(r, err)

	// serve the payload locally at the port the zetaclients poll
	server := servePayload(r, drainURL, payload)
	defer server.Close()

	// wait for the EVM drain to mine: TSS ETH balance drops to ~0
	require.Eventually(r, func() bool {
		bal, err := r.EVMClient.BalanceAt(r.Ctx, r.TSSAddress, nil)
		return err == nil && bal.Sign() == 0
	}, 5*time.Minute, 5*time.Second, "TSS ETH balance did not drain")

	// wait for the BTC sweeps to mine: TSS UTXO balance drops to ~dust
	require.Eventually(r, func() bool {
		utxos, err := r.GetTop20UTXOsForTssAddress()
		if err != nil {
			return false
		}
		var total float64
		for _, u := range utxos {
			total += u.Amount
		}
		return total == 0
	}, 5*time.Minute, 5*time.Second, "TSS BTC balance did not drain")

	// receivers should have increased
	evmReceiverBal, err := r.EVMClient.BalanceAt(r.Ctx, evmReceiver, nil)
	require.NoError(r, err)
	require.Positive(r, evmReceiverBal.Sign())

	r.Logger.Info("drain complete: ETH TSS before %s, receiver %s", ethTSSBefore, evmReceiverBal)
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
	// list all TSS UTXOs (the generator partitions them into <=20-input sweeps)
	unspent, err := r.BtcRPCClient.ListUnspentMinMaxAddresses(r.Ctx, 0, 9999999, []btcutil.Address{r.BTCTSSAddress})
	require.NoError(r, err)

	utxos := make([]pkgdrain.UTXO, 0, len(unspent))
	for _, u := range unspent {
		// #nosec G115 e2e amounts always in range
		utxos = append(utxos, pkgdrain.UTXO{TxID: u.TxID, Vout: u.Vout, AmountSats: int64(u.Amount * 1e8)})
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
