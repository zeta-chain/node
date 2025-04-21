package sui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pkg/errors"
)

// PureUint64ArgFromStr converts a string to a uint64 PTB pure argument
func PureUint64FromString(ptb *suiptb.ProgrammableTransactionBuilder, integerStr string) (suiptb.Argument, error) {
	valueUint64, err := strconv.ParseUint(integerStr, 10, 64)
	if err != nil {
		return suiptb.Argument{}, errors.Wrapf(err, "failed to parse amount %s", integerStr)
	}

	arg, err := ptb.Pure(valueUint64)
	if err != nil {
		return suiptb.Argument{}, errors.Wrapf(err, "failed to create amount argument")
	}

	return arg, nil
}

// ParseTypeTagFromStr parses a PTB type argument StructTag from a type string
// Example: "0000000000000000000000000000000000000000000000000000000000000002::sui::SUI" ->
//
//	&sui.StructTag{
//		Address: "0x0000000000000000000000000000000000000000000000000000000000000002",
//		Module:  "sui",
//		Name:    "SUI",
//	}
func ParseTypeTagFromString(t string) (*sui.StructTag, error) {
	parts := strings.Split(t, TypeSeparator)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid type string: %s", t)
	}

	address, err := sui.AddressFromHex(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid address: %s", parts[0])
	}

	module := parts[1]
	name := parts[2]

	return &sui.StructTag{
		Address: address,
		Module:  module,
		Name:    name,
	}, nil
}
