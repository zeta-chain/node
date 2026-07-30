package cli

import (
	"context"
	"math/big"
	"os"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/cmd/zetatool/clients"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	pkgchains "github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/drain"
	"github.com/zeta-chain/node/pkg/migration"
	"github.com/zeta-chain/node/pkg/rpc"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	btccommon "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

// throwaway receivers: these tests only build and inspect drain txs, never broadcast them.
const (
	liveEVMReceiver = "0x000000000000000000000000000000000000dEaD"
	// bc1q.../tb1q... are the BIP173 bech32 P2WPKH example addresses for mainnet/testnet.
	liveBTCReceiverMainnet = "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4"
	liveBTCReceiverTestnet = "tb1qw508d6qejxtdg4y5r3zarvary0c5xw7kxpjzsx"
)

func requireLiveTest(t *testing.T) {
	if os.Getenv("ZETATOOL_LIVE_TEST") == "" {
		t.Skip("set ZETATOOL_LIVE_TEST=1 to run live-RPC drain tests")
	}
}

// liveNetwork returns the network the live tests run against: testnet by default, or the value
// of ZETATOOL_LIVE_NETWORK when set (used to inspect a real BTC sweep on mainnet).
func liveNetwork(t *testing.T) string {
	network := os.Getenv("ZETATOOL_LIVE_NETWORK")
	if network == "" {
		return config.NetworkTestnet
	}
	switch network {
	case config.NetworkTestnet, config.NetworkMainnet:
		return network
	default:
		t.Fatalf("unsupported ZETATOOL_LIVE_NETWORK %q: use testnet or mainnet", network)
		return ""
	}
}

// liveState resolves the current TSS state on the selected network, exactly like
// payloadGenerator.generate does.
func liveState(t *testing.T) (
	context.Context,
	*config.Config,
	rpc.Clients,
	*observertypes.QueryGetTssAddressByFinalizedHeightResponse,
	[]pkgchains.Chain,
	int64,
) {
	ctx := context.Background()
	network := liveNetwork(t)
	t.Logf("live network: %s", network)

	cfg, err := config.GetConfigByNetwork(network, "")
	require.NoError(t, err)

	zetacore, err := rpc.NewCometBFTClients(cfg.ZetaChainRPC)
	require.NoError(t, err)

	btcChainID, err := clients.GetBTCChainID(network)
	require.NoError(t, err)

	tss, err := currentTSS(ctx, zetacore)
	require.NoError(t, err)
	t.Logf("current TSS pubkey=%s finalizedHeight=%d", tss.TssPubkey, tss.FinalizedZetaHeight)

	tssAddrRes, err := zetacore.Observer.GetTssAddressByFinalizedHeight(
		ctx,
		&observertypes.QueryGetTssAddressByFinalizedHeightRequest{
			FinalizedZetaHeight: tss.FinalizedZetaHeight,
			BitcoinChainId:      btcChainID,
		},
	)
	require.NoError(t, err)

	supported, err := zetacore.Observer.SupportedChains(ctx, &observertypes.QuerySupportedChains{})
	require.NoError(t, err)

	return ctx, cfg, zetacore, tssAddrRes, supported.Chains, btcChainID
}

func TestLiveEVMDrain(t *testing.T) {
	requireLiveTest(t)

	ctx, cfg, zetacore, tssAddrRes, chains, _ := liveState(t)
	tssAddr := ethcommon.HexToAddress(tssAddrRes.Eth)
	t.Logf("TSS EVM address: %s", tssAddr.Hex())

	buffer, ok := new(big.Int).SetString(migration.BufferAmountEVM, 10)
	require.True(t, ok)

	for _, chain := range chains {
		if !chain.IsExternal || chain.Vm != pkgchains.Vm_evm {
			continue
		}
		rpcURL := getRPCForChain(cfg, chain)
		if rpcURL == "" {
			t.Logf("chain %d (%s): skipped, RPC not configured", chain.ChainId, chain.Name)
			continue
		}

		balance, err := clients.GetEVMBalance(ctx, rpcURL, tssAddr)
		if err != nil {
			t.Logf("chain %d (%s): skipped, balance error: %v", chain.ChainId, chain.Name, err)
			continue
		}
		if balance.Sign() == 0 {
			t.Logf("chain %d (%s): skipped, zero balance", chain.ChainId, chain.Name)
			continue
		}

		median, err := medianGasPrice(ctx, zetacore, chain.ChainId)
		if err != nil {
			t.Logf("chain %d (%s): skipped, gas price error: %v", chain.ChainId, chain.Name, err)
			continue
		}
		nonce, err := clients.GetEVMNonce(ctx, rpcURL, tssAddr)
		if err != nil {
			t.Logf("chain %d (%s): skipped, nonce error: %v", chain.ChainId, chain.Name, err)
			continue
		}

		tx, err := drain.GenerateEVMTx(drain.EVMInput{
			ChainID:        chain.ChainId,
			To:             liveEVMReceiver,
			Balance:        sdkmath.NewUintFromString(balance.String()),
			MedianGasPrice: median,
			Nonce:          nonce,
		})
		if err != nil {
			t.Logf("chain %d (%s): skipped, build error: %v", chain.ChainId, chain.Name, err)
			continue
		}

		t.Logf(
			"chain %d (%s): balance=%s to=%s amount=%s gasPrice=%s gasLimit=%d nonce=%d",
			chain.ChainId, chain.Name, balance.String(), tx.To, tx.Amount, tx.GasPrice, tx.GasLimit, tx.Nonce,
		)

		// amount == balance - (gasLimit*gasPrice + BufferAmountEVM)
		gasPrice, ok := new(big.Int).SetString(tx.GasPrice, 10)
		require.True(t, ok)
		fee := new(big.Int).Mul(new(big.Int).SetUint64(tx.GasLimit), gasPrice)
		fee.Add(fee, buffer)
		expected := new(big.Int).Sub(balance, fee)
		amount, ok := new(big.Int).SetString(tx.Amount, 10)
		require.True(t, ok)
		require.Equal(t, expected.String(), amount.String())
		require.GreaterOrEqual(t, amount.Sign(), 0)
	}
}

func TestLiveBTCDrain(t *testing.T) {
	requireLiveTest(t)

	ctx, _, _, tssAddrRes, _, btcChainID := liveState(t)

	// The TSS BTC address is derived from the pubkey and returned regardless of whether BTC is an
	// active supported chain, so we query it directly by address+chainID via mempool.space.
	btcAddr := tssAddrRes.Btc
	if btcAddr == "" {
		t.Skipf("no TSS BTC address derived on %s", liveNetwork(t))
	}
	t.Logf("TSS BTC address: %s", btcAddr)

	netParams, err := pkgchains.BitcoinNetParamsFromChainID(btcChainID)
	require.NoError(t, err)
	btcReceiver := liveBTCReceiverTestnet
	if netParams.Name == chaincfg.MainNetParams.Name {
		btcReceiver = liveBTCReceiverMainnet
	}
	receiver, err := btcutil.DecodeAddress(btcReceiver, netParams)
	require.NoError(t, err)

	balance, err := clients.GetBTCBalance(ctx, btcAddr, btcChainID)
	require.NoError(t, err)
	t.Logf("TSS BTC balance: %.8f BTC", balance)

	utxos, err := clients.GetBTCUtxos(ctx, btcAddr, btcChainID)
	require.NoError(t, err)
	if len(utxos) == 0 {
		t.Skipf(
			"TSS BTC address %s on %s has no UTXOs; try ZETATOOL_LIVE_NETWORK=mainnet",
			btcAddr, liveNetwork(t),
		)
	}

	var totalSats int64
	drainUTXOs := make([]drain.UTXO, len(utxos))
	for i, u := range utxos {
		totalSats += u.Value
		drainUTXOs[i] = drain.UTXO{TxID: u.TxID, Vout: u.Vout, AmountSats: u.Value}
	}
	t.Logf("fetched %d UTXOs totalling %d sats", len(utxos), totalSats)

	feeRate := migration.BTCConservativeFeeRate
	txs, err := drain.GenerateBTCTxs(drain.BTCInput{
		ChainID: btcChainID,
		To:      receiver,
		FeeRate: feeRate,
		UTXOs:   drainUTXOs,
	})
	require.NoError(t, err)

	var sweptInputs int
	var sweptSats, totalFees int64
	for _, tx := range txs {
		var sumInputs int64
		for _, in := range tx.Inputs {
			sumInputs += in.AmountSats
		}
		sweptInputs += len(tx.Inputs)
		sweptSats += sumInputs
		totalFees += tx.FeeSats
		t.Logf(
			"chain %d: to=%s inputs=%d outputSats=%d feeSats=%d",
			tx.ChainID, tx.To, len(tx.Inputs), tx.OutputSats, tx.FeeSats,
		)
		// miner fee only (no RBF / nonce reserve), right-sized to the input count
		wantSize, err := btccommon.EstimateOutboundSize(int64(len(tx.Inputs)), []btcutil.Address{receiver})
		require.NoError(t, err)
		require.Equal(t, feeRate*wantSize, tx.FeeSats)
		require.Equal(t, sumInputs-tx.FeeSats, tx.OutputSats)
		require.GreaterOrEqual(t, tx.OutputSats, int64(constant.BTCWithdrawalDustAmount))
		require.LessOrEqual(t, tx.FeeSats, sumInputs/drain.MaxBTCFeeFraction) // poller fee cap
	}

	// uneconomical dust groups are skipped, so swept UTXOs may be fewer than fetched
	t.Logf(
		"swept %d UTXOs (%d sats) in %d txs; skipped %d UTXOs (%d sats); total miner fees %d sats",
		sweptInputs, sweptSats, len(txs), len(utxos)-sweptInputs, totalSats-sweptSats, totalFees,
	)
}
