package solana_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gagliardetto/solana-go"
	"github.com/zeta-chain/node/pkg/chains"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
)

func Test_MsgWithdrawHash(t *testing.T) {
	t.Run("should pass for archived inbound, receipt and cctx", func(t *testing.T) {
		// #nosec G115 always positive
		chainID := uint64(chains.SolanaLocalnet.ChainId)
		nonce := uint64(0)
		amount := uint64(1336000)
		to := solana.MustPublicKeyFromBase58("37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ")

		wantHash := "a20cddb3f888f4064ced892a477101f45469a8c50f783b966d3fec2455887c05"
		wantHashBytes, err := hex.DecodeString(wantHash)
		require.NoError(t, err)

		// create new withdraw message
		hash := contracts.NewMsgWithdraw(chainID, nonce, amount, to).Hash()
		require.True(t, bytes.Equal(hash[:], wantHashBytes))
	})
}
