package memo_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/memo"
)

func Test_GetLenBytes(t *testing.T) {
	// Define table-driven test cases
	tests := []struct {
		name        string
		encodeFmt   memo.EncodingFormat
		expectedLen int
		expectErr   bool
	}{
		{
			name:        "compact short",
			encodeFmt:   memo.EncodingFmtCompactShort,
			expectedLen: 1,
		},
		{
			name:        "compact long",
			encodeFmt:   memo.EncodingFmtCompactLong,
			expectedLen: 2,
		},
		{
			name:        "non-compact encoding format",
			encodeFmt:   memo.EncodingFmtABI,
			expectedLen: 0,
			expectErr:   true,
		},
	}

	// Loop through each test case
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			length, err := memo.GetLenBytes(tc.encodeFmt)

			// Check if error is expected
			if tc.expectErr {
				require.Error(t, err)
				require.Equal(t, 0, length)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedLen, length)
			}
		})
	}
}

func Test_GetCodec(t *testing.T) {
	// Define table-driven test cases
	tests := []struct {
		name      string
		encodeFmt memo.EncodingFormat
		errMsg    string
	}{
		{
			name:      "should get ABI codec",
			encodeFmt: memo.EncodingFmtABI,
		},
		{
			name:      "should get compact codec",
			encodeFmt: memo.EncodingFmtCompactShort,
		},
		{
			name:      "should get compact codec",
			encodeFmt: memo.EncodingFmtCompactLong,
		},
		{
			name:      "should fail to get codec",
			encodeFmt: 0b0011,
			errMsg:    "invalid encoding format",
		},
	}

	// Loop through each test case
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			codec, err := memo.GetCodec(tc.encodeFmt)
			if tc.errMsg != "" {
				require.Error(t, err)
				require.Nil(t, codec)
			} else {
				require.NoError(t, err)
				require.NotNil(t, codec)
			}
		})
	}
}
