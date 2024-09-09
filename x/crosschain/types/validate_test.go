package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestValidateAddressForChain(t *testing.T) {
	// test for eth chain
	require.Error(t, types.ValidateAddressForChain("0x123", chains.Goerli.ChainId))
	require.Error(t, types.ValidateAddressForChain("", chains.Goerli.ChainId))
	require.Error(t, types.ValidateAddressForChain("%%%%", chains.Goerli.ChainId))
	require.NoError(
		t,
		types.ValidateAddressForChain("0x792c127Fa3AC1D52F904056Baf1D9257391e7D78", chains.Goerli.ChainId),
	)

	// test for btc chain
	require.NoError(
		t,
		types.ValidateAddressForChain(
			"bc1p4scddlkkuw9486579autxumxmkvuphm5pz4jvf7f6pdh50p2uzqstawjt9",
			chains.BitcoinMainnet.ChainId,
		),
	)
	require.NoError(
		t,
		types.ValidateAddressForChain("327z4GyFM8Y8DiYfasGKQWhRK4MvyMSEgE", chains.BitcoinMainnet.ChainId),
	)
	require.NoError(
		t,
		types.ValidateAddressForChain("1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3", chains.BitcoinMainnet.ChainId),
	)
	require.Error(
		t,
		types.ValidateAddressForChain("bcrt1qs758ursh4q9z627kt3pp5yysm78ddny6txaqgw", chains.BitcoinMainnet.ChainId),
	)
	require.Error(t, types.ValidateAddressForChain("", chains.BitcoinRegtest.ChainId))
	require.NoError(
		t,
		types.ValidateAddressForChain("bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu", chains.BitcoinMainnet.ChainId),
	)
	require.NoError(
		t,
		types.ValidateAddressForChain("bcrt1qs758ursh4q9z627kt3pp5yysm78ddny6txaqgw", chains.BitcoinRegtest.ChainId),
	)

	// test for zeta chain
	require.NoError(
		t,
		types.ValidateAddressForChain("bcrt1qs758ursh4q9z627kt3pp5yysm78ddny6txaqgw", chains.ZetaChainMainnet.ChainId),
	)
	require.NoError(
		t,
		types.ValidateAddressForChain("0x792c127Fa3AC1D52F904056Baf1D9257391e7D78", chains.ZetaChainMainnet.ChainId),
	)
}

func TestValidateCCTXIndex(t *testing.T) {
	require.NoError(t, types.ValidateCCTXIndex("0x84bd5c9922b63c52d8a9ca686e0a57ff978150b71be0583514d01c27aa341910"))
	require.NoError(t, types.ValidateCCTXIndex(sample.ZetaIndex(t)))
	require.Error(t, types.ValidateCCTXIndex("0"))
	require.Error(t, types.ValidateCCTXIndex("0x70e967acFcC17c3941E87562161406d41676FD83"))
}

func TestValidateHashForChain(t *testing.T) {
	require.NoError(
		t,
		types.ValidateHashForChain(
			"0x84bd5c9922b63c52d8a9ca686e0a57ff978150b71be0583514d01c27aa341910",
			chains.Goerli.ChainId,
		),
	)
	require.Error(t, types.ValidateHashForChain("", chains.Goerli.ChainId))
	require.Error(
		t,
		types.ValidateHashForChain(
			"a0fa5a82f106fb192e4c503bfa8d54b2de20a821e09338094ab825cc9b275059",
			chains.Goerli.ChainId,
		),
	)
	require.NoError(
		t,
		types.ValidateHashForChain(
			"15b7880f5d236e857a5e8f043ce9d56f5ef01e1c3f2a786baf740fc0bb7a22a3",
			chains.BitcoinMainnet.ChainId,
		),
	)
	require.NoError(
		t,
		types.ValidateHashForChain(
			"a0fa5a82f106fb192e4c503bfa8d54b2de20a821e09338094ab825cc9b275059",
			chains.BitcoinTestnet.ChainId,
		),
	)
	require.Error(
		t,
		types.ValidateHashForChain(
			"0x84bd5c9922b63c52d8a9ca686e0a57ff978150b71be0583514d01c27aa341910",
			chains.BitcoinMainnet.ChainId,
		),
	)
}
