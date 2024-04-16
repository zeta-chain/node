// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: crosschain/gas_price.proto

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

type GasPrice struct {
	Creator     string   `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Index       string   `protobuf:"bytes,2,opt,name=index,proto3" json:"index,omitempty"`
	ChainId     int64    `protobuf:"varint,3,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Signers     []string `protobuf:"bytes,4,rep,name=signers,proto3" json:"signers,omitempty"`
	BlockNums   []uint64 `protobuf:"varint,5,rep,packed,name=block_nums,json=blockNums,proto3" json:"block_nums,omitempty"`
	Prices      []uint64 `protobuf:"varint,6,rep,packed,name=prices,proto3" json:"prices,omitempty"`
	MedianIndex uint64   `protobuf:"varint,7,opt,name=median_index,json=medianIndex,proto3" json:"median_index,omitempty"`
}

func (m *GasPrice) Reset()         { *m = GasPrice{} }
func (m *GasPrice) String() string { return proto.CompactTextString(m) }
func (*GasPrice) ProtoMessage()    {}
func (*GasPrice) Descriptor() ([]byte, []int) {
	return fileDescriptor_a9c78c67aa323583, []int{0}
}
func (m *GasPrice) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GasPrice) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GasPrice.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GasPrice) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GasPrice.Merge(m, src)
}
func (m *GasPrice) XXX_Size() int {
	return m.Size()
}
func (m *GasPrice) XXX_DiscardUnknown() {
	xxx_messageInfo_GasPrice.DiscardUnknown(m)
}

var xxx_messageInfo_GasPrice proto.InternalMessageInfo

func (m *GasPrice) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *GasPrice) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *GasPrice) GetChainId() int64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

func (m *GasPrice) GetSigners() []string {
	if m != nil {
		return m.Signers
	}
	return nil
}

func (m *GasPrice) GetBlockNums() []uint64 {
	if m != nil {
		return m.BlockNums
	}
	return nil
}

func (m *GasPrice) GetPrices() []uint64 {
	if m != nil {
		return m.Prices
	}
	return nil
}

func (m *GasPrice) GetMedianIndex() uint64 {
	if m != nil {
		return m.MedianIndex
	}
	return 0
}

func init() {
	proto.RegisterType((*GasPrice)(nil), "crosschain.GasPrice")
}

func init() { proto.RegisterFile("crosschain/gas_price.proto", fileDescriptor_a9c78c67aa323583) }

var fileDescriptor_a9c78c67aa323583 = []byte{
	// 268 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xc1, 0x4a, 0xc4, 0x30,
	0x14, 0x45, 0x27, 0xb6, 0xd3, 0xce, 0x3c, 0x5d, 0x05, 0x91, 0x28, 0x18, 0xaa, 0xab, 0x6e, 0x9c,
	0x22, 0xfe, 0x81, 0x1b, 0x19, 0x04, 0x91, 0x2e, 0xdd, 0x94, 0x34, 0x0d, 0x9d, 0xa0, 0x6d, 0x4a,
	0x5e, 0x0b, 0xa3, 0x5f, 0xe1, 0x67, 0xe9, 0x6e, 0x96, 0x2e, 0xa5, 0xfd, 0x11, 0x69, 0x3a, 0xa2,
	0xbb, 0x77, 0xce, 0x85, 0xcb, 0xe3, 0xc2, 0x99, 0xb4, 0x06, 0x51, 0x6e, 0x84, 0xae, 0x93, 0x52,
	0x60, 0xd6, 0x58, 0x2d, 0xd5, 0xaa, 0xb1, 0xa6, 0x35, 0x14, 0xfe, 0xb2, 0xcb, 0x4f, 0x02, 0x8b,
	0x3b, 0x81, 0x8f, 0x63, 0x4c, 0x19, 0x84, 0xd2, 0x2a, 0xd1, 0x1a, 0xcb, 0x48, 0x44, 0xe2, 0x65,
	0xfa, 0x8b, 0xf4, 0x18, 0xe6, 0xba, 0x2e, 0xd4, 0x96, 0x1d, 0x38, 0x3f, 0x01, 0x3d, 0x85, 0x85,
	0x6b, 0xc9, 0x74, 0xc1, 0xbc, 0x88, 0xc4, 0x5e, 0x1a, 0x3a, 0x5e, 0x17, 0x63, 0x15, 0xea, 0xb2,
	0x56, 0x16, 0x99, 0x1f, 0x79, 0x63, 0xd5, 0x1e, 0xe9, 0x39, 0x40, 0xfe, 0x62, 0xe4, 0x73, 0x56,
	0x77, 0x15, 0xb2, 0x79, 0xe4, 0xc5, 0x7e, 0xba, 0x74, 0xe6, 0xa1, 0xab, 0x90, 0x9e, 0x40, 0xe0,
	0x7e, 0x45, 0x16, 0xb8, 0x68, 0x4f, 0xf4, 0x02, 0x8e, 0x2a, 0x55, 0x68, 0x51, 0x67, 0xd3, 0x23,
	0x61, 0x44, 0x62, 0x3f, 0x3d, 0x9c, 0xdc, 0x7a, 0x54, 0xb7, 0xf7, 0x1f, 0x3d, 0x27, 0xbb, 0x9e,
	0x93, 0xef, 0x9e, 0x93, 0xf7, 0x81, 0xcf, 0x76, 0x03, 0x9f, 0x7d, 0x0d, 0x7c, 0xf6, 0x74, 0x5d,
	0xea, 0x76, 0xd3, 0xe5, 0x2b, 0x69, 0xaa, 0xe4, 0x4d, 0xb5, 0xe2, 0x6a, 0x1a, 0x66, 0x3c, 0xa5,
	0xb1, 0x2a, 0xd9, 0x26, 0xff, 0xe6, 0x6a, 0x5f, 0x1b, 0x85, 0x79, 0xe0, 0xb6, 0xba, 0xf9, 0x09,
	0x00, 0x00, 0xff, 0xff, 0x46, 0x57, 0x36, 0xa6, 0x49, 0x01, 0x00, 0x00,
}

