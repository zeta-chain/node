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
	t.Run("invalid string", func(t *testing.T) {
		_, err := DecodeBtcAddress("�U�ڷ���i߭����꿚�l", BtcTestNetChain().ChainId)
		require.ErrorContains(t, err, "runtime error: index out of range")
	})
	t.Run("invalid chain", func(t *testing.T) {
		_, err := DecodeBtcAddress("14CEjTd5ci3228J45GdnGeUKLSSeCWUQxK", 0)
		require.ErrorContains(t, err, "is not a bitcoin chain")
	})
	t.Run("invalid checksum", func(t *testing.T) {
		_, err := DecodeBtcAddress("tb1qy9pqmk2pd9sv63g27jt8r657wy0d9uee4x2dt2", BtcTestNetChain().ChainId)
		require.ErrorContains(t, err, "invalid checksum")
	})
	t.Run("valid legacy main-net address address incorrect params TestNet", func(t *testing.T) {
		_, err := DecodeBtcAddress("1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3", BtcTestNetChain().ChainId)
		require.ErrorContains(t, err, "decode address failed")
	})
	t.Run("valid legacy main-net address address incorrect params RegTestNet", func(t *testing.T) {
		_, err := DecodeBtcAddress("1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3", BtcRegtestChain().ChainId)
		require.ErrorContains(t, err, "decode address failed")
	})

	t.Run("valid legacy main-net address address correct params", func(t *testing.T) {
		_, err := DecodeBtcAddress("1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3", BtcMainnetChain().ChainId)
		require.NoError(t, err)
	})
	t.Run("valid legacy testnet address with correct params", func(t *testing.T) {
		_, err := DecodeBtcAddress("n2TCLD16i8SNjwPCcgGBkTEeG6CQAcYTN1", BtcTestNetChain().ChainId)
		require.NoError(t, err)
	})

	t.Run("non legacy valid address with incorrect params", func(t *testing.T) {
		_, err := DecodeBtcAddress("bcrt1qy9pqmk2pd9sv63g27jt8r657wy0d9uee4x2dt2", BtcMainnetChain().ChainId)
		require.ErrorContains(t, err, "address is not for network main-net")
	})
	t.Run("non legacy valid address with correct params", func(t *testing.T) {
		_, err := DecodeBtcAddress("bcrt1qy9pqmk2pd9sv63g27jt8r657wy0d9uee4x2dt2", BtcTestNetChain().ChainId)
		require.NoError(t, err)
	})
}
