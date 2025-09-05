package solana_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/gagliardetto/solana-go"
	"github.com/zeta-chain/node/pkg/chains"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/testutil"
)

func Test_MsgWithdrawHash(t *testing.T) {
	t.Run("should calculate expected hash", func(t *testing.T) {
		// ARRANGE
		// #nosec G115 always positive
		chainID := uint64(chains.SolanaLocalnet.ChainId)
		nonce := uint64(0)
		amount := uint64(1336000)
		to := solana.MustPublicKeyFromBase58("37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ")

		wantHash := "ab576990e7a36adcaeed559d434680308e18cadc60e9c6a50d2891282274dbb0"
		wantHashBytes := testutil.HexToBytes(t, wantHash)

		// ACT
		// create new withdraw message
		hash := contracts.NewMsgWithdraw(chainID, nonce, amount, to).Hash()

		// ASSERT
		require.EqualValues(t, hash[:], wantHashBytes)
	})
}

func Test_MsgWhitelistHash(t *testing.T) {
	t.Run("should calculate expected hash", func(t *testing.T) {
		// ARRANGE
		// #nosec G115 always positive
		chainID := uint64(chains.SolanaLocalnet.ChainId)
		nonce := uint64(0)
		whitelistCandidate := solana.MustPublicKeyFromBase58("37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ")
		whitelistEntry := solana.MustPublicKeyFromBase58("2kJndCL9NBR36ySiQ4bmArs4YgWQu67LmCDfLzk5Gb7s")

		wantHash := "04de341bab0307cdcf9f3ff23c38461e17c33f98338163eeb7a5f50e6b12d052"
		wantHashBytes := testutil.HexToBytes(t, wantHash)

		// ACT
		// create new withdraw message
		hash := contracts.NewMsgWhitelist(whitelistCandidate, whitelistEntry, chainID, nonce).Hash()

		// ASSERT
		require.EqualValues(t, hash[:], wantHashBytes)
	})
}

func Test_MsgWithdrawSPLHash(t *testing.T) {
	t.Run("should calculate expected hash", func(t *testing.T) {
		// ARRANGE
		// #nosec G115 always positive
		chainID := uint64(chains.SolanaLocalnet.ChainId)
		nonce := uint64(0)
		amount := uint64(1336000)
		mintAccount := solana.MustPublicKeyFromBase58("AS48jKNQsDGkEdDvfwu1QpqjtqbCadrAq9nGXjFmdX3Z")
		to := solana.MustPublicKeyFromBase58("37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ")
		toAta, _, err := solana.FindAssociatedTokenAddress(to, mintAccount)
		require.NoError(t, err)

		wantHash := "628d9011b2784ed2376fcc1ba57e49c0e9dfc860d10d379c8d228c8972c9cd65"
		wantHashBytes := testutil.HexToBytes(t, wantHash)

		// ACT
		// create new withdraw message
		hash := contracts.NewMsgWithdrawSPL(chainID, nonce, amount, 8, mintAccount, to, toAta).Hash()

		// ASSERT
		require.EqualValues(t, hash[:], wantHashBytes)
	})
}

func Test_MsgIncrementNonceHash(t *testing.T) {
	t.Run("should calculate expected hash", func(t *testing.T) {
		// ARRANGE
		// #nosec G115 always positive
		chainID := uint64(chains.SolanaLocalnet.ChainId)
		nonce := uint64(0)
		amount := uint64(1336000)

		wantHash := "7acd44905f89cc36b4542b103a25f2110d25ca469dee854fbf342f1f3ee8ba90"
		wantHashBytes := testutil.HexToBytes(t, wantHash)

		// ACT
		// create new increment nonce message
		hash := contracts.NewMsgIncrementNonce(chainID, nonce, amount).Hash()

		// ASSERT
		require.EqualValues(t, hash[:], wantHashBytes)
	})
}

func Test_MsgExecuteHash(t *testing.T) {
	t.Run("should calculate expected hash", func(t *testing.T) {
		// ARRANGE
		// #nosec G115 always positive
		chainID := uint64(chains.SolanaLocalnet.ChainId)
		nonce := uint64(0)
		amount := uint64(1336000)
		to := solana.MustPublicKeyFromBase58("37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ")
		sender := common.HexToAddress("0x42bd6E2ce4CDb2F58Ed0A0E427F011A0645D5E33")

		wantHash := "ff0737262c010614dead18a4ed152ca38bd92aae694f08ba611cea307f3d92d9"
		wantHashBytes := testutil.HexToBytes(t, wantHash)

		// ACT
		// create new execute message
		hash := contracts.NewMsgExecute(chainID, nonce, amount, to, sender.Hex(), []byte("hello"), contracts.ExecuteTypeCall, []*solana.AccountMeta{}).
			Hash()

		// ASSERT
		require.EqualValues(t, hash[:], wantHashBytes)
	})
}

