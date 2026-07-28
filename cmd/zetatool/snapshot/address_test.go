package snapshot

import (
	"encoding/hex"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/stretchr/testify/require"
)

func TestCanonical(t *testing.T) {
	// ARRANGE
	raw := make([]byte, addressLen)
	for i := range raw {
		raw[i] = byte(i + 1)
	}
	want := "0x" + hex.EncodeToString(raw)

	acc, err := bech32.ConvertAndEncode("zeta", raw)
	require.NoError(t, err)
	val, err := bech32.ConvertAndEncode("zetavaloper", raw)
	require.NoError(t, err)
	hexAddr := want

	// ACT
	gotAcc, errAcc := Canonical(acc)
	gotVal, errVal := Canonical(val)
	gotHex, errHex := Canonical(hexAddr)

	// ASSERT
	require.NoError(t, errAcc)
	require.NoError(t, errVal)
	require.NoError(t, errHex)
	require.Equal(t, want, gotAcc)
	require.Equal(t, want, gotVal, "validator operator and account address must collapse to the same key")
	require.Equal(t, want, gotHex)
}

func TestCanonicalErrors(t *testing.T) {
	// ARRANGE / ACT / ASSERT
	_, err := Canonical("")
	require.Error(t, err)

	_, err = Canonical("0xdeadbeef") // too short
	require.Error(t, err)

	_, err = Canonical("not-a-real-address")
	require.Error(t, err)
}
