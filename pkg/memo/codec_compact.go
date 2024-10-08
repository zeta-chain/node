package memo

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var _ Codec = (*CodecCompact)(nil)

// CodecCompact is a coder/decoder for compact encoded memo fields
type CodecCompact struct {
	// lenBytes is the number of bytes used to encode the length of the data
	lenBytes int

	// args contains the list of arguments
	args []CodecArg
}

// NewCodecCompact creates a new compact codec
func NewCodecCompact(encodingFmt uint8) (*CodecCompact, error) {
	lenBytes, err := GetLenBytes(encodingFmt)
	if err != nil {
		return nil, err
	}

	return &CodecCompact{
		lenBytes: lenBytes,
		args:     make([]CodecArg, 0),
	}, nil
}

// AddArguments adds a list of arguments to the codec
func (c *CodecCompact) AddArguments(args ...CodecArg) {
	c.args = append(c.args, args...)
}

// PackArguments packs the arguments into the compact encoded data
func (c *CodecCompact) PackArguments() ([]byte, error) {
	data := make([]byte, 0)

	// pack according to argument type
	for _, arg := range c.args {
		switch arg.Type {
		case ArgTypeBytes:
			dataBytes, err := c.packBytes(arg.Arg)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to pack bytes argument: %s", arg.Name)
			}
			data = append(data, dataBytes...)
		case ArgTypeAddress:
			dateAddress, err := c.packAddress(arg.Arg)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to pack address argument: %s", arg.Name)
			}
			data = append(data, dateAddress...)
		case ArgTypeString:
			dataString, err := c.packString(arg.Arg)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to pack string argument: %s", arg.Name)
			}
			data = append(data, dataString...)
		default:
			return nil, fmt.Errorf("unsupported argument (%s) type: %s", arg.Name, arg.Type)
		}
	}

	return data, nil
}

// UnpackArguments unpacks the compact encoded data into the output arguments
func (c *CodecCompact) UnpackArguments(data []byte) error {
	// unpack according to argument type
	offset := 0
	for _, arg := range c.args {
		switch arg.Type {
		case ArgTypeBytes:
			bytesRead, err := c.unpackBytes(data[offset:], arg.Arg)
			if err != nil {
				return errors.Wrapf(err, "failed to unpack bytes argument: %s", arg.Name)
			}
			offset += bytesRead
		case ArgTypeAddress:
			bytesRead, err := c.unpackAddress(data[offset:], arg.Arg)
			if err != nil {
				return errors.Wrapf(err, "failed to unpack address argument: %s", arg.Name)
			}
			offset += bytesRead
		case ArgTypeString:
			bytesRead, err := c.unpackString(data[offset:], arg.Arg)
			if err != nil {
				return errors.Wrapf(err, "failed to unpack string argument: %s", arg.Name)
			}
			offset += bytesRead
		default:
			return fmt.Errorf("unsupported argument (%s) type: %s", arg.Name, arg.Type)
		}
	}

	// ensure all data is consumed
	if offset != len(data) {
		return fmt.Errorf("consumed bytes (%d) != total bytes (%d)", offset, len(data))
	}

	return nil
}

// packLength packs the length of the data into the compact format
func (c *CodecCompact) packLength(length int) ([]byte, error) {
	data := make([]byte, c.lenBytes)

	switch c.lenBytes {
	case LenBytesShort:
		if length > math.MaxUint8 {
			return nil, fmt.Errorf("data length %d exceeds %d bytes", length, math.MaxUint8)
		}
		data[0] = uint8(length)
	case LenBytesLong:
		if length > math.MaxUint16 {
			return nil, fmt.Errorf("data length %d exceeds %d bytes", length, math.MaxUint16)
		}
		binary.LittleEndian.PutUint16(data, uint16(length))
	}
	return data, nil
}

// packAddress packs argument of type 'address'.
func (c *CodecCompact) packAddress(arg interface{}) ([]byte, error) {
	// type assertion
	address, ok := arg.(common.Address)
	if !ok {
		return nil, fmt.Errorf("argument is not of type common.Address")
	}

	return address.Bytes(), nil
}

// packBytes packs argument of type 'bytes'.
func (c *CodecCompact) packBytes(arg interface{}) ([]byte, error) {
	// type assertion
	bytes, ok := arg.([]byte)
	if !ok {
		return nil, fmt.Errorf("argument is not of type []byte")
	}

	// pack length of the data
	data, err := c.packLength(len(bytes))
	if err != nil {
		return nil, errors.Wrap(err, "failed to pack length of bytes")
	}

	// append the data
	data = append(data, bytes...)
	return data, nil
}

// packString packs argument of type 'string'.
func (c *CodecCompact) packString(arg interface{}) ([]byte, error) {
	// type assertion
	str, ok := arg.(string)
	if !ok {
		return nil, fmt.Errorf("argument is not of type string")
	}

	// pack length of the data
	data, err := c.packLength(len([]byte(str)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to pack length of string")
	}

	// append the string
	data = append(data, []byte(str)...)
	return data, nil
}

// unpackLength returns the length of the data encoded in the compact format
func (c *CodecCompact) unpackLength(data []byte) (int, error) {
	if len(data) < c.lenBytes {
		return 0, fmt.Errorf("expected %d bytes to decode length, got %d", c.lenBytes, len(data))
	}

	// decode length of the data
	length := 0
	switch c.lenBytes {
	case LenBytesShort:
		length = int(data[0])
	case LenBytesLong:
		// convert little-endian bytes to integer
		length = int(binary.LittleEndian.Uint16(data[:2]))
	}

	// ensure remaining data is long enough
	if len(data) < c.lenBytes+length {
		return 0, fmt.Errorf("expected %d bytes, got %d", length, len(data)-c.lenBytes)
	}

	return length, nil
}

// unpackAddress unpacks argument of type 'address'.
func (c *CodecCompact) unpackAddress(data []byte, output interface{}) (int, error) {
	// type assertion
	pAddress, ok := output.(*common.Address)
	if !ok {
		return 0, fmt.Errorf("argument is not of type *common.Address")
	}

	// ensure remaining data >= 20 bytes
	if len(data) < common.AddressLength {
		return 0, fmt.Errorf("expected address, got %d bytes", len(data))
	}
	*pAddress = common.BytesToAddress((data[:20]))

	return common.AddressLength, nil
}

// unpackBytes unpacks argument of type 'bytes' and returns the number of bytes read.
func (c *CodecCompact) unpackBytes(data []byte, output interface{}) (int, error) {
	// type assertion
	pSlice, ok := output.(*[]byte)
	if !ok {
		return 0, fmt.Errorf("argument is not of type *[]byte")
	}

	// unpack length
	dataLen, err := c.unpackLength(data)
	if err != nil {
		return 0, errors.Wrap(err, "failed to unpack length of bytes")
	}

	// make a copy of the data
	*pSlice = make([]byte, dataLen)
	copy(*pSlice, data[c.lenBytes:c.lenBytes+dataLen])

	return c.lenBytes + dataLen, nil
}

// unpackString unpacks argument of type 'string' and returns the number of bytes read.
func (c *CodecCompact) unpackString(data []byte, output interface{}) (int, error) {
	// type assertion
	pString, ok := output.(*string)
	if !ok {
		return 0, fmt.Errorf("argument is not of type *string")
	}

	// unpack length
	strLen, err := c.unpackLength(data)
	if err != nil {
		return 0, errors.Wrap(err, "failed to unpack length of string")
	}

	// make a copy of the string
	*pString = string(data[c.lenBytes : c.lenBytes+strLen])

	return c.lenBytes + strLen, nil
}
