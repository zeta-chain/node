// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: zetachain/zetacore/crosschain/rate_limiter_flags.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	coin "github.com/zeta-chain/node/pkg/coin"
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

type RateLimiterFlags struct {
	Enabled bool `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	// window in blocks
	Window int64 `protobuf:"varint,2,opt,name=window,proto3" json:"window,omitempty"`
	// rate in azeta per block
	Rate github_com_cosmos_cosmos_sdk_types.Uint `protobuf:"bytes,3,opt,name=rate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"rate"`
	// conversion in azeta per token
	Conversions []Conversion `protobuf:"bytes,4,rep,name=conversions,proto3" json:"conversions"`
}

func (m *RateLimiterFlags) Reset()         { *m = RateLimiterFlags{} }
func (m *RateLimiterFlags) String() string { return proto.CompactTextString(m) }
func (*RateLimiterFlags) ProtoMessage()    {}
func (*RateLimiterFlags) Descriptor() ([]byte, []int) {
	return fileDescriptor_9c435f4c2dabc0eb, []int{0}
}
func (m *RateLimiterFlags) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RateLimiterFlags) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RateLimiterFlags.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RateLimiterFlags) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RateLimiterFlags.Merge(m, src)
}
func (m *RateLimiterFlags) XXX_Size() int {
	return m.Size()
}
func (m *RateLimiterFlags) XXX_DiscardUnknown() {
	xxx_messageInfo_RateLimiterFlags.DiscardUnknown(m)
}

var xxx_messageInfo_RateLimiterFlags proto.InternalMessageInfo

func (m *RateLimiterFlags) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

func (m *RateLimiterFlags) GetWindow() int64 {
	if m != nil {
		return m.Window
	}
	return 0
}

func (m *RateLimiterFlags) GetConversions() []Conversion {
	if m != nil {
		return m.Conversions
	}
	return nil
}

type Conversion struct {
	Zrc20 string                                 `protobuf:"bytes,1,opt,name=zrc20,proto3" json:"zrc20,omitempty"`
	Rate  github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=rate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"rate"`
}

func (m *Conversion) Reset()         { *m = Conversion{} }
func (m *Conversion) String() string { return proto.CompactTextString(m) }
func (*Conversion) ProtoMessage()    {}
func (*Conversion) Descriptor() ([]byte, []int) {
	return fileDescriptor_9c435f4c2dabc0eb, []int{1}
}
func (m *Conversion) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Conversion) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Conversion.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Conversion) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Conversion.Merge(m, src)
}
func (m *Conversion) XXX_Size() int {
	return m.Size()
}
func (m *Conversion) XXX_DiscardUnknown() {
	xxx_messageInfo_Conversion.DiscardUnknown(m)
}

var xxx_messageInfo_Conversion proto.InternalMessageInfo

func (m *Conversion) GetZrc20() string {
	if m != nil {
		return m.Zrc20
	}
	return ""
}

type AssetRate struct {
	ChainId  int64                                  `protobuf:"varint,1,opt,name=chainId,proto3" json:"chainId,omitempty"`
	Asset    string                                 `protobuf:"bytes,2,opt,name=asset,proto3" json:"asset,omitempty"`
	Decimals uint32                                 `protobuf:"varint,3,opt,name=decimals,proto3" json:"decimals,omitempty"`
	CoinType coin.CoinType                          `protobuf:"varint,4,opt,name=coin_type,json=coinType,proto3,enum=zetachain.zetacore.pkg.coin.CoinType" json:"coin_type,omitempty"`
	Rate     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=rate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"rate"`
}

func (m *AssetRate) Reset()         { *m = AssetRate{} }
func (m *AssetRate) String() string { return proto.CompactTextString(m) }
func (*AssetRate) ProtoMessage()    {}
func (*AssetRate) Descriptor() ([]byte, []int) {
	return fileDescriptor_9c435f4c2dabc0eb, []int{2}
}
func (m *AssetRate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AssetRate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AssetRate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AssetRate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AssetRate.Merge(m, src)
}
func (m *AssetRate) XXX_Size() int {
	return m.Size()
}
func (m *AssetRate) XXX_DiscardUnknown() {
	xxx_messageInfo_AssetRate.DiscardUnknown(m)
}

