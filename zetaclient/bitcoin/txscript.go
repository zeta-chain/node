package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/cosmos/btcutil/base58"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/bitcoin"
	"golang.org/x/crypto/ripemd160"

	"github.com/pkg/errors"
)

const (
	// P2TR script type
	ScriptTypeP2TR = "witness_v1_taproot"

	// P2WSH script type
	ScriptTypeP2WSH = "witness_v0_scripthash"

	// P2WPKH script type
	ScriptTypeP2WPKH = "witness_v0_keyhash"

	// P2SH script type
	ScriptTypeP2SH = "scripthash"

	// P2PKH script type
	ScriptTypeP2PKH = "pubkeyhash"
)

// PayToAddrScript creates a new script to pay a transaction output to a the
// specified address.
func PayToAddrScript(addr btcutil.Address) ([]byte, error) {
	switch addr := addr.(type) {
	case *bitcoin.AddressTaproot:
		return bitcoin.PayToWitnessTaprootScript(addr.ScriptAddress())
	default:
		return txscript.PayToAddrScript(addr)
	}
}

// DecodeVoutP2TR decodes receiver and amount from P2TR output
func DecodeVoutP2TR(vout btcjson.Vout, net *chaincfg.Params) (string, error) {
	// check tx script type
	if vout.ScriptPubKey.Type != ScriptTypeP2TR {
		return "", fmt.Errorf("want scriptPubKey type witness_v1_taproot, got %s", vout.ScriptPubKey.Type)
	}
	// decode P2TR scriptPubKey [OP_1 0x20 <32-byte-hash>]
	scriptPubKey := vout.ScriptPubKey.Hex
	decodedScriptPubKey, err := hex.DecodeString(scriptPubKey)
	if err != nil {
		return "", errors.Wrapf(err, "error decoding scriptPubKey %s", scriptPubKey)
	}
	if len(decodedScriptPubKey) != 34 ||
		decodedScriptPubKey[0] != txscript.OP_1 ||
		decodedScriptPubKey[1] != 0x20 {
		return "", fmt.Errorf("invalid P2TR scriptPubKey: %s", scriptPubKey)
	}
	witnessProg := decodedScriptPubKey[2:]
	receiverAddress, err := bitcoin.NewAddressTaproot(witnessProg, net)
	if err != nil { // should never happen
		return "", errors.Wrapf(err, "error getting receiver from scriptPubKey %s", scriptPubKey)
	}
	return receiverAddress.EncodeAddress(), nil
}

// DecodeVoutP2WSH decodes receiver and amount from P2WSH output
func DecodeVoutP2WSH(vout btcjson.Vout, net *chaincfg.Params) (string, error) {
	// check tx script type
	if vout.ScriptPubKey.Type != ScriptTypeP2WSH {
		return "", fmt.Errorf("want scriptPubKey type witness_v0_scripthash, got %s", vout.ScriptPubKey.Type)
	}
	// decode P2WSH scriptPubKey [OP_0 0x20 <32-byte-hash>]
	scriptPubKey := vout.ScriptPubKey.Hex
	decodedScriptPubKey, err := hex.DecodeString(scriptPubKey)
	if err != nil {
		return "", errors.Wrapf(err, "error decoding scriptPubKey %s", scriptPubKey)
	}
	if len(decodedScriptPubKey) != 34 ||
		decodedScriptPubKey[0] != txscript.OP_0 ||
		decodedScriptPubKey[1] != 0x20 {
		return "", fmt.Errorf("invalid P2WSH scriptPubKey: %s", scriptPubKey)
	}
	witnessProg := decodedScriptPubKey[2:]
	receiverAddress, err := btcutil.NewAddressWitnessScriptHash(witnessProg, net)
	if err != nil { // should never happen
		return "", errors.Wrapf(err, "error getting receiver from scriptPubKey %s", scriptPubKey)
	}
	return receiverAddress.EncodeAddress(), nil
}

