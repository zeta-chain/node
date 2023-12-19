package common

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	zcommon "github.com/zeta-chain/zetacore/common/cosmos"
)

// GetTssAddrEVM returns the ethereum address of the tss pubkey
func GetTssAddrEVM(tssPubkey string) (ethcommon.Address, error) {
	var keyAddr ethcommon.Address
	pubk, err := zcommon.GetPubKeyFromBech32(zcommon.Bech32PubKeyTypeAccPub, tssPubkey)
	if err != nil {
		return keyAddr, err
	}
	//keyAddrBytes := pubk.EVMAddress().Bytes()
	pubk.Bytes()
	decompresspubkey, err := crypto.DecompressPubkey(pubk.Bytes())
	if err != nil {
		return keyAddr, err
	}

	keyAddr = crypto.PubkeyToAddress(*decompresspubkey)

	return keyAddr, nil
}

// GetTssAddrBTC returns the bitcoin address of the tss pubkey
func GetTssAddrBTC(tssPubkey string, bitcoinParams *chaincfg.Params) (string, error) {
	addrWPKH, err := getKeyAddrBTCWitnessPubkeyHash(tssPubkey, bitcoinParams)
	if err != nil {
		return "", err
	}

	return addrWPKH.EncodeAddress(), nil
}

func getKeyAddrBTCWitnessPubkeyHash(tssPubkey string, bitcoinParams *chaincfg.Params) (*btcutil.AddressWitnessPubKeyHash, error) {
	pubk, err := zcommon.GetPubKeyFromBech32(zcommon.Bech32PubKeyTypeAccPub, tssPubkey)
	if err != nil {
		return nil, err
	}
	addr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubk.Bytes()), bitcoinParams)
	if err != nil {
		return nil, err
	}
	return addr, nil
}
