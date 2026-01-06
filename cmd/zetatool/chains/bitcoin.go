package chains

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	cosmosmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/cmd/zetatool/config"
	zetacontext "github.com/zeta-chain/node/cmd/zetatool/context"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/memo"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	zetaclientObserver "github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
)

func BitcoinBallotIdentifier(
	ctx *zetacontext.Context,
	btcClient *client.Client,
	params *chaincfg.Params,
	tss string,
	feeRateMultiplier float64,
	txHash string,
	senderChainID int64,
	zetacoreChainID int64,
	confirmationCount uint64,
) (cctxIdentifier string, isConfirmed bool, err error) {
	var (
		goCtx = ctx.GetContext()
	)

	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return
	}
	tx, err := btcClient.GetRawTransactionVerbose(goCtx, hash)
	if err != nil {
		return
	}

	if tx.Confirmations >= confirmationCount {
		isConfirmed = true
	}

	blockHash, err := chainhash.NewHashFromStr(tx.BlockHash)
	if err != nil {
		return
	}

	blockVb, err := btcClient.GetBlockVerbose(goCtx, blockHash)
	if err != nil {
		return
	}

	event, err := zetaclientObserver.GetBtcEventWithWitness(
		goCtx,
		btcClient,
		*tx,
		tss,
		uint64(blockVb.Height), // #nosec G115 always positive
		feeRateMultiplier,
		zerolog.New(zerolog.Nop()),
		params,
		common.CalcDepositorFee,
	)
	if err != nil {
		return
	}
	if event == nil {
		err = fmt.Errorf("no event built for btc sent to TSS")
		return
	}

	cctxIdentifier, err = identifierFromBtcEvent(event, senderChainID, zetacoreChainID)
	return
}

func identifierFromBtcEvent(event *zetaclientObserver.BTCInboundEvent,
	senderChainID int64,
	zetacoreChainID int64) (cctxIdentifier string, err error) {
	err = event.DecodeMemoBytes(senderChainID)
	if err != nil {
		return
	}

	amountSats, err := common.GetSatoshis(event.Value)
	if err != nil {
		return
	}
	amountInt := big.NewInt(amountSats)

	var msg *crosschaintypes.MsgVoteInbound
	switch event.MemoStd {
	case nil:
		{
			msg = voteFromLegacyMemo(event, amountInt, senderChainID, zetacoreChainID)
		}
	default:
		{
			msg = voteFromStdMemo(event, amountInt, senderChainID, zetacoreChainID)
		}
	}
	if msg == nil {
		return
	}

	cctxIdentifier = msg.Digest()
	return
}

// NewInboundVoteFromLegacyMemo creates a MsgVoteInbound message for inbound that uses legacy memo
func voteFromLegacyMemo(
	event *zetaclientObserver.BTCInboundEvent,
	amountSats *big.Int,
	senderChainID int64,
	zetacoreChainID int64,
) *crosschaintypes.MsgVoteInbound {
	message := hex.EncodeToString(event.MemoBytes)

	return crosschaintypes.NewMsgVoteInbound(
		"",
		event.FromAddress,
		senderChainID,
		event.FromAddress,
		event.ToAddress,
		zetacoreChainID,
		cosmosmath.NewUintFromBigInt(amountSats),
		message,
		event.TxHash,
		event.BlockNumber,
		0,
		coin.CoinType_Gas,
		"",
		0,
		crosschaintypes.ProtocolContractVersion_V2,
		false, // not relevant for v1
		crosschaintypes.InboundStatus_SUCCESS,
		crosschaintypes.ConfirmationMode_SAFE,
		crosschaintypes.WithCrossChainCall(len(event.MemoBytes) > 0),
	)
}

