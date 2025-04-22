package sample

import (
	"crypto/ecdsa"
	cryptoed25519 "crypto/ed25519"
	"math/big"
	"math/rand"
	"strconv"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/coming-chat/go-sui/v2/account"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/cosmos"
	"github.com/zeta-chain/node/pkg/crypto"
)

func PubKeySet() *crypto.PubKeySet {
	pubKeySet := crypto.PubKeySet{
		Secp256k1: crypto.PubKey(secp256k1.GenPrivKey().PubKey().Bytes()),
		Ed25519:   crypto.PubKey(ed25519.GenPrivKey().PubKey().String()),
	}
	return &pubKeySet
}

func Ed25519PrivateKeyFromRand(r *rand.Rand) (*ed25519.PrivKey, error) {
	randomBytes := make([]byte, 32)
	_, err := r.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	return ed25519.GenPrivKeyFromSecret(randomBytes), nil
}

// PubKeyString returns a sample public key string
func PubKeyString() string {
	priKey := ed25519.GenPrivKey()
	s, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, priKey.PubKey())
	if err != nil {
		panic(err)
	}
	pubkey, err := crypto.NewPubKey(s)
	if err != nil {
		panic(err)
	}
	return pubkey.String()
}

func PubkeyStringFromRand(r *rand.Rand) (string, error) {
	priKey, err := Ed25519PrivateKeyFromRand(r)
	if err != nil {
		return "", err
	}
	s, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, priKey.PubKey())
	if err != nil {
		return "", err
	}
	pubkey, err := crypto.NewPubKey(s)
	if err != nil {
		return "", err
	}
	return pubkey.String(), nil
}

// PrivKeyAddressPair returns a private key, address pair
func PrivKeyAddressPair() (*ed25519.PrivKey, sdk.AccAddress) {
	privKey := ed25519.GenPrivKey()
	addr := privKey.PubKey().Address()

	return privKey, sdk.AccAddress(addr)
}

// EthAddress returns a sample ethereum address
func EthAddress() ethcommon.Address {
	return ethcommon.BytesToAddress(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).Bytes())
}

func EthAddressFromRand(r *rand.Rand) ethcommon.Address {
	return ethcommon.BytesToAddress(sdk.AccAddress(PubKey(r).Address()).Bytes())
}

// BTCAddressP2WPKH returns a sample Bitcoin Pay-to-Witness-Public-Key-Hash (P2WPKH) address
func BTCAddressP2WPKH(t *testing.T, r *rand.Rand, net *chaincfg.Params) *btcutil.AddressWitnessPubKeyHash {
	privateKey, err := secp.GeneratePrivateKeyFromRand(r)
	require.NoError(t, err)

	pubKeyHash := btcutil.Hash160(privateKey.PubKey().SerializeCompressed())
	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, net)
	require.NoError(t, err)

	return addr
}

// BTCAddressP2WPKHScript returns a pkscript for a sample btc P2WPKH address
func BTCAddressP2WPKHScript(t *testing.T, r *rand.Rand, net *chaincfg.Params) []byte {
	addr := BTCAddressP2WPKH(t, r, net)
	script, err := txscript.PayToAddrScript(addr)
	require.NoError(t, err)
	return script
}

// SolanaPrivateKey returns a sample solana private key
func SolanaPrivateKey(t *testing.T) solana.PrivateKey {
	privKey, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	return privKey
}

func SolanaPrivateKeyFromRand(r *rand.Rand) (solana.PrivateKey, error) {
	pub, priv, err := cryptoed25519.GenerateKey(r)
	if err != nil {
		return nil, err
	}
	var publicKey cryptoed25519.PublicKey
	copy(publicKey[:], pub)
	return solana.PrivateKey(priv), nil
}

// SolanaAddress returns a sample solana address
func SolanaAddress(t *testing.T) string {
	privKey, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	return privKey.PublicKey().String()
}