var xxx_messageInfo_AssetRate proto.InternalMessageInfo

func (m *AssetRate) GetChainId() int64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

func (m *AssetRate) GetAsset() string {
	if m != nil {
		return m.Asset
	}
	return ""
}

func (m *AssetRate) GetDecimals() uint32 {
	if m != nil {
		return m.Decimals
	}
	return 0
}

func (m *AssetRate) GetCoinType() coin.CoinType {
	if m != nil {
		return m.CoinType
	}
	return coin.CoinType_Zeta
}

func init() {
	proto.RegisterType((*RateLimiterFlags)(nil), "zetachain.zetacore.crosschain.RateLimiterFlags")
	proto.RegisterType((*Conversion)(nil), "zetachain.zetacore.crosschain.Conversion")
	proto.RegisterType((*AssetRate)(nil), "zetachain.zetacore.crosschain.AssetRate")
}

func init() {
	proto.RegisterFile("zetachain/zetacore/crosschain/rate_limiter_flags.proto", fileDescriptor_9c435f4c2dabc0eb)
}

var fileDescriptor_9c435f4c2dabc0eb = []byte{
	// 437 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0x41, 0x8b, 0x13, 0x31,
	0x14, 0x6e, 0xb6, 0xdd, 0xb5, 0xcd, 0xa2, 0x48, 0x58, 0x64, 0x28, 0x38, 0x3b, 0x14, 0x5c, 0xc7,
	0x43, 0x13, 0xa9, 0xe0, 0xdd, 0xa9, 0x28, 0x82, 0x17, 0x83, 0x5e, 0xbc, 0x94, 0x34, 0x93, 0x9d,
	0x0d, 0xed, 0x24, 0xc3, 0x24, 0xba, 0xee, 0xfe, 0x0a, 0x7f, 0xd6, 0x1e, 0xf7, 0x28, 0x22, 0xab,
	0xb4, 0x7f, 0x44, 0x92, 0x4c, 0x67, 0x7b, 0x28, 0x22, 0x7b, 0x69, 0xdf, 0x37, 0xbc, 0xef, 0x7d,
	0xdf, 0xfb, 0xf2, 0xe0, 0xcb, 0x4b, 0x61, 0x19, 0x3f, 0x63, 0x52, 0x11, 0x5f, 0xe9, 0x5a, 0x10,
	0x5e, 0x6b, 0x63, 0xc2, 0xb7, 0x9a, 0x59, 0x31, 0x5b, 0xca, 0x52, 0x5a, 0x51, 0xcf, 0x4e, 0x97,
	0xac, 0x30, 0xb8, 0xaa, 0xb5, 0xd5, 0xe8, 0x71, 0xcb, 0xc3, 0x1b, 0x1e, 0xbe, 0xe5, 0x0d, 0x8f,
	0x0a, 0x5d, 0x68, 0xdf, 0x49, 0x5c, 0x15, 0x48, 0xc3, 0x93, 0x1d, 0x62, 0xd5, 0xa2, 0x20, 0x5c,
	0x4b, 0xe5, 0x7f, 0x42, 0xdf, 0xe8, 0x17, 0x80, 0x0f, 0x29, 0xb3, 0xe2, 0x7d, 0x10, 0x7e, 0xe3,
	0x74, 0x51, 0x04, 0xef, 0x09, 0xc5, 0xe6, 0x4b, 0x91, 0x47, 0x20, 0x01, 0x69, 0x9f, 0x6e, 0x20,
	0x7a, 0x04, 0x0f, 0xce, 0xa5, 0xca, 0xf5, 0x79, 0xb4, 0x97, 0x80, 0xb4, 0x4b, 0x1b, 0x84, 0xa6,
	0xb0, 0xe7, 0xfc, 0x47, 0xdd, 0x04, 0xa4, 0x83, 0x8c, 0x5c, 0xdd, 0x1c, 0x77, 0x7e, 0xde, 0x1c,
	0x3f, 0x2d, 0xa4, 0x3d, 0xfb, 0x32, 0xc7, 0x5c, 0x97, 0x84, 0x6b, 0x53, 0x6a, 0xd3, 0xfc, 0x8d,
	0x4d, 0xbe, 0x20, 0xf6, 0xa2, 0x12, 0x06, 0x7f, 0x92, 0xca, 0x52, 0x4f, 0x46, 0x1f, 0xe0, 0x21,
	0xd7, 0xea, 0xab, 0xa8, 0x8d, 0xd4, 0xca, 0x44, 0xbd, 0xa4, 0x9b, 0x1e, 0x4e, 0x9e, 0xe1, 0x7f,
	0xae, 0x8f, 0xa7, 0x2d, 0x23, 0xeb, 0x39, 0x59, 0xba, 0x3d, 0x63, 0x74, 0x0a, 0xe1, 0x6d, 0x03,
	0x3a, 0x82, 0xfb, 0x97, 0x35, 0x9f, 0x3c, 0xf7, 0x5b, 0x0d, 0x68, 0x00, 0x28, 0x6b, 0xbc, 0xef,
	0x79, 0xef, 0xb8, 0xf1, 0x7e, 0xf2, 0x1f, 0xde, 0x5f, 0x0b, 0x1e, 0xac, 0x8f, 0x7e, 0x03, 0x38,
	0x78, 0x65, 0x8c, 0xb0, 0x2e, 0x4b, 0x97, 0x9f, 0x37, 0xf7, 0x2e, 0xe4, 0xd7, 0xa5, 0x1b, 0xe8,
	0x1c, 0x30, 0xd7, 0x16, 0xc4, 0x68, 0x00, 0x68, 0x08, 0xfb, 0xb9, 0xe0, 0xb2, 0x64, 0x4b, 0xe3,
	0x13, 0xbc, 0x4f, 0x5b, 0x8c, 0x32, 0x38, 0x70, 0xcf, 0x35, 0x73, 0x8a, 0x51, 0x2f, 0x01, 0xe9,
	0x83, 0xc9, 0x93, 0x5d, 0x91, 0x54, 0x8b, 0x02, 0xfb, 0x77, 0x9d, 0x6a, 0xa9, 0x3e, 0x5e, 0x54,
	0x82, 0xf6, 0x79, 0x53, 0xb5, 0x1b, 0xee, 0xdf, 0x7d, 0xc3, 0xec, 0xed, 0xd5, 0x2a, 0x06, 0xd7,
	0xab, 0x18, 0xfc, 0x59, 0xc5, 0xe0, 0xfb, 0x3a, 0xee, 0x5c, 0xaf, 0xe3, 0xce, 0x8f, 0x75, 0xdc,
	0xf9, 0x3c, 0xde, 0x9a, 0xe3, 0xec, 0x8c, 0xc3, 0xd9, 0x29, 0x9d, 0x0b, 0xf2, 0x6d, 0xfb, 0xc2,
	0xfd, 0xc8, 0xf9, 0x81, 0x3f, 0xbc, 0x17, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xc7, 0x00, 0x01,
	0xd6, 0x0f, 0x03, 0x00, 0x00,
}

