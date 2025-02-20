package sui

import (
	"encoding/base64"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCrypto(t *testing.T) {
	t.Run("Digest", func(t *testing.T) {
		// ARRANGE
		// Given a tx (imagine client.MoveCall(...) result)
		//
		// Data was generated using sui cli:
		// sui client transfer-sui --to "0xac5bceec1b789ff840d7d4e6ce4ce61c90d190a7f8c4f4ddf0bff6ee2413c33c" \
		// --sui-coin-object-id "0x0466a9a57add505b7b85ac485054f9b71f574f4504d9c70acd8f73ef11e0dc30" \
		// --gas-budget 500000 --serialize-unsigned-transaction
		//
		// https://docs.sui.io/concepts/cryptography/transaction-auth/offline-signing#sign
		tx := models.TxnMetaData{
			TxBytes: "AAABACCsW87sG3if+EDX1ObOTOYckNGQp/jE9N3wv/buJBPDPAEBAQABAACsW87sG3if+EDX1ObOTOYckNGQp/jE9N3wv/buJBPDPAEEZqmlet1QW3uFrEhQVPm3H1dPRQTZxwrNj3PvEeDcMPCkyhwAAAAAICNNoyg5v4obnoVYDWw0XhxB6Tq/b+OPXnJKPc2+QM5QrFvO7Bt4n/hA19TmzkzmHJDRkKf4xPTd8L/27iQTwzzuAgAAAAAAACChBwAAAAAAAA==",
		}

		// Given expected digest based on SUI cli:
		// https://docs.sui.io/concepts/cryptography/transaction-auth/offline-signing#sign
		// sui keytool sign --address "..." --data "$txBytesBase64" --json | jq ".digest"
		const expectedDigestBase64 = "A1NY74R1IScWR/GPtOMNHVY/RwTNzWHlUbOkwp3911g="

		// ACT
		digest, err := Digest(tx)

		digestBase64 := base64.StdEncoding.EncodeToString(digest[:])

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, expectedDigestBase64, digestBase64)
	})

	t.Run("AddressFromPubKeyECDSA", func(t *testing.T) {
		// `$> sui keytool generate secp256k1`
		for _, tt := range []struct {
			pubKey  string
			address string
		}{
			{
				pubKey:  "AQJz6a5yi6Wtf55atMWlW/ZA4Xdd6lJKC22u3Xi/h9yeBw==",
				address: "0xccf49bfb6c8159f5e53c80f5b6ecf748e4af89c8c10c27d24302207b2bc97744",
			},
			{
				pubKey:  "AQKUgO1kyhheTjbzYYhP67nxDD1UZwEhqkLyX1KRBm1xTQ==",
				address: "0x2dc141f8a8d8a3fe397054f538dcc8207fd5edf4a1415c54b7d5a4dc124d9b3e",
			},
			{
				pubKey:  "AQIgwiNQwm529+fvaKW/n5ITbaQVUToZq+ZIpNjjOw7Spw==",
				address: "0x17012be22c34ad1396f8af272b2e7b0edb529b3441912bd532faf874bf2c9262",
			},
		} {
			// ARRANGE
			pubKeyBytes, err := base64.StdEncoding.DecodeString(tt.pubKey)
			require.NoError(t, err)

			// type_flag + compression_flag + 32bytes
			require.Equal(t, 1+1+32, len(pubKeyBytes))

			pk, err := crypto.DecompressPubkey(pubKeyBytes[1:])
			require.NoError(t, err)

			// ACT
			addr := AddressFromPubKeyECDSA(pk)

			// ASSERT
			assert.Equal(t, tt.address, addr)
		}
	})
}
