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

	chainInfo, err := cli.ReadChainInfoFromFile(fs, "types/chain_info.json")
	require.NoError(t, err)

	require.Len(t, chainInfo.Chains, 2)
	require.EqualValues(t, chains.Chain{
		ChainId:     42,
		Network:     chains.Network_eth,
		NetworkType: chains.NetworkType_mainnet,
		Vm:          chains.Vm_no_vm,
		Consensus:   chains.Consensus_ethereum,
		IsExternal:  false,
		CctxGateway: chains.CCTXGateway_zevm,
		Name:        "foo",
	}, chainInfo.Chains[0])
	require.EqualValues(t, chains.Chain{
		ChainId:     84,
		Network:     chains.Network_zeta,
		NetworkType: chains.NetworkType_testnet,
		Vm:          chains.Vm_evm,
		Consensus:   chains.Consensus_tendermint,
		IsExternal:  true,
		CctxGateway: chains.CCTXGateway_observers,
		Name:        "bar",
	}, chainInfo.Chains[1])
}
