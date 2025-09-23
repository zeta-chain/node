package encoder

import (
	"testing"

	"cosmossdk.io/math"

	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/ton"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestEncoder(t *testing.T) {
	t.Run("encode and decode transaction", func(t *testing.T) {
		accountID := sample.GenerateTONAccountID()
		deposit := toncontracts.Deposit{
			Sender:    sample.GenerateTONAccountID(),
			Amount:    math.NewUint(10),
			Recipient: sample.EthAddress(),
		}
		tx := sample.TONDeposit(t, accountID, deposit)

		lt, hash, err := DecodeTx(EncodeTx(tx))
		require.NoError(t, err)
		require.Equal(t, tx.Lt, lt)
		require.Equal(t, tx.Hash().Hex(), hash.Hex())
	})

	t.Run("encode and decode hash", func(t *testing.T) {
		lt1 := uint64(123)
		hash1 := ton.Bits256(sample.Hash())

		lt2, hash2, err := DecodeTx(EncodeHash(lt1, hash1))
		require.NoError(t, err)
		require.Equal(t, lt1, lt2)
		require.Equal(t, hash1.Hex(), hash2.Hex())
	})

	t.Run("decode invalid inputs", func(t *testing.T) {
		inputs := []struct {
			value string
			err   error
		}{
			{"", ErrInvalidString},
			{"::", ErrInvalidString},
			{"1:a:3", ErrInvalidString},
			{":", ErrInvalidLT},
			{"-1:abc", ErrInvalidLT},
			{"1:", ErrInvalidHash},
			{"1:abc", ErrInvalidHash},
		}
		for _, input := range inputs {
			_, _, err := DecodeTx(input.value)
			require.Error(t, err)
			require.ErrorIs(t, err, input.err, "failed input: %q", input.value)
		}
	})
}
