package chains

import (
	"encoding/hex"
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
		require.Equal(t, addrStr, addr.String())
		require.Equal(t, addrStr, addr.EncodeAddress())
		require.True(t, addr.IsForNet(&chaincfg.MainNetParams))
	}
	{
		// should parse testnet taproot address
		addrStr := "tb1pzeclkt6upu8xwuksjcz36y4q56dd6jw5r543eu8j8238yaxpvcvq7t8f33"
		addr, err := DecodeTaprootAddress(addrStr)
		require.Nil(t, err)
		require.Equal(t, addrStr, addr.String())
		require.Equal(t, addrStr, addr.EncodeAddress())
		require.True(t, addr.IsForNet(&chaincfg.TestNet3Params))
	}
	{
		// should parse regtest taproot address
		addrStr := "bcrt1pqqqsyqcyq5rqwzqfpg9scrgwpugpzysnzs23v9ccrydpk8qarc0sj9hjuh"
		addr, err := DecodeTaprootAddress(addrStr)
		require.Nil(t, err)
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
	{
		// should create correct taproot address from given witness program
		// these hex string comes from link
		// https://mempool.space/tx/41f7cbaaf9a8d378d09ee86de32eebef455225520cb71015cc9a7318fb42e326
		witnessProg, err := hex.DecodeString("af06f3d4c726a3b952f2f648d86398af5ddfc3df27aa531d97203987add8db2e")
		require.Nil(t, err)
		addr, err := NewAddressTaproot(witnessProg[:], &chaincfg.MainNetParams)
		require.Nil(t, err)
		require.Equal(t, addr.EncodeAddress(), "bc1p4ur084x8y63mj5hj7eydscuc4awals7ly749x8vhyquc0twcmvhquspa5c")
	}
	{
		// should give correct ScriptAddress for taproot address
		// example comes from
		// https://blockstream.info/tx/09298a2f32f5267f419aeaf8a58c4807dcf6cac3edb59815a3b129cd8f1219b0?expand
		addrStr := "bc1p6pls9gpm24g8ntl37pajpjtuhd3y08hs5rnf9a4n0wq595hwdh9suw7m2h"
		addr, err := DecodeTaprootAddress(addrStr)
		require.Nil(t, err)
		require.Equal(t, "d07f02a03b555079aff1f07b20c97cbb62479ef0a0e692f6b37b8142d2ee6dcb", hex.EncodeToString(addr.ScriptAddress()))
	}
}
