package inbound_ballot

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	cosmosmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/rpc"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	zetaclientObserver "github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
)

func btcInboundBallotIdentifier(
	cfg config.Config,
	zetacoreClient rpc.Clients,
	inboundHash string,
	inboundChain chains.Chain,
	zetaChainID int64) (string, error) {
	params, err := chains.BitcoinNetParamsFromChainID(inboundChain.ChainId)
	if err != nil {
		return "", fmt.Errorf("unable to get bitcoin net params from chain id: %s", err)
	}

	connCfg := &rpcclient.ConnConfig{
		Host:         cfg.BtcHost,
		User:         cfg.BtcUser,
		Pass:         cfg.BtcPassword,
		HTTPPostMode: true,
		DisableTLS:   true,
		Params:       params.Name,
	}
	rpcClient, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return "", fmt.Errorf("error creating rpc client: %s", err)
	}

	err = rpcClient.Ping()
	if err != nil {
		return "", fmt.Errorf("error ping the bitcoin server: %s", err)
	}
	res, err := zetacoreClient.Observer.GetTssAddress(context.Background(), &types.QueryGetTssAddressRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to get tss address %s", err.Error())
	}
	tssBtcAddress := res.GetBtc()

	return bitcoinBallotIdentifier(
		rpcClient,
		params,
		tssBtcAddress,
		inboundHash,
		inboundChain.ChainId,
		zetaChainID,
	)
}

func bitcoinBallotIdentifier(
	btcClient *rpcclient.Client,
	params *chaincfg.Params,
	tss string,
	txHash string,
	senderChainID int64,
	zetacoreChainID int64) (string, error) {
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return "", err
	}

	tx, err := btcClient.GetRawTransactionVerbose(hash)
	if err != nil {
		return "", err
	}

	blockHash, err := chainhash.NewHashFromStr(tx.BlockHash)
	if err != nil {
		return "", err
	}

	blockVb, err := btcClient.GetBlockVerboseTx(blockHash)
	if err != nil {
		return "", err
	}

	if len(blockVb.Tx) <= 1 {
		return "", fmt.Errorf("block %d has no transactions", blockVb.Height)
	}
	// #nosec G115 always positive

	event, err := zetaclientObserver.GetBtcEvent(
		btcClient,
		*tx,
		tss,
		uint64(blockVb.Height),
		zerolog.New(zerolog.Nop()),
		params,
		common.CalcDepositorFee,
	)
	if err != nil {
		return "", fmt.Errorf("error getting btc event: %s", err)
	}

	if event == nil {
		return "", fmt.Errorf("no event built for btc sent to TSS")
	}

	return identifierFromBtcEvent(event, senderChainID, zetacoreChainID)
}

func identifierFromBtcEvent(event *zetaclientObserver.BTCInboundEvent,
	senderChainID int64,
	zetacoreChainID int64) (string, error) {
	// decode event memo bytes
	err := event.DecodeMemoBytes(senderChainID)
	if err != nil {
		return "", fmt.Errorf("error decoding memo bytes: %s", err)
	}

	// convert the amount to integer (satoshis)
	amountSats, err := common.GetSatoshis(event.Value)
	if err != nil {
		return "", fmt.Errorf("error converting amount to satoshis: %s", err)
	}
	amountInt := big.NewInt(amountSats)

	switch event.MemoStd {
	case nil:
		{
			msg := voteFromLegacyMemo(event, amountInt, senderChainID, zetacoreChainID)
			return msg.Digest(), nil
		}
	default:
		{
			msg := voteFromStdMemo(event, amountInt, senderChainID, zetacoreChainID)
			return msg.Digest(), nil
		}
	}
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
		crosschaintypes.WithRevertOptions(revertOptions),
	)
}
