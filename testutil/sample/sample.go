package sample

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccAddress returns a sample account address
func AccAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr).String()
}

// PrivKeyAddressPair returns a private key, address pair
func PrivKeyAddressPair() (*ed25519.PrivKey, sdk.AccAddress) {
	privKey := ed25519.GenPrivKey()
	addr := privKey.PubKey().Address()

	return privKey, sdk.AccAddress(addr)
}
