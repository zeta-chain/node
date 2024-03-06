package bitcoin

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
)

func TestAddressTaproot(t *testing.T) {
	{
		// should parse mainnet taproot address
		addrStr := "bc1p4ur084x8y63mj5hj7eydscuc4awals7ly749x8vhyquc0twcmvhquspa5c"
		addr, err := DecodeTaprootAddress(addrStr)
		require.Nil(t, err)
		t.Log(addr.String())
		require.Equal(t, addrStr, addr.String())
		require.Equal(t, addrStr, addr.EncodeAddress())
		require.True(t, addr.IsForNet(&chaincfg.MainNetParams))
	}
	{
		// should parse testnet taproot address
		addrStr := "tb1pzeclkt6upu8xwuksjcz36y4q56dd6jw5r543eu8j8238yaxpvcvq7t8f33"
		addr, err := DecodeTaprootAddress(addrStr)
		require.Nil(t, err)
		t.Log(addr.String())
		require.Equal(t, addrStr, addr.String())
		require.Equal(t, addrStr, addr.EncodeAddress())
		require.True(t, addr.IsForNet(&chaincfg.TestNet3Params))
	}
	{
		// should parse regtest taproot address
		addrStr := "bcrt1pqqqsyqcyq5rqwzqfpg9scrgwpugpzysnzs23v9ccrydpk8qarc0sj9hjuh"
		addr, err := DecodeTaprootAddress(addrStr)
		require.Nil(t, err)
		t.Log(addr.String())
		require.Equal(t, addrStr, addr.String())
		require.Equal(t, addrStr, addr.EncodeAddress())
		require.True(t, addr.IsForNet(&chaincfg.RegressionNetParams))
	}

	{
		// should fail to parse invalid taproot address
		// should parse mainnet taproot address
		addrStr := "bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu"
		_, err := DecodeTaprootAddress(addrStr)
		require.Error(t, err)
	}
	{
		var witnessProg [32]byte
		for i := 0; i < 32; i++ {
			witnessProg[i] = byte(i)
		}
		_, err := newAddressTaproot("bcrt", witnessProg[:])
		require.Nil(t, err)
		//t.Logf("addr: %v", addr)
	}

}
