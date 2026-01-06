package chains

import (
	"errors"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/require"
	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

func TestAddress(t *testing.T) {
	addr := NewAddress("bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	require.EqualValuesf(t, NoAddress, addr, "address string should be empty")
	require.True(t, addr.IsEmpty())

	addr = NewAddress("bogus")
	require.EqualValuesf(t, NoAddress, addr, "address string should be empty")
	require.True(t, addr.IsEmpty())

	addr = NewAddress("0x90f2b1ae50e6018230e90a33f98c7844a0ab635a")
	require.EqualValuesf(
		t,
		"0x90f2b1ae50e6018230e90a33f98c7844a0ab635a",
		addr.String(),
		"address string should be equal",
	)
	require.False(t, addr.IsEmpty())

	addr2 := NewAddress("0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5")
	require.EqualValuesf(
		t,
		"0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5",
		addr2.String(),
		"address string should be equal",
	)
	require.False(t, addr.IsEmpty())

	require.False(t, addr.Equals(addr2))
	require.True(t, addr.Equals(addr))
}

func TestDecodeBtcAddress(t *testing.T) {
	t.Run("invalid string", func(t *testing.T) {
		_, err := DecodeBtcAddress("�U�ڷ���i߭����꿚�l", BitcoinTestnet.ChainId)
		require.ErrorContains(t, err, "decode address failed: decoded address is of unknown format")
	})
	t.Run("invalid chain", func(t *testing.T) {
		_, err := DecodeBtcAddress("14CEjTd5ci3228J45GdnGeUKLSSeCWUQxK", 0)
		require.ErrorContains(t, err, "is not a bitcoin chain")
	})
	t.Run("invalid checksum", func(t *testing.T) {
		_, err := DecodeBtcAddress("tb1qy9pqmk2pd9sv63g27jt8r657wy0d9uee4x2dt2", BitcoinTestnet.ChainId)
		require.ErrorContains(t, err, "invalid checksum")
	})
	t.Run("valid legacy main-net address address incorrect params TestNet", func(t *testing.T) {
		_, err := DecodeBtcAddress("1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3", BitcoinTestnet.ChainId)
		require.ErrorContains(t, err, "decode address failed")
	})
	t.Run("valid legacy main-net address address incorrect params RegTestNet", func(t *testing.T) {
		_, err := DecodeBtcAddress("1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3", BitcoinRegtest.ChainId)
		require.ErrorContains(t, err, "decode address failed")
	})

	t.Run("valid legacy main-net address address correct params", func(t *testing.T) {
		_, err := DecodeBtcAddress("1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3", BitcoinMainnet.ChainId)
		require.NoError(t, err)
	})
	t.Run("valid legacy testnet address with correct params", func(t *testing.T) {
		_, err := DecodeBtcAddress("n2TCLD16i8SNjwPCcgGBkTEeG6CQAcYTN1", BitcoinTestnet.ChainId)
		require.NoError(t, err)
	})

	t.Run("non legacy valid address with incorrect params", func(t *testing.T) {
		_, err := DecodeBtcAddress("bcrt1qy9pqmk2pd9sv63g27jt8r657wy0d9uee4x2dt2", BitcoinMainnet.ChainId)
		require.ErrorContains(t, err, "not for network mainnet")
	})
	t.Run("non legacy valid address with correct params", func(t *testing.T) {
		_, err := DecodeBtcAddress("bcrt1qy9pqmk2pd9sv63g27jt8r657wy0d9uee4x2dt2", BitcoinRegtest.ChainId)
		require.NoError(t, err)
	})

	t.Run("taproot address with correct params", func(t *testing.T) {
		_, err := DecodeBtcAddress(
			"bc1p4ur084x8y63mj5hj7eydscuc4awals7ly749x8vhyquc0twcmvhquspa5c",
			BitcoinMainnet.ChainId,
		)
		require.NoError(t, err)
	})
	t.Run("taproot address with incorrect params", func(t *testing.T) {
		_, err := DecodeBtcAddress(
			"bc1p4ur084x8y63mj5hj7eydscuc4awals7ly749x8vhyquc0twcmvhquspa5c",
			BitcoinTestnet.ChainId,
		)
		require.ErrorContains(t, err, "not for network testnet")
	})
}

func Test_DecodeSolanaWalletAddress(t *testing.T) {
	tests := []struct {
		name         string
		address      string
		want         string
		errorMessage string
	}{
		{
			name:         "should decode a valid Solana wallet address",
			address:      "DCAK36VfExkPdAkYUQg6ewgxyinvcEyPLyHjRbmveKFw",
			want:         "DCAK36VfExkPdAkYUQg6ewgxyinvcEyPLyHjRbmveKFw",
			errorMessage: "",
		},
		{
			name:         "should fail to decode an address with invalid length",
			address:      "DCAK36VfExkPdAkYUQg6ewgxyinvcEyPLyHjRbmveK",
			want:         "",
			errorMessage: "error decoding solana wallet address",
		},
		{
			name:         "should fail to decode an invalid Base58 address",
			address:      "DCAK36VfExkPdAkYUQg6ewgxyinvcEyPLyHjRbmveKFl", // contains invalid character 'l'
			want:         "",
			errorMessage: "error decoding solana wallet address",
		},
		{
			name:         "should fail to decode a program derived address (not on ed25519 curve)",
			address:      "9dcAyYG4bawApZocwZSyJBi9Mynf5EuKAJfifXdfkqik",
			want:         "",
			errorMessage: "is not on ed25519 curve",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk, err := DecodeSolanaWalletAddress(tt.address)
			if tt.errorMessage != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMessage)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, pk.String())
		})
	}
}