// DecodeVoutP2WPKH decodes receiver and amount from P2WPKH output
func DecodeVoutP2WPKH(vout btcjson.Vout, net *chaincfg.Params) (string, error) {
	// check tx script type
	if vout.ScriptPubKey.Type != ScriptTypeP2WPKH {
		return "", fmt.Errorf("want scriptPubKey type witness_v0_keyhash, got %s", vout.ScriptPubKey.Type)
	}
	// decode P2WPKH scriptPubKey [OP_0 0x14 <20-byte-hash>]
	scriptPubKey := vout.ScriptPubKey.Hex
	decodedScriptPubKey, err := hex.DecodeString(scriptPubKey)
	if err != nil {
		return "", errors.Wrapf(err, "error decoding scriptPubKey %s", scriptPubKey)
	}
	if len(decodedScriptPubKey) != 22 ||
		decodedScriptPubKey[0] != txscript.OP_0 ||
		decodedScriptPubKey[1] != 0x14 {
		return "", fmt.Errorf("invalid P2WPKH scriptPubKey: %s", scriptPubKey)
	}
	witnessProg := decodedScriptPubKey[2:]
	receiverAddress, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, net)
	if err != nil { // should never happen
		return "", errors.Wrapf(err, "error getting receiver from scriptPubKey %s", scriptPubKey)
	}
	return receiverAddress.EncodeAddress(), nil
}

// DecodeVoutP2SH decodes receiver address from P2SH output
func DecodeVoutP2SH(vout btcjson.Vout, net *chaincfg.Params) (string, error) {
	// check tx script type
	if vout.ScriptPubKey.Type != ScriptTypeP2SH {
		return "", fmt.Errorf("want scriptPubKey type scripthash, got %s", vout.ScriptPubKey.Type)
	}
	// decode P2SH scriptPubKey [OP_HASH160 0x14 <20-byte-hash> OP_EQUAL]
	scriptPubKey := vout.ScriptPubKey.Hex
	decodedScriptPubKey, err := hex.DecodeString(scriptPubKey)
	if err != nil {
		return "", errors.Wrapf(err, "error decoding scriptPubKey %s", scriptPubKey)
	}
	if len(decodedScriptPubKey) != 23 ||
		decodedScriptPubKey[0] != txscript.OP_HASH160 ||
		decodedScriptPubKey[1] != 0x14 ||
		decodedScriptPubKey[22] != txscript.OP_EQUAL {
		return "", fmt.Errorf("invalid P2SH scriptPubKey: %s", scriptPubKey)
	}
	scriptHash := decodedScriptPubKey[2:22]
	return EncodeAddress(scriptHash, net.ScriptHashAddrID), nil
}

// DecodeVoutP2PKH decodes receiver address from P2PKH output
func DecodeVoutP2PKH(vout btcjson.Vout, net *chaincfg.Params) (string, error) {
	// check tx script type
	if vout.ScriptPubKey.Type != ScriptTypeP2PKH {
		return "", fmt.Errorf("want scriptPubKey type pubkeyhash, got %s", vout.ScriptPubKey.Type)
	}
	// decode P2PKH scriptPubKey [OP_DUP OP_HASH160 0x14 <20-byte-hash> OP_EQUALVERIFY OP_CHECKSIG]
	scriptPubKey := vout.ScriptPubKey.Hex
	decodedScriptPubKey, err := hex.DecodeString(scriptPubKey)
	if err != nil {
		return "", errors.Wrapf(err, "error decoding scriptPubKey %s", scriptPubKey)
	}
	if len(decodedScriptPubKey) != 25 ||
		decodedScriptPubKey[0] != txscript.OP_DUP ||
		decodedScriptPubKey[1] != txscript.OP_HASH160 ||
		decodedScriptPubKey[2] != 0x14 ||
		decodedScriptPubKey[23] != txscript.OP_EQUALVERIFY ||
		decodedScriptPubKey[24] != txscript.OP_CHECKSIG {
		return "", fmt.Errorf("invalid P2PKH scriptPubKey: %s", scriptPubKey)
	}
	pubKeyHash := decodedScriptPubKey[3:23]
	return EncodeAddress(pubKeyHash, net.PubKeyHashAddrID), nil
}

