package memo

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
)

const (
	// ABIAlignment is the number of bytes used to align the ABI encoded data
	ABIAlignment = 32

	// selectorLength is the length of the selector in bytes
	selectorLength = 4

	// codecMethod is the name of the codec method
	codecMethod = "codec"

	// codecMethodABIString is the ABI string template for codec method
	codecMethodABIString = `[{"name":"codec", "inputs":[%s], "outputs":[%s], "type":"function"}]`
)

var _ Codec = (*CodecABI)(nil)

// CodecABI is a coder/decoder for ABI encoded memo fields
type CodecABI struct {
	// abiTypes contains the ABI types of the arguments
	abiTypes []string

	// abiArgs contains the ABI arguments to be packed or unpacked into
	abiArgs []interface{}
}

// NewCodecABI creates a new ABI codec
func NewCodecABI() *CodecABI {
	return &CodecABI{
		abiTypes: make([]string, 0),
		abiArgs:  make([]interface{}, 0),
	}
}

// AddArguments adds a list of arguments to the codec
func (c *CodecABI) AddArguments(args ...CodecArg) {
	for _, arg := range args {
		typeJSON := fmt.Sprintf(`{"type":"%s"}`, arg.Type)
		c.abiTypes = append(c.abiTypes, typeJSON)
		c.abiArgs = append(c.abiArgs, arg.Arg)
	}
}

// PackArguments packs the arguments into the ABI encoded data
func (c *CodecABI) PackArguments() ([]byte, error) {
	// get parsed ABI based on the inputs
	parsedABI, err := c.parsedABI()
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse ABI string")
	}

	// pack the arguments
	data, err := parsedABI.Pack(codecMethod, c.abiArgs...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to pack ABI arguments")
	}

	return data[selectorLength:], nil
}

// UnpackArguments unpacks the ABI encoded data into the output arguments
func (c *CodecABI) UnpackArguments(data []byte) error {
	// get parsed ABI based on the inputs
	parsedABI, err := c.parsedABI()
	if err != nil {
		return errors.Wrap(err, "failed to parse ABI string")
	}

	// unpack data into outputs
	err = parsedABI.UnpackIntoInterface(&c.abiArgs, codecMethod, data)
	if err != nil {
		return errors.Wrap(err, "failed to unpack ABI encoded data")
	}

	return nil
}

// parsedABI builds a parsed ABI based on the inputs
func (c *CodecABI) parsedABI() (abi.ABI, error) {
	typeList := strings.Join(c.abiTypes, ",")
	abiString := fmt.Sprintf(codecMethodABIString, typeList, typeList)
	return abi.JSON(strings.NewReader(abiString))
}