func voteFromStdMemo(
	event *zetaclientObserver.BTCInboundEvent,
	amountSats *big.Int,
	senderChainID int64,
	zetacoreChainID int64,
) *crosschaintypes.MsgVoteInbound {
	revertOptions := crosschaintypes.RevertOptions{
		RevertAddress: event.MemoStd.RevertOptions.RevertAddress,
	}

	isCrosschainCall := event.MemoStd.OpCode == memo.OpCodeCall || event.MemoStd.OpCode == memo.OpCodeDepositAndCall

	return crosschaintypes.NewMsgVoteInbound(
		"",
		event.FromAddress,
		senderChainID,
		event.FromAddress,
		event.ToAddress,
		zetacoreChainID,
		cosmosmath.NewUintFromBigInt(amountSats),
		hex.EncodeToString(event.MemoStd.Payload),
		event.TxHash,
		event.BlockNumber,
		0,
		coin.CoinType_Gas,
		"",
		0,
		crosschaintypes.ProtocolContractVersion_V2,
		false, // not relevant for v1
		event.Status,
		crosschaintypes.ConfirmationMode_SAFE,
		crosschaintypes.WithRevertOptions(revertOptions),
		crosschaintypes.WithCrossChainCall(isCrosschainCall),
	)
}

const (
	mempoolAddressAPIMainnet  = "https://mempool.space/api/address/%s"
	mempoolAddressAPITestnet3 = "https://mempool.space/testnet/api/address/%s"
	mempoolAddressAPISignet   = "https://mempool.space/signet/api/address/%s"
	mempoolAddressAPITestnet4 = "https://mempool.space/testnet4/api/address/%s"
	satoshisPerBitcoin        = 100_000_000
	httpClientTimeout         = 30 * time.Second
)

// BTCAddressStats represents the response from mempool.space address API
type BTCAddressStats struct {
	Address    string `json:"address"`
	ChainStats struct {
		FundedTxoCount int   `json:"funded_txo_count"`
		FundedTxoSum   int64 `json:"funded_txo_sum"`
		SpentTxoCount  int   `json:"spent_txo_count"`
		SpentTxoSum    int64 `json:"spent_txo_sum"`
		TxCount        int   `json:"tx_count"`
	} `json:"chain_stats"`
	MempoolStats struct {
		FundedTxoCount int   `json:"funded_txo_count"`
		FundedTxoSum   int64 `json:"funded_txo_sum"`
		SpentTxoCount  int   `json:"spent_txo_count"`
		SpentTxoSum    int64 `json:"spent_txo_sum"`
		TxCount        int   `json:"tx_count"`
	} `json:"mempool_stats"`
}

// GetBTCBalance fetches the BTC balance for a given address using mempool.space API
// Returns the balance in BTC (not satoshis)
func GetBTCBalance(ctx context.Context, address string, chainID int64) (float64, error) {
	apiURL := getMempoolAddressAPIURL(chainID, address)
	if apiURL == "" {
		return 0, fmt.Errorf("unsupported Bitcoin chain ID: %d", chainID)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	btcClient := &http.Client{Timeout: httpClientTimeout}
	resp, err := btcClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch address stats: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("mempool.space API returned status %d", resp.StatusCode)
	}

	var stats BTCAddressStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return 0, fmt.Errorf("failed to decode address stats: %w", err)
	}

	balanceSatoshis := stats.ChainStats.FundedTxoSum - stats.ChainStats.SpentTxoSum

	return float64(balanceSatoshis) / satoshisPerBitcoin, nil
}

// getMempoolAddressAPIURL returns the mempool.space address API URL for the given chain ID
func getMempoolAddressAPIURL(chainID int64, address string) string {
	switch chainID {
	case 8332: // Bitcoin mainnet
		return fmt.Sprintf(mempoolAddressAPIMainnet, address)
	case 18332: // Bitcoin testnet3
		return fmt.Sprintf(mempoolAddressAPITestnet3, address)
	case 18333: // Bitcoin signet
		return fmt.Sprintf(mempoolAddressAPISignet, address)
	case 18334: // Bitcoin testnet4
		return fmt.Sprintf(mempoolAddressAPITestnet4, address)
	default:
		return ""
	}
}

// GetBTCChainID returns the Bitcoin chain ID for the given network
func GetBTCChainID(network string) int64 {
	switch network {
	case config.NetworkMainnet:
		return 8332
	case config.NetworkTestnet:
		return 18332
	case config.NetworkLocalnet:
		return 18444
	default:
		panic("invalid network")
	}
}
