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

		wantHash := "aa609ef9480303e8d743f6e36fe1bea0cc56b8d27dcbd8220846125c1181b681"
		wantHashBytes, err := hex.DecodeString(wantHash)
		require.NoError(t, err)

		// create new withdraw message
		hash := contracts.NewMsgWithdraw(chainID, nonce, amount, to).Hash()
		require.True(t, bytes.Equal(hash[:], wantHashBytes))
	})
}