func Test_MsgExecuteSPLHash(t *testing.T) {
	t.Run("should calculate expected hash", func(t *testing.T) {
		// ARRANGE
		// #nosec G115 always positive
		chainID := uint64(chains.SolanaLocalnet.ChainId)
		nonce := uint64(0)
		amount := uint64(1336000)
		mintAccount := solana.MustPublicKeyFromBase58("AS48jKNQsDGkEdDvfwu1QpqjtqbCadrAq9nGXjFmdX3Z")
		to := solana.MustPublicKeyFromBase58("37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ")
		toAta, _, err := solana.FindAssociatedTokenAddress(to, mintAccount)
		require.NoError(t, err)
		sender := common.HexToAddress("0x42bd6E2ce4CDb2F58Ed0A0E427F011A0645D5E33")

		wantHash := "8ed9d12294399703611bf6b4d131aa0cdd9b1c2ca7e2586238c47929766ed0b9"
		wantHashBytes := testutil.HexToBytes(t, wantHash)

		// ACT
		// create new execute message
		hash := contracts.NewMsgExecuteSPL(chainID, nonce, amount, 8, mintAccount, to, toAta, sender.Hex(), []byte("hello"), contracts.ExecuteTypeCall, []*solana.AccountMeta{}).
			Hash()

		// ASSERT
		require.EqualValues(t, hash[:], wantHashBytes)
	})
}

func Test_MsgExecuteRevertHash(t *testing.T) {
	t.Run("should calculate expected hash", func(t *testing.T) {
		// ARRANGE
		// #nosec G115 always positive
		chainID := uint64(chains.SolanaLocalnet.ChainId)
		nonce := uint64(0)
		amount := uint64(1336000)
		to := solana.MustPublicKeyFromBase58("37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ")
		sender := solana.MustPublicKeyFromBase58("CVoPuE3EMu6QptGHLx7mDGb2ZgASJRQ5BcTvmhZNJd8A")

		wantHash := "8a84b88736c9b677b1d859b212bd07f9808f8ca07682b3585eb1b6c5f2f6f4dd"
		wantHashBytes := testutil.HexToBytes(t, wantHash)

		// ACT
		// create new execute message
		hash := contracts.NewMsgExecute(chainID, nonce, amount, to, sender.String(), []byte("hello"), contracts.ExecuteTypeRevert, []*solana.AccountMeta{}).
			Hash()

		// ASSERT
		require.EqualValues(t, hash[:], wantHashBytes)
	})
}
func Test_MsgExecuteSPLRevertHash(t *testing.T) {
	t.Run("should calculate expected hash", func(t *testing.T) {
		// ARRANGE
		// #nosec G115 always positive
		chainID := uint64(chains.SolanaLocalnet.ChainId)
		nonce := uint64(0)
		amount := uint64(1336000)
		mintAccount := solana.MustPublicKeyFromBase58("AS48jKNQsDGkEdDvfwu1QpqjtqbCadrAq9nGXjFmdX3Z")
		to := solana.MustPublicKeyFromBase58("37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ")
		toAta, _, err := solana.FindAssociatedTokenAddress(to, mintAccount)
		require.NoError(t, err)
		sender := solana.MustPublicKeyFromBase58("CVoPuE3EMu6QptGHLx7mDGb2ZgASJRQ5BcTvmhZNJd8A")

		wantHash := "147f471f794e0227653f69c40d9ac0ab278ac8549de41d9229e9b6ca14735c57"
		wantHashBytes := testutil.HexToBytes(t, wantHash)

		// ACT
		// create new execute message
		hash := contracts.NewMsgExecuteSPL(chainID, nonce, amount, 8, mintAccount, to, toAta, sender.String(), []byte("hello"), contracts.ExecuteTypeRevert, []*solana.AccountMeta{}).
			Hash()

		// ASSERT
		require.EqualValues(t, hash[:], wantHashBytes)
	})
}
