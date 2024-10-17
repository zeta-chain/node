package memo

// ArgType is the enum for types supported by the codec
type ArgType string

// Define all the types supported by the codec
const (
	ArgTypeBytes   ArgType = "bytes"
	ArgTypeString  ArgType = "string"
	ArgTypeAddress ArgType = "address"
)

// CodecArg represents a codec argument
type CodecArg struct {
	Name string
	Type ArgType
	Arg  interface{}
}

// NewArg create a new codec argument
func NewArg(name string, argType ArgType, arg interface{}) CodecArg {
	return CodecArg{
		Name: name,
		Type: argType,
		Arg:  arg,
	}
}

// ArgReceiver wraps the receiver address in a CodecArg
func ArgReceiver(arg interface{}) CodecArg {
	return NewArg("receiver", ArgTypeAddress, arg)
}

// ArgPayload wraps the payload in a CodecArg
func ArgPayload(arg interface{}) CodecArg {
	return NewArg("payload", ArgTypeBytes, arg)
}

// ArgRevertAddress wraps the revert address in a CodecArg
func ArgRevertAddress(arg interface{}) CodecArg {
	return NewArg("revertAddress", ArgTypeString, arg)
}

// ArgAbortAddress wraps the abort address in a CodecArg
func ArgAbortAddress(arg interface{}) CodecArg {
	return NewArg("abortAddress", ArgTypeAddress, arg)
}

// ArgRevertMessage wraps the revert message in a CodecArg
func ArgRevertMessage(arg interface{}) CodecArg {
	return NewArg("revertMessage", ArgTypeBytes, arg)
}
