package common

// #nosec G507 ripemd160 required for bitcoin address encoding
import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/cosmos/btcutil/base58"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ripemd160"

	"github.com/zeta-chain/node/pkg/chains"
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
	receiverAddress, err := btcutil.NewAddressTaproot(witnessProg, net)
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
func DecodeOpReturnMemo(scriptHex string) ([]byte, bool, error) {
	// decode hex script
	scriptBytes, err := hex.DecodeString(scriptHex)
	if err != nil {
		return nil, false, errors.Wrapf(err, "error decoding script hex: %s", scriptHex)
	}

	// skip non-OP_RETURN script
	// OP_RETURN script has to be at least 2 bytes: [OP_RETURN + dataLen]
	if len(scriptBytes) < 2 || scriptBytes[0] != txscript.OP_RETURN {
		return nil, false, nil
	}

	// extract appended data in the OP_RETURN script
	var memoBytes []byte
	var memoSize = scriptBytes[1]
	switch {
	case memoSize < txscript.OP_PUSHDATA1:
		// memo size has to match the actual data
		if int(memoSize) != (len(scriptBytes) - 2) {
			return nil, false, fmt.Errorf("memo size mismatch: %d != %d", memoSize, (len(scriptBytes) - 2))
		}
		memoBytes = scriptBytes[2:]
	case memoSize == txscript.OP_PUSHDATA1:
		// when data size >= OP_PUSHDATA1 (76), Bitcoin uses 2 bytes to represent the length: [OP_PUSHDATA1 + dataLen]
		// see: https://github.com/btcsuite/btcd/blob/master/txscript/scriptbuilder.go#L183
		if len(scriptBytes) < 3 {
			return nil, false, fmt.Errorf("script too short: %s", scriptHex)
		}
		memoSize = scriptBytes[2]

		// memo size has to match the actual data
		if int(memoSize) != (len(scriptBytes) - 3) {
			return nil, false, fmt.Errorf("memo size mismatch: %d != %d", memoSize, (len(scriptBytes) - 3))
		}
		memoBytes = scriptBytes[3:]
	default:
		// should never happen
		// OP_RETURN script won't carry more than 80 bytes
		return nil, false, fmt.Errorf("invalid OP_RETURN script: %s", scriptHex)
	}

	return memoBytes, true, nil
}

// DecodeScript decodes memo wrapped in an inscription like script in witness
// returns (memo, found, error)
//
// Note: the format of the script is following that of "inscription" defined in ordinal theory.
// However, to separate from inscription (as this use case is not an NFT), simplifications are made.
// The bitcoin envelope script is as follows:
// OP_DATA_32 <32 byte of public key> OP_CHECKSIG
// OP_FALSE
// OP_IF
//
//	OP_PUSH 0x...
//	OP_PUSH 0x...
//
// OP_ENDIF
// There are no content-type or any other attributes, it's just raw bytes.
func DecodeScript(script []byte) ([]byte, bool, error) {
	t := txscript.MakeScriptTokenizer(0, script)

	if err := checkInscriptionEnvelope(&t); err != nil {
		return nil, false, errors.Wrap(err, "checkInscriptionEnvelope: unable to check the envelope")
	}

	memoBytes, err := decodeInscriptionPayload(&t)
	if err != nil {
		return nil, false, errors.Wrap(err, "decodeInscriptionPayload: unable to decode the payload")
	}

	return memoBytes, true, nil
}

// EncodeAddress returns a human-readable payment address given a ripemd160 hash
// and netID which encodes the bitcoin network and address type. It is used
// in both pay-to-pubkey-hash (P2PKH) and pay-to-script-hash (P2SH) address
// encoding.
// Note: this function is a copy of the function in btcutil/address.go
func EncodeAddress(hash160 []byte, netID byte) string {
	// Format is 1 byte for a network and address class (i.e. P2PKH vs
	// P2SH), 20 bytes for a RIPEMD160 hash, and 4 bytes of checksum.
	return base58.CheckEncode(hash160[:ripemd160.Size], netID)
}

// DecodeSenderFromScript decodes sender from a given script
func DecodeSenderFromScript(pkScript []byte, net *chaincfg.Params) (string, error) {
	scriptHex := hex.EncodeToString(pkScript)

	// decode sender address from according to script type
	switch {
	case IsPkScriptP2TR(pkScript):
		return DecodeScriptP2TR(scriptHex, net)
	case IsPkScriptP2WSH(pkScript):
		return DecodeScriptP2WSH(scriptHex, net)
	case IsPkScriptP2WPKH(pkScript):
		return DecodeScriptP2WPKH(scriptHex, net)
	case IsPkScriptP2SH(pkScript):
		return DecodeScriptP2SH(scriptHex, net)
	case IsPkScriptP2PKH(pkScript):
		return DecodeScriptP2PKH(scriptHex, net)
	default:
		// sender address not found, return nil and move on to the next tx
		return "", nil
	}
}

// DecodeTSSVout decodes receiver and amount from a given TSS vout
func DecodeTSSVout(vout btcjson.Vout, receiverExpected btcutil.Address, chain chains.Chain) (string, int64, error) {
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

	// parse receiver address from vout
	var receiverVout string
	switch receiverExpected.(type) {
	case *btcutil.AddressTaproot:
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
		return "", 0, fmt.Errorf("unsupported receiver address type: %T", receiverExpected)
	}
	if err != nil {
		return "", 0, errors.Wrap(err, "error decoding TSS vout")
	}

	return receiverVout, amount, nil
}

func decodeInscriptionPayload(t *txscript.ScriptTokenizer) ([]byte, error) {
	if !t.Next() || t.Opcode() != txscript.OP_FALSE {
		return nil, fmt.Errorf("OP_FALSE not found")
	}

	if !t.Next() || t.Opcode() != txscript.OP_IF {
		return nil, fmt.Errorf("OP_IF not found")
	}

	memo := make([]byte, 0)
	var next byte
	for t.Next() {
		next = t.Opcode()
		if next == txscript.OP_ENDIF {
			return memo, nil
		}
		if next < txscript.OP_DATA_1 || next > txscript.OP_PUSHDATA4 {
			return nil, fmt.Errorf("expecting data push, found %d", next)
		}
		memo = append(memo, t.Data()...)
	}
	if t.Err() != nil {
		return nil, t.Err()
	}
	return nil, fmt.Errorf("should contain more data, but script ended")
}

// checkInscriptionEnvelope decodes the envelope for the script monitoring. The format is
// OP_PUSHBYTES_32 <32 bytes> OP_CHECKSIG <Content>
func checkInscriptionEnvelope(t *txscript.ScriptTokenizer) error {
	if !t.Next() || t.Opcode() != txscript.OP_DATA_32 {
		return fmt.Errorf("public key not found: %v", t.Err())
	}

	if !t.Next() || t.Opcode() != txscript.OP_CHECKSIG {
		return fmt.Errorf("OP_CHECKSIG not found: %v", t.Err())
	}

	return nil
}
