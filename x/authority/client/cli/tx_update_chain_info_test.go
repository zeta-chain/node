package cli_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/testdata"
	"github.com/zeta-chain/node/x/authority/client/cli"
)

func TestReadChainInfoFromFile(t *testing.T) {
	fs := testdata.TypesFiles

	chain, err := cli.ReadChainFromFile(fs, "types/chain.json")
	require.NoError(t, err)

	require.EqualValues(t, chains.Chain{
		ChainId:     42,
		Network:     chains.Network_eth,
		NetworkType: chains.NetworkType_mainnet,
		Vm:          chains.Vm_no_vm,
		Consensus:   chains.Consensus_ethereum,
		IsExternal:  false,
		CctxGateway: chains.CCTXGateway_zevm,
		Name:        "foo",
	}, chain)
}
