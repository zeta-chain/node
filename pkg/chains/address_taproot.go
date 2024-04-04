package chains

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
)

// taproot address type

type AddressSegWit struct {
	hrp            string
	witnessVersion byte
	witnessProgram []byte
}

type AddressTaproot struct {
	AddressSegWit
}

var _ btcutil.Address = &AddressTaproot{}

// NewAddressTaproot returns a new AddressTaproot.
func NewAddressTaproot(witnessProg []byte,
	net *chaincfg.Params) (*AddressTaproot, error) {

	return newAddressTaproot(net.Bech32HRPSegwit, witnessProg)
}

// newAddressWitnessScriptHash is an internal helper function to create an
// AddressWitnessScriptHash with a known human-readable part, rather than
// looking it up through its parameters.
func newAddressTaproot(hrp string, witnessProg []byte) (*AddressTaproot, error) {
	// Check for valid program length for witness version 1, which is 32
	// for P2TR.
	if len(witnessProg) != 32 {
		return nil, errors.New("witness program must be 32 bytes for " +
			"p2tr")
	}

	addr := &AddressTaproot{
		AddressSegWit{
			hrp:            strings.ToLower(hrp),
			witnessVersion: 0x01,
			witnessProgram: witnessProg,
		},
	}

	return addr, nil
}

// EncodeAddress returns the bech32 (or bech32m for SegWit v1) string encoding
// of an AddressSegWit.
//
// NOTE: This method is part of the Address interface.
func (a AddressSegWit) EncodeAddress() string {
	str, err := encodeSegWitAddress(
		a.hrp, a.witnessVersion, a.witnessProgram[:],
	)
	if err != nil {
		return ""
	}
	return str
}

// encodeSegWitAddress creates a bech32 (or bech32m for SegWit v1) encoded
// address string representation from witness version and witness program.
func encodeSegWitAddress(hrp string, witnessVersion byte, witnessProgram []byte) (string, error) {
	// Group the address bytes into 5 bit groups, as this is what is used to
	// encode each character in the address string.
	converted, err := bech32.ConvertBits(witnessProgram, 8, 5, true)
	if err != nil {
		return "", err
	}

	// Concatenate the witness version and program, and encode the resulting
	// bytes using bech32 encoding.
	combined := make([]byte, len(converted)+1)
	combined[0] = witnessVersion
	copy(combined[1:], converted)

	var bech string
	switch witnessVersion {
	case 0:
		bech, err = bech32.Encode(hrp, combined)

	case 1:
		bech, err = bech32.EncodeM(hrp, combined)

	default:
		return "", fmt.Errorf("unsupported witness version %d",
			witnessVersion)
	}
	if err != nil {
		return "", err
	}

	// Check validity by decoding the created address.
	_, version, program, err := decodeSegWitAddress(bech)
	if err != nil {
		return "", fmt.Errorf("invalid segwit address: %v", err)
	}

	if version != witnessVersion || !bytes.Equal(program, witnessProgram) {
		return "", fmt.Errorf("invalid segwit address")
	}

	return bech, nil
}

// decodeSegWitAddress parses a bech32 encoded segwit address string and
// returns the witness version and witness program byte representation.
func decodeSegWitAddress(address string) (string, byte, []byte, error) {
	// Decode the bech32 encoded address.
	hrp, data, bech32version, err := bech32.DecodeGeneric(address)
	if err != nil {
		return "", 0, nil, err
	}

	// The first byte of the decoded address is the witness version, it must
	// exist.
	if len(data) < 1 {
		return "", 0, nil, fmt.Errorf("no witness version")
	}

	// ...and be <= 16.
	version := data[0]
	if version > 16 {
		return "", 0, nil, fmt.Errorf("invalid witness version: %v", version)
	}

	// The remaining characters of the address returned are grouped into
	// words of 5 bits. In order to restore the original witness program
	// bytes, we'll need to regroup into 8 bit words.
	regrouped, err := bech32.ConvertBits(data[1:], 5, 8, false)
	if err != nil {
		return "", 0, nil, err
	}

	// The regrouped data must be between 2 and 40 bytes.
	if len(regrouped) < 2 || len(regrouped) > 40 {
		return "", 0, nil, fmt.Errorf("invalid data length")
	}

	// For witness version 0, address MUST be exactly 20 or 32 bytes.
	if version == 0 && len(regrouped) != 20 && len(regrouped) != 32 {
		return "", 0, nil, fmt.Errorf("invalid data length for witness "+
			"version 0: %v", len(regrouped))
	}

	// For witness version 0, the bech32 encoding must be used.
	if version == 0 && bech32version != bech32.Version0 {
		return "", 0, nil, fmt.Errorf("invalid checksum expected bech32 " +
			"encoding for address with witness version 0")
	}

	// For witness version 1, the bech32m encoding must be used.
	if version == 1 && bech32version != bech32.VersionM {
		return "", 0, nil, fmt.Errorf("invalid checksum expected bech32m " +
			"encoding for address with witness version 1")
	}

	return hrp, version, regrouped, nil
}

// ScriptAddress returns the witness program for this address.
//
// NOTE: This method is part of the Address interface.
func (a *AddressSegWit) ScriptAddress() []byte {
	return a.witnessProgram[:]
}

// IsForNet returns whether the AddressSegWit is associated with the passed
// bitcoin network.
//
// NOTE: This method is part of the Address interface.
func (a *AddressSegWit) IsForNet(net *chaincfg.Params) bool {
	return a.hrp == net.Bech32HRPSegwit
}

// String returns a human-readable string for the AddressWitnessPubKeyHash.
// This is equivalent to calling EncodeAddress, but is provided so the type
// can be used as a fmt.Stringer.
//
// NOTE: This method is part of the Address interface.
func (a *AddressSegWit) String() string {
	return a.EncodeAddress()
}

// DecodeTaprootAddress decodes taproot address only and returns error on non-taproot address
func DecodeTaprootAddress(addr string) (*AddressTaproot, error) {
	hrp, version, program, err := decodeSegWitAddress(addr)
	if err != nil {
		return nil, err
	}
	if version != 1 {
		return nil, errors.New("invalid witness version; taproot address must be version 1")
	}
	return &AddressTaproot{
		AddressSegWit{
			hrp:            hrp,
			witnessVersion: version,
			witnessProgram: program,
		},
	}, nil
}

// PayToWitnessTaprootScript creates a new script to pay to a version 1
// (taproot) witness program. The passed hash is expected to be valid.
func PayToWitnessTaprootScript(rawKey []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_1).AddData(rawKey).Script()
}
