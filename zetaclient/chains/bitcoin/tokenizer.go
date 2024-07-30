package bitcoin

import (
	"encoding/binary"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
)

func newScriptTokenizer(script []byte) scriptTokenizer {
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
	offset int
	op     byte
	data   []byte
	err    error
}

// Done returns true when either all opcodes have been exhausted or a parse
// failure was encountered and therefore the state has an associated error.
func (t *scriptTokenizer) Done() bool {
	return t.err != nil || t.offset >= len(t.script)
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

		// The length should be: int(op) - txscript.OP_DATA_1 + 2, i.e. op is txscript.OP_DATA_10, that means
		// the data length should be 10, which is txscript.OP_DATA_10 - txscript.OP_DATA_1 + 1.
		// Here, 2 instead of 1 because `script` also includes the opcode which means it contains one more byte.
		// Since txscript.OP_DATA_1 is 1, then length is just int(op) - 1 + 2 = int(op) + 1
		length := int(op) + 1
		if len(script) < length {
			t.err = fmt.Errorf("opcode %d detected, but script only %d bytes remaining", op, len(script))
			return false
		}

		// Move the offset forward and set the opcode and data accordingly.
		t.offset += length
		t.op = op
		t.data = script[1:length]
		return true

	case op > txscript.OP_PUSHDATA4:
		t.err = fmt.Errorf("unexpected op code %d", op)
		return false

	// Data pushes with parsed lengths -- OP_PUSHDATA{1,2,4}.
	default:
		var length int
		switch op {
		case txscript.OP_PUSHDATA1:
			length = 1
		case txscript.OP_PUSHDATA2:
			length = 2
		case txscript.OP_PUSHDATA4:
			length = 4
		default:
			t.err = fmt.Errorf("unexpected op code %d", op)
			return false
		}

		script := t.script[t.offset+1:]
		if len(script) < length {
			t.err = fmt.Errorf("opcode %d requires %d bytes, only %d remaining", op, length, len(script))
			return false
		}

		// Next -length bytes are little endian length of data.
		var dataLen int
		switch length {
		case 1:
			dataLen = int(script[0])
		case 2:
			dataLen = int(binary.LittleEndian.Uint16(script[:length]))
		case 4:
			dataLen = int(binary.LittleEndian.Uint32(script[:length]))
		default:
			t.err = fmt.Errorf("invalid opcode length %d", length)
			return false
		}

		// Move to the beginning of the data.
		script = script[length:]

		// Disallow entries that do not fit script or were sign extended.
		if dataLen > len(script) || dataLen < 0 {
			t.err = fmt.Errorf("opcode %d pushes %d bytes, only %d remaining", op, dataLen, len(script))
			return false
		}

		// Move the offset forward and set the opcode and data accordingly.
		// 1 is the opcode size, which is just 1 byte. int(op) is the opcode value,
		// it should not be mixed with the size.
		t.offset += 1 + length + dataLen
		t.op = op
		t.data = script[:dataLen]
		return true
	}
}
