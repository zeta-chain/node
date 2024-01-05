package common

import (
	"testing"

	"github.com/stretchr/testify/require"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

func TestAddress(t *testing.T) {
	addr := NewAddress("bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	require.EqualValuesf(t, NoAddress, addr, "address string should be empty")

	addr = NewAddress("bogus")
	require.EqualValuesf(t, NoAddress, addr, "address string should be empty")

	addr = NewAddress("0x90f2b1ae50e6018230e90a33f98c7844a0ab635a")
	require.EqualValuesf(t, "0x90f2b1ae50e6018230e90a33f98c7844a0ab635a", addr.String(), "address string should be equal")
}

func TestDecodeBtcAddress(t *testing.T) {
	// �U�ڷ���i߭����꿚�l
	// 14CEjTd5ci3228J45GdnGeUKLSSeCWUQxK
	t.Run("invalid string", func(t *testing.T) {
		_, err := DecodeBtcAddress("�U�ڷ���i߭����꿚�l", 18332)
		require.ErrorContains(t, err, "runtime error: index out of range")
	})
	t.Run("invalid chain", func(t *testing.T) {
		_, err := DecodeBtcAddress("14CEjTd5ci3228J45GdnGeUKLSSeCWUQxK", 0)
		require.ErrorContains(t, err, "is not a Bitcoin chain")
	})
	t.Run("invalid checksum", func(t *testing.T) {
		_, err := DecodeBtcAddress("tb1qy9pqmk2pd9sv63g27jt8r657wy0d9uee4x2dt2", 18332)
		require.ErrorContains(t, err, "invalid checksum")
	})
	t.Run("valid address", func(t *testing.T) {
		_, err := DecodeBtcAddress("bcrt1qy9pqmk2pd9sv63g27jt8r657wy0d9uee4x2dt2", 18444)
		require.NoError(t, err)
	})

}