func Test_IsBtcAddressSupported_P2TR(t *testing.T) {
	tests := []struct {
		name      string
		addr      string
		chainId   int64
		supported bool
	}{
		{
			// https://mempool.space/tx/259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7
			name:      "mainnet taproot address",
			addr:      "bc1p4scddlkkuw9486579autxumxmkvuphm5pz4jvf7f6pdh50p2uzqstawjt9",
			chainId:   BitcoinMainnet.ChainId,
			supported: true,
		},
		{
			name:      "mainnet taproot address with uppercase",
			addr:      strings.ToUpper("bc1p4scddlkkuw9486579autxumxmkvuphm5pz4jvf7f6pdh50p2uzqstawjt9"),
			chainId:   BitcoinMainnet.ChainId,
			supported: true,
		},
		{
			// https://mempool.space/testnet/tx/24991bd2fdc4f744bf7bbd915d4915925eecebdae249f81e057c0a6ffb700ab9
			name:      "testnet taproot address",
			addr:      "tb1p7qqaucx69xtwkx7vwmhz03xjmzxxpy3hk29y7q06mt3k6a8sehhsu5lacw",
			chainId:   BitcoinTestnet.ChainId,
			supported: true,
		},
		{
			name:      "regtest taproot address",
			addr:      "bcrt1pqqqsyqcyq5rqwzqfpg9scrgwpugpzysnzs23v9ccrydpk8qarc0sj9hjuh",
			chainId:   BitcoinRegtest.ChainId,
			supported: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// it should be a taproot address
			addr, err := DecodeBtcAddress(tt.addr, tt.chainId)
			require.NoError(t, err)
			_, ok := addr.(*btcutil.AddressTaproot)
			require.True(t, ok)

			// it should be supported
			require.NoError(t, err)
			supported := IsBtcAddressSupported(addr)
			require.Equal(t, tt.supported, supported)
		})
	}
}

