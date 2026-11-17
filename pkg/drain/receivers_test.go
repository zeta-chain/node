package drain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/drain"
)

func TestReceiversValidate(t *testing.T) {
	tests := []struct {
		name    string
		r       drain.Receivers
		wantErr bool
	}{
		{"valid", drain.Receivers{EVM: "0x74D6F908a320Fed7E1c0002eBa7996C4376A8071", BTC: "bcrt1qzyfpx9q4zct3sxg6rvwp68slyqsjygeyuwdjcu"}, false},
		{"zero evm", drain.Receivers{EVM: "0x0000000000000000000000000000000000000000", BTC: "bcrt1qx"}, true},
		{"empty evm", drain.Receivers{EVM: "", BTC: "bcrt1qx"}, true},
		{"unset evm", drain.Receivers{EVM: "UNSET", BTC: "bcrt1qx"}, true},
		{"bad evm", drain.Receivers{EVM: "not-an-address", BTC: "bcrt1qx"}, true},
		{"empty btc", drain.Receivers{EVM: "0x74D6F908a320Fed7E1c0002eBa7996C4376A8071", BTC: ""}, true},
		{"valid override", drain.Receivers{EVM: "0x74D6F908a320Fed7E1c0002eBa7996C4376A8071", BTC: "bcrt1qx", EVMByChain: map[int64]string{8453: "0x919F207971b9Dea34150D095169388d25e167F1A"}}, false},
		{"non-checksummed override", drain.Receivers{EVM: "0x74D6F908a320Fed7E1c0002eBa7996C4376A8071", BTC: "bcrt1qx", EVMByChain: map[int64]string{8453: "0x919f207971b9dea34150d095169388d25e167f1a"}}, true},
		{"zero override", drain.Receivers{EVM: "0x74D6F908a320Fed7E1c0002eBa7996C4376A8071", BTC: "bcrt1qx", EVMByChain: map[int64]string{8453: "0x0000000000000000000000000000000000000000"}}, true},
		{"bad override", drain.Receivers{EVM: "0x74D6F908a320Fed7E1c0002eBa7996C4376A8071", BTC: "bcrt1qx", EVMByChain: map[int64]string{8453: "UNSET"}}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.r.Validate()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestResolveAnchorsTestnetUnset(t *testing.T) {
	// testnet anchors resolve but are the UNSET sentinel until configured, so Validate fails closed.
	_, receivers, err := drain.ResolveAnchors(drain.NetworkTestnet)
	require.NoError(t, err)
	require.Error(t, receivers.Validate())
}

func TestResolveAnchorsMainnetPerChainEVM(t *testing.T) {
	_, receivers, err := drain.ResolveAnchors(drain.NetworkMainnet)
	require.NoError(t, err)
	require.NoError(t, receivers.Validate())

	// Pin the reviewed mainnet anchors: chains already drained keep the original receiver,
	// every other EVM chain goes to the Gnosis Safe.
	const (
		eoa  = "0x0a538985123729f48D70DBCaE82a7f47f1CbA8f8"
		safe = "0x919F207971b9Dea34150D095169388d25e167F1A"
	)
	for chainID, want := range map[int64]string{1: eoa, 56: eoa, 137: eoa, 8453: safe, 42161: safe, 43114: safe} {
		require.Equal(t, want, receivers.EVMFor(chainID), "chain %d", chainID)
	}
	require.Equal(t, "bc1qkl02lqffhmpf5hn3khnetc0ay9yyk4eajmlefy", receivers.BTC)
	require.Len(t, receivers.EVMByChain, 3)
}