func SolanaAddressFromRand(r *rand.Rand) (string, error) {
	privKey, err := SolanaPrivateKeyFromRand(r)
	if err != nil {
		return "", err
	}
	return privKey.PublicKey().String(), nil
}

// SolanaSignature returns a sample solana signature
func SolanaSignature(t *testing.T) solana.Signature {
	// Generate a random keypair
	keypair, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)

	// Generate a random message to sign
	// #nosec G404 test purpose - weak randomness is not an issue here
	r := rand.New(rand.NewSource(900))
	message := StringRandom(r, 64)

	// Sign the message with the private key
	signature, err := keypair.Sign([]byte(message))
	require.NoError(t, err)

	return signature
}

// SuiAddress returns a sample sui address
func SuiAddress(t require.TestingT) string {
	privateKey := ed25519.GenPrivKey()

	// create a new account with ed25519 scheme
	scheme, err := sui_types.NewSignatureScheme(0)
	require.NoError(t, err)
	acc := account.NewAccount(scheme, privateKey.GetKey().Seed())

	return acc.Address
}

// SuiDigest returns a sample sui digest
func SuiDigest(t *testing.T) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	require.NoError(t, err)

	return base58.Encode(randomBytes)
}

// Hash returns a sample hash
func Hash() ethcommon.Hash {
	return ethcommon.BytesToHash(EthAddress().Bytes())
}

// Hash returns a sample hash
func HashFromRand(r *rand.Rand) ethcommon.Hash {
	return ethcommon.BytesToHash(EthAddressFromRand(r).Bytes())
}

// BtcHash returns a sample btc hash
func BtcHash() chainhash.Hash {
	return chainhash.Hash(Hash())
}

// PubKey returns a sample account PubKey
func PubKey(r *rand.Rand) cryptotypes.PubKey {
	seed := []byte(strconv.Itoa(r.Int()))
	return ed25519.GenPrivKeyFromSecret(seed).PubKey()
}

// Bech32AccAddress returns a sample account address
func Bech32AccAddress() sdk.AccAddress {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr)
}

// AccAddress returns a sample account address in string
func AccAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr).String()
}

// AccAddressFromRand returns a sample account address in string
func AccAddressFromRand(r *rand.Rand) string {
	pk := PubKey(r)
	addr := pk.Address()
	return sdk.AccAddress(addr).String()
}

// ValAddress returns a sample validator operator address
func ValAddress(r *rand.Rand) sdk.ValAddress {
	return sdk.ValAddress(PubKey(r).Address())
}

// EthTx returns a sample ethereum transaction with the associated tx data bytes
func EthTx(t *testing.T, chainID int64, to ethcommon.Address, nonce uint64) (*ethtypes.Transaction, []byte) {
	tx := ethtypes.NewTx(&ethtypes.DynamicFeeTx{
		ChainID:   big.NewInt(chainID),
		Nonce:     nonce,
		GasTipCap: nil,
		GasFeeCap: nil,
		Gas:       21000,
		To:        &to,
		Value:     big.NewInt(5),
		Data:      nil,
	})

	txBytes, err := tx.MarshalBinary()
	require.NoError(t, err)

	return tx, txBytes
}

// EthTxSigned returns a sample signed ethereum transaction with the address of the sender
func EthTxSigned(
	t *testing.T,
	chainID int64,
	to ethcommon.Address,
	nonce uint64,
) (*ethtypes.Transaction, []byte, ethcommon.Address) {
	tx, _ := EthTx(t, chainID, to, nonce)

	// generate a private key and get address
	privateKey, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	fromAddress := ethcrypto.PubkeyToAddress(*publicKeyECDSA)
	require.True(t, ok)

	// sign the transaction
	signer := ethtypes.NewLondonSigner(tx.ChainId())
	signedTx, err := ethtypes.SignTx(tx, signer, privateKey)
	require.NoError(t, err)

	txBytes, err := signedTx.MarshalBinary()
	require.NoError(t, err)

	return signedTx, txBytes, fromAddress
}
