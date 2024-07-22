package bitcoin

// #nosec G507 ripemd160 required for bitcoin address encoding
import (
	"bytes"
	"encoding/binary"
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
	t := makeScriptTokenizer(script)

	if err := checkInscriptionEnvelope(&t); err != nil {
		return nil, false, err
	}

	memoBytes, err := decodeInscriptionPayload(&t)
	if err != nil {
		return nil, false, errors.Wrap(err, "unable to decode the payload")
	}

	return memoBytes, true, nil
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

// decodeInscriptionPayload checks the envelope for the script monitoring. The format is
// OP_FALSE
// OP_IF
//
//		OP_PUSHDATA_N ...
//	 ...
//		OP_PUSHDATA_N ...
//
// OP_ENDIF
//
// Note: the total data pushed will always be more than 80 bytes and within the btc transaction size limit.
func decodeInscriptionPayload(t *scriptTokenizer) ([]byte, error) {
	if !t.Next() || t.Opcode() != txscript.OP_FALSE {
		return nil, fmt.Errorf("OP_FALSE not found")
	}

	if !t.Next() || t.Opcode() != txscript.OP_IF {
		return nil, fmt.Errorf("OP_IF not found")
	}

	memo := make([]byte, 0)
	var next byte
	for {
		if !t.Next() {
			if t.Err() != nil {
				return nil, t.Err()
			}
			return nil, fmt.Errorf("should contain more data, but script ended")
		}

		next = t.Opcode()

		if next == txscript.OP_ENDIF {
			break
		}

		if next < txscript.OP_DATA_1 || next > txscript.OP_PUSHDATA4 {
			return nil, fmt.Errorf("expecting data push, found %d", next)
		}

		memo = append(memo, t.Data()...)
	}

	return memo, nil
}

// checkInscriptionEnvelope decodes the envelope for the script monitoring. The format is
// OP_PUSHBYTES_32 <32 bytes> OP_CHECKSIG <Content>
func checkInscriptionEnvelope(t *scriptTokenizer) error {
	if !t.Next() || t.Opcode() != txscript.OP_DATA_32 {
		return fmt.Errorf("cannot obtain public key bytes")
	}

	if !t.Next() || t.Opcode() != txscript.OP_CHECKSIG {
		return fmt.Errorf("cannot parse OP_CHECKSIG")
	}

	return nil
}

func makeScriptTokenizer(script []byte) scriptTokenizer {
	return scriptTokenizer{
		script: script,
		offset: 0,
	}
}

// scriptTokenizer is supposed to be replaced by txscript.ScriptTokenizer. However,
// it seems currently the btcsuite version does not have ScriptTokenizer. A simplified
// version of that is implemented here. This is fully compatible with txscript.ScriptTokenizer
// one should consider upgrading txscript and remove this implementation
type scriptTokenizer struct {
	script []byte
	offset int32
	op     byte
	data   []byte
	err    error
}

// Done returns true when either all opcodes have been exhausted or a parse
// failure was encountered and therefore the state has an associated error.
func (t *scriptTokenizer) Done() bool {
	return t.err != nil || t.offset >= int32(len(t.script))
}

// Data returns the data associated with the most recently successfully parsed
// opcode.
func (t *scriptTokenizer) Data() []byte {
	return t.data
}

// Err returns any errors currently associated with the tokenizer.  This will
// only be non-nil in the case a parsing error was encountered.
func (t *scriptTokenizer) Err() error {
	return t.err
}

// Opcode returns the current opcode associated with the tokenizer.
func (t *scriptTokenizer) Opcode() byte {
	return t.op
}

// Next attempts to parse the next opcode and returns whether or not it was
// successful.  It will not be successful if invoked when already at the end of
// the script, a parse failure is encountered, or an associated error already
// exists due to a previous parse failure.
//
// In the case of a true return, the parsed opcode and data can be obtained with
// the associated functions and the offset into the script will either point to
// the next opcode or the end of the script if the final opcode was parsed.
//
// In the case of a false return, the parsed opcode and data will be the last
// successfully parsed values (if any) and the offset into the script will
// either point to the failing opcode or the end of the script if the function
// was invoked when already at the end of the script.
//
// Invoking this function when already at the end of the script is not
// considered an error and will simply return false.
func (t *scriptTokenizer) Next() bool {
	if t.Done() {
		return false
	}

	op := t.script[t.offset]

	// Only the following op_code will be encountered:
	// OP_PUSHDATA*, OP_DATA_*, OP_CHECKSIG, OP_IF, OP_ENDIF, OP_FALSE
	switch {
	// No additional data.  Note that some of the opcodes, notably OP_1NEGATE,
	// OP_0, and OP_[1-16] represent the data themselves.
	case op == txscript.OP_FALSE || op == txscript.OP_IF || op == txscript.OP_CHECKSIG || op == txscript.OP_ENDIF:
		t.offset++
		t.op = op
		t.data = nil
		return true

	// Data pushes of specific lengths -- OP_DATA_[1-75].
	case op >= txscript.OP_DATA_1 && op <= txscript.OP_DATA_75:
		script := t.script[t.offset:]

		// add 2 instead of 1 because script includes the opcode as well
		length := int32(op) - txscript.OP_DATA_1 + 2
		if int32(len(script)) < length {
			t.err = fmt.Errorf("opcode %d requires %d bytes, but script only "+
				"has %d remaining", op, length, len(script))
			return false
		}

		// Move the offset forward and set the opcode and data accordingly.
		t.offset += length
		t.op = op
		t.data = script[1:length]

		return true
	case op > txscript.OP_PUSHDATA4:
		t.err = fmt.Errorf("unexpected op code")
		return false
	// Data pushes with parsed lengths -- OP_PUSHDATA{1,2,4}.
	default:
		var length int32
		switch op {
		case txscript.OP_PUSHDATA1:
			length = 1
		case txscript.OP_PUSHDATA2:
			length = 2
		default:
			length = 4
		}

		script := t.script[t.offset+1:]
		if int32(len(script)) < length {
			t.err = fmt.Errorf("opcode %d requires %d bytes, but script only "+
				"has %d remaining", op, length, len(script))
			return false
		}

		// Next -length bytes are little endian length of data.
		var dataLen int32
		switch length {
		case 1:
			dataLen = int32(script[0])
		case 2:
			dataLen = int32(binary.LittleEndian.Uint16(script[:2]))
		case 4:
			dataLen = int32(binary.LittleEndian.Uint32(script[:4]))
		default:
			t.err = fmt.Errorf("invalid opcode length %d", length)
			return false
		}

		// Move to the beginning of the data.
		script = script[length:]

		// Disallow entries that do not fit script or were sign extended.
		if dataLen > int32(len(script)) || dataLen < 0 {
			t.err = fmt.Errorf("opcode %d pushes %d bytes, but script only "+
				"has %d remaining", op, dataLen, len(script))
			return false
		}

		// Move the offset forward and set the opcode and data accordingly.
		t.offset += 1 + int32(length) + dataLen
		t.op = op
		t.data = script[:dataLen]
		return true
	}
}
