// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: zetachain/zetacore/observer/node_account.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	crypto "github.com/zeta-chain/node/pkg/crypto"
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

type NodeStatus int32

const (
	NodeStatus_Unknown     NodeStatus = 0
	NodeStatus_Whitelisted NodeStatus = 1
	NodeStatus_Standby     NodeStatus = 2
	NodeStatus_Ready       NodeStatus = 3
	NodeStatus_Active      NodeStatus = 4
	NodeStatus_Disabled    NodeStatus = 5
)

var NodeStatus_name = map[int32]string{
	0: "Unknown",
	1: "Whitelisted",
	2: "Standby",
	3: "Ready",
	4: "Active",
	5: "Disabled",
}

var NodeStatus_value = map[string]int32{
	"Unknown":     0,
	"Whitelisted": 1,
	"Standby":     2,
	"Ready":       3,
	"Active":      4,
	"Disabled":    5,
}

func (x NodeStatus) String() string {
	return proto.EnumName(NodeStatus_name, int32(x))
}

func (NodeStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_67bb97178fb2bc84, []int{0}
}

type NodeAccount struct {
	Operator       string            `protobuf:"bytes,1,opt,name=operator,proto3" json:"operator,omitempty"`
	GranteeAddress string            `protobuf:"bytes,2,opt,name=granteeAddress,proto3" json:"granteeAddress,omitempty"`
	GranteePubkey  *crypto.PubKeySet `protobuf:"bytes,3,opt,name=granteePubkey,proto3" json:"granteePubkey,omitempty"`
	NodeStatus     NodeStatus        `protobuf:"varint,4,opt,name=nodeStatus,proto3,enum=zetachain.zetacore.observer.NodeStatus" json:"nodeStatus,omitempty"`
}

func (m *NodeAccount) Reset()         { *m = NodeAccount{} }
func (m *NodeAccount) String() string { return proto.CompactTextString(m) }
func (*NodeAccount) ProtoMessage()    {}
func (*NodeAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_67bb97178fb2bc84, []int{0}
}
func (m *NodeAccount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NodeAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NodeAccount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NodeAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeAccount.Merge(m, src)
}
func (m *NodeAccount) XXX_Size() int {
	return m.Size()
}
func (m *NodeAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeAccount.DiscardUnknown(m)
}

var xxx_messageInfo_NodeAccount proto.InternalMessageInfo

func (m *NodeAccount) GetOperator() string {
	if m != nil {
		return m.Operator
	}
	return ""
}

func (m *NodeAccount) GetGranteeAddress() string {
	if m != nil {
		return m.GranteeAddress
	}
	return ""
}

func (m *NodeAccount) GetGranteePubkey() *crypto.PubKeySet {
	if m != nil {
		return m.GranteePubkey
	}
	return nil
}

func (m *NodeAccount) GetNodeStatus() NodeStatus {
	if m != nil {
		return m.NodeStatus
	}
	return NodeStatus_Unknown
}

func init() {
	proto.RegisterEnum("zetachain.zetacore.observer.NodeStatus", NodeStatus_name, NodeStatus_value)
	proto.RegisterType((*NodeAccount)(nil), "zetachain.zetacore.observer.NodeAccount")
}

func init() {
	proto.RegisterFile("zetachain/zetacore/observer/node_account.proto", fileDescriptor_67bb97178fb2bc84)
}

