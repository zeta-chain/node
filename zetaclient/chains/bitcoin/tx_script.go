package bitcoin

// #nosec G507 ripemd160 required for bitcoin address encoding
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
	"github.com/pkg/errors"
	"golang.org/x/crypto/ripemd160"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/constant"
)

const (
	// LengthScriptP2TR is the lenth of P2TR script [OP_1 0x20 <32-byte-hash>]
	LengthScriptP2TR = 34

	// LengthScriptP2WSH is the length of P2WSH script [OP_0 0x20 <32-byte-hash>]
	LengthScriptP2WSH = 34

	// LengthScriptP2WPKH is the length of P2WPKH script [OP_0 0x14 <20-byte-hash>]
	LengthScriptP2WPKH = 22

	// LengthScriptP2SH is the length of P2SH script [OP_HASH160 0x14 <20-byte-hash> OP_EQUAL]
	LengthScriptP2SH = 23

	// LengthScriptP2PKH is the length of P2PKH script [OP_DUP OP_HASH160 0x14 <20-byte-hash> OP_EQUALVERIFY OP_CHECKSIG]
	LengthScriptP2PKH = 25
)

// PayToAddrScript creates a new script to pay a transaction output to a the
// specified address.
func PayToAddrScript(addr btcutil.Address) ([]byte, error) {
	switch addr := addr.(type) {
	case *chains.AddressTaproot:
		return chains.PayToWitnessTaprootScript(addr.ScriptAddress())
	default:
		return txscript.PayToAddrScript(addr)
	}
}

// IsPkScriptP2TR checks if the given script is a P2TR script
func IsPkScriptP2TR(script []byte) bool {
	return len(script) == LengthScriptP2TR && script[0] == txscript.OP_1 && script[1] == 0x20
}

// IsPkScriptP2WSH checks if the given script is a P2WSH script
func IsPkScriptP2WSH(script []byte) bool {
	return len(script) == LengthScriptP2WSH && script[0] == txscript.OP_0 && script[1] == 0x20
}

// IsPkScriptP2WPKH checks if the given script is a P2WPKH script
func IsPkScriptP2WPKH(script []byte) bool {
	return len(script) == LengthScriptP2WPKH && script[0] == txscript.OP_0 && script[1] == 0x14
}

// IsPkScriptP2SH checks if the given script is a P2SH script
func IsPkScriptP2SH(script []byte) bool {
	return len(script) == LengthScriptP2SH &&
		script[0] == txscript.OP_HASH160 &&
		script[1] == 0x14 &&
		script[22] == txscript.OP_EQUAL
}

// IsPkScriptP2PKH checks if the given script is a P2PKH script
func IsPkScriptP2PKH(script []byte) bool {
	return len(script) == LengthScriptP2PKH &&
		script[0] == txscript.OP_DUP &&
		script[1] == txscript.OP_HASH160 &&
		script[2] == 0x14 &&
		script[23] == txscript.OP_EQUALVERIFY &&
		script[24] == txscript.OP_CHECKSIG
}

// DecodeScriptP2TR decodes address from P2TR script
func DecodeScriptP2TR(scriptHex string, net *chaincfg.Params) (string, error) {
	script, err := hex.DecodeString(scriptHex)
	if err != nil {
		return "", errors.Wrapf(err, "error decoding script %s", scriptHex)
	}
	if !IsPkScriptP2TR(script) {
		return "", fmt.Errorf("invalid P2TR script: %s", scriptHex)
	}

	witnessProg := script[2:]
	receiverAddress, err := chains.NewAddressTaproot(witnessProg, net)
	if err != nil { // should never happen
		return "", errors.Wrapf(err, "error getting address from script %s", scriptHex)
	}

	return receiverAddress.EncodeAddress(), nil
}

// DecodeScriptP2WSH decodes address from P2WSH script
func DecodeScriptP2WSH(scriptHex string, net *chaincfg.Params) (string, error) {
	script, err := hex.DecodeString(scriptHex)
	if err != nil {
		return "", errors.Wrapf(err, "error decoding script: %s", scriptHex)
	}
	if !IsPkScriptP2WSH(script) {
		return "", fmt.Errorf("invalid P2WSH script: %s", scriptHex)
	}

	witnessProg := script[2:]
	receiverAddress, err := btcutil.NewAddressWitnessScriptHash(witnessProg, net)
	if err != nil { // should never happen
		return "", errors.Wrapf(err, "error getting receiver from script: %s", scriptHex)
	}

	return receiverAddress.EncodeAddress(), nil
}

