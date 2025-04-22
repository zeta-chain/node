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

		wantHash := "7391cf357fd80e7cb3d2a9758932fdea4988d03c87210a2632f03b467728d199"
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

		wantHash := "d90f9640faecd76509b4e88fa7d18f130918130b8666179f1597f185a828d3a5"
		wantHashBytes := testutil.HexToBytes(t, wantHash)

		// ACT
		// create new execute message
		hash := contracts.NewMsgExecuteSPL(chainID, nonce, amount, 8, mintAccount, to, toAta, sender, []byte("hello"), []*solana.AccountMeta{}).
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
		sender := common.HexToAddress("0x42bd6E2ce4CDb2F58Ed0A0E427F011A0645D5E33")

		wantHash := "A55A5E8E302D5BA9A4C2DCDE54225F45F5E20081873AA7F6A3A361DB20527E31"
		wantHashBytes := testutil.HexToBytes(t, wantHash)

		// ACT
		// create new execute message
		hash := contracts.NewMsgExecute(chainID, nonce, amount, to, sender.Hex(), []byte("hello"), contracts.ExecuteTypeRevert, []*solana.AccountMeta{}).
			Hash()

		// ASSERT
		require.EqualValues(t, hash[:], wantHashBytes)
	})
}
