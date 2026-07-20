package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
)

// NewDrainPayloadCMD creates the command that builds and signs an emergency drain payload.
func NewDrainPayloadCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drain-payload <chain>",
		Short: "Build and sign an emergency TSS drain payload (EVM + BTC)",
		Long: `Build the fully-resolved, signed drain payload that moves all native TSS funds to
the hardcoded safe wallet. The operator supplies only the trigger height; balances,
median gas prices, nonces and UTXOs are derived from the configured RPCs.

The chain argument selects the network (mainnet/testnet/localnet) the same way as
tss-balances. The signed JSON payload is printed to stdout.`,
		Args: cobra.ExactArgs(1),
		RunE: runDrainPayload,
	}

	cmd.Flags().Int64(FlagTriggerHeight, 0, "zeta block height at which clients fire (required)")
	cmd.Flags().Bool(FlagFinal, false, "mark the payload as final (clients only sign final payloads)")
	cmd.Flags().String(FlagSigningKey, "", "hex-encoded secp256k1 operator private key (required)")
	cmd.Flags().Uint64(FlagSeq, 0, "monotonic payload version")
	cmd.Flags().Int64(FlagFeeRate, conservativeFeeRate, "BTC fee rate in sat/vB")

	return cmd
}

func runDrainPayload(cmd *cobra.Command, args []string) error {
	chain, err := zetatoolcommon.ResolveChain(args[0])
	if err != nil {
		return fmt.Errorf("failed to resolve chain %q: %w", args[0], err)
	}
	network := zetatoolcommon.NetworkTypeFromChain(chain)

	triggerHeight, err := cmd.Flags().GetInt64(FlagTriggerHeight)
	if err != nil {
		return err
	}
	if triggerHeight <= 0 {
		return fmt.Errorf("--%s is required and must be positive", FlagTriggerHeight)
	}

	final, err := cmd.Flags().GetBool(FlagFinal)
	if err != nil {
		return err
	}
	seq, err := cmd.Flags().GetUint64(FlagSeq)
	if err != nil {
		return err
	}
	feeRate, err := cmd.Flags().GetInt64(FlagFeeRate)
	if err != nil {
		return err
	}

	signingKeyHex, err := cmd.Flags().GetString(FlagSigningKey)
	if err != nil {
		return err
	}
	priv, err := ethcrypto.HexToECDSA(strings.TrimPrefix(signingKeyHex, "0x"))
	if err != nil {
		return fmt.Errorf("invalid --%s: %w", FlagSigningKey, err)
	}

	configFile, err := cmd.Flags().GetString(config.FlagConfig)
	if err != nil {
		return err
	}
	cfg, err := config.GetConfigByNetwork(network, configFile)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}
	if cfg.ZetaChainRPC == "" {
		return fmt.Errorf("ZetaChainRPC is not configured for network %s", network)
	}

	receivers, err := drain.ReceiverForNetwork(drainNetwork(network))
	if err != nil {
		return err
	}

	ctx := context.Background()
	zetacoreClient, err := rpc.NewCometBFTClients(cfg.ZetaChainRPC)
	if err != nil {
		return fmt.Errorf("failed to create zetacore client: %w", err)
	}

	tss, err := currentTSS(ctx, zetacoreClient)
	if err != nil {
		return err
	}

	btcChainID, err := clients.GetBTCChainID(network)
	if err != nil {
		return fmt.Errorf("failed to get BTC chain ID: %w", err)
	}
	tssAddrRes, err := zetacoreClient.Observer.GetTssAddressByFinalizedHeight(
		ctx,
		&observertypes.QueryGetTssAddressByFinalizedHeightRequest{
			FinalizedZetaHeight: tss.FinalizedZetaHeight,
			BitcoinChainId:      btcChainID,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to get TSS addresses: %w", err)
	}

	supportedChains, err := zetacoreClient.Observer.SupportedChains(ctx, &observertypes.QuerySupportedChains{})
	if err != nil {
		return fmt.Errorf("failed to get supported chains: %w", err)
	}

	evmTxs, err := buildEVMTxs(ctx, cfg, zetacoreClient, supportedChains.Chains, ethcommon.HexToAddress(tssAddrRes.Eth), receivers.EVM)
	if err != nil {
		return err
	}

	btcTxs, err := buildBTCTxs(ctx, cfg, supportedChains.Chains, tssAddrRes.Btc, receivers.BTC, feeRate)
	if err != nil {
		return err
	}

	payload, err := drain.BuildPayload(triggerHeight, seq, final, evmTxs, btcTxs, priv)
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

// currentTSS returns the TSS with the highest finalized height.
func currentTSS(ctx context.Context, c rpc.Clients) (observertypes.TSS, error) {
	res, err := c.Observer.TssHistory(ctx, &observertypes.QueryTssHistoryRequest{})
	if err != nil {
		return observertypes.TSS{}, fmt.Errorf("failed to fetch TSS history: %w", err)
	}
	if len(res.TssList) == 0 {
		return observertypes.TSS{}, fmt.Errorf("no TSS entries found")
	}
	current := res.TssList[0]
	for _, t := range res.TssList {
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
	var evmTxs []draintx.EVMTx
	for _, chain := range supportedChains {
		if !chain.IsExternal || chain.Vm != pkgchains.Vm_evm {
			continue
		}
		rpcURL := getRPCForChain(cfg, chain)
		if rpcURL == "" {
			continue
		}

		balance, err := clients.GetEVMBalance(ctx, rpcURL, tssAddr)
		if err != nil {
			return nil, fmt.Errorf("get balance for chain %d: %w", chain.ChainId, err)
		}
		if balance.Sign() == 0 {
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
			// insufficient funds to cover the fee — skip this chain rather than aborting.
			fmt.Printf("skipping chain %d: %v\n", chain.ChainId, err)
			continue
		}
		evmTxs = append(evmTxs, tx)
	}
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
	for _, u := range unspent {
		sats, err := btccommon.GetSatoshis(u.Amount)
		if err != nil {
			return nil, fmt.Errorf("convert utxo amount: %w", err)
		}
		utxos = append(utxos, drain.UTXO{TxID: u.TxID, Vout: u.Vout, AmountSats: sats})
	}

	return drain.GenerateBTCTxs(drain.BTCInput{
		ChainID: btcChain.ChainId,
		To:      receiverAddr,
		FeeRate: feeRate,
		UTXOs:   utxos,
	})
}

// medianGasPrice queries zetacore for the median gas price of a chain.
func medianGasPrice(ctx context.Context, c rpc.Clients, chainID int64) (sdkmath.Uint, error) {
	res, err := c.Crosschain.GasPrice(ctx, &crosschaintypes.QueryGetGasPriceRequest{
		Index: strconv.FormatInt(chainID, 10),
	})
	if err != nil {
		return sdkmath.ZeroUint(), fmt.Errorf("get gas price for chain %d: %w", chainID, err)
	}
	gp := res.GasPrice
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
