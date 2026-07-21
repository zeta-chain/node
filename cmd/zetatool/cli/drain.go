package cli

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcutil"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/cmd/zetatool/clients"
	zetatoolcommon "github.com/zeta-chain/node/cmd/zetatool/common"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	pkgchains "github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/drain"
	"github.com/zeta-chain/node/pkg/draintx"
	"github.com/zeta-chain/node/pkg/rpc"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	btccommon "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

const (
	// FlagTriggerHeight is the zeta block height at which clients fire.
	FlagTriggerHeight = "trigger-height"
	// FlagFinal marks the payload as the single authoritative version clients may sign.
	FlagFinal = "final"
	// FlagSigningKey is the hex-encoded secp256k1 operator private key that signs the payload.
	FlagSigningKey = "signing-key"
	// FlagSeq is the monotonic payload version (observability only).
	FlagSeq = "seq"
	// FlagFeeRate is the BTC fee rate in sat/vB pinned into the sweeps.
	FlagFeeRate = "fee-rate"
	// FlagServe runs the draft->freeze->final cron, serving the payload over HTTP.
	FlagServe = "serve"
	// FlagServeAddr is the address the serve-mode HTTP server binds.
	FlagServeAddr = "serve-addr"
	// FlagFreezeWindow K: publish the single final once currentHeight >= triggerHeight - K.
	FlagFreezeWindow = "freeze-window"
	// FlagInterval is how often the cron recomputes and republishes drafts.
	FlagInterval = "interval"
)

// NewDrainPayloadCMD creates the command that builds and signs an emergency drain payload.
func NewDrainPayloadCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drain-payload <chain>",
		Short: "Build and sign an emergency TSS drain payload (EVM + BTC)",
		Long: `Build the fully-resolved, signed drain payload that moves all native TSS funds to
the hardcoded safe wallet. The operator supplies only the trigger height; balances,
median gas prices, nonces and UTXOs are derived from the configured RPCs.

Without --serve it prints a single signed payload to stdout. With --serve it runs the
draft->freeze->final cron: it recomputes and republishes drafts (final:false) every
--interval over HTTP, and once the zeta height reaches (trigger-height - freeze-window)
it publishes exactly one final:true payload and stops.

The chain argument selects the network (mainnet/testnet/localnet) the same way as
tss-balances.`,
		Args: cobra.ExactArgs(1),
		RunE: runDrainPayload,
	}

	cmd.Flags().Int64(FlagTriggerHeight, 0, "zeta block height at which clients fire (required)")
	cmd.Flags().Bool(FlagFinal, false, "mark the one-shot payload as final (ignored with --serve)")
	cmd.Flags().String(FlagSigningKey, "", "hex-encoded secp256k1 operator private key (required)")
	cmd.Flags().Uint64(FlagSeq, 0, "monotonic payload version (one-shot mode)")
	cmd.Flags().Int64(FlagFeeRate, conservativeFeeRate, "BTC fee rate in sat/vB")
	cmd.Flags().Bool(FlagServe, false, "run the draft->freeze->final cron and serve over HTTP")
	cmd.Flags().String(FlagServeAddr, ":8899", "address the serve-mode HTTP server binds")
	cmd.Flags().Int64(FlagFreezeWindow, 20, "blocks before trigger-height at which to freeze and publish the final")
	cmd.Flags().Duration(FlagInterval, 10*time.Second, "serve-mode republish interval")

	return cmd
}

