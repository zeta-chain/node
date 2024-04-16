// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: crosschain/out_tx_tracker.proto

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

type TxHashList struct {
	TxHash   string `protobuf:"bytes,1,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	TxSigner string `protobuf:"bytes,2,opt,name=tx_signer,json=txSigner,proto3" json:"tx_signer,omitempty"`
	Proved   bool   `protobuf:"varint,3,opt,name=proved,proto3" json:"proved,omitempty"`
}

func (m *TxHashList) Reset()         { *m = TxHashList{} }
func (m *TxHashList) String() string { return proto.CompactTextString(m) }
func (*TxHashList) ProtoMessage()    {}
func (*TxHashList) Descriptor() ([]byte, []int) {
	return fileDescriptor_5638c11005e4d36d, []int{0}
}
func (m *TxHashList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxHashList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxHashList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxHashList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxHashList.Merge(m, src)
}
func (m *TxHashList) XXX_Size() int {
	return m.Size()
}
func (m *TxHashList) XXX_DiscardUnknown() {
	xxx_messageInfo_TxHashList.DiscardUnknown(m)
}

var xxx_messageInfo_TxHashList proto.InternalMessageInfo

func (m *TxHashList) GetTxHash() string {
	if m != nil {
		return m.TxHash
	}
	return ""
}

func (m *TxHashList) GetTxSigner() string {
	if m != nil {
		return m.TxSigner
	}
	return ""
}

func (m *TxHashList) GetProved() bool {
	if m != nil {
		return m.Proved
	}
	return false
}

type OutTxTracker struct {
	Index    string        `protobuf:"bytes,1,opt,name=index,proto3" json:"index,omitempty"`
	ChainId  int64         `protobuf:"varint,2,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Nonce    uint64        `protobuf:"varint,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	HashList []*TxHashList `protobuf:"bytes,4,rep,name=hash_list,json=hashList,proto3" json:"hash_list,omitempty"`
}

func (m *OutTxTracker) Reset()         { *m = OutTxTracker{} }
func (m *OutTxTracker) String() string { return proto.CompactTextString(m) }
func (*OutTxTracker) ProtoMessage()    {}
func (*OutTxTracker) Descriptor() ([]byte, []int) {
	return fileDescriptor_5638c11005e4d36d, []int{1}
}
func (m *OutTxTracker) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OutTxTracker) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OutTxTracker.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OutTxTracker) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OutTxTracker.Merge(m, src)
}
func (m *OutTxTracker) XXX_Size() int {
	return m.Size()
}
func (m *OutTxTracker) XXX_DiscardUnknown() {
	xxx_messageInfo_OutTxTracker.DiscardUnknown(m)
}

var xxx_messageInfo_OutTxTracker proto.InternalMessageInfo

func (m *OutTxTracker) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *OutTxTracker) GetChainId() int64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

func (m *OutTxTracker) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *OutTxTracker) GetHashList() []*TxHashList {
	if m != nil {
		return m.HashList
	}
	return nil
}

func init() {
	proto.RegisterType((*TxHashList)(nil), "crosschain.TxHashList")
	proto.RegisterType((*OutTxTracker)(nil), "crosschain.OutTxTracker")
}

func init() { proto.RegisterFile("crosschain/out_tx_tracker.proto", fileDescriptor_5638c11005e4d36d) }

var fileDescriptor_5638c11005e4d36d = []byte{
	// 288 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xb1, 0x4e, 0xc3, 0x30,
	0x14, 0x45, 0x6b, 0x5a, 0xda, 0xd4, 0x30, 0x59, 0xa8, 0x04, 0x21, 0x99, 0xaa, 0x53, 0x17, 0x12,
	0x41, 0xff, 0x80, 0x09, 0x04, 0x12, 0x52, 0xe8, 0xd4, 0xc5, 0x4a, 0x1d, 0xab, 0xb6, 0x80, 0x38,
	0xb2, 0x5f, 0x90, 0xe1, 0x13, 0x98, 0xf8, 0x2c, 0xc6, 0x8e, 0x8c, 0x28, 0xf9, 0x11, 0x14, 0x27,
	0x52, 0xd8, 0xde, 0xd1, 0xb3, 0xef, 0xbb, 0xf7, 0xe2, 0x0b, 0x6e, 0xb4, 0xb5, 0x5c, 0xa6, 0x2a,
	0x8f, 0x75, 0x09, 0x0c, 0x1c, 0x03, 0x93, 0xf2, 0x67, 0x61, 0xa2, 0xc2, 0x68, 0xd0, 0x04, 0xf7,
	0x0f, 0x16, 0x1b, 0x8c, 0xd7, 0xee, 0x36, 0xb5, 0xf2, 0x41, 0x59, 0x20, 0xa7, 0x78, 0x02, 0x8e,
	0xc9, 0xd4, 0xca, 0x10, 0xcd, 0xd1, 0x72, 0x9a, 0x8c, 0xc1, 0x2f, 0xc9, 0x39, 0x9e, 0x82, 0x63,
	0x56, 0xed, 0x72, 0x61, 0xc2, 0x03, 0xbf, 0x0a, 0xc0, 0x3d, 0x79, 0x26, 0x33, 0x3c, 0x2e, 0x8c,
	0x7e, 0x13, 0x59, 0x38, 0x9c, 0xa3, 0x65, 0x90, 0x74, 0xb4, 0xf8, 0x44, 0xf8, 0xf8, 0xb1, 0x84,
	0xb5, 0x5b, 0xb7, 0xe7, 0xc9, 0x09, 0x3e, 0x54, 0x79, 0x26, 0x5c, 0x27, 0xde, 0x02, 0x39, 0xc3,
	0x81, 0xf7, 0xc2, 0x54, 0xe6, 0xa5, 0x87, 0xc9, 0xc4, 0xf3, 0x5d, 0xd6, 0x7c, 0xc8, 0x75, 0xce,
	0x85, 0x17, 0x1e, 0x25, 0x2d, 0x90, 0x15, 0x9e, 0x36, 0x16, 0xd9, 0x8b, 0xb2, 0x10, 0x8e, 0xe6,
	0xc3, 0xe5, 0xd1, 0xf5, 0x2c, 0xea, 0x33, 0x45, 0x7d, 0xa0, 0x24, 0x90, 0xdd, 0x74, 0x73, 0xff,
	0x5d, 0x51, 0xb4, 0xaf, 0x28, 0xfa, 0xad, 0x28, 0xfa, 0xaa, 0xe9, 0x60, 0x5f, 0xd3, 0xc1, 0x4f,
	0x4d, 0x07, 0x9b, 0xab, 0x9d, 0x02, 0x59, 0x6e, 0x23, 0xae, 0x5f, 0xe3, 0x0f, 0x01, 0xe9, 0x65,
	0x5b, 0x5d, 0x33, 0x72, 0x6d, 0x44, 0xec, 0xe2, 0x7f, 0x85, 0xc2, 0x7b, 0x21, 0xec, 0x76, 0xec,
	0x8b, 0x5c, 0xfd, 0x05, 0x00, 0x00, 0xff, 0xff, 0x89, 0xfa, 0x69, 0x30, 0x6b, 0x01, 0x00, 0x00,
}