func (m *GasPrice) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GasPrice) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GasPrice) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MedianIndex != 0 {
		i = encodeVarintGasPrice(dAtA, i, uint64(m.MedianIndex))
		i--
		dAtA[i] = 0x38
	}
	if len(m.Prices) > 0 {
		dAtA2 := make([]byte, len(m.Prices)*10)
		var j1 int
		for _, num := range m.Prices {
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		i -= j1
		copy(dAtA[i:], dAtA2[:j1])
		i = encodeVarintGasPrice(dAtA, i, uint64(j1))
		i--
		dAtA[i] = 0x32
	}
	if len(m.BlockNums) > 0 {
		dAtA4 := make([]byte, len(m.BlockNums)*10)
		var j3 int
		for _, num := range m.BlockNums {
			for num >= 1<<7 {
				dAtA4[j3] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j3++
			}
			dAtA4[j3] = uint8(num)
			j3++
		}
		i -= j3
		copy(dAtA[i:], dAtA4[:j3])
		i = encodeVarintGasPrice(dAtA, i, uint64(j3))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Signers) > 0 {
		for iNdEx := len(m.Signers) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Signers[iNdEx])
			copy(dAtA[i:], m.Signers[iNdEx])
			i = encodeVarintGasPrice(dAtA, i, uint64(len(m.Signers[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if m.ChainId != 0 {
		i = encodeVarintGasPrice(dAtA, i, uint64(m.ChainId))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Index) > 0 {
		i -= len(m.Index)
		copy(dAtA[i:], m.Index)
		i = encodeVarintGasPrice(dAtA, i, uint64(len(m.Index)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintGasPrice(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGasPrice(dAtA []byte, offset int, v uint64) int {
	offset -= sovGasPrice(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GasPrice) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovGasPrice(uint64(l))
	}
	l = len(m.Index)
	if l > 0 {
		n += 1 + l + sovGasPrice(uint64(l))
	}
	if m.ChainId != 0 {
		n += 1 + sovGasPrice(uint64(m.ChainId))
	}
	if len(m.Signers) > 0 {
		for _, s := range m.Signers {
			l = len(s)
			n += 1 + l + sovGasPrice(uint64(l))
		}
	}
	if len(m.BlockNums) > 0 {
		l = 0
		for _, e := range m.BlockNums {
			l += sovGasPrice(uint64(e))
		}
		n += 1 + sovGasPrice(uint64(l)) + l
	}
	if len(m.Prices) > 0 {
		l = 0
		for _, e := range m.Prices {
			l += sovGasPrice(uint64(e))
		}
		n += 1 + sovGasPrice(uint64(l)) + l
	}
	if m.MedianIndex != 0 {
		n += 1 + sovGasPrice(uint64(m.MedianIndex))
	}
	return n
}

func sovGasPrice(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGasPrice(x uint64) (n int) {
	return sovGasPrice(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GasPrice) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGasPrice
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
			return fmt.Errorf("proto: GasPrice: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GasPrice: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGasPrice
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
				return ErrInvalidLengthGasPrice
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGasPrice
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGasPrice
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
				return ErrInvalidLengthGasPrice
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGasPrice
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Index = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			m.ChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGasPrice
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChainId |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signers", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGasPrice
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
				return ErrInvalidLengthGasPrice
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGasPrice
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signers = append(m.Signers, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 5:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGasPrice
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.BlockNums = append(m.BlockNums, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGasPrice
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthGasPrice
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthGasPrice
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.BlockNums) == 0 {
					m.BlockNums = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGasPrice
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.BlockNums = append(m.BlockNums, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockNums", wireType)
			}
		case 6:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGasPrice
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Prices = append(m.Prices, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGasPrice
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthGasPrice
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthGasPrice
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.Prices) == 0 {
					m.Prices = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGasPrice
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Prices = append(m.Prices, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Prices", wireType)
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MedianIndex", wireType)
			}
			m.MedianIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGasPrice
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MedianIndex |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGasPrice(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGasPrice
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
func skipGasPrice(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGasPrice
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
					return 0, ErrIntOverflowGasPrice
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
					return 0, ErrIntOverflowGasPrice
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
				return 0, ErrInvalidLengthGasPrice
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGasPrice
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGasPrice
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGasPrice        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGasPrice          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGasPrice = fmt.Errorf("proto: unexpected end of group")
)
