package sui

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

// Payload sample extracted from https://github.com/zeta-chain/example-contracts/pull/250
// npx ts-node sui/setup/encodeCallArgs.ts \
// "$TOKEN_TYPE" \
// "$CONFIG,$POOL,$PARTNER,$CLOCK" \
// "$MESSAGE"
const formattedPayloadSample = "00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000001c00000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000503078623131326633373062633865336261366534356164316139353436363030393966633365366465326132303364663964323665313161613064383730663633353a3a746f6b656e3a3a544f4b454e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000457dd7b5841300199ac87b420ddeb48229523e76af423b4fce37da0cb78604408bab1a2d90ea585eab574932e1b3467ff1d5d3f2aee55fed304f963ca2b9209ebee6f1f44d24a8bf7268d82425d6e7bd8b9c48d11b2119b20756ee150c8e24ac3039ce62b538a0d0fca21c3c3a5b99adf519d55e534c536568fbcca40ee61fb7e00000000000000000000000000000000000000000000000000000000000000203573924024f4a7ff8e6755cb2d9fdeef69bdb65329f081d21b0b6ab37a265d06"

func TestCallPayload_UnpackABI(t *testing.T) {
	t.Run("can parse a correctly formatted payload", func(t *testing.T) {
		// ARRANGE
		payload, err := hex.DecodeString(formattedPayloadSample)
		require.NoError(t, err)

		var cp CallPayload

		// ACT
		err = cp.UnpackABI(payload)

		// ASSERT
		require.NoError(t, err)
		require.EqualValues(t, []string{
			"0xb112f370bc8e3ba6e45ad1a954660099fc3e6de2a203df9d26e11aa0d870f635::token::TOKEN",
		}, cp.TypeArgs)
		require.EqualValues(t, []string{
			"0x57dd7b5841300199ac87b420ddeb48229523e76af423b4fce37da0cb78604408",
			"0xbab1a2d90ea585eab574932e1b3467ff1d5d3f2aee55fed304f963ca2b9209eb",
			"0xee6f1f44d24a8bf7268d82425d6e7bd8b9c48d11b2119b20756ee150c8e24ac3",
			"0x039ce62b538a0d0fca21c3c3a5b99adf519d55e534c536568fbcca40ee61fb7e",
		}, cp.ObjectIDs)
		expectedMessage, err := hex.DecodeString("3573924024f4a7ff8e6755cb2d9fdeef69bdb65329f081d21b0b6ab37a265d06")
		require.NoError(t, err)
		require.EqualValues(t, expectedMessage, cp.Message)
	})

	t.Run("unable to unpack ABI encoded payload", func(t *testing.T) {
		payload, err := hex.DecodeString("deadbeef")
		require.NoError(t, err)

		// ACT
		var cp CallPayload
		err = cp.UnpackABI(payload)

		// ASSERT
		require.ErrorIs(t, err, ErrInvalidPayload)
		require.ErrorContains(t, err, "unable to unpack ABI encoded payload (deadbeef):")
	})
}

func TestFormatWithdrawAndCallPayload(t *testing.T) {
	t.Run("can format a payload", func(t *testing.T) {
		// ARRANGE
		typeArgs := []string{
			"0xb112f370bc8e3ba6e45ad1a954660099fc3e6de2a203df9d26e11aa0d870f635::token::TOKEN",
		}
		objects := []string{
			"0x57dd7b5841300199ac87b420ddeb48229523e76af423b4fce37da0cb78604408",
			"0xbab1a2d90ea585eab574932e1b3467ff1d5d3f2aee55fed304f963ca2b9209eb",
			"0xee6f1f44d24a8bf7268d82425d6e7bd8b9c48d11b2119b20756ee150c8e24ac3",
			"0x039ce62b538a0d0fca21c3c3a5b99adf519d55e534c536568fbcca40ee61fb7e",
		}
		message, err := hex.DecodeString("3573924024f4a7ff8e6755cb2d9fdeef69bdb65329f081d21b0b6ab37a265d06")
		require.NoError(t, err)

		cp := NewCallPayload(typeArgs, objects, message)

		// ACT
		payload, err := cp.PackABI()

		// ASSERT
		require.NoError(t, err)
		require.EqualValues(t, formattedPayloadSample, hex.EncodeToString(payload))
	})

	t.Run("invalid object", func(t *testing.T) {
		// ARRANGE
		typeArgs := []string{
			"0xb112f370bc8e3ba6e45ad1a954660099fc3e6de2a203df9d26e11aa0d870f635::token::TOKEN",
		}
		objects := []string{
			"0x57dd7b5841300199ac87b420ddeb48229523e76af423b4fce37da0cb78604408",
			"0xbab1a2d90ea585eab574932e1b3467ff1d5d3f2aee55fed304f963ca2b9209eb",
			"0xee6f1f44d24a8bf7268d82425d6e7bd8b9c48d11b2119b20756ee150c8e24a",
			"0x039ce62b538a0d0fca21c3c3a5b99adf519d55e534c536568fbcca40ee61fb7e",
		}
		message, err := hex.DecodeString("3573924024f4a7ff8e6755cb2d9fdeef69bdb65329f081d21b0b6ab37a265d06")
		require.NoError(t, err)

		cp := NewCallPayload(typeArgs, objects, message)

		// ACT
		_, err = cp.PackABI()

		// ASSERT
		require.Error(t, err)
	})
}