func (m *TxHashList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxHashList) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxHashList) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Proved {
		i--
		if m.Proved {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if len(m.TxSigner) > 0 {
		i -= len(m.TxSigner)
		copy(dAtA[i:], m.TxSigner)
		i = encodeVarintOutTxTracker(dAtA, i, uint64(len(m.TxSigner)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.TxHash) > 0 {
		i -= len(m.TxHash)
		copy(dAtA[i:], m.TxHash)
		i = encodeVarintOutTxTracker(dAtA, i, uint64(len(m.TxHash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *OutTxTracker) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OutTxTracker) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OutTxTracker) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.HashList) > 0 {
		for iNdEx := len(m.HashList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.HashList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintOutTxTracker(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.Nonce != 0 {
		i = encodeVarintOutTxTracker(dAtA, i, uint64(m.Nonce))
		i--
		dAtA[i] = 0x18
	}
	if m.ChainId != 0 {
		i = encodeVarintOutTxTracker(dAtA, i, uint64(m.ChainId))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Index) > 0 {
		i -= len(m.Index)
		copy(dAtA[i:], m.Index)
		i = encodeVarintOutTxTracker(dAtA, i, uint64(len(m.Index)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintOutTxTracker(dAtA []byte, offset int, v uint64) int {
	offset -= sovOutTxTracker(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TxHashList) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.TxHash)
	if l > 0 {
		n += 1 + l + sovOutTxTracker(uint64(l))
	}
	l = len(m.TxSigner)
	if l > 0 {
		n += 1 + l + sovOutTxTracker(uint64(l))
	}
	if m.Proved {
		n += 2
	}
	return n
}

func (m *OutTxTracker) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Index)
	if l > 0 {
		n += 1 + l + sovOutTxTracker(uint64(l))
	}
	if m.ChainId != 0 {
		n += 1 + sovOutTxTracker(uint64(m.ChainId))
	}
	if m.Nonce != 0 {
		n += 1 + sovOutTxTracker(uint64(m.Nonce))
	}
	if len(m.HashList) > 0 {
		for _, e := range m.HashList {
			l = e.Size()
			n += 1 + l + sovOutTxTracker(uint64(l))
		}
	}
	return n
}

func sovOutTxTracker(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozOutTxTracker(x uint64) (n int) {
	return sovOutTxTracker(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TxHashList) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOutTxTracker
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
			return fmt.Errorf("proto: TxHashList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxHashList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutTxTracker
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
				return ErrInvalidLengthOutTxTracker
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOutTxTracker
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxSigner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutTxTracker
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
				return ErrInvalidLengthOutTxTracker
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOutTxTracker
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxSigner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Proved", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutTxTracker
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
			m.Proved = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipOutTxTracker(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthOutTxTracker
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
func (m *OutTxTracker) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOutTxTracker
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
			return fmt.Errorf("proto: OutTxTracker: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OutTxTracker: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutTxTracker
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
				return ErrInvalidLengthOutTxTracker
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOutTxTracker
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Index = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			m.ChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutTxTracker
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
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			m.Nonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutTxTracker
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Nonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HashList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutTxTracker
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
				return ErrInvalidLengthOutTxTracker
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOutTxTracker
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HashList = append(m.HashList, &TxHashList{})
			if err := m.HashList[len(m.HashList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipOutTxTracker(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthOutTxTracker
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
func skipOutTxTracker(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowOutTxTracker
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
					return 0, ErrIntOverflowOutTxTracker
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
					return 0, ErrIntOverflowOutTxTracker
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
				return 0, ErrInvalidLengthOutTxTracker
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupOutTxTracker
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthOutTxTracker
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthOutTxTracker        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowOutTxTracker          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupOutTxTracker = fmt.Errorf("proto: unexpected end of group")
)
