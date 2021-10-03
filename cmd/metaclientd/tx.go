package metaclientd

import (
	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/rs/zerolog/log"
)

// Post Txin to Metachain, with signature of the signer.
// MetaChain takes this as a vote to the PostTxIn.
func (b *MetachainBridge) PostTxIn(fromAddress string, toAddress string,sourceAsset string, sourceAmount uint64, mBurnt uint64, destAsset string, txHash string, blockHeight uint64) (string, error) {
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
