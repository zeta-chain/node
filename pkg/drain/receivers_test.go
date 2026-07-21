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

func TestResolveAnchorsTestnetUnset(t *testing.T) {
	// testnet anchors resolve but are the UNSET sentinel until configured, so Validate fails closed.
	_, receivers, err := drain.ResolveAnchors(drain.NetworkTestnet)
	require.NoError(t, err)
	require.Error(t, receivers.Validate())
}
