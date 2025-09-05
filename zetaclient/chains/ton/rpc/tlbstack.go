package rpc

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tonkeeper/tongo/tlb"
)

const (
	typeTinyInt = "VmStkTinyInt"
	typeInt     = "VmStkInt"
)

// toncenter api uses toncenter/pytonlib which relies on toncenter/tvm_valuetypes.
// by looking at its sources, we can partially mimic the logic
// https://github.com/toncenter/tvm_valuetypes/blob/55b910782eceee5824bc01c3d280905c1432be9d/tvm_valuetypes/cell.py#L437-L445
//
// it supports: "num" (hex) and "cell" (base64) for input arguments. however we only support "num"
// due to awful encoding logic in tvm_valuetypes. Not an issue, since observer-signer doesn't need this feature,
// only e2e tests need it.
func marshalStack(stack tlb.VmStack) ([][]any, error) {
	items := make([][]any, len(stack))

	for i, arg := range stack {
		switch {
		case arg.SumType == typeTinyInt:
			items[i] = []any{"num", fmt.Sprintf("%d", arg.VmStkTinyInt)}

		case arg.SumType == typeInt:
			bi := big.Int(arg.VmStkInt)
			items[i] = []any{"num", "0x" + bi.Text(16)}

		default:
			return nil, errors.Errorf("unsupported argument type: %s", arg.SumType)
		}
	}

	return items, nil
}

func parseGetMethodResponse(res json.RawMessage) (uint32, tlb.VmStack, error) {
	items := gjson.GetManyBytes(res, "exit_code", "stack")

	exitCode, err := hexToInt(items[0].String())
	if err != nil {
		return 0, tlb.VmStack{}, errors.Wrapf(err, "unable to parse exit code")
	}

	stack := tlb.VmStack{}

	for _, arg := range items[1].Array() {
		pair := arg.Array()
		if len(pair) != 2 {
			return 0, tlb.VmStack{}, errors.Errorf("expected 2 items in pair, got %d", len(pair))
		}

		if pair[0].String() != "num" {
			return 0, tlb.VmStack{}, errors.Errorf("only num is supported, got %s", pair[0].String())
		}

		num, err := hexToInt(pair[1].String())
		if err != nil {
			return 0, tlb.VmStack{}, errors.Wrapf(err, "unable to parse num")
		}

		stack = append(stack, tlb.VmStackValue{
			SumType:      "VmStkTinyInt",
			VmStkTinyInt: num,
		})
	}

	// #nosec G115 always in range
	return uint32(exitCode), stack, nil
}

func hexToInt(s string) (int64, error) {
	s = strings.TrimPrefix(s, "0x")

	return strconv.ParseInt(s, 16, 64)
}