func runDrainPayload(cmd *cobra.Command, args []string) error {
	gen, opts, err := setupGenerator(cmd, args[0])
	if err != nil {
		return err
	}

	if opts.serve {
		return serveDrain(cmd.Context(), gen, opts)
	}

	payload, err := gen.generate(context.Background(), opts.seq, opts.final)
	if err != nil {
		return fmt.Errorf("failed to build payload: %w", err)
	}
	out, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

// drainOptions holds the parsed command flags.
type drainOptions struct {
	serve       bool
	serveAddr   string
	final       bool
	seq         uint64
	interval    time.Duration
	freezeK     int64
	triggerHigh int64
}

// setupGenerator parses flags and builds a payloadGenerator + options.
func setupGenerator(cmd *cobra.Command, chainArg string) (*payloadGenerator, drainOptions, error) {
	var opts drainOptions

	chain, err := zetatoolcommon.ResolveChain(chainArg)
	if err != nil {
		return nil, opts, fmt.Errorf("failed to resolve chain %q: %w", chainArg, err)
	}
	network := zetatoolcommon.NetworkTypeFromChain(chain)

	opts.triggerHigh = must(cmd.Flags().GetInt64(FlagTriggerHeight))
	if opts.triggerHigh <= 0 {
		return nil, opts, fmt.Errorf("--%s is required and must be positive", FlagTriggerHeight)
	}
	opts.final = must(cmd.Flags().GetBool(FlagFinal))
	opts.seq = must(cmd.Flags().GetUint64(FlagSeq))
	opts.serve = must(cmd.Flags().GetBool(FlagServe))
	opts.serveAddr = must(cmd.Flags().GetString(FlagServeAddr))
	opts.interval = must(cmd.Flags().GetDuration(FlagInterval))
	opts.freezeK = must(cmd.Flags().GetInt64(FlagFreezeWindow))
	feeRate := must(cmd.Flags().GetInt64(FlagFeeRate))

	priv, err := ethcrypto.HexToECDSA(strings.TrimPrefix(must(cmd.Flags().GetString(FlagSigningKey)), "0x"))
	if err != nil {
		return nil, opts, fmt.Errorf("invalid --%s: %w", FlagSigningKey, err)
	}

	cfg, err := config.GetConfigByNetwork(network, must(cmd.Flags().GetString(config.FlagConfig)))
	if err != nil {
		return nil, opts, fmt.Errorf("failed to get config: %w", err)
	}
	if cfg.ZetaChainRPC == "" {
		return nil, opts, fmt.Errorf("ZetaChainRPC is not configured for network %s", network)
	}

	receivers, err := drain.ReceiverForNetwork(drainNetwork(network))
	if err != nil {
		return nil, opts, err
	}
	if err := receivers.Validate(); err != nil {
		return nil, opts, fmt.Errorf("drain receivers not configured for %s: %w", network, err)
	}

	zetacoreClient, err := rpc.NewCometBFTClients(cfg.ZetaChainRPC)
	if err != nil {
		return nil, opts, fmt.Errorf("failed to create zetacore client: %w", err)
	}

	btcChainID, err := clients.GetBTCChainID(network)
	if err != nil {
		return nil, opts, fmt.Errorf("failed to get BTC chain ID: %w", err)
	}

	return &payloadGenerator{
		cfg:           cfg,
		zetacore:      zetacoreClient,
		btcChainID:    btcChainID,
		triggerHeight: opts.triggerHigh,
		evmReceiver:   receivers.EVM,
		btcReceiver:   receivers.BTC,
		feeRate:       feeRate,
		priv:          priv,
	}, opts, nil
}

// payloadGenerator builds a signed payload from live chain state; reused per cron tick.
type payloadGenerator struct {
	cfg           *config.Config
	zetacore      rpc.Clients
	btcChainID    int64
	triggerHeight int64
	evmReceiver   string
	btcReceiver   string
	feeRate       int64
	priv          *ecdsa.PrivateKey
}

// generate resolves live balances/gas/nonces/UTXOs and returns a signed payload.
func (g *payloadGenerator) generate(ctx context.Context, seq uint64, final bool) (draintx.Payload, error) {
	tss, err := currentTSS(ctx, g.zetacore)
	if err != nil {
		return draintx.Payload{}, err
	}
	tssAddrRes, err := g.zetacore.Observer.GetTssAddressByFinalizedHeight(
		ctx,
		&observertypes.QueryGetTssAddressByFinalizedHeightRequest{
			FinalizedZetaHeight: tss.FinalizedZetaHeight,
			BitcoinChainId:      g.btcChainID,
		},
	)
	if err != nil {
		return draintx.Payload{}, fmt.Errorf("failed to get TSS addresses: %w", err)
	}
	supportedChains, err := g.zetacore.Observer.SupportedChains(ctx, &observertypes.QuerySupportedChains{})
	if err != nil {
		return draintx.Payload{}, fmt.Errorf("failed to get supported chains: %w", err)
	}

	evmTxs, err := buildEVMTxs(
		ctx,
		g.cfg,
		g.zetacore,
		supportedChains.Chains,
		ethcommon.HexToAddress(tssAddrRes.Eth),
		g.evmReceiver,
	)
	if err != nil {
		return draintx.Payload{}, err
	}
	btcTxs, err := buildBTCTxs(ctx, g.cfg, supportedChains.Chains, tssAddrRes.Btc, g.btcReceiver, g.feeRate)
	if err != nil {
		return draintx.Payload{}, err
	}
	return drain.BuildPayload(g.triggerHeight, seq, final, evmTxs, btcTxs, g.priv)
}

// serveDrain runs the draft->freeze->final cron, serving the payload over HTTP.
func serveDrain(ctx context.Context, gen *payloadGenerator, opts drainOptions) error {
	if ctx == nil {
		ctx = context.Background()
	}
	server := drain.NewPayloadServer()
	if err := server.Start(opts.serveAddr); err != nil {
		return fmt.Errorf("failed to start payload server: %w", err)
	}
	defer server.Close()
	fmt.Fprintf(os.Stderr, "serving drain payload at %s (trigger=%d, freeze-window=%d)\n",
		server.URL(), opts.triggerHigh, opts.freezeK)

	isFinalTime := func(ctx context.Context) (bool, error) {
		height, err := zetaHeight(ctx, gen.zetacore)
		if err != nil {
			return false, err
		}
		return height >= opts.triggerHigh-opts.freezeK, nil
	}

	return drain.RunCron(ctx, opts.interval, gen.generate, server.Publish, isFinalTime)
}

// zetaHeight returns the latest zeta block height.
func zetaHeight(ctx context.Context, c rpc.Clients) (int64, error) {
	res, err := c.Crosschain.LastZetaHeight(ctx, &crosschaintypes.QueryLastZetaHeightRequest{})
	if err != nil {
		return 0, fmt.Errorf("get zeta height: %w", err)
	}
	return res.Height, nil
}

// currentTSS returns the TSS with the highest finalized height.
func currentTSS(ctx context.Context, c rpc.Clients) (observertypes.TSS, error) {
	res, err := c.Observer.TssHistory(ctx, &observertypes.QueryTssHistoryRequest{})
	if err != nil {
		return observertypes.TSS{}, fmt.Errorf("failed to fetch TSS history: %w", err)
	}
	return latestTSS(res.TssList)
}

// latestTSS picks the TSS entry with the highest finalized zeta height.
func latestTSS(list []observertypes.TSS) (observertypes.TSS, error) {
	if len(list) == 0 {
		return observertypes.TSS{}, fmt.Errorf("no TSS entries found")
	}
	current := list[0]
	for _, t := range list {
		if t.FinalizedZetaHeight > current.FinalizedZetaHeight {
			current = t
		}
	}
	return current, nil
}

func buildEVMTxs(
	ctx context.Context,
	cfg *config.Config,
	zetacoreClient rpc.Clients,
	supportedChains []pkgchains.Chain,
	tssAddr ethcommon.Address,
	receiver string,
) ([]draintx.EVMTx, error) {
	var (
		evmTxs   []draintx.EVMTx
		included []int64
	)
	for _, chain := range supportedChains {
		if !chain.IsExternal || chain.Vm != pkgchains.Vm_evm {
			continue
		}
		rpcURL := getRPCForChain(cfg, chain)
		if rpcURL == "" {
			fmt.Fprintf(os.Stderr, "WARN chain %d skipped: RPC not configured\n", chain.ChainId)
			continue
		}

		balance, err := clients.GetEVMBalance(ctx, rpcURL, tssAddr)
		if err != nil {
			return nil, fmt.Errorf("get balance for chain %d: %w", chain.ChainId, err)
		}
		if balance.Sign() == 0 {
			fmt.Fprintf(os.Stderr, "WARN chain %d skipped: zero balance\n", chain.ChainId)
			continue
		}

		median, err := medianGasPrice(ctx, zetacoreClient, chain.ChainId)
		if err != nil {
			return nil, err
		}
		nonce, err := clients.GetEVMNonce(ctx, rpcURL, tssAddr)
		if err != nil {
			return nil, fmt.Errorf("get nonce for chain %d: %w", chain.ChainId, err)
		}

		tx, err := drain.GenerateEVMTx(drain.EVMInput{
			ChainID:        chain.ChainId,
			To:             receiver,
			Balance:        sdkmath.NewUintFromString(balance.String()),
			MedianGasPrice: median,
			Nonce:          nonce,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARN chain %d skipped: %v\n", chain.ChainId, err)
			continue
		}
		evmTxs = append(evmTxs, tx)
		included = append(included, chain.ChainId)
	}
	fmt.Fprintf(os.Stderr, "EVM drain: %d chains included %v\n", len(included), included)
	return evmTxs, nil
}

func buildBTCTxs(
	ctx context.Context,
	cfg *config.Config,
	supportedChains []pkgchains.Chain,
	tssBtcAddr string,
	receiver string,
	feeRate int64,
) ([]draintx.BTCTx, error) {
	var btcChain *pkgchains.Chain
	for i := range supportedChains {
		if supportedChains[i].IsExternal && supportedChains[i].Vm == pkgchains.Vm_no_vm {
			btcChain = &supportedChains[i]
			break
		}
	}
	if btcChain == nil {
		return nil, nil
	}

	netParams, err := pkgchains.BitcoinNetParamsFromChainID(btcChain.ChainId)
	if err != nil {
		return nil, fmt.Errorf("bitcoin net params: %w", err)
	}
	receiverAddr, err := btcutil.DecodeAddress(receiver, netParams)
	if err != nil {
		return nil, fmt.Errorf("invalid BTC receiver %q: %w", receiver, err)
	}
	tssAddr, err := btcutil.DecodeAddress(tssBtcAddr, netParams)
	if err != nil {
		return nil, fmt.Errorf("invalid TSS BTC address %q: %w", tssBtcAddr, err)
	}

	btcAdapter, err := clients.NewBitcoinClientAdapter(cfg, *btcChain, zerolog.Nop())
	if err != nil {
		return nil, fmt.Errorf("create bitcoin client: %w", err)
	}
	unspent, err := btcAdapter.ListUnspentByAddress(ctx, tssAddr)
	if err != nil {
		return nil, fmt.Errorf("list unspent: %w", err)
	}

	utxos := make([]drain.UTXO, 0, len(unspent))
	var totalSats int64
	for _, u := range unspent {
		sats, err := btccommon.GetSatoshis(u.Amount)
		if err != nil {
			return nil, fmt.Errorf("convert utxo amount: %w", err)
		}
		totalSats += sats
		utxos = append(utxos, drain.UTXO{TxID: u.TxID, Vout: u.Vout, AmountSats: sats})
	}
	fmt.Fprintf(os.Stderr, "BTC drain: %d UTXOs (%d sats) for chain %d\n", len(utxos), totalSats, btcChain.ChainId)

	txs, err := drain.GenerateBTCTxs(drain.BTCInput{
		ChainID: btcChain.ChainId,
		To:      receiverAddr,
		FeeRate: feeRate,
		UTXOs:   utxos,
	})
	if err != nil {
		return nil, err
	}

	// surface dust that GenerateBTCTxs skipped as uneconomical, so it isn't dropped silently.
	var sweptInputs int
	var sweptSats int64
	for _, tx := range txs {
		sweptInputs += len(tx.Inputs)
		for _, in := range tx.Inputs {
			sweptSats += in.AmountSats
		}
	}
	fmt.Fprintf(os.Stderr, "BTC drain: swept %d UTXOs (%d sats) in %d txs, skipped %d UTXOs (%d sats) as uneconomical\n",
		sweptInputs, sweptSats, len(txs), len(utxos)-sweptInputs, totalSats-sweptSats)

	return txs, nil
}

// medianGasPrice queries zetacore for the median gas price of a chain.
func medianGasPrice(ctx context.Context, c rpc.Clients, chainID int64) (sdkmath.Uint, error) {
	res, err := c.Crosschain.GasPrice(ctx, &crosschaintypes.QueryGetGasPriceRequest{
		Index: strconv.FormatInt(chainID, 10),
	})
	if err != nil {
		return sdkmath.ZeroUint(), fmt.Errorf("get gas price for chain %d: %w", chainID, err)
	}
	return pickMedian(res.GasPrice, chainID)
}

// pickMedian extracts the median gas price from a GasPrice record.
func pickMedian(gp *crosschaintypes.GasPrice, chainID int64) (sdkmath.Uint, error) {
	if gp == nil || len(gp.Prices) == 0 || gp.MedianIndex >= uint64(len(gp.Prices)) {
		return sdkmath.ZeroUint(), fmt.Errorf("no median gas price for chain %d", chainID)
	}
	return sdkmath.NewUint(gp.Prices[gp.MedianIndex]), nil
}

// drainNetwork maps a zetatool network string to a drain receiver network.
func drainNetwork(network string) string {
	switch network {
	case config.NetworkMainnet:
		return drain.NetworkMainnet
	case config.NetworkLocalnet:
		return drain.NetworkLocalnet
	default:
		return drain.NetworkTestnet
	}
}

// must panics on a flag-parse error; flags are statically defined so this cannot fail.
func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
