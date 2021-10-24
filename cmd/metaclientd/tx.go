package metaclientd

import (
	"context"
	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/rs/zerolog/log"
)

// Post Txin to Metachain, with signature of the signer.
// MetaChain takes this as a vote to the PostTxIn.
func (b *MetachainBridge) PostTxIn(fromAddress string, toAddress string, sourceAsset string, sourceAmount uint64, mBurnt uint64, destAsset string, txHash string, blockHeight uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgCreateTxinVoter(
		signerAddress,
		txHash, sourceAsset, sourceAmount, mBurnt, destAsset,
		fromAddress, toAddress, blockHeight,
	)

	metaTxHash, err := b.Broadcast(msg)
	if err != nil {
		log.Err(err).Msg("PostTxIn broadcast fail")
		return "", err
	}
	return metaTxHash, nil
}

// Post Txin to Metachain, with signature of the signer.
// MetaChain takes this as a vote to the PostTxIn.
func (b *MetachainBridge) PostTxoutConfirmation(txoutId uint64, txHash string, mMint uint64, destinationAsset string, destinationAmount uint64, toAddress string, blockHeight uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgTxoutConfirmationVoter(signerAddress, txoutId, txHash, mMint, destinationAsset, destinationAmount, toAddress, blockHeight)
	metaTxHash, err := b.Broadcast(msg)
	if err != nil {
		log.Err(err).Msg("PostTxoutConfirmation broadcast fail")
		return "", err
	}
	return metaTxHash, nil
}

// Get all current Txout from MetaCore
func (b *MetachainBridge) GetAllTxout() ([]*types.Txout, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.TxoutAll(context.Background(), &types.QueryAllTxoutRequest{})
	if err != nil {
		log.Error().Err(err).Msg("query TxoutAll error")
		return nil, err
	}
	return resp.Txout, nil
}

func (b *MetachainBridge) PostSend(sender string, senderChain string, receiver string, receiverChain string, mBurnt string, mMint string, message string, inTxHash string, inBlockHeight uint64) (string, error) {
	signerAddress := b.keys.GetSignerInfo().GetAddress().String()
	msg := types.NewMsgSendVoter(signerAddress, sender, senderChain, receiver, receiverChain, mBurnt, mMint, message, inTxHash, inBlockHeight)

	metaTxHash, err := b.Broadcast(msg)
	if err != nil {
		log.Err(err).Msg("PostSend broadcast fail")
		return "", err
	}
	return metaTxHash, nil
}

func (b *MetachainBridge) GetAllSend() ([]*types.Send, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.SendAll(context.Background(), &types.QueryAllSendRequest{})
	if err != nil {
		log.Error().Err(err).Msg("query SendAll error")
		return nil, err
	}
	return resp.Send, nil
}