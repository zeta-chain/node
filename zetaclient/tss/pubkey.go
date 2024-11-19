package tss

import (
	"crypto/ecdsa"
	"crypto/elliptic"

	"github.com/btcsuite/btcd/btcutil"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/cosmos"
)

// PubKey represents TSS public key in various formats.
type PubKey struct {
	cosmosPubKey cryptotypes.PubKey
	ecdsaPubKey  *ecdsa.PublicKey
}

// NewPubKeyFromBech32 creates a new PubKey from a bech32 address.
// Example: `zetapub1addwnpepq2fdhcmfyv07s86djjca835l4f2n2ta0c7le6vnl508mseca2s9g6slj0gm`
func NewPubKeyFromBech32(bech32 string) (PubKey, error) {
	if bech32 == "" {
		return PubKey{}, errors.New("empty bech32 address")
	}

	cosmosPubKey, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, bech32)
	if err != nil {
		return PubKey{}, errors.Wrap(err, "unable to GetPubKeyFromBech32")
	}

	pubKey, err := crypto.DecompressPubkey(cosmosPubKey.Bytes())
	if err != nil {
		return PubKey{}, errors.Wrap(err, "unable to DecompressPubkey")
	}

	crypto.FromECDSAPub(pubKey)

	return PubKey{
		cosmosPubKey: cosmosPubKey,
		ecdsaPubKey:  pubKey,
	}, nil
}

// Bytes marshals pubKey to bytes either as compressed or uncompressed slice.
//
// In ECDSA, a compressed pubKey includes only the X and a parity bit for the Y,
// allowing the full Y to be reconstructed using the elliptic curve equation,
// thus reducing the key size while maintaining the ability to fully recover the pubKey.
func (k PubKey) Bytes(compress bool) []byte {
	pk := k.ecdsaPubKey
	if compress {
		return elliptic.MarshalCompressed(pk.Curve, pk.X, pk.Y)
	}

	return crypto.FromECDSAPub(pk)
}

// Bech32String returns the bech32 string of the public key.
// Example: `zetapub1addwnpepq2fdhcmfyv07s86djjca835l4f2n2ta0c7le6vnl508mseca2s9g6slj0gm`
func (k PubKey) Bech32String() string {
	v, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, k.cosmosPubKey)

	// should not happen as we only set k.cosmosPubKey from the constructor
	if err != nil {
		panic("PubKey.Bech32String: " + err.Error())
	}

	return v
}

// AddressBTC returns the bitcoin address of the public key.
func (k PubKey) AddressBTC(chainID int64) (*btcutil.AddressWitnessPubKeyHash, error) {
	return bitcoinP2WPKH(k.Bytes(true), chainID)
}

// AddressEVM returns the ethereum address of the public key.
func (k PubKey) AddressEVM() eth.Address {
	return crypto.PubkeyToAddress(*k.ecdsaPubKey)
}

// bitcoinP2WPKH returns P2WPKH (pay to witness pub key hash) address from the compressed pub key.
func bitcoinP2WPKH(pkCompressed []byte, chainID int64) (*btcutil.AddressWitnessPubKeyHash, error) {
	params, err := chains.BitcoinNetParamsFromChainID(chainID)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get btc net params")
	}

	hash := btcutil.Hash160(pkCompressed)

	return btcutil.NewAddressWitnessPubKeyHash(hash, params)
}