func (m *RateLimiterFlags) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RateLimiterFlags) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RateLimiterFlags) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Conversions) > 0 {
		for iNdEx := len(m.Conversions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Conversions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintRateLimiterFlags(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	{
		size := m.Rate.Size()
		i -= size
		if _, err := m.Rate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintRateLimiterFlags(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.Window != 0 {
		i = encodeVarintRateLimiterFlags(dAtA, i, uint64(m.Window))
		i--
		dAtA[i] = 0x10
	}
	if m.Enabled {
		i--
		if m.Enabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Conversion) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Conversion) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Conversion) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Rate.Size()
		i -= size
		if _, err := m.Rate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintRateLimiterFlags(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Zrc20) > 0 {
		i -= len(m.Zrc20)
		copy(dAtA[i:], m.Zrc20)
		i = encodeVarintRateLimiterFlags(dAtA, i, uint64(len(m.Zrc20)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *AssetRate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AssetRate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AssetRate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Rate.Size()
		i -= size
		if _, err := m.Rate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintRateLimiterFlags(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if m.CoinType != 0 {
		i = encodeVarintRateLimiterFlags(dAtA, i, uint64(m.CoinType))
		i--
		dAtA[i] = 0x20
	}
	if m.Decimals != 0 {
		i = encodeVarintRateLimiterFlags(dAtA, i, uint64(m.Decimals))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Asset) > 0 {
		i -= len(m.Asset)
		copy(dAtA[i:], m.Asset)
		i = encodeVarintRateLimiterFlags(dAtA, i, uint64(len(m.Asset)))
		i--
		dAtA[i] = 0x12
	}
	if m.ChainId != 0 {
		i = encodeVarintRateLimiterFlags(dAtA, i, uint64(m.ChainId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintRateLimiterFlags(dAtA []byte, offset int, v uint64) int {
	offset -= sovRateLimiterFlags(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *RateLimiterFlags) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Enabled {
		n += 2
	}
	if m.Window != 0 {
		n += 1 + sovRateLimiterFlags(uint64(m.Window))
	}
	l = m.Rate.Size()
	n += 1 + l + sovRateLimiterFlags(uint64(l))
	if len(m.Conversions) > 0 {
		for _, e := range m.Conversions {
			l = e.Size()
			n += 1 + l + sovRateLimiterFlags(uint64(l))
		}
	}
	return n
}

func (m *Conversion) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Zrc20)
	if l > 0 {
		n += 1 + l + sovRateLimiterFlags(uint64(l))
	}
	l = m.Rate.Size()
	n += 1 + l + sovRateLimiterFlags(uint64(l))
	return n
}

func (m *AssetRate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ChainId != 0 {
		n += 1 + sovRateLimiterFlags(uint64(m.ChainId))
	}
	l = len(m.Asset)
	if l > 0 {
		n += 1 + l + sovRateLimiterFlags(uint64(l))
	}
	if m.Decimals != 0 {
		n += 1 + sovRateLimiterFlags(uint64(m.Decimals))
	}
	if m.CoinType != 0 {
		n += 1 + sovRateLimiterFlags(uint64(m.CoinType))
	}
	l = m.Rate.Size()
	n += 1 + l + sovRateLimiterFlags(uint64(l))
	return n
}

func sovRateLimiterFlags(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRateLimiterFlags(x uint64) (n int) {
	return sovRateLimiterFlags(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RateLimiterFlags) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRateLimiterFlags
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
			return fmt.Errorf("proto: RateLimiterFlags: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RateLimiterFlags: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Enabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRateLimiterFlags
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Enabled = bool(v != 0)
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Window", wireType)
			}
			m.Window = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRateLimiterFlags
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Window |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRateLimiterFlags
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
				return ErrInvalidLengthRateLimiterFlags
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRateLimiterFlags
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Rate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Conversions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRateLimiterFlags
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRateLimiterFlags
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRateLimiterFlags
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Conversions = append(m.Conversions, Conversion{})
			if err := m.Conversions[len(m.Conversions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRateLimiterFlags(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRateLimiterFlags
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
func (m *Conversion) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRateLimiterFlags
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
			return fmt.Errorf("proto: Conversion: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Conversion: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Zrc20", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRateLimiterFlags
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
				return ErrInvalidLengthRateLimiterFlags
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRateLimiterFlags
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Zrc20 = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRateLimiterFlags
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
				return ErrInvalidLengthRateLimiterFlags
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRateLimiterFlags
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Rate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRateLimiterFlags(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRateLimiterFlags
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
func (m *AssetRate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRateLimiterFlags
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
			return fmt.Errorf("proto: AssetRate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AssetRate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			m.ChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRateLimiterFlags
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
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Asset", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRateLimiterFlags
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
				return ErrInvalidLengthRateLimiterFlags
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRateLimiterFlags
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Asset = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Decimals", wireType)
			}
			m.Decimals = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRateLimiterFlags
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Decimals |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CoinType", wireType)
			}
			m.CoinType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRateLimiterFlags
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CoinType |= coin.CoinType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRateLimiterFlags
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
				return ErrInvalidLengthRateLimiterFlags
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRateLimiterFlags
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Rate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRateLimiterFlags(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRateLimiterFlags
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
func skipRateLimiterFlags(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRateLimiterFlags
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
					return 0, ErrIntOverflowRateLimiterFlags
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
					return 0, ErrIntOverflowRateLimiterFlags
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
				return 0, ErrInvalidLengthRateLimiterFlags
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupRateLimiterFlags
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthRateLimiterFlags
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthRateLimiterFlags        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRateLimiterFlags          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupRateLimiterFlags = fmt.Errorf("proto: unexpected end of group")
)
