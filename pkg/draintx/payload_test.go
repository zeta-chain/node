package draintx

import (
	"testing"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func samplePayload() Payload {
	return Payload{
		TriggerZetaHeight: 1_234_567,
		Seq:               12,
		Final:             true,
		Network:           "mainnet",
		EVMTxs: []EVMTx{
			{ChainID: 1, To: "0x1111111111111111111111111111111111111111", Nonce: 42, Amount: "1000000000000000000", GasPrice: "250000", GasLimit: 21000},
		},
		BTCTxs: []BTCTx{
			{
				ChainID:    8332,
				To:         "bc1qsafe",
				OutputSats: 99_000_000,
				FeeSats:    1_000_000,
				Inputs: []BTCInput{
					{TxID: "aa", Vout: 0, AmountSats: 50_000_000},
					{TxID: "bb", Vout: 1, AmountSats: 50_000_000},
				},
			},
		},
	}
}

func TestSignVerifyRoundTrip(t *testing.T) {
	// ARRANGE
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	pubCompressed := ethcrypto.CompressPubkey(&priv.PublicKey)
	pubUncompressed := ethcrypto.FromECDSAPub(&priv.PublicKey)
	p := samplePayload()

	// ACT
	require.NoError(t, p.Sign(priv))

	// ASSERT: verifies against both compressed and uncompressed forms
	require.NoError(t, p.Verify(pubCompressed))
	require.NoError(t, p.Verify(pubUncompressed))
}

func TestVerifyRejectsWrongKey(t *testing.T) {
	// ARRANGE
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	other, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	p := samplePayload()
	require.NoError(t, p.Sign(priv))

	// ACT
	err = p.Verify(ethcrypto.CompressPubkey(&other.PublicKey))

	// ASSERT
	require.Error(t, err)
}

func TestVerifyRejectsMutation(t *testing.T) {
	// ARRANGE
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	pub := ethcrypto.CompressPubkey(&priv.PublicKey)

	t.Run("mutated amount", func(t *testing.T) {
		p := samplePayload()
		require.NoError(t, p.Sign(priv))
		p.EVMTxs[0].Amount = "999"
		require.Error(t, p.Verify(pub))
	})

	t.Run("mutated btc input", func(t *testing.T) {
		p := samplePayload()
		require.NoError(t, p.Sign(priv))
		p.BTCTxs[0].Inputs[0].AmountSats = 1
		require.Error(t, p.Verify(pub))
	})

	t.Run("mutated final flag", func(t *testing.T) {
		p := samplePayload()
		require.NoError(t, p.Sign(priv))
		p.Final = false
		require.Error(t, p.Verify(pub))
	})

	t.Run("mutated network", func(t *testing.T) {
		p := samplePayload()
		require.NoError(t, p.Sign(priv))
		p.Network = "testnet"
		require.Error(t, p.Verify(pub))
	})
}

func TestCanonicalBytesStable(t *testing.T) {
	// ARRANGE
	p := samplePayload()

	// ACT
	first := p.canonicalBytes()
	second := p.canonicalBytes()

	// ASSERT
	require.Equal(t, first, second)
	// signature is excluded from canonical bytes
	p.Signature = "0xdeadbeef"
	require.Equal(t, first, p.canonicalBytes())
}

func TestCanonicalBytesFieldOrderMatters(t *testing.T) {
	// ARRANGE: two payloads differing only by swapped input order must differ
	a := samplePayload()
	b := samplePayload()
	b.BTCTxs[0].Inputs[0], b.BTCTxs[0].Inputs[1] = b.BTCTxs[0].Inputs[1], b.BTCTxs[0].Inputs[0]

	// ASSERT
	require.NotEqual(t, a.canonicalBytes(), b.canonicalBytes())
}

func TestCanonicalBytesNetworkMatters(t *testing.T) {
	// ARRANGE: two payloads differing only by network must produce different canonical bytes
	a := samplePayload()
	b := samplePayload()
	b.Network = "testnet"

	// ASSERT
	require.NotEqual(t, a.canonicalBytes(), b.canonicalBytes())
}
