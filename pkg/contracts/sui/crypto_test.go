package sui

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/cosmos/cosmos-sdk/types/bech32"
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

	t.Run("SerializeSignatureECDSA", func(t *testing.T) {
		// ARRANGE
		// Given a pubKey
		enc, _ := hex.DecodeString(
			"04760c4460e5336ac9bbd87952a3c7ec4363fc0a97bd31c86430806e287b437fd1b01abc6e1db640cf3106b520344af1d58b00b57823db3e1407cbc433e1b6d04d",
		)
		pubKey, err := crypto.UnmarshalPubkey(enc)
		require.NoError(t, err)

		// Given signature
		signature := [65]byte{0, 1, 3}

		// ACT
		res, err := SerializeSignatureECDSA(signature, pubKey)

		// ASSERT
		require.NoError(t, err)

		// Check signature
		resBin, err := base64.StdEncoding.DecodeString(res)
		require.NoError(t, err)
		require.Equal(t, (1 + 64 + 33), len(resBin))

		// ACT 2
		pubKey2, signature2, err := DeserializeSignatureECDSA(res)

		// ASSERT 2
		require.NoError(t, err)
		assert.True(t, pubKey2.Equal(pubKey))
		assert.Equal(t, signature[:64], signature2[:])
	})

	t.Run("PrivateKeySecp256k1FromHex", func(t *testing.T) {
		for _, tt := range []struct {
			privKeyHex             string
			privKeyBech32Secp256k1 string
			errMsg                 string
		}{
			{
				privKeyHex:             "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263",
				privKeyBech32Secp256k1: "suiprivkey1q8v8htmm7mw9vz39yktx0rqjus0h695zsdlstv5agydu8au2utpxxgjwf3h",
			},
			{
				privKeyHex: "invalid",
				errMsg:     "failed to decode private key hex",
			},
			{
				privKeyHex: "abcdef",
				errMsg:     "invalid private key length",
			},
		} {
			t.Run(tt.privKeyHex, func(t *testing.T) {
				privKey, err := PrivateKeyBech32Secp256k1FromHex(tt.privKeyHex)
				if tt.errMsg != "" {
					require.Empty(t, privKey)
					require.Contains(t, err.Error(), tt.errMsg)
				} else {
					require.NoError(t, err)
					require.Equal(t, tt.privKeyBech32Secp256k1, privKey)
				}
			})
		}
	})

	t.Run("SignerSecp256k1", func(t *testing.T) {
		for _, tt := range []struct {
			privKey string
			address string
		}{
			{
				privKey: "suiprivkey1q8h7ejwfcdn6gc2x9kddtd9vld3kvuvtr9gdtj9qcf7nqnw94tvl79cwfq4",
				address: "0x68f6d05fd44366bd431405e8ea32e2d2f8e330d98e0089c9447ecfbbdf6e236f",
			},
			{
				privKey: "suiprivkey1qxghtp2vncr94s8h7ctvgf58gy27l9nch75ka2jh37qr90xuyxhlwk5khxc",
				address: "0x8ec6f13affbf48c73550567f2a1cb8f05474c0133ebf745f91a2f3b971c1ec9a",
			},
			{
				privKey: "suiprivkey1q99wkv3fj352cn97yu5r9jwqhcvyyk6t9scwrftyjgqgflandfgc70r74hg",
				address: "0xa0f6b839f7945065ebdd030cec8e6e30d01a74c8cb31b1945909fd426c2cef80",
			},
		} {
			t.Run(tt.privKey, func(t *testing.T) {
				// ARRANGE
				// Given a private key
				_, privateKeyBytes, err := bech32.DecodeAndConvert(tt.privKey)
				require.NoError(t, err)
				require.Equal(t, byte(flagSecp256k1), privateKeyBytes[0])

				// Given signer (pk imported w/o flag byte)
				signer := NewSignerSecp256k1(privateKeyBytes[1:])

				// ACT
				// Check signer's Sui address
				address := signer.Address()

				// Sign some stub tx
				// We don't have a good way outside e2e to verify the signature is correct,
				// but let's exercise it anyway
				const exampleBase64 = "ZXhhbXBsZQo="
				_, errSign := signer.SignTxBlock(models.TxnMetaData{
					TxBytes: exampleBase64,
				})

				// ASSERT
				require.Equal(t, tt.address, address)
				require.NoError(t, errSign)
			})
		}
	})
}
