// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: crosschain/in_tx_hash_to_cctx.proto

package types

import (
	fmt "fmt"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type InTxHashToCctx struct {
	InTxHash  string   `protobuf:"bytes,1,opt,name=in_tx_hash,json=inTxHash,proto3" json:"in_tx_hash,omitempty"`
	CctxIndex []string `protobuf:"bytes,2,rep,name=cctx_index,json=cctxIndex,proto3" json:"cctx_index,omitempty"`
}

func (m *InTxHashToCctx) Reset()         { *m = InTxHashToCctx{} }
func (m *InTxHashToCctx) String() string { return proto.CompactTextString(m) }
func (*InTxHashToCctx) ProtoMessage()    {}
func (*InTxHashToCctx) Descriptor() ([]byte, []int) {
	return fileDescriptor_67ee1b8208d56a23, []int{0}
}
func (m *InTxHashToCctx) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InTxHashToCctx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InTxHashToCctx.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *InTxHashToCctx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InTxHashToCctx.Merge(m, src)
}
func (m *InTxHashToCctx) XXX_Size() int {
	return m.Size()
}
func (m *InTxHashToCctx) XXX_DiscardUnknown() {
	xxx_messageInfo_InTxHashToCctx.DiscardUnknown(m)
}

var xxx_messageInfo_InTxHashToCctx proto.InternalMessageInfo

func (m *InTxHashToCctx) GetInTxHash() string {
	if m != nil {
		return m.InTxHash
	}
	return ""
}

func (m *InTxHashToCctx) GetCctxIndex() []string {
	if m != nil {
		return m.CctxIndex
	}
	return nil
}

func init() {
	proto.RegisterType((*InTxHashToCctx)(nil), "crosschain.InTxHashToCctx")
}

func init() {
	proto.RegisterFile("crosschain/in_tx_hash_to_cctx.proto", fileDescriptor_67ee1b8208d56a23)
}

var fileDescriptor_67ee1b8208d56a23 = []byte{
	// 194 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4e, 0x2e, 0xca, 0x2f,
	0x2e, 0x4e, 0xce, 0x48, 0xcc, 0xcc, 0xd3, 0xcf, 0xcc, 0x8b, 0x2f, 0xa9, 0x88, 0xcf, 0x48, 0x2c,
	0xce, 0x88, 0x2f, 0xc9, 0x8f, 0x4f, 0x4e, 0x2e, 0xa9, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0xe2, 0x42, 0x28, 0x52, 0xf2, 0xe5, 0xe2, 0xf3, 0xcc, 0x0b, 0xa9, 0xf0, 0x48, 0x2c, 0xce, 0x08,
	0xc9, 0x77, 0x4e, 0x2e, 0xa9, 0x10, 0x92, 0xe1, 0xe2, 0x42, 0xe8, 0x94, 0x60, 0x54, 0x60, 0xd4,
	0xe0, 0x0c, 0xe2, 0xc8, 0x84, 0xaa, 0x11, 0x92, 0xe5, 0xe2, 0x02, 0x99, 0x14, 0x9f, 0x99, 0x97,
	0x92, 0x5a, 0x21, 0xc1, 0xa4, 0xc0, 0xac, 0xc1, 0x19, 0xc4, 0x09, 0x12, 0xf1, 0x04, 0x09, 0x38,
	0x79, 0x9f, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x13, 0x1e,
	0xcb, 0x31, 0x5c, 0x78, 0x2c, 0xc7, 0x70, 0xe3, 0xb1, 0x1c, 0x43, 0x94, 0x61, 0x7a, 0x66, 0x49,
	0x46, 0x69, 0x92, 0x5e, 0x72, 0x7e, 0xae, 0x7e, 0x55, 0x6a, 0x49, 0xa2, 0x2e, 0xc4, 0x91, 0x20,
	0x66, 0x72, 0x7e, 0x51, 0xaa, 0x7e, 0x85, 0x3e, 0x92, 0xd3, 0x4b, 0x2a, 0x0b, 0x52, 0x8b, 0x93,
	0xd8, 0xc0, 0xce, 0x35, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xe0, 0x52, 0xbf, 0xdc, 0xd5, 0x00,
	0x00, 0x00,
}

func (m *InTxHashToCctx) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InTxHashToCctx) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *InTxHashToCctx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.CctxIndex) > 0 {
		for iNdEx := len(m.CctxIndex) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.CctxIndex[iNdEx])
			copy(dAtA[i:], m.CctxIndex[iNdEx])
			i = encodeVarintInTxHashToCctx(dAtA, i, uint64(len(m.CctxIndex[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.InTxHash) > 0 {
		i -= len(m.InTxHash)
		copy(dAtA[i:], m.InTxHash)
		i = encodeVarintInTxHashToCctx(dAtA, i, uint64(len(m.InTxHash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintInTxHashToCctx(dAtA []byte, offset int, v uint64) int {
	offset -= sovInTxHashToCctx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *InTxHashToCctx) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.InTxHash)
	if l > 0 {
		n += 1 + l + sovInTxHashToCctx(uint64(l))
	}
	if len(m.CctxIndex) > 0 {
		for _, s := range m.CctxIndex {
			l = len(s)
			n += 1 + l + sovInTxHashToCctx(uint64(l))
		}
	}
	return n
}

func sovInTxHashToCctx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozInTxHashToCctx(x uint64) (n int) {
	return sovInTxHashToCctx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *InTxHashToCctx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInTxHashToCctx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: InTxHashToCctx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InTxHashToCctx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InTxHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInTxHashToCctx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInTxHashToCctx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInTxHashToCctx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.InTxHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CctxIndex", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInTxHashToCctx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthInTxHashToCctx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthInTxHashToCctx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CctxIndex = append(m.CctxIndex, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipInTxHashToCctx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthInTxHashToCctx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipInTxHashToCctx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowInTxHashToCctx
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowInTxHashToCctx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowInTxHashToCctx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthInTxHashToCctx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupInTxHashToCctx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthInTxHashToCctx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthInTxHashToCctx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowInTxHashToCctx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupInTxHashToCctx = fmt.Errorf("proto: unexpected end of group")
)
