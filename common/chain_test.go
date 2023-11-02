package common

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChain_EncodeAddress(t *testing.T) {
	type fields struct {
		ChainName ChainName
		ChainId   int32
	}

	tests := []struct {
		name    string
		chain   Chain
		b       []byte
		want    string
		wantErr bool
	}{
		{
			name: "should error if b is not a valid address on the network",
			chain: Chain{
				ChainName: ChainName_btc_testnet,
				ChainId:   18332,
			},
			b:       []byte("bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c"),
			want:    "",
			wantErr: true,
		},
		{
			name: "should pass if b is a valid address on the network",
			chain: Chain{
				ChainName: ChainName_btc_mainnet,
				ChainId:   8332,
			},
			b:       []byte("bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c"),
			want:    "bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c",
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			s, err := tc.chain.EncodeAddress(tc.b)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.Equal(t, tc.want, s)
		})
	}
}