func Test_IsBtcAddressSupported_P2WSH(t *testing.T) {
	tests := []struct {
		name      string
		addr      string
		chainId   int64
		supported bool
	}{
		{
			// https://mempool.space/tx/791bb9d16f7ab05f70a116d18eaf3552faf77b9d5688699a480261424b4f7e53
			name:      "mainnet P2WSH address",
			addr:      "bc1qqv6pwn470vu0tssdfha4zdk89v3c8ch5lsnyy855k9hcrcv3evequdmjmc",
			chainId:   BitcoinMainnet.ChainId,
			supported: true,
		},
		{
			name:      "mainnet P2WSH address with uppercase",
			addr:      strings.ToUpper("bc1qqv6pwn470vu0tssdfha4zdk89v3c8ch5lsnyy855k9hcrcv3evequdmjmc"),
			chainId:   BitcoinMainnet.ChainId,
			supported: true,
		},
		{
			// https://mempool.space/testnet/tx/78fac3f0d4c0174c88d21c4bb1e23a8f007e890c6d2cfa64c97389ead16c51ed
			name:      "testnet P2WSH address",
			addr:      "tb1quhassyrlj43qar0mn0k5sufyp6mazmh2q85lr6ex8ehqfhxpzsksllwrsu",
			chainId:   BitcoinTestnet.ChainId,
			supported: true,
		},
		{
			name:      "regtest P2WSH address",
			addr:      "bcrt1qm9mzhyky4w853ft2ms6dtqdyyu3z2tmrq8jg8xglhyuv0dsxzmgs2f0sqy",
			chainId:   BitcoinRegtest.ChainId,
			supported: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// it should be a P2WSH address
			addr, err := DecodeBtcAddress(tt.addr, tt.chainId)
			require.NoError(t, err)
			_, ok := addr.(*btcutil.AddressWitnessScriptHash)
			require.True(t, ok)

			// it should be supported
			supported := IsBtcAddressSupported(addr)
			require.Equal(t, tt.supported, supported)
		})
	}
}

func Test_IsBtcAddressSupported_P2WPKH(t *testing.T) {
	tests := []struct {
		name      string
		addr      string
		chainId   int64
		supported bool
	}{
		{
			// https://mempool.space/tx/5d09d232bfe41c7cb831bf53fc2e4029ab33a99087fd5328a2331b52ff2ebe5b
			name:      "mainnet P2WPKH address",
			addr:      "bc1qaxf82vyzy8y80v000e7t64gpten7gawewzu42y",
			chainId:   BitcoinMainnet.ChainId,
			supported: true,
		},
		{
			name:      "mainnet P2WPKH address with uppercase",
			addr:      strings.ToUpper("bc1qaxf82vyzy8y80v000e7t64gpten7gawewzu42y"),
			chainId:   BitcoinMainnet.ChainId,
			supported: true,
		},
		{
			// https://mempool.space/testnet/tx/508b4d723c754bad001eae9b7f3c12377d3307bd5b595c27fd8a90089094f0e9
			name:      "testnet P2WPKH address",
			addr:      "tb1q6rufg6myrxurdn0h57d2qhtm9zfmjw2mzcm05q",
			chainId:   BitcoinTestnet.ChainId,
			supported: true,
		},
		{
			name:      "regtest P2WPKH address",
			addr:      "bcrt1qy9pqmk2pd9sv63g27jt8r657wy0d9uee4x2dt2",
			chainId:   BitcoinRegtest.ChainId,
			supported: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// it should be a P2WPKH address
			addr, err := DecodeBtcAddress(tt.addr, tt.chainId)
			require.NoError(t, err)
			_, ok := addr.(*btcutil.AddressWitnessPubKeyHash)
			require.True(t, ok)

			// it should be supported
			supported := IsBtcAddressSupported(addr)
			require.Equal(t, tt.supported, supported)
		})
	}
}

