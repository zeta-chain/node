// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pkg/crypto/crypto.proto

package crypto

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

// PubKeySet contains two pub keys , secp256k1 and ed25519
type PubKeySet struct {
	Secp256k1 PubKey `protobuf:"bytes,1,opt,name=secp256k1,proto3,casttype=PubKey" json:"secp256k1,omitempty"`
	Ed25519   PubKey `protobuf:"bytes,2,opt,name=ed25519,proto3,casttype=PubKey" json:"ed25519,omitempty"`
}

func (m *PubKeySet) Reset()         { *m = PubKeySet{} }
func (m *PubKeySet) String() string { return proto.CompactTextString(m) }
func (*PubKeySet) ProtoMessage()    {}
func (*PubKeySet) Descriptor() ([]byte, []int) {
	return fileDescriptor_5643a513c82df681, []int{0}
}
func (m *PubKeySet) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PubKeySet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PubKeySet.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PubKeySet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PubKeySet.Merge(m, src)
}
func (m *PubKeySet) XXX_Size() int {
	return m.Size()
}
func (m *PubKeySet) XXX_DiscardUnknown() {
	xxx_messageInfo_PubKeySet.DiscardUnknown(m)
}

var xxx_messageInfo_PubKeySet proto.InternalMessageInfo

func (m *PubKeySet) GetSecp256k1() PubKey {
	if m != nil {
		return m.Secp256k1
	}
	return ""
}

func (m *PubKeySet) GetEd25519() PubKey {
	if m != nil {
		return m.Ed25519
	}
	return ""
}

func init() {
	proto.RegisterType((*PubKeySet)(nil), "crypto.PubKeySet")
}

func init() { proto.RegisterFile("pkg/crypto/crypto.proto", fileDescriptor_5643a513c82df681) }

var fileDescriptor_5643a513c82df681 = []byte{
	// 195 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2f, 0xc8, 0x4e, 0xd7,
	0x4f, 0x2e, 0xaa, 0x2c, 0x28, 0xc9, 0x87, 0x52, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0x6c,
	0x10, 0x9e, 0x94, 0x48, 0x7a, 0x7e, 0x7a, 0x3e, 0x58, 0x48, 0x1f, 0xc4, 0x82, 0xc8, 0x2a, 0x65,
	0x70, 0x71, 0x06, 0x94, 0x26, 0x79, 0xa7, 0x56, 0x06, 0xa7, 0x96, 0x08, 0x99, 0x72, 0x71, 0x16,
	0xa7, 0x26, 0x17, 0x18, 0x99, 0x9a, 0x65, 0x1b, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x3a, 0x89,
	0x3f, 0xba, 0x27, 0xcf, 0x19, 0x0c, 0x13, 0xfc, 0x75, 0x4f, 0x9e, 0x0d, 0xa2, 0x3c, 0x08, 0xa1,
	0x52, 0x48, 0x85, 0x8b, 0x3d, 0x35, 0xc5, 0xc8, 0xd4, 0xd4, 0xd0, 0x52, 0x82, 0x09, 0xac, 0x89,
	0x0b, 0x49, 0x1d, 0x4c, 0xca, 0xc9, 0xf9, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f,
	0x3c, 0x92, 0x63, 0x9c, 0xf0, 0x58, 0x8e, 0xe1, 0xc2, 0x63, 0x39, 0x86, 0x1b, 0x8f, 0xe5, 0x18,
	0xa2, 0x34, 0xd3, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0xab, 0x52, 0x4b,
	0x12, 0x75, 0x93, 0x33, 0x12, 0x33, 0xf3, 0xc0, 0xcc, 0xe4, 0xfc, 0xa2, 0x54, 0x7d, 0x84, 0xcf,
	0x92, 0xd8, 0xc0, 0xae, 0x36, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xe0, 0x44, 0x43, 0x35, 0xee,
	0x00, 0x00, 0x00,
}

func (m *PubKeySet) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PubKeySet) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PubKeySet) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Ed25519) > 0 {
		i -= len(m.Ed25519)
		copy(dAtA[i:], m.Ed25519)
		i = encodeVarintCrypto(dAtA, i, uint64(len(m.Ed25519)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Secp256k1) > 0 {
		i -= len(m.Secp256k1)
		copy(dAtA[i:], m.Secp256k1)
		i = encodeVarintCrypto(dAtA, i, uint64(len(m.Secp256k1)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintCrypto(dAtA []byte, offset int, v uint64) int {
	offset -= sovCrypto(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PubKeySet) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Secp256k1)
	if l > 0 {
		n += 1 + l + sovCrypto(uint64(l))
	}
	l = len(m.Ed25519)
	if l > 0 {
		n += 1 + l + sovCrypto(uint64(l))
	}
	return n
}

func sovCrypto(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCrypto(x uint64) (n int) {
	return sovCrypto(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PubKeySet) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCrypto
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
			return fmt.Errorf("proto: PubKeySet: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PubKeySet: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Secp256k1", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCrypto
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
				return ErrInvalidLengthCrypto
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCrypto
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Secp256k1 = PubKey(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ed25519", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCrypto
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
				return ErrInvalidLengthCrypto
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCrypto
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ed25519 = PubKey(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCrypto(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCrypto
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
func skipCrypto(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCrypto
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
					return 0, ErrIntOverflowCrypto
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
					return 0, ErrIntOverflowCrypto
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
				return 0, ErrInvalidLengthCrypto
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCrypto
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCrypto
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCrypto        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCrypto          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCrypto = fmt.Errorf("proto: unexpected end of group")
)
