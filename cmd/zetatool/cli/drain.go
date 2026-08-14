package cli

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
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
	"github.com/zeta-chain/node/pkg/migration"
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
	// FlagOnlyChains restricts the drain to these chain IDs (comma-separated); empty = all.
	FlagOnlyChains = "only-chains"
	// FlagExcludeChains drops these chain IDs from the drain (comma-separated).
	FlagExcludeChains = "exclude-chains"
	// FlagEVMMaxAmount caps the per-chain EVM transfer (wei) for a small-value rehearsal.
	FlagEVMMaxAmount = "evm-max-amount"
	// FlagBTCMaxSats caps the total BTC swept (sats) for a small-value rehearsal.
	FlagBTCMaxSats = "btc-max-sats"
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

--evm-max-amount and --btc-max-sats cap the value moved so the whole path (payload ->
TSS ceremony -> broadcast) can be rehearsed with a small amount before committing the
full balance. They are rehearsal-only: a capped payload leaves the remainder at the TSS
address, so the real drain is a second run without the caps at a higher trigger height.
Wait for the rehearsal txs to CONFIRM before the real payload freezes: the pinned nonce
is the confirmed one, and a still-pending tx makes the poller reject that chain for the
whole firing window.

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
	cmd.Flags().String(FlagOnlyChains, "", "comma-separated chain IDs to drain (default: all with funds)")
	cmd.Flags().String(FlagExcludeChains, "", "comma-separated chain IDs to skip")
	cmd.Flags().String(
		FlagEVMMaxAmount,
		"",
		"rehearsal only: cap the per-chain EVM transfer at this many wei (default: full balance)",
	)
	cmd.Flags().Int64(
		FlagBTCMaxSats,
		0,
		"rehearsal only: cap the total BTC swept at this many sats, by selecting a subset of UTXOs. "+
			"Whole UTXOs only, so a cap below the smallest sweepable UTXO sweeps nothing (reported) (default: all)",
	)

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
	if feeRate <= 0 {
		return nil, opts, fmt.Errorf("--%s must be positive, got %d", FlagFeeRate, feeRate)
	}
	if feeRate < migration.BTCConservativeFeeRate {
		fmt.Fprintf(
			os.Stderr,
			"WARN --%s %d is below the conservative default %d sat/vB; set it from the live mempool so sweeps confirm\n",
			FlagFeeRate,
			feeRate,
			migration.BTCConservativeFeeRate,
		)
	}

	filter, err := newChainFilter(
		must(cmd.Flags().GetString(FlagOnlyChains)),
		must(cmd.Flags().GetString(FlagExcludeChains)),
	)
	if err != nil {
		return nil, opts, err
	}

	evmMaxAmount, err := parseEVMMaxAmount(must(cmd.Flags().GetString(FlagEVMMaxAmount)))
	if err != nil {
		return nil, opts, err
	}
	btcMaxSats := must(cmd.Flags().GetInt64(FlagBTCMaxSats))
	if btcMaxSats < 0 {
		return nil, opts, fmt.Errorf("--%s must not be negative, got %d", FlagBTCMaxSats, btcMaxSats)
	}

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
		network:       drainNetwork(network),
		evmReceiver:   receivers.EVM,
		btcReceiver:   receivers.BTC,
		feeRate:       feeRate,
		filter:        filter,
		evmMaxAmount:  evmMaxAmount,
		btcMaxSats:    btcMaxSats,
		priv:          priv,
	}, opts, nil
}

// parseEVMMaxAmount parses the --evm-max-amount flag. Empty means no cap; the value is a decimal
// wei amount, matching how amounts are carried in the payload.
func parseEVMMaxAmount(s string) (sdkmath.Uint, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return sdkmath.ZeroUint(), nil
	}
	v, ok := new(big.Int).SetString(s, 10)
	if !ok || v.Sign() < 0 {
		return sdkmath.ZeroUint(), fmt.Errorf("invalid --%s %q: want a decimal wei amount", FlagEVMMaxAmount, s)
	}
	if v.Sign() == 0 {
		return sdkmath.ZeroUint(), fmt.Errorf(
			"--%s must be positive; omit the flag to drain the full balance",
			FlagEVMMaxAmount,
		)
	}
	// sdkmath.Uint is bounded at 256 bits and its constructor panics past that, which would take
	// down zetatool on an extra-zeros typo instead of reporting a bad flag.
	if v.BitLen() > 256 {
		return sdkmath.ZeroUint(), fmt.Errorf(
			"--%s %q exceeds the maximum 256-bit wei amount",
			FlagEVMMaxAmount,
			s,
		)
	}
	return sdkmath.NewUintFromBigInt(v), nil
}

