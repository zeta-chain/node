package cli_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/testdata"
	"github.com/zeta-chain/node/x/authority/client/cli"
)

func TestReadChainInfoFromFile(t *testing.T) {
	t.Run("successfully read file", func(t *testing.T) {
		fs := testdata.TypesFiles

		chain, err := cli.ReadChainFromFile(fs, "types/chain.json")
		require.NoError(t, err)

		require.EqualValues(t, chains.Chain{
			ChainId:     1,
			Network:     chains.Network_zeta,
			NetworkType: chains.NetworkType_devnet,
			Vm:          chains.Vm_svm,
			Consensus:   chains.Consensus_solana_consensus,
			IsExternal:  true,
			CctxGateway: chains.CCTXGateway_zevm,
			Name:        "testchain",
		}, chain)
	})

	t.Run("file not found", func(t *testing.T) {
		fs := testdata.TypesFiles

		_, err := cli.ReadChainFromFile(fs, "types/chain_not_found.json")
		require.Error(t, err)
	})
}
