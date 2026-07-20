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

func TestResolveAnchorsLocalnetEnvOverride(t *testing.T) {
	t.Setenv(drain.EnvLocalnetPubKey, "0x0284bf7562262bbd6940085748f3be6afa52ae317155181ece31b66351ccffa4b0")
	t.Setenv(drain.EnvLocalnetEVMReceiver, "0x971d9a4763D845F4346D39292b849C567184D201")
	t.Setenv(drain.EnvLocalnetBTCReceiver, "bcrt1q5zs69gay5kn2029f4246etdw47ctrv4nvzy38a")

	pub, receivers, err := drain.ResolveAnchors(drain.NetworkLocalnet)
	require.NoError(t, err)
	require.Len(t, pub, 33)
	require.Equal(t, "0x971d9a4763D845F4346D39292b849C567184D201", receivers.EVM)
	require.Equal(t, "bcrt1q5zs69gay5kn2029f4246etdw47ctrv4nvzy38a", receivers.BTC)
	require.NoError(t, receivers.Validate())
}

func TestResolveAnchorsTestnetIgnoresEnv(t *testing.T) {
	// env overrides must NOT apply to testnet/mainnet — the compiled anchors win.
	t.Setenv(drain.EnvLocalnetEVMReceiver, "0x971d9a4763D845F4346D39292b849C567184D201")

	_, receivers, err := drain.ResolveAnchors(drain.NetworkTestnet)
	require.NoError(t, err)
	require.NotEqual(t, "0x971d9a4763D845F4346D39292b849C567184D201", receivers.EVM)
	// testnet is unset by default -> fails closed
	require.Error(t, receivers.Validate())
}