// warnIfCapped makes a rehearsal payload impossible to mistake for the real drain: the caps only
// show up as smaller numbers inside the payload, which is easy to miss when the operator is
// reading it under pressure.
//
// The wording deliberately avoids claiming nothing is drained: a cap is an upper bound, so any
// chain whose drainable balance already sits below it is swept in full by a "rehearsal".
func warnIfCapped(w io.Writer, evmMaxAmount sdkmath.Uint, btcMaxSats int64) {
	if evmMaxAmount.IsZero() && btcMaxSats == 0 {
		return
	}
	evmCap := "uncapped"
	if !evmMaxAmount.IsZero() {
		evmCap = evmMaxAmount.String() + " wei per chain"
	}
	btcCap := "uncapped"
	if btcMaxSats > 0 {
		btcCap = fmt.Sprintf("%d sats total", btcMaxSats)
	}
	fmt.Fprintf(
		w,
		"WARN REHEARSAL PAYLOAD: value is capped (evm: %s, btc: %s). This does NOT drain the TSS, "+
			"except on chains already below the cap, which are swept in full. "+
			"Re-run without the caps at a higher trigger height for the real drain\n",
		evmCap,
		btcCap,
	)
}

// payloadGenerator builds a signed payload from live chain state; reused per cron tick.
type payloadGenerator struct {
	cfg           *config.Config
	zetacore      rpc.Clients
	btcChainID    int64
	triggerHeight int64
	network       string
	evmReceiver   string
	btcReceiver   string
	feeRate       int64
	filter        chainFilter
	// evmMaxAmount / btcMaxSats are the rehearsal caps; zero means uncapped (a real drain).
	evmMaxAmount sdkmath.Uint
	btcMaxSats   int64
	priv         *ecdsa.PrivateKey
}

// chainFilter restricts which chains the drain covers. Zero value (both maps nil) allows all.
type chainFilter struct {
	only    map[int64]bool
	exclude map[int64]bool
}

// newChainFilter parses the --only-chains / --exclude-chains flag values. The two are mutually
// exclusive; empty means no restriction.
func newChainFilter(only, exclude string) (chainFilter, error) {
	if only != "" && exclude != "" {
		return chainFilter{}, fmt.Errorf("--%s and --%s are mutually exclusive", FlagOnlyChains, FlagExcludeChains)
	}
	onlyIDs, err := parseChainIDs(only)
	if err != nil {
		return chainFilter{}, fmt.Errorf("invalid --%s: %w", FlagOnlyChains, err)
	}
	excludeIDs, err := parseChainIDs(exclude)
	if err != nil {
		return chainFilter{}, fmt.Errorf("invalid --%s: %w", FlagExcludeChains, err)
	}
	return chainFilter{only: onlyIDs, exclude: excludeIDs}, nil
}

// allow reports whether chainID should be drained under this filter.
func (f chainFilter) allow(chainID int64) bool {
	if len(f.only) > 0 {
		return f.only[chainID]
	}
	return !f.exclude[chainID]
}

// parseChainIDs parses a comma-separated list of int64 chain IDs into a set.
func parseChainIDs(s string) (map[int64]bool, error) {
	if s == "" {
		return nil, nil
	}
	ids := make(map[int64]bool)
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid chain id %q: %w", part, err)
		}
		ids[id] = true
	}
	return ids, nil
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
		g.filter,
		g.evmMaxAmount,
	)
	if err != nil {
		return draintx.Payload{}, err
	}
	btcTxs, err := buildBTCTxs(
		ctx,
		g.cfg,
		g.btcChainID,
		tssAddrRes.Btc,
		g.btcReceiver,
		g.feeRate,
		g.filter,
		g.btcMaxSats,
	)
	if err != nil {
		return draintx.Payload{}, err
	}

	// A final payload is the one clients sign; signing an empty one lets the poller "succeed"
	// while moving nothing. A draft with zero txs is a fine transient state, and a subset
	// (EVM-only or BTC-only) is valid — only the fully-empty final is refused.
	if final && len(evmTxs) == 0 && len(btcTxs) == 0 {
		return draintx.Payload{}, fmt.Errorf(
			"refusing to sign an empty final drain payload: no chains have drainable funds",
		)
	}

	// Emitted per payload, not once at startup: in --serve mode the cron rebuilds every
	// --interval, and a startup-only banner would scroll off behind the per-tick chain logs
	// within a minute. This is the only in-terminal signal that separates a rehearsal from the
	// real drain, so it stays next to the payload it describes.
	warnIfCapped(os.Stderr, g.evmMaxAmount, g.btcMaxSats)

	return drain.BuildPayload(g.triggerHeight, seq, final, g.network, evmTxs, btcTxs, g.priv)
}