var fileDescriptor_67bb97178fb2bc84 = []byte{
	// 379 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xcd, 0xaa, 0xd3, 0x40,
	0x14, 0xc7, 0x33, 0xfd, 0xb2, 0x9d, 0x68, 0x0d, 0x83, 0x8b, 0x10, 0x21, 0x14, 0x17, 0x1a, 0x2a,
	0x4e, 0xa0, 0x3e, 0x41, 0x45, 0x71, 0x21, 0x94, 0x92, 0x22, 0x82, 0x1b, 0x99, 0xc9, 0x1c, 0xd2,
	0x90, 0x3a, 0x13, 0x26, 0x93, 0x6a, 0x7c, 0x0a, 0x1f, 0xc2, 0x85, 0x8f, 0xe2, 0xb2, 0x4b, 0x97,
	0xd2, 0xee, 0x7c, 0x0a, 0x49, 0xd2, 0x0f, 0xef, 0xa5, 0xdc, 0xd5, 0x9c, 0x39, 0xe7, 0xf7, 0xe7,
	0x9c, 0x39, 0xff, 0xc1, 0xf4, 0x1b, 0x18, 0x16, 0xaf, 0x59, 0x2a, 0xc3, 0x26, 0x52, 0x1a, 0x42,
	0xc5, 0x0b, 0xd0, 0x5b, 0xd0, 0xa1, 0x54, 0x02, 0x3e, 0xb1, 0x38, 0x56, 0xa5, 0x34, 0x34, 0xd7,
	0xca, 0x28, 0xf2, 0xf8, 0xcc, 0xd3, 0x13, 0x4f, 0x4f, 0xbc, 0xf7, 0x28, 0x51, 0x89, 0x6a, 0xb8,
	0xb0, 0x8e, 0x5a, 0x89, 0x37, 0xbd, 0xd2, 0x22, 0xcf, 0x92, 0x30, 0xd6, 0x55, 0x6e, 0xd4, 0xf1,
	0x68, 0xd9, 0x27, 0x7f, 0x11, 0xb6, 0x17, 0x4a, 0xc0, 0xbc, 0x6d, 0x4a, 0x3c, 0x3c, 0x54, 0x39,
	0x68, 0x66, 0x94, 0x76, 0xd1, 0x04, 0x05, 0xa3, 0xe8, 0x7c, 0x27, 0x4f, 0xf1, 0x38, 0xd1, 0x4c,
	0x1a, 0x80, 0xb9, 0x10, 0x1a, 0x8a, 0xc2, 0xed, 0x34, 0xc4, 0xad, 0x2c, 0x59, 0xe0, 0x07, 0xc7,
	0xcc, 0xb2, 0xe4, 0x19, 0x54, 0x6e, 0x77, 0x82, 0x02, 0x7b, 0x16, 0xd0, 0x2b, 0x4f, 0xc9, 0xb3,
	0x84, 0x1e, 0x07, 0x5a, 0x96, 0xfc, 0x1d, 0x54, 0x2b, 0x30, 0xd1, 0x4d, 0x39, 0x79, 0x8b, 0x71,
	0xbd, 0x98, 0x95, 0x61, 0xa6, 0x2c, 0xdc, 0xde, 0x04, 0x05, 0xe3, 0xd9, 0x33, 0x7a, 0xc7, 0x5e,
	0xe8, 0xe2, 0x8c, 0x47, 0xff, 0x49, 0xa7, 0x1c, 0xe3, 0x4b, 0x85, 0xd8, 0xf8, 0xde, 0x7b, 0x99,
	0x49, 0xf5, 0x45, 0x3a, 0x16, 0x79, 0x88, 0xed, 0x0f, 0xeb, 0xd4, 0xc0, 0x26, 0x2d, 0x0c, 0x08,
	0x07, 0xd5, 0xd5, 0x95, 0x61, 0x52, 0xf0, 0xca, 0xe9, 0x90, 0x11, 0xee, 0x47, 0xc0, 0x44, 0xe5,
	0x74, 0x09, 0xc6, 0x83, 0x79, 0x6c, 0xd2, 0x2d, 0x38, 0x3d, 0x72, 0x1f, 0x0f, 0x5f, 0xa7, 0x05,
	0xe3, 0x1b, 0x10, 0x4e, 0xdf, 0xeb, 0xfd, 0xfc, 0xe1, 0xa3, 0x57, 0x6f, 0x7e, 0xed, 0x7d, 0xb4,
	0xdb, 0xfb, 0xe8, 0xcf, 0xde, 0x47, 0xdf, 0x0f, 0xbe, 0xb5, 0x3b, 0xf8, 0xd6, 0xef, 0x83, 0x6f,
	0x7d, 0x7c, 0x9e, 0xa4, 0x66, 0x5d, 0x72, 0x1a, 0xab, 0xcf, 0x8d, 0x2f, 0x2f, 0x5a, 0x8b, 0xea,
	0xf9, 0xc2, 0xaf, 0x97, 0x3f, 0x60, 0xaa, 0x1c, 0x0a, 0x3e, 0x68, 0xec, 0x79, 0xf9, 0x2f, 0x00,
	0x00, 0xff, 0xff, 0x76, 0x26, 0x93, 0x70, 0x2f, 0x02, 0x00, 0x00,
}

func (m *NodeAccount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NodeAccount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *NodeAccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.NodeStatus != 0 {
		i = encodeVarintNodeAccount(dAtA, i, uint64(m.NodeStatus))
		i--
		dAtA[i] = 0x20
	}
	if m.GranteePubkey != nil {
		{
			size, err := m.GranteePubkey.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintNodeAccount(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.GranteeAddress) > 0 {
		i -= len(m.GranteeAddress)
		copy(dAtA[i:], m.GranteeAddress)
		i = encodeVarintNodeAccount(dAtA, i, uint64(len(m.GranteeAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Operator) > 0 {
		i -= len(m.Operator)
		copy(dAtA[i:], m.Operator)
		i = encodeVarintNodeAccount(dAtA, i, uint64(len(m.Operator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintNodeAccount(dAtA []byte, offset int, v uint64) int {
	offset -= sovNodeAccount(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *NodeAccount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Operator)
	if l > 0 {
		n += 1 + l + sovNodeAccount(uint64(l))
	}
	l = len(m.GranteeAddress)
	if l > 0 {
		n += 1 + l + sovNodeAccount(uint64(l))
	}
	if m.GranteePubkey != nil {
		l = m.GranteePubkey.Size()
		n += 1 + l + sovNodeAccount(uint64(l))
	}
	if m.NodeStatus != 0 {
		n += 1 + sovNodeAccount(uint64(m.NodeStatus))
	}
	return n
}

func sovNodeAccount(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozNodeAccount(x uint64) (n int) {
	return sovNodeAccount(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *NodeAccount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNodeAccount
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
			return fmt.Errorf("proto: NodeAccount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NodeAccount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Operator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNodeAccount
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
				return ErrInvalidLengthNodeAccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthNodeAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Operator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GranteeAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNodeAccount
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
				return ErrInvalidLengthNodeAccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthNodeAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GranteeAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GranteePubkey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNodeAccount
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
				return ErrInvalidLengthNodeAccount
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthNodeAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.GranteePubkey == nil {
				m.GranteePubkey = &crypto.PubKeySet{}
			}
			if err := m.GranteePubkey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NodeStatus", wireType)
			}
			m.NodeStatus = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNodeAccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NodeStatus |= NodeStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipNodeAccount(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthNodeAccount
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
func skipNodeAccount(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowNodeAccount
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
					return 0, ErrIntOverflowNodeAccount
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
					return 0, ErrIntOverflowNodeAccount
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
				return 0, ErrInvalidLengthNodeAccount
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupNodeAccount
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthNodeAccount
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthNodeAccount        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowNodeAccount          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupNodeAccount = fmt.Errorf("proto: unexpected end of group")
)