// DecodeScriptP2WPKH decodes address from P2WPKH script
func DecodeScriptP2WPKH(scriptHex string, net *chaincfg.Params) (string, error) {
	script, err := hex.DecodeString(scriptHex)
	if err != nil {
		return "", errors.Wrapf(err, "error decoding script: %s", scriptHex)
	}
	if !IsPkScriptP2WPKH(script) {
		return "", fmt.Errorf("invalid P2WPKH script: %s", scriptHex)
	}

	witnessProg := script[2:]
	receiverAddress, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, net)
	if err != nil { // should never happen
		return "", errors.Wrapf(err, "error getting receiver from script: %s", scriptHex)
	}

	return receiverAddress.EncodeAddress(), nil
}

// DecodeScriptP2SH decodes address from P2SH script
func DecodeScriptP2SH(scriptHex string, net *chaincfg.Params) (string, error) {
	script, err := hex.DecodeString(scriptHex)
	if err != nil {
		return "", errors.Wrapf(err, "error decoding script: %s", scriptHex)
	}
	if !IsPkScriptP2SH(script) {
		return "", fmt.Errorf("invalid P2SH script: %s", scriptHex)
	}

	scriptHash := script[2:22]

	return EncodeAddress(scriptHash, net.ScriptHashAddrID), nil
}

// DecodeScriptP2PKH decodes address from P2PKH script
func DecodeScriptP2PKH(scriptHex string, net *chaincfg.Params) (string, error) {
	script, err := hex.DecodeString(scriptHex)
	if err != nil {
		return "", errors.Wrapf(err, "error decoding script: %s", scriptHex)
	}
	if !IsPkScriptP2PKH(script) {
		return "", fmt.Errorf("invalid P2PKH script: %s", scriptHex)
	}

	pubKeyHash := script[3:23]

	return EncodeAddress(pubKeyHash, net.PubKeyHashAddrID), nil
}

// DecodeOpReturnMemo decodes memo from OP_RETURN script
// returns (memo, found, error)
func DecodeOpReturnMemo(scriptHex string, txid string) ([]byte, bool, error) {
	if len(scriptHex) >= 4 && scriptHex[:2] == "6a" { // OP_RETURN
		memoSize, err := strconv.ParseInt(scriptHex[2:4], 16, 32)
		if err != nil {
			return nil, false, errors.Wrapf(err, "error decoding memo size: %s", scriptHex)
		}
		if int(memoSize) != (len(scriptHex)-4)/2 {
			return nil, false, fmt.Errorf("memo size mismatch: %d != %d", memoSize, (len(scriptHex)-4)/2)
		}

		memoBytes, err := hex.DecodeString(scriptHex[4:])
		if err != nil {
			return nil, false, errors.Wrapf(err, "error hex decoding memo: %s", scriptHex)
		}
		if bytes.Equal(memoBytes, []byte(constant.DonationMessage)) {
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
func DecodeTSSVout(vout btcjson.Vout, receiverExpected string, chain chains.Chain) (string, int64, error) {
	// parse amount
	amount, err := GetSatoshis(vout.Value)
	if err != nil {
		return "", 0, errors.Wrap(err, "error getting satoshis")
	}

	// get btc chain params
	chainParams, err := chains.GetBTCChainParams(chain.ChainId)
	if err != nil {
		return "", 0, errors.Wrapf(err, "error GetBTCChainParams for chain %d", chain.ChainId)
	}

	// decode cctx receiver address
	addr, err := chains.DecodeBtcAddress(receiverExpected, chain.ChainId)
	if err != nil {
		return "", 0, errors.Wrapf(err, "error decoding receiver %s", receiverExpected)
	}

	// parse receiver address from vout
	var receiverVout string
	switch addr.(type) {
	case *chains.AddressTaproot:
		receiverVout, err = DecodeScriptP2TR(vout.ScriptPubKey.Hex, chainParams)
	case *btcutil.AddressWitnessScriptHash:
		receiverVout, err = DecodeScriptP2WSH(vout.ScriptPubKey.Hex, chainParams)
	case *btcutil.AddressWitnessPubKeyHash:
		receiverVout, err = DecodeScriptP2WPKH(vout.ScriptPubKey.Hex, chainParams)
	case *btcutil.AddressScriptHash:
		receiverVout, err = DecodeScriptP2SH(vout.ScriptPubKey.Hex, chainParams)
	case *btcutil.AddressPubKeyHash:
		receiverVout, err = DecodeScriptP2PKH(vout.ScriptPubKey.Hex, chainParams)
	default:
		return "", 0, fmt.Errorf("unsupported receiver address type: %T", addr)
	}
	if err != nil {
		return "", 0, errors.Wrap(err, "error decoding TSS vout")
	}

	return receiverVout, amount, nil
}
