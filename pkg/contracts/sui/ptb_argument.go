package sui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pkg/errors"
)

// PureUint64FromString converts a string to a uint64 and creates a PTB pure argument
func PureUint64FromString(
	ptb *suiptb.ProgrammableTransactionBuilder,
	integerStr string,
) (arg suiptb.Argument, value uint64, err error) {
	value, err = strconv.ParseUint(integerStr, 10, 64)
	if err != nil {
		return suiptb.Argument{}, 0, errors.Wrapf(err, "failed to parse amount %s", integerStr)
	}

	arg, err = ptb.Pure(value)
	if err != nil {
		return suiptb.Argument{}, 0, errors.Wrapf(err, "failed to create amount argument")
	}

	return arg, value, nil
}

// TypeTagFromString creates a PTB type argument StructTag from a type string
// Example: "0000000000000000000000000000000000000000000000000000000000000002::sui::SUI" ->
//
//	&sui.StructTag{
//		Address: "0x0000000000000000000000000000000000000000000000000000000000000002",
//		Module:  "sui",
//		Name:    "SUI",
//	}
func TypeTagFromString(t string) (tag sui.StructTag, err error) {
	parts := strings.Split(t, typeSeparator)
	if len(parts) != 3 {
		return tag, fmt.Errorf("invalid type string: %s", t)
	}

	address, err := sui.AddressFromHex(parts[0])
	if err != nil {
		return tag, errors.Wrapf(err, "invalid address: %s", parts[0])
	}

	module := parts[1]
	name := parts[2]

	return sui.StructTag{
		Address: address,
		Module:  module,
		Name:    name,
	}, nil
}
