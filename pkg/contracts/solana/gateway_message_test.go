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

func Test_MsgWhitelistHash(t *testing.T) {
	t.Run("should pass for archived inbound, receipt and cctx", func(t *testing.T) {
		// #nosec G115 always positive
		chainID := uint64(chains.SolanaLocalnet.ChainId)
		nonce := uint64(0)
		whitelistCandidate := solana.MustPublicKeyFromBase58("37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ")
		whitelistEntry := solana.MustPublicKeyFromBase58("2kJndCL9NBR36ySiQ4bmArs4YgWQu67LmCDfLzk5Gb7s")

		wantHash := "cde8fa3ab24b50320db1c47f30492e789177d28e76208176f0a52b8ed54ce2dd"
		wantHashBytes, err := hex.DecodeString(wantHash)
		require.NoError(t, err)

		// create new withdraw message
		hash := contracts.NewMsgWhitelist(whitelistCandidate, whitelistEntry, chainID, nonce).Hash()
		require.True(t, bytes.Equal(hash[:], wantHashBytes))
	})
}

func Test_MsgWithdrawSPLHash(t *testing.T) {
	t.Run("should pass for archived inbound, receipt and cctx", func(t *testing.T) {
		// #nosec G115 always positive
		chainID := uint64(chains.SolanaLocalnet.ChainId)
		nonce := uint64(0)
		amount := uint64(1336000)
		mintAccount := solana.MustPublicKeyFromBase58("AS48jKNQsDGkEdDvfwu1QpqjtqbCadrAq9nGXjFmdX3Z")
		to := solana.MustPublicKeyFromBase58("37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ")
		toAta, _, err := solana.FindAssociatedTokenAddress(to, mintAccount)
		require.NoError(t, err)

		wantHash := "87fa5c0ed757c6e1ea9d8976537eaf7868bc1f1bbf55ab198a01645d664fe0ae"
		wantHashBytes, err := hex.DecodeString(wantHash)
		require.NoError(t, err)

		// create new withdraw message
		hash := contracts.NewMsgWithdrawSPL(chainID, nonce, amount, 8, mintAccount, to, toAta).Hash()
		require.True(t, bytes.Equal(hash[:], wantHashBytes))
	})
}