func TestFormatAndParseWithdrawAndCallPayload(t *testing.T) {
	decodeHex := func(s string) []byte {
		b, err := hex.DecodeString(s)
		require.NoError(t, err)
		return b
	}

	tests := []struct {
		name        string
		callPayload CallPayload
	}{
		{
			name: "sample payload",
			callPayload: CallPayload{
				TypeArgs: []string{
					"0xb112f370bc8e3ba6e45ad1a954660099fc3e6de2a203df9d26e11aa0d870f635::token::TOKEN",
				},
				ObjectIDs: []string{
					"0x57dd7b5841300199ac87b420ddeb48229523e76af423b4fce37da0cb78604408",
					"0xbab1a2d90ea585eab574932e1b3467ff1d5d3f2aee55fed304f963ca2b9209eb",
					"0xee6f1f44d24a8bf7268d82425d6e7bd8b9c48d11b2119b20756ee150c8e24ac3",
					"0x039ce62b538a0d0fca21c3c3a5b99adf519d55e534c536568fbcca40ee61fb7e",
				},
				Message: decodeHex("3573924024f4a7ff8e6755cb2d9fdeef69bdb65329f081d21b0b6ab37a265d06"),
			},
		},
		{
			name: "no argument types",
			callPayload: CallPayload{
				TypeArgs: []string{},
				ObjectIDs: []string{
					"0x57dd7b5841300199ac87b420ddeb48229523e76af423b4fce37da0cb78604408",
					"0xbab1a2d90ea585eab574932e1b3467ff1d5d3f2aee55fed304f963ca2b9209eb",
					"0xee6f1f44d24a8bf7268d82425d6e7bd8b9c48d11b2119b20756ee150c8e24ac3",
					"0x039ce62b538a0d0fca21c3c3a5b99adf519d55e534c536568fbcca40ee61fb7e",
				},
				Message: decodeHex("3573924024f4a7ff8e6755cb2d9fdeef69bdb65329f081d21b0b6ab37a265d06"),
			},
		},
		{
			name: "no objects",
			callPayload: CallPayload{
				TypeArgs: []string{
					"0xb112f370bc8e3ba6e45ad1a954660099fc3e6de2a203df9d26e11aa0d870f635::token::TOKEN",
				},
				ObjectIDs: []string{},
				Message:   decodeHex("3573924024f4a7ff8e6755cb2d9fdeef69bdb65329f081d21b0b6ab37a265d06"),
			},
		},
		{
			name: "empty message",
			callPayload: CallPayload{
				TypeArgs: []string{
					"0xb112f370bc8e3ba6e45ad1a954660099fc3e6de2a203df9d26e11aa0d870f635::token::TOKEN",
				},
				ObjectIDs: []string{
					"0x57dd7b5841300199ac87b420ddeb48229523e76af423b4fce37da0cb78604408",
					"0xbab1a2d90ea585eab574932e1b3467ff1d5d3f2aee55fed304f963ca2b9209eb",
					"0xee6f1f44d24a8bf7268d82425d6e7bd8b9c48d11b2119b20756ee150c8e24ac3",
					"0x039ce62b538a0d0fca21c3c3a5b99adf519d55e534c536568fbcca40ee61fb7e",
				},
				Message: []byte{},
			},
		},
		{
			name: "empty payload",
			callPayload: CallPayload{
				TypeArgs:  []string{},
				ObjectIDs: []string{},
				Message:   []byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT
			payload, err := tt.callPayload.PackABI()
			require.NoError(t, err)

			var cp CallPayload
			require.NoError(t, cp.UnpackABI(payload))

			// ASSERT
			require.EqualValues(t, tt.callPayload, cp)
		})
	}
}
