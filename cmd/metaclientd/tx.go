package metaclientd

import (
	"context"
	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/rs/zerolog/log"
)


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