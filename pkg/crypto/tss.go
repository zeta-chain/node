package crypto

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	zetasui "github.com/zeta-chain/node/pkg/contracts/sui"
	"github.com/zeta-chain/node/pkg/cosmos"
)

// GetTSSAddrEVM returns the ethereum address of the tss pubkey
func GetTSSAddrEVM(tssPubkey string) (ethcommon.Address, error) {
	var keyAddr ethcommon.Address
	pubk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, tssPubkey)
	if err != nil {
		return keyAddr, err
	}
	decompressPubKey, err := crypto.DecompressPubkey(pubk.Bytes())
	if err != nil {
		return keyAddr, err
	}

	keyAddr = crypto.PubkeyToAddress(*decompressPubKey)

	return keyAddr, nil
}

// GetTSSAddrBTC returns the bitcoin address of the tss pubkey
func GetTSSAddrBTC(tssPubkey string, bitcoinParams *chaincfg.Params) (string, error) {
	addrWPKH, err := getKeyAddrBTCWitnessPubkeyHash(tssPubkey, bitcoinParams)
	if err != nil {
		return "", err
	}

	return addrWPKH.EncodeAddress(), nil
}

// GetTSSAddrSui returns the sui address of the tss pubkey
func GetTSSAddrSui(tssPubkey string) (string, error) {
	pubk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, tssPubkey)
	if err != nil {
		return "", err
	}
	decompressPubKey, err := crypto.DecompressPubkey(pubk.Bytes())
	if err != nil {
		return "", err
	}

	return zetasui.AddressFromPubKeyECDSA(decompressPubKey), nil
}

func getKeyAddrBTCWitnessPubkeyHash(
	tssPubkey string,
	bitcoinParams *chaincfg.Params,
) (*btcutil.AddressWitnessPubKeyHash, error) {
	pubk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, tssPubkey)
	if err != nil {
		return nil, err
	}
	addr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubk.Bytes()), bitcoinParams)
	if err != nil {
		return nil, err
	}
	return addr, nil
}
