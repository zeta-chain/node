// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: zetachain/zetacore/pkg/proofs/proofs.proto

package proofs

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	bitcoin "github.com/zeta-chain/node/pkg/proofs/bitcoin"
	ethereum "github.com/zeta-chain/node/pkg/proofs/ethereum"
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

type BlockHeader struct {
	Height     int64  `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Hash       []byte `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	ParentHash []byte `protobuf:"bytes,3,opt,name=parent_hash,json=parentHash,proto3" json:"parent_hash,omitempty"`
	ChainId    int64  `protobuf:"varint,4,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	// chain specific header
	Header HeaderData `protobuf:"bytes,5,opt,name=header,proto3" json:"header"`
}

func (m *BlockHeader) Reset()         { *m = BlockHeader{} }
func (m *BlockHeader) String() string { return proto.CompactTextString(m) }
func (*BlockHeader) ProtoMessage()    {}
func (*BlockHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_874830d2276ded66, []int{0}
}
func (m *BlockHeader) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BlockHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BlockHeader.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BlockHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockHeader.Merge(m, src)
}
func (m *BlockHeader) XXX_Size() int {
	return m.Size()
}
func (m *BlockHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockHeader.DiscardUnknown(m)
}

var xxx_messageInfo_BlockHeader proto.InternalMessageInfo

func (m *BlockHeader) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *BlockHeader) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *BlockHeader) GetParentHash() []byte {
	if m != nil {
		return m.ParentHash
	}
	return nil
}

func (m *BlockHeader) GetChainId() int64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

func (m *BlockHeader) GetHeader() HeaderData {
	if m != nil {
		return m.Header
	}
	return HeaderData{}
}

type HeaderData struct {
	// Types that are valid to be assigned to Data:
	//
	//	*HeaderData_EthereumHeader
	//	*HeaderData_BitcoinHeader
	Data isHeaderData_Data `protobuf_oneof:"data"`
}

func (m *HeaderData) Reset()         { *m = HeaderData{} }
func (m *HeaderData) String() string { return proto.CompactTextString(m) }
func (*HeaderData) ProtoMessage()    {}
func (*HeaderData) Descriptor() ([]byte, []int) {
	return fileDescriptor_874830d2276ded66, []int{1}
}
func (m *HeaderData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HeaderData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HeaderData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HeaderData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeaderData.Merge(m, src)
}
func (m *HeaderData) XXX_Size() int {
	return m.Size()
}
func (m *HeaderData) XXX_DiscardUnknown() {
	xxx_messageInfo_HeaderData.DiscardUnknown(m)
}

var xxx_messageInfo_HeaderData proto.InternalMessageInfo

type isHeaderData_Data interface {
	isHeaderData_Data()
	MarshalTo([]byte) (int, error)
	Size() int
}

type HeaderData_EthereumHeader struct {
	EthereumHeader []byte `protobuf:"bytes,1,opt,name=ethereum_header,json=ethereumHeader,proto3,oneof" json:"ethereum_header,omitempty"`
}
type HeaderData_BitcoinHeader struct {
	BitcoinHeader []byte `protobuf:"bytes,2,opt,name=bitcoin_header,json=bitcoinHeader,proto3,oneof" json:"bitcoin_header,omitempty"`
}

func (*HeaderData_EthereumHeader) isHeaderData_Data() {}
func (*HeaderData_BitcoinHeader) isHeaderData_Data()  {}

func (m *HeaderData) GetData() isHeaderData_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *HeaderData) GetEthereumHeader() []byte {
	if x, ok := m.GetData().(*HeaderData_EthereumHeader); ok {
		return x.EthereumHeader
	}
	return nil
}

func (m *HeaderData) GetBitcoinHeader() []byte {
	if x, ok := m.GetData().(*HeaderData_BitcoinHeader); ok {
		return x.BitcoinHeader
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*HeaderData) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*HeaderData_EthereumHeader)(nil),
		(*HeaderData_BitcoinHeader)(nil),
	}
}

type Proof struct {
	// Types that are valid to be assigned to Proof:
	//
	//	*Proof_EthereumProof
	//	*Proof_BitcoinProof
	Proof isProof_Proof `protobuf_oneof:"proof"`
}

func (m *Proof) Reset()         { *m = Proof{} }
func (m *Proof) String() string { return proto.CompactTextString(m) }
func (*Proof) ProtoMessage()    {}
func (*Proof) Descriptor() ([]byte, []int) {
	return fileDescriptor_874830d2276ded66, []int{2}
}
func (m *Proof) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Proof) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Proof.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Proof) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Proof.Merge(m, src)
}
func (m *Proof) XXX_Size() int {
	return m.Size()
}
func (m *Proof) XXX_DiscardUnknown() {
	xxx_messageInfo_Proof.DiscardUnknown(m)
}

var xxx_messageInfo_Proof proto.InternalMessageInfo

type isProof_Proof interface {
	isProof_Proof()
	MarshalTo([]byte) (int, error)
	Size() int
}

type Proof_EthereumProof struct {
	EthereumProof *ethereum.Proof `protobuf:"bytes,1,opt,name=ethereum_proof,json=ethereumProof,proto3,oneof" json:"ethereum_proof,omitempty"`
}
type Proof_BitcoinProof struct {
	BitcoinProof *bitcoin.Proof `protobuf:"bytes,2,opt,name=bitcoin_proof,json=bitcoinProof,proto3,oneof" json:"bitcoin_proof,omitempty"`
}

func (*Proof_EthereumProof) isProof_Proof() {}
func (*Proof_BitcoinProof) isProof_Proof()  {}

func (m *Proof) GetProof() isProof_Proof {
	if m != nil {
		return m.Proof
	}
	return nil
}

func (m *Proof) GetEthereumProof() *ethereum.Proof {
	if x, ok := m.GetProof().(*Proof_EthereumProof); ok {
		return x.EthereumProof
	}
	return nil
}

func (m *Proof) GetBitcoinProof() *bitcoin.Proof {
	if x, ok := m.GetProof().(*Proof_BitcoinProof); ok {
		return x.BitcoinProof
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Proof) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Proof_EthereumProof)(nil),
		(*Proof_BitcoinProof)(nil),
	}
}

func init() {
	proto.RegisterType((*BlockHeader)(nil), "zetachain.zetacore.pkg.proofs.BlockHeader")
	proto.RegisterType((*HeaderData)(nil), "zetachain.zetacore.pkg.proofs.HeaderData")
	proto.RegisterType((*Proof)(nil), "zetachain.zetacore.pkg.proofs.Proof")
}

func init() {
	proto.RegisterFile("zetachain/zetacore/pkg/proofs/proofs.proto", fileDescriptor_874830d2276ded66)
}

var fileDescriptor_874830d2276ded66 = []byte{
	// 398 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0x41, 0xab, 0xd3, 0x40,
	0x10, 0xce, 0xf6, 0xa5, 0x7d, 0x32, 0xa9, 0x4f, 0x58, 0x44, 0xe2, 0x03, 0xd3, 0x52, 0x10, 0x5b,
	0xb1, 0x09, 0xb4, 0x78, 0x16, 0x82, 0x60, 0xbd, 0x49, 0x04, 0x0f, 0x5e, 0xc2, 0x36, 0x59, 0xb3,
	0xa1, 0x36, 0x1b, 0xd2, 0xed, 0xc5, 0x5f, 0xe1, 0x1f, 0xd2, 0x73, 0x8f, 0x3d, 0x7a, 0x12, 0x69,
	0xff, 0x88, 0x64, 0x76, 0x37, 0x7a, 0x6a, 0x4f, 0x3b, 0x3b, 0xf3, 0xcd, 0xf7, 0xcd, 0x0c, 0x1f,
	0xbc, 0xfc, 0xc6, 0x15, 0xcb, 0x04, 0x2b, 0xab, 0x08, 0x23, 0xd9, 0xf0, 0xa8, 0xde, 0x14, 0x51,
	0xdd, 0x48, 0xf9, 0x65, 0x67, 0x9e, 0xb0, 0x6e, 0xa4, 0x92, 0xf4, 0x59, 0x87, 0x0d, 0x2d, 0x36,
	0xac, 0x37, 0x45, 0xa8, 0x41, 0xf7, 0x8f, 0x0b, 0x59, 0x48, 0x44, 0x46, 0x6d, 0xa4, 0x9b, 0xee,
	0x97, 0x97, 0x05, 0xd6, 0xa5, 0xca, 0x64, 0x59, 0xd9, 0xd7, 0x34, 0xbd, 0xbe, 0xdc, 0xc4, 0x95,
	0xe0, 0x0d, 0xdf, 0x6f, 0xbb, 0x40, 0xb7, 0x4d, 0x7e, 0x12, 0xf0, 0xe2, 0xaf, 0x32, 0xdb, 0xac,
	0x38, 0xcb, 0x79, 0x43, 0x9f, 0xc0, 0x40, 0xf0, 0xb2, 0x10, 0xca, 0x27, 0x63, 0x32, 0xbd, 0x49,
	0xcc, 0x8f, 0x52, 0x70, 0x05, 0xdb, 0x09, 0xbf, 0x37, 0x26, 0xd3, 0x61, 0x82, 0x31, 0x1d, 0x81,
	0x57, 0xb3, 0x86, 0x57, 0x2a, 0xc5, 0xd2, 0x0d, 0x96, 0x40, 0xa7, 0x56, 0x2d, 0xe0, 0x29, 0x3c,
	0xc0, 0x89, 0xd2, 0x32, 0xf7, 0x5d, 0xa4, 0xbb, 0xc5, 0xff, 0xfb, 0x9c, 0xbe, 0x6b, 0x75, 0x5a,
	0x45, 0xbf, 0x3f, 0x26, 0x53, 0x6f, 0x31, 0x0b, 0x2f, 0x5e, 0x2a, 0xd4, 0xe3, 0xbd, 0x65, 0x8a,
	0xc5, 0xee, 0xe1, 0xf7, 0xc8, 0x49, 0x4c, 0xfb, 0x44, 0x00, 0xfc, 0xab, 0xd1, 0x19, 0x3c, 0xb2,
	0x0b, 0xa6, 0x86, 0xbf, 0xdd, 0x63, 0xb8, 0x72, 0x92, 0x3b, 0x5b, 0x30, 0x9b, 0xbe, 0x80, 0x3b,
	0x73, 0x41, 0x8b, 0xec, 0x19, 0xe4, 0x43, 0x93, 0xd7, 0xc0, 0x78, 0x00, 0x6e, 0xce, 0x14, 0x9b,
	0xfc, 0x20, 0xd0, 0xff, 0xd0, 0x4e, 0x43, 0x3f, 0x41, 0x47, 0x96, 0xe2, 0x7c, 0x28, 0xe2, 0x2d,
	0xe6, 0x57, 0x96, 0xe8, 0x6e, 0x8f, 0x34, 0xad, 0x92, 0xcd, 0x68, 0xde, 0x8f, 0x60, 0xa5, 0x0d,
	0x6d, 0x0f, 0x69, 0x5f, 0x5d, 0xa1, 0xb5, 0x46, 0xb0, 0xac, 0x43, 0x93, 0xc0, 0x7f, 0x7c, 0x0b,
	0x7d, 0xc4, 0xc5, 0x6f, 0x0e, 0xa7, 0x80, 0x1c, 0x4f, 0x01, 0xf9, 0x73, 0x0a, 0xc8, 0xf7, 0x73,
	0xe0, 0x1c, 0xcf, 0x81, 0xf3, 0xeb, 0x1c, 0x38, 0x9f, 0x9f, 0x17, 0xa5, 0x12, 0xfb, 0x75, 0x98,
	0xc9, 0x2d, 0x9a, 0x67, 0xae, 0x7d, 0x54, 0xc9, 0xfc, 0x7f, 0x0f, 0xad, 0x07, 0x68, 0x99, 0xe5,
	0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xea, 0x34, 0x1b, 0x28, 0x01, 0x03, 0x00, 0x00,
}

func (m *BlockHeader) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BlockHeader) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BlockHeader) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Header.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintProofs(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if m.ChainId != 0 {
		i = encodeVarintProofs(dAtA, i, uint64(m.ChainId))
		i--
		dAtA[i] = 0x20
	}
	if len(m.ParentHash) > 0 {
		i -= len(m.ParentHash)
		copy(dAtA[i:], m.ParentHash)
		i = encodeVarintProofs(dAtA, i, uint64(len(m.ParentHash)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintProofs(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0x12
	}
	if m.Height != 0 {
		i = encodeVarintProofs(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *HeaderData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HeaderData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HeaderData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Data != nil {
		{
			size := m.Data.Size()
			i -= size
			if _, err := m.Data.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	return len(dAtA) - i, nil
}

func (m *HeaderData_EthereumHeader) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HeaderData_EthereumHeader) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.EthereumHeader != nil {
		i -= len(m.EthereumHeader)
		copy(dAtA[i:], m.EthereumHeader)
		i = encodeVarintProofs(dAtA, i, uint64(len(m.EthereumHeader)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func (m *HeaderData_BitcoinHeader) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HeaderData_BitcoinHeader) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.BitcoinHeader != nil {
		i -= len(m.BitcoinHeader)
		copy(dAtA[i:], m.BitcoinHeader)
		i = encodeVarintProofs(dAtA, i, uint64(len(m.BitcoinHeader)))
		i--
		dAtA[i] = 0x12
	}
	return len(dAtA) - i, nil
}
func (m *Proof) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Proof) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Proof) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Proof != nil {
		{
			size := m.Proof.Size()
			i -= size
			if _, err := m.Proof.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	return len(dAtA) - i, nil
}

func (m *Proof_EthereumProof) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Proof_EthereumProof) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.EthereumProof != nil {
		{
			size, err := m.EthereumProof.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProofs(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func (m *Proof_BitcoinProof) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Proof_BitcoinProof) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.BitcoinProof != nil {
		{
			size, err := m.BitcoinProof.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProofs(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	return len(dAtA) - i, nil
}
func encodeVarintProofs(dAtA []byte, offset int, v uint64) int {
	offset -= sovProofs(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BlockHeader) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovProofs(uint64(m.Height))
	}
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovProofs(uint64(l))
	}
	l = len(m.ParentHash)
	if l > 0 {
		n += 1 + l + sovProofs(uint64(l))
	}
	if m.ChainId != 0 {
		n += 1 + sovProofs(uint64(m.ChainId))
	}
	l = m.Header.Size()
	n += 1 + l + sovProofs(uint64(l))
	return n
}

func (m *HeaderData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Data != nil {
		n += m.Data.Size()
	}
	return n
}

func (m *HeaderData_EthereumHeader) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.EthereumHeader != nil {
		l = len(m.EthereumHeader)
		n += 1 + l + sovProofs(uint64(l))
	}
	return n
}
func (m *HeaderData_BitcoinHeader) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.BitcoinHeader != nil {
		l = len(m.BitcoinHeader)
		n += 1 + l + sovProofs(uint64(l))
	}
	return n
}
func (m *Proof) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Proof != nil {
		n += m.Proof.Size()
	}
	return n
}

func (m *Proof_EthereumProof) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.EthereumProof != nil {
		l = m.EthereumProof.Size()
		n += 1 + l + sovProofs(uint64(l))
	}
	return n
}
func (m *Proof_BitcoinProof) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.BitcoinProof != nil {
		l = m.BitcoinProof.Size()
		n += 1 + l + sovProofs(uint64(l))
	}
	return n
}

func sovProofs(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozProofs(x uint64) (n int) {
	return sovProofs(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BlockHeader) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProofs
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
			return fmt.Errorf("proto: BlockHeader: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BlockHeader: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProofs
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProofs
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthProofs
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthProofs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = append(m.Hash[:0], dAtA[iNdEx:postIndex]...)
			if m.Hash == nil {
				m.Hash = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ParentHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProofs
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthProofs
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthProofs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ParentHash = append(m.ParentHash[:0], dAtA[iNdEx:postIndex]...)
			if m.ParentHash == nil {
				m.ParentHash = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			m.ChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProofs
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
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Header", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProofs
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
				return ErrInvalidLengthProofs
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProofs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Header.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProofs(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProofs
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
func (m *HeaderData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProofs
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
			return fmt.Errorf("proto: HeaderData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HeaderData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EthereumHeader", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProofs
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthProofs
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthProofs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := make([]byte, postIndex-iNdEx)
			copy(v, dAtA[iNdEx:postIndex])
			m.Data = &HeaderData_EthereumHeader{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BitcoinHeader", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProofs
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthProofs
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthProofs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := make([]byte, postIndex-iNdEx)
			copy(v, dAtA[iNdEx:postIndex])
			m.Data = &HeaderData_BitcoinHeader{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProofs(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProofs
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
func (m *Proof) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProofs
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
			return fmt.Errorf("proto: Proof: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Proof: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EthereumProof", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProofs
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
				return ErrInvalidLengthProofs
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProofs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &ethereum.Proof{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Proof = &Proof_EthereumProof{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BitcoinProof", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProofs
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
				return ErrInvalidLengthProofs
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProofs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &bitcoin.Proof{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Proof = &Proof_BitcoinProof{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProofs(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProofs
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
func skipProofs(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProofs
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
					return 0, ErrIntOverflowProofs
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
					return 0, ErrIntOverflowProofs
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
				return 0, ErrInvalidLengthProofs
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupProofs
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthProofs
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthProofs        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProofs          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupProofs = fmt.Errorf("proto: unexpected end of group")
)
