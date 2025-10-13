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

// remainingAccounts to be used in execute tests
var remainingAccounts = []*solana.AccountMeta{
	solana.NewAccountMeta(solana.MustPublicKeyFromBase58("4C2kkMnqXMfPJ8PPK5v6TCg42k5z16f2kzqVocXoSDcq"), false, false),
	solana.NewAccountMeta(solana.MustPublicKeyFromBase58("3tqcVCVz5Jztwnku1H9zpvjaWshSpHakMuX4xUJQuuuA"), false, false),
	solana.NewAccountMeta(solana.MustPublicKeyFromBase58("7duGsuv6nB3yr15EuWuHEDD7rWovpAnjuveXJ5ySZuFV"), false, false),
}

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

		wantHash := "c971ad67a1727edd0c918bf48c0859c0aa43c8bff32d533b7517d4da25aa9c01"
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

		wantHash := "f7cab9856e40b3d6381a69aa4854c90c4c7affe00d86dfba747599359f07126b"
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

		wantHash := "a930d636213955b4c604199b735433452f5b9ba2a9d2869caca8b1c1af9218f2"
		wantHashBytes := testutil.HexToBytes(t, wantHash)

		// ACT
		// create new execute message
		hash := contracts.NewMsgExecute(chainID, nonce, amount, to, sender.Hex(), []byte("hello"), contracts.ExecuteTypeCall, remainingAccounts, nil, nil).
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

		wantHash := "a245ba0c28fd5ca344f4747d3249e94e2fee0a883cd3acb5ac85d041556df9c0"
		wantHashBytes := testutil.HexToBytes(t, wantHash)

		// ACT
		// create new execute message
		hash := contracts.NewMsgExecuteSPL(chainID, nonce, amount, 8, mintAccount, to, toAta, sender.Hex(), []byte("hello"), contracts.ExecuteTypeCall, remainingAccounts, nil, nil).
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

		wantHash := "04ac292b75cda6b32c7b4a5761877d2e7f23cf22e4b97945eb5b959b7e64c0a5"
		wantHashBytes := testutil.HexToBytes(t, wantHash)

		// ACT
		// create new execute message
		hash := contracts.NewMsgExecute(chainID, nonce, amount, to, sender.String(), []byte("hello"), contracts.ExecuteTypeRevert, remainingAccounts, nil, nil).
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

		wantHash := "bdd0f1218cea32d81bf5626f0031ebccef477e843a648a4c3eb4612a234604ee"
		wantHashBytes := testutil.HexToBytes(t, wantHash)

		// ACT
		// create new execute message
		hash := contracts.NewMsgExecuteSPL(chainID, nonce, amount, 8, mintAccount, to, toAta, sender.String(), []byte("hello"), contracts.ExecuteTypeRevert, remainingAccounts, nil, nil).
			Hash()

		// ASSERT
		require.EqualValues(t, hash[:], wantHashBytes)
	})
}
