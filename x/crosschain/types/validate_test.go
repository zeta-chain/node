package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

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