// DecodeVoutMemoP2WPKH decodes memo from P2WPKH output
// returns (memo, found, error)
func DecodeVoutMemoP2WPKH(vout btcjson.Vout, txid string) ([]byte, bool, error) {
	script := vout.ScriptPubKey.Hex
	if len(script) >= 4 && script[:2] == "6a" { // OP_RETURN
		memoSize, err := strconv.ParseInt(script[2:4], 16, 32)
		if err != nil {
			return nil, false, errors.Wrapf(err, "error decoding memo size: %s", script)
		}
		if int(memoSize) != (len(script)-4)/2 {
			return nil, false, fmt.Errorf("memo size mismatch: %d != %d", memoSize, (len(script)-4)/2)
		}
		memoBytes, err := hex.DecodeString(script[4:])
		if err != nil {
			return nil, false, errors.Wrapf(err, "error hex decoding memo: %s", script)
		}
		if bytes.Equal(memoBytes, []byte(common.DonationMessage)) {
			return nil, false, fmt.Errorf("donation tx: %s", txid)
		}
		return memoBytes, true, nil
	}
	return nil, false, nil
}

// EncodeAddress returns a human-readable payment address given a ripemd160 hash
// and netID which encodes the bitcoin network and address type.  It is used
// in both pay-to-pubkey-hash (P2PKH) and pay-to-script-hash (P2SH) address
// encoding.
// Note: this function is a copy of the function in btcutil/address.go
func EncodeAddress(hash160 []byte, netID byte) string {
	// Format is 1 byte for a network and address class (i.e. P2PKH vs
	// P2SH), 20 bytes for a RIPEMD160 hash, and 4 bytes of checksum.
	return base58.CheckEncode(hash160[:ripemd160.Size], netID)
}

// DecodeTSSVout decodes receiver and amount from a given TSS vout
func DecodeTSSVout(vout btcjson.Vout, receiverExpected string, chain common.Chain) (string, int64, error) {
	// parse amount
	amount, err := GetSatoshis(vout.Value)
	if err != nil {
		return "", 0, errors.Wrap(err, "error getting satoshis")
	}
	// get btc chain params
	chainParams, err := common.GetBTCChainParams(chain.ChainId)
	if err != nil {
		return "", 0, errors.Wrapf(err, "error GetBTCChainParams for chain %d", chain.ChainId)
	}
	// decode cctx receiver address
	addr, err := common.DecodeBtcAddress(receiverExpected, chain.ChainId)
	if err != nil {
		return "", 0, errors.Wrapf(err, "error decoding receiver %s", receiverExpected)
	}
	// parse receiver address from vout
	var receiverVout string
	switch addr.(type) {
	case *bitcoin.AddressTaproot:
		receiverVout, err = DecodeVoutP2TR(vout, chainParams)
	case *btcutil.AddressWitnessScriptHash:
		receiverVout, err = DecodeVoutP2WSH(vout, chainParams)
	case *btcutil.AddressWitnessPubKeyHash:
		receiverVout, err = DecodeVoutP2WPKH(vout, chainParams)
	case *btcutil.AddressScriptHash:
		receiverVout, err = DecodeVoutP2SH(vout, chainParams)
	case *btcutil.AddressPubKeyHash:
		receiverVout, err = DecodeVoutP2PKH(vout, chainParams)
	default:
		return "", 0, fmt.Errorf("unsupported receiver address type: %T", addr)
	}
	if err != nil {
		return "", 0, errors.Wrap(err, "error decoding TSS vout")
	}
	return receiverVout, amount, nil
}
