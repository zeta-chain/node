package memo_test

import (
	"testing"

	"github.com/test-go/testify/require"
	"github.com/zeta-chain/node/pkg/memo"
)

func Test_GetLenBytes(t *testing.T) {
	// Define table-driven test cases
	tests := []struct {
		name        string
		encodingFmt uint8
		expectedLen int
		expectErr   bool
	}{
		{
			name:        "compact short",
			encodingFmt: memo.EncodingFmtCompactShort,
			expectedLen: 1,
		},
		{
			name:        "compact long",
			encodingFmt: memo.EncodingFmtCompactLong,
			expectedLen: 2,
		},
		{
			name:        "non-compact encoding format",
			encodingFmt: memo.EncodingFmtABI,
			expectedLen: 0,
			expectErr:   true,
		},
	}

	// Loop through each test case
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			length, err := memo.GetLenBytes(tc.encodingFmt)

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
		name        string
		encodingFmt uint8
		errMsg      string
	}{
		{
			name:        "should get ABI codec",
			encodingFmt: memo.EncodingFmtABI,
		},
		{
			name:        "should get compact codec",
			encodingFmt: memo.EncodingFmtCompactShort,
		},
		{
			name:        "should get compact codec",
			encodingFmt: memo.EncodingFmtCompactLong,
		},
		{
			name:        "should fail to get codec",
			encodingFmt: 0b11,
			errMsg:      "unsupported",
		},
	}

	// Loop through each test case
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			codec, err := memo.GetCodec(tc.encodingFmt)
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
