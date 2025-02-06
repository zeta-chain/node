package inbound

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	cosmosmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/rpc"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	zetaclientObserver "github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	zetaclientConfig "github.com/zeta-chain/node/zetaclient/config"
)

func btcInboundBallotIdentifier(
	ctx context.Context,
	cfg config.Config,
	zetacoreClient rpc.Clients,
	inboundHash string,
	inboundChain chains.Chain,
	zetaChainID int64,
	logger zerolog.Logger) (string, error) {
	params, err := chains.BitcoinNetParamsFromChainID(inboundChain.ChainId)
	if err != nil {
		return "", fmt.Errorf("unable to get bitcoin net params from chain id: %w", err)
	}

	connCfg := zetaclientConfig.BTCConfig{
		RPCUsername: cfg.BtcUser,
		RPCPassword: cfg.BtcPassword,
		RPCHost:     cfg.BtcHost,
		RPCParams:   params.Name,
	}

	rpcClient, err := client.New(connCfg, inboundChain.ChainId, logger)
	if err != nil {
		return "", fmt.Errorf("unable to create rpc client: %w", err)
	}

	err = rpcClient.Ping(ctx)
	if err != nil {
		return "", fmt.Errorf("error ping the bitcoin server: %w", err)
	}
	res, err := zetacoreClient.Observer.GetTssAddress(context.Background(), &types.QueryGetTssAddressRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to get tss address: %w", err)
	}
	tssBtcAddress := res.GetBtc()

	chainParams, err := zetacoreClient.GetChainParamsForChainID(context.Background(), inboundChain.ChainId)
	if err != nil {
		return "", fmt.Errorf("failed to get chain params: %w", err)
	}

	return bitcoinBallotIdentifier(
		ctx,
		rpcClient,
		params,
		tssBtcAddress,
		inboundHash,
		inboundChain.ChainId,
		zetaChainID,
		chainParams.ConfirmationParams.SafeInboundCount,
	)
}

func bitcoinBallotIdentifier(
	ctx context.Context,
	btcClient *client.Client,
	params *chaincfg.Params,
	tss string,
	txHash string,
	senderChainID int64,
	zetacoreChainID int64,
	confirmationCount uint64) (string, error) {
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return "", err
	}
	confirmationMessage := ""
	tx, err := btcClient.GetRawTransactionVerbose(ctx, hash)
	if err != nil {
		return "", err
	}
	if tx.Confirmations < confirmationCount {
		confirmationMessage = fmt.Sprintf("tx might not be confirmed on chain: %d", senderChainID)
	}

	blockHash, err := chainhash.NewHashFromStr(tx.BlockHash)
	if err != nil {
		return "", err
	}

	blockVb, err := btcClient.GetBlockVerbose(ctx, blockHash)
	if err != nil {
		return "", err
	}

	event, err := zetaclientObserver.GetBtcEvent(
		ctx,
		btcClient,
		*tx,
		tss,
		uint64(blockVb.Height), // #nosec G115 always positive
		zerolog.New(zerolog.Nop()),
		params,
		common.CalcDepositorFee,
	)
	if err != nil {
		return "", fmt.Errorf("error getting btc event: %w", err)
	}
	if event == nil {
		return "", fmt.Errorf("no event built for btc sent to TSS")
	}

	return identifierFromBtcEvent(event, senderChainID, zetacoreChainID, confirmationMessage)
}

func identifierFromBtcEvent(event *zetaclientObserver.BTCInboundEvent,
	senderChainID int64,
	zetacoreChainID int64, confirmationMessage string) (string, error) {
	// decode event memo bytes
	err := event.DecodeMemoBytes(senderChainID)
	if err != nil {
		return "", fmt.Errorf("error decoding memo bytes: %w", err)
	}

	// convert the amount to integer (satoshis)
	amountSats, err := common.GetSatoshis(event.Value)
	if err != nil {
		return "", fmt.Errorf("error converting amount to satoshis: %w", err)
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
		return "", fmt.Errorf("failed to create vote message")
	}

	index := msg.Digest()
	if confirmationMessage != "" {
		return fmt.Sprintf("ballot identifier: %s warning: %s", index, confirmationMessage), nil
	}
	return fmt.Sprintf("ballot identifier: %s", msg.Digest()), nil
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
		crosschaintypes.ProtocolContractVersion_V1,
		false, // not relevant for v1
		crosschaintypes.InboundStatus_SUCCESS,
	)
}

func voteFromStdMemo(
	event *zetaclientObserver.BTCInboundEvent,
	amountSats *big.Int,
	senderChainID int64,
	zetacoreChainID int64,
) *crosschaintypes.MsgVoteInbound {
	// zetacore will create a revert outbound that points to the custom revert address.
	revertOptions := crosschaintypes.RevertOptions{
		RevertAddress: event.MemoStd.RevertOptions.RevertAddress,
	}

	// make a legacy message so that zetacore can process it as V1
	msgBytes := append(event.MemoStd.Receiver.Bytes(), event.MemoStd.Payload...)
	message := hex.EncodeToString(msgBytes)

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
		crosschaintypes.ProtocolContractVersion_V1,
		false, // not relevant for v1
		crosschaintypes.InboundStatus_SUCCESS,
		crosschaintypes.WithRevertOptions(revertOptions),
	)
}