func Test_IsBtcAddressSupported_P2SH(t *testing.T) {
	tests := []struct {
		name      string
		addr      string
		chainId   int64
		supported bool
	}{
		{
			// https://mempool.space/tx/fd68c8b4478686ca6f5ae4c28eaab055490650dbdaa6c2c8e380a7e075958a21
			name:      "mainnet P2SH address",
			addr:      "327z4GyFM8Y8DiYfasGKQWhRK4MvyMSEgE",
			chainId:   BitcoinMainnet.ChainId,
			supported: true,
		},
		{
			// https://mempool.space/testnet/tx/0c8c8f94817e0288a5273f5c971adaa3cee18a895c3ec8544785dddcd96f3848
			name:      "testnet P2SH address 1",
			addr:      "2N6AoUj3KPS7wNGZXuCckh8YEWcSYNsGbqd",
			chainId:   BitcoinTestnet.ChainId,
			supported: true,
		},
		{
			// https://mempool.space/testnet/tx/b5e074c5e021fcbd91ea14b1db29dfe5d14e1a6e046039467bf6ada7f8cc01b3
			name:      "testnet P2SH address 2",
			addr:      "2MwbFpRpZWv4zREjbdLB9jVW3Q8xonpVeyE",
			chainId:   BitcoinTestnet.ChainId,
			supported: true,
		},
		{
			name:      "testnet P2SH address 1 should also be supported in regtest",
			addr:      "2N6AoUj3KPS7wNGZXuCckh8YEWcSYNsGbqd",
			chainId:   BitcoinRegtest.ChainId,
			supported: true,
		},
		{
			name:      "testnet P2SH address 2 should also be supported in regtest",
			addr:      "2MwbFpRpZWv4zREjbdLB9jVW3Q8xonpVeyE",
			chainId:   BitcoinRegtest.ChainId,
			supported: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// it should be a P2SH address
			addr, err := DecodeBtcAddress(tt.addr, tt.chainId)
			require.NoError(t, err)
			_, ok := addr.(*btcutil.AddressScriptHash)
			require.True(t, ok)

			// it should be supported
			supported := IsBtcAddressSupported(addr)
			require.Equal(t, tt.supported, supported)
		})
	}
}

func Test_IsBtcAddressSupported_P2PKH(t *testing.T) {
	tests := []struct {
		name      string
		addr      string
		chainId   int64
		supported bool
	}{
		{
			// https://mempool.space/tx/9c741de6e17382b7a9113fc811e3558981a35a360e3d1262a6675892c91322ca
			name:      "mainnet P2PKH address 1",
			addr:      "1FueivsE338W2LgifJ25HhTcVJ7CRT8kte",
			chainId:   BitcoinMainnet.ChainId,
			supported: true,
		},
		{
			// https://mempool.space/testnet/tx/1e3974386f071de7f65cabb57346c1a22ec9b3e211a96928a98149673f681237
			name:      "testnet P2PKH address 1",
			addr:      "mxpYha3UJKUgSwsAz2qYRqaDSwAkKZ3YEY",
			chainId:   BitcoinTestnet.ChainId,
			supported: true,
		},
		{
			// https://mempool.space/testnet/tx/e48459f372727f2253b0ea8c71ded83e8270873b8a044feb3435fc7a799a648f
			name:      "testnet P2PKH address 2",
			addr:      "n1gXcqxmzwqHmqmgobe1XXuJaweSu69tZz",
			chainId:   BitcoinTestnet.ChainId,
			supported: true,
		},
		{
			name:      "testnet P2PKH address should also be supported in regtest",
			addr:      "mxpYha3UJKUgSwsAz2qYRqaDSwAkKZ3YEY",
			chainId:   BitcoinRegtest.ChainId,
			supported: true,
		},
		{
			name:      "testnet P2PKH address should also be supported in regtest",
			addr:      "n1gXcqxmzwqHmqmgobe1XXuJaweSu69tZz",
			chainId:   BitcoinRegtest.ChainId,
			supported: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// it should be a P2PKH address
			addr, err := DecodeBtcAddress(tt.addr, tt.chainId)
			require.NoError(t, err)
			_, ok := addr.(*btcutil.AddressPubKeyHash)
			require.True(t, ok)

			// it should be supported
			supported := IsBtcAddressSupported(addr)
			require.Equal(t, tt.supported, supported)
		})
	}
}

func TestConvertRecoverToError(t *testing.T) {
	t.Run("recover with string", func(t *testing.T) {
		err := ConvertRecoverToError("error occurred")
		require.Error(t, err)
		require.Equal(t, "error occurred", err.Error())
	})

	t.Run("recover with error", func(t *testing.T) {
		originalErr := errors.New("original error")
		err := ConvertRecoverToError(originalErr)
		require.Error(t, err)
		require.Equal(t, originalErr, err)
	})

	t.Run("recover with non-string and non-error", func(t *testing.T) {
		err := ConvertRecoverToError(12345)
		require.Error(t, err)
		require.Equal(t, "12345", err.Error())
	})

	t.Run("recover with nil", func(t *testing.T) {
		err := ConvertRecoverToError(nil)
		require.Error(t, err)
		require.Equal(t, "<nil>", err.Error())
	})
}