// serveDrain runs the draft->freeze->final cron, serving the payload over HTTP.
func serveDrain(ctx context.Context, gen *payloadGenerator, opts drainOptions) error {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	// A failing tick is transient (zetacore RPC) and retried, so surface it and keep serving.
	onError := func(err error) {
		fmt.Fprintf(os.Stderr, "ERROR drain payload publish failed, retrying next tick: %v\n", err)
	}

	if err := drain.RunCron(ctx, opts.interval, gen.generate, server.Publish, isFinalTime, onError); err != nil {
		if ctx.Err() != nil {
			return nil // interrupted before the final was published
		}
		return err
	}

	// The final has been published (at H-K). Closing now would leave clients unable to fetch it
	// during the [H, H+window) firing window, so keep serving until the operator tears down.
	fmt.Fprintf(
		os.Stderr,
		"final drain payload published; serving until interrupted (Ctrl-C) — keep running through the firing window (trigger height %d)\n",
		opts.triggerHigh,
	)
	<-ctx.Done()
	return nil
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
	filter chainFilter,
	maxAmount sdkmath.Uint,
) ([]draintx.EVMTx, error) {
	evmTxs := make([]draintx.EVMTx, 0, len(supportedChains))
	included := make([]int64, 0, len(supportedChains))
	skipped := make([]int64, 0, len(supportedChains))
	for _, chain := range supportedChains {
		if !chain.IsExternal || chain.Vm != pkgchains.Vm_evm {
			continue
		}
		if !filter.allow(chain.ChainId) {
			fmt.Fprintf(os.Stderr, "WARN chain %d skipped: excluded by chain filter\n", chain.ChainId)
			continue
		}
		rpcURL := getRPCForChain(cfg, chain)
		if rpcURL == "" {
			fmt.Fprintf(os.Stderr, "WARN chain %d skipped: RPC not configured\n", chain.ChainId)
			continue
		}

		// A per-chain RPC failure skips that chain loudly rather than aborting the whole payload:
		// one unreachable endpoint must not deny the drain of every other chain.
		balance, err := clients.GetEVMBalance(ctx, rpcURL, tssAddr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR chain %d skipped: get balance failed: %v\n", chain.ChainId, err)
			skipped = append(skipped, chain.ChainId)
			continue
		}
		if balance.Sign() == 0 {
			fmt.Fprintf(os.Stderr, "WARN chain %d skipped: zero balance\n", chain.ChainId)
			continue
		}

		median, err := medianGasPrice(ctx, zetacoreClient, chain.ChainId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR chain %d skipped: get median gas price failed: %v\n", chain.ChainId, err)
			skipped = append(skipped, chain.ChainId)
			continue
		}
		nonce, err := clients.GetEVMNonce(ctx, rpcURL, tssAddr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR chain %d skipped: get nonce failed: %v\n", chain.ChainId, err)
			skipped = append(skipped, chain.ChainId)
			continue
		}
		warnIfNonceNotQuiesced(ctx, os.Stderr, rpcURL, tssAddr, chain.ChainId, nonce)

		tx, err := drain.GenerateEVMTx(drain.EVMInput{
			ChainID:        chain.ChainId,
			To:             receiver,
			Balance:        sdkmath.NewUintFromString(balance.String()),
			MedianGasPrice: median,
			Nonce:          nonce,
			MaxAmount:      maxAmount,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARN chain %d skipped: %v\n", chain.ChainId, err)
			skipped = append(skipped, chain.ChainId)
			continue
		}
		evmTxs = append(evmTxs, tx)
		included = append(included, chain.ChainId)
	}
	fmt.Fprintf(os.Stderr, "EVM drain: %d chains included %v\n", len(included), included)
	fmt.Fprintf(os.Stderr, "EVM drain: skipped chains %v\n", skipped)
	return evmTxs, nil
}

func buildBTCTxs(
	ctx context.Context,
	cfg *config.Config,
	btcChainID int64,
	tssBtcAddr string,
	receiver string,
	feeRate int64,
	filter chainFilter,
	maxSats int64,
) ([]draintx.BTCTx, error) {
	// The TSS BTC address is derived from the pubkey and keyed to btcChainID; the on-chain
	// SupportedChains list can name a different (or no) BTC chain, so we never scan it here —
	// draining the wrong chain would leave real funds behind.
	if tssBtcAddr == "" {
		return nil, nil
	}
	if !filter.allow(btcChainID) {
		fmt.Fprintf(os.Stderr, "WARN chain %d skipped: excluded by chain filter\n", btcChainID)
		return nil, nil
	}

	netParams, err := pkgchains.BitcoinNetParamsFromChainID(btcChainID)
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

	// A per-chain RPC failure skips BTC loudly rather than aborting the whole payload, so an
	// unreachable BTC endpoint doesn't deny the EVM drain.
	btcAdapter, err := clients.NewBitcoinClientAdapter(cfg, pkgchains.Chain{ChainId: btcChainID}, zerolog.Nop())
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR chain %d skipped: create bitcoin client failed: %v\n", btcChainID, err)
		return nil, nil
	}
	unspent, err := btcAdapter.ListUnspentByAddress(ctx, tssAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR chain %d skipped: list unspent failed: %v\n", btcChainID, err)
		return nil, nil
	}

	utxos := make([]drain.UTXO, 0, len(unspent))
	var totalSats int64
	for _, u := range unspent {
		sats, err := btccommon.GetSatoshis(u.Amount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR chain %d skipped: convert utxo amount failed: %v\n", btcChainID, err)
			return nil, nil
		}
		totalSats += sats
		utxos = append(utxos, drain.UTXO{TxID: u.TxID, Vout: u.Vout, AmountSats: sats})
	}
	fmt.Fprintf(os.Stderr, "BTC drain: %d UTXOs (%d sats) for chain %d\n", len(utxos), totalSats, btcChainID)

	txs, err := drain.GenerateBTCTxs(drain.BTCInput{
		ChainID: btcChainID,
		To:      receiverAddr,
		FeeRate: feeRate,
		UTXOs:   utxos,
		MaxSats: maxSats,
	})
	if err != nil {
		// Same policy as every other BTC failure in this function: skip BTC loudly rather than
		// abort, so a BTC-side problem never denies the EVM drain. Only the capped path can get
		// here (the uncapped one skips uneconomical groups instead of erroring), and it would take
		// a receiver address whose output size can't be estimated — impossible for the hardcoded
		// anchors — but the two paths should not diverge on how they fail.
		fmt.Fprintf(os.Stderr, "ERROR chain %d skipped: build sweeps failed: %v\n", btcChainID, err)
		return nil, nil
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
	// With a cap the leftover is mostly deliberate (out-of-cap), not uneconomical dust; saying
	// "uneconomical" there would send the operator hunting a problem that isn't one.
	reason := "as uneconomical"
	if maxSats > 0 {
		reason = fmt.Sprintf("as out-of-cap or uneconomical (rehearsal cap %d sats)", maxSats)
	}
	// A capped run that sweeps nothing is the failure mode that matters: BTC is the leg with no
	// testnet coverage, so an empty BTC section must not pass for "the rehearsal covered BTC".
	// Say what cap would actually work instead of leaving it to be inferred from a missing tx.
	if maxSats > 0 && len(txs) == 0 {
		reportUnviableBTCCap(os.Stderr, utxos, maxSats, feeRate, receiverAddr, btcChainID)
	}
	fmt.Fprintf(
		os.Stderr,
		"BTC drain: swept %d UTXOs (%d sats) in %d txs, skipped %d UTXOs (%d sats) %s\n",
		sweptInputs,
		sweptSats,
		len(txs),
		len(utxos)-sweptInputs,
		totalSats-sweptSats,
		reason,
	)

	return txs, nil
}

// warnIfNonceNotQuiesced flags a chain whose pending nonce has run ahead of its confirmed one.
//
// The payload pins the confirmed nonce, but the poller's executeEVM compares the pinned value
// against the *pending* nonce and treats "pinned already consumed" as a hard stop — that chain
// then drops out for the whole firing window and needs a republish at a higher height to recover.
// So an unmined TSS tx at freeze time silently costs a chain its drain. A rehearsal run creates
// exactly that condition, which makes waiting for the rehearsal txs to confirm a gate on the real
// drain rather than housekeeping.
//
// A failure to read the pending nonce is not fatal: the confirmed nonce is already in hand and
// this is advisory, so it degrades to a warning rather than dropping the chain.
func warnIfNonceNotQuiesced(
	ctx context.Context,
	w io.Writer,
	rpcURL string,
	tssAddr ethcommon.Address,
	chainID int64,
	pinnedNonce uint64,
) {
	pending, err := clients.GetEVMPendingNonce(ctx, rpcURL, tssAddr)
	if err != nil {
		fmt.Fprintf(w, "WARN chain %d: cannot verify nonce quiescence: %v\n", chainID, err)
		return
	}
	reportNonceState(w, chainID, pinnedNonce, pending)
}

// reportNonceState is the decision half of warnIfNonceNotQuiesced, split out so the comparison
// and its wording are testable without an RPC endpoint.
func reportNonceState(w io.Writer, chainID int64, pinnedNonce, pending uint64) {
	switch {
	case pending == pinnedNonce:
		return
	case pending < pinnedNonce:
		// The pending nonce can only lag the confirmed one if this endpoint's view is stale —
		// a load-balanced RPC answering the two calls from different backends will do it. Report
		// it as the RPC problem it is: subtracting here would underflow uint64 and claim ~1.8e19
		// txs in flight, and the "wait for confirmations" advice would be wrong anyway.
		fmt.Fprintf(
			w,
			"WARN chain %d: inconsistent nonce view: pending nonce %d is behind the confirmed nonce %d. "+
				"The RPC is likely load-balanced across backends; confirm the pinned nonce against a "+
				"single trusted endpoint before signing\n",
			chainID,
			pending,
			pinnedNonce,
		)
		return
	}
	fmt.Fprintf(
		w,
		"WARN chain %d: NONCE NOT QUIESCED: pinning confirmed nonce %d but %d tx(s) are still in "+
			"flight (pending nonce %d). If this payload is signed as-is the poller will reject it "+
			"and this chain will not drain. Wait for the in-flight tx(s) to confirm\n",
		chainID,
		pinnedNonce,
		pending-pinnedNonce,
		pending,
	)
}

// reportUnviableBTCCap explains why a capped run produced no BTC sweep, and what cap would.
//
// A UTXO wallet cannot always honour a small cap: the sweep spends whole UTXOs, so a viable
// rehearsal needs a single UTXO that is at once large enough to out-earn its own fee and small
// enough to fit the cap. If the wallet's UTXOs are all far above the cap and the rest is dust,
// no cap in between exists — the operator needs to know that now, not after concluding BTC was
// covered.
func reportUnviableBTCCap(
	w io.Writer,
	utxos []drain.UTXO,
	maxSats, feeRate int64,
	receiver btcutil.Address,
	chainID int64,
) {
	minViable, err := drain.MinViableSweepSats(feeRate, receiver)
	if err != nil {
		fmt.Fprintf(w, "ERROR chain %d: cannot compute the minimum viable sweep: %v\n", chainID, err)
		return
	}

	// the smallest UTXO that would make a viable one-input sweep on its own; the operator can
	// raise the cap to it, and if none exists the wallet's shape rules out a small BTC rehearsal
	var smallestViable int64
	for _, u := range utxos {
		if u.AmountSats >= minViable && (smallestViable == 0 || u.AmountSats < smallestViable) {
			smallestViable = u.AmountSats
		}
	}

	fmt.Fprintf(
		w,
		"ERROR chain %d: REHEARSAL SWEPT NO BTC. --%s %d admits no economical sweep at %d sat/vB "+
			"(a sweep needs at least %d sats to out-earn its own fee)\n",
		chainID,
		FlagBTCMaxSats,
		maxSats,
		feeRate,
		minViable,
	)
	if smallestViable > 0 {
		fmt.Fprintf(
			w,
			"ERROR chain %d: raise --%s to at least %d (the smallest UTXO that can be swept alone) "+
				"or lower --%s, otherwise this run does NOT rehearse BTC\n",
			chainID,
			FlagBTCMaxSats,
			smallestViable,
			FlagFeeRate,
		)
		return
	}
	fmt.Fprintf(
		w,
		"ERROR chain %d: no single UTXO reaches %d sats at this fee rate, so no cap can rehearse BTC "+
			"on these holdings — BTC can only be drained uncapped\n",
		chainID,
		minViable,
	)
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
