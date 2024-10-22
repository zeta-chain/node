package ton

import (
	"bytes"

	"cosmossdk.io/math"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

// MarshalTLB encodes entity to BOC
func MarshalTLB(v tlb.MarshalerTLB) (*boc.Cell, error) {
	cell := boc.NewCell()

	if err := v.MarshalTLB(cell, &tlb.Encoder{}); err != nil {
		return nil, err
	}

	return cell, nil
}

// UnmarshalTLB decodes entity from BOC
func UnmarshalTLB(t tlb.UnmarshalerTLB, cell *boc.Cell) error {
	return t.UnmarshalTLB(cell, &tlb.Decoder{})
}

// UnmarshalSnakeCell decodes TLB cell to []byte using snake-cell encoding
func UnmarshalSnakeCell(cell *boc.Cell) ([]byte, error) {
	var sd tlb.SnakeData

	if err := UnmarshalTLB(&sd, cell); err != nil {
		return nil, err
	}

	cd := boc.BitString(sd)

	// TLB operates with bits, so we (might) need to trim some "leftovers" (null chars)
	return bytes.Trim(cd.Buffer(), "\x00"), nil
}

// MarshalSnakeCell encodes []byte to TLB using snake-cell encoding
func MarshalSnakeCell(data []byte) (*boc.Cell, error) {
	b := boc.NewCell()

	wrapped := tlb.Bytes(data)
	if err := wrapped.MarshalTLB(b, &tlb.Encoder{}); err != nil {
		return nil, err
	}

	return b, nil
}

// UnmarshalEVMAddress decodes eth.Address from BOC
func UnmarshalEVMAddress(cell *boc.Cell) (eth.Address, error) {
	const evmAddrBits = 20 * 8

	s, err := cell.ReadBits(evmAddrBits)
	if err != nil {
		return eth.Address{}, err
	}

	return eth.BytesToAddress(s.Buffer()), nil
}

func GramsToUint(g tlb.Grams) math.Uint {
	return math.NewUint(uint64(g))
}

func ErrCollect(errs ...error) error {
	for i, err := range errs {
		if err != nil {
			return errors.Wrapf(err, "error at index %d", i)
		}
	}

	return nil
}
