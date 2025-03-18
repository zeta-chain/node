package ton

import (
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

	bitString := boc.BitString(sd)
	n := bitString.BitsAvailableForRead() / 8

	return bitString.ReadBytes(n)
}

// MarshalSnakeCell encodes []byte to TLB using snake-cell encoding
func MarshalSnakeCell(data []byte) (*boc.Cell, error) {
	bs := boc.NewBitString(len(data) * 8)

	if err := bs.WriteBytes(data); err != nil {
		return nil, err
	}

	return MarshalTLB(tlb.SnakeData(bs))
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
