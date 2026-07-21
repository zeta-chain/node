package cli

import (
	"context"
	"math/big"
	"os"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcutil"
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
	// tb1q... is the BIP173 testnet bech32 example address; valid on testnet3.
	liveBTCReceiver = "tb1qw508d6qejxtdg4y5r3zarvary0c5xw7kxpjzsx"
)

func requireLiveTest(t *testing.T) {
	if os.Getenv("ZETATOOL_LIVE_TEST") == "" {
		t.Skip("set ZETATOOL_LIVE_TEST=1 to run live-RPC drain tests")
	}
}

// liveTestnet resolves the current testnet TSS state exactly like payloadGenerator.generate does.
func liveTestnet(t *testing.T) (
	context.Context,
	*config.Config,
	rpc.Clients,
	*observertypes.QueryGetTssAddressByFinalizedHeightResponse,
	[]pkgchains.Chain,
	int64,
) {
	ctx := context.Background()
	cfg := config.TestnetConfig()

	zetacore, err := rpc.NewCometBFTClients(cfg.ZetaChainRPC)
	require.NoError(t, err)

	btcChainID, err := clients.GetBTCChainID(config.NetworkTestnet)
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

	ctx, cfg, zetacore, tssAddrRes, chains, _ := liveTestnet(t)
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

	ctx, _, _, tssAddrRes, chains, btcChainID := liveTestnet(t)

	var btcChain *pkgchains.Chain
	for i := range chains {
		if chains[i].IsExternal && chains[i].Vm == pkgchains.Vm_no_vm {
			btcChain = &chains[i]
			break
		}
	}
	if btcChain == nil {
		t.Skip("no external BTC chain in supported chains")
	}

	btcAddr := tssAddrRes.Btc
	t.Logf("TSS BTC address: %s", btcAddr)

	netParams, err := pkgchains.BitcoinNetParamsFromChainID(btcChain.ChainId)
	require.NoError(t, err)
	receiver, err := btcutil.DecodeAddress(liveBTCReceiver, netParams)
	require.NoError(t, err)

	utxos, err := clients.GetBTCUtxos(ctx, btcAddr, btcChainID)
	require.NoError(t, err)
	if len(utxos) == 0 {
		t.Skipf("TSS BTC address %s has no UTXOs", btcAddr)
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
		ChainID: btcChain.ChainId,
		To:      receiver,
		FeeRate: feeRate,
		UTXOs:   drainUTXOs,
	})
	require.NoError(t, err)

	expectedFee := migration.BTCConservativeFeeRate*btccommon.OutboundBytesMax +
		migration.BTCReservedRBFFeeSats + migration.BTCNonceMarkBufferSats

	for _, tx := range txs {
		var sumInputs int64
		for _, in := range tx.Inputs {
			sumInputs += in.AmountSats
		}
		t.Logf(
			"chain %d: to=%s inputs=%d outputSats=%d feeSats=%d",
			tx.ChainID, tx.To, len(tx.Inputs), tx.OutputSats, tx.FeeSats,
		)
		require.Equal(t, sumInputs-tx.FeeSats, tx.OutputSats)
		require.Equal(t, expectedFee, tx.FeeSats)
		require.GreaterOrEqual(t, tx.OutputSats, int64(constant.BTCWithdrawalDustAmount))
	}
}
