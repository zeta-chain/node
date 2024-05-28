// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: zetachain/zetacore/authority/authorization.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
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

type Authorization struct {
	// The URL of the message that needs to be authorized
	MsgUrl string `protobuf:"bytes,1,opt,name=msg_url,json=msgUrl,proto3" json:"msg_url,omitempty"`
	// The policy that is authorized to access the message
	AuthorizedPolicy PolicyType `protobuf:"varint,2,opt,name=authorized_policy,json=authorizedPolicy,proto3,enum=zetachain.zetacore.authority.PolicyType" json:"authorized_policy,omitempty"`
}

func (m *Authorization) Reset()         { *m = Authorization{} }
func (m *Authorization) String() string { return proto.CompactTextString(m) }
func (*Authorization) ProtoMessage()    {}
func (*Authorization) Descriptor() ([]byte, []int) {
	return fileDescriptor_b7303e09de7c755a, []int{0}
}
func (m *Authorization) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Authorization) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Authorization.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Authorization) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Authorization.Merge(m, src)
}
func (m *Authorization) XXX_Size() int {
	return m.Size()
}
func (m *Authorization) XXX_DiscardUnknown() {
	xxx_messageInfo_Authorization.DiscardUnknown(m)
}

var xxx_messageInfo_Authorization proto.InternalMessageInfo

func (m *Authorization) GetMsgUrl() string {
	if m != nil {
		return m.MsgUrl
	}
	return ""
}

func (m *Authorization) GetAuthorizedPolicy() PolicyType {
	if m != nil {
		return m.AuthorizedPolicy
	}
	return PolicyType_groupEmergency
}

// AuthorizationList holds the list of authorizations on zetachain
type AuthorizationList struct {
	Authorizations []Authorization `protobuf:"bytes,1,rep,name=authorizations,proto3" json:"authorizations"`
}

func (m *AuthorizationList) Reset()         { *m = AuthorizationList{} }
func (m *AuthorizationList) String() string { return proto.CompactTextString(m) }
func (*AuthorizationList) ProtoMessage()    {}
func (*AuthorizationList) Descriptor() ([]byte, []int) {
	return fileDescriptor_b7303e09de7c755a, []int{1}
}
func (m *AuthorizationList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AuthorizationList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AuthorizationList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AuthorizationList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthorizationList.Merge(m, src)
}
func (m *AuthorizationList) XXX_Size() int {
	return m.Size()
}
func (m *AuthorizationList) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthorizationList.DiscardUnknown(m)
}

var xxx_messageInfo_AuthorizationList proto.InternalMessageInfo

func (m *AuthorizationList) GetAuthorizations() []Authorization {
	if m != nil {
		return m.Authorizations
	}
	return nil
}

func init() {
	proto.RegisterType((*Authorization)(nil), "zetachain.zetacore.authority.Authorization")
	proto.RegisterType((*AuthorizationList)(nil), "zetachain.zetacore.authority.AuthorizationList")
}

func init() {
	proto.RegisterFile("zetachain/zetacore/authority/authorization.proto", fileDescriptor_b7303e09de7c755a)
}

var fileDescriptor_b7303e09de7c755a = []byte{
	// 271 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0xa8, 0x4a, 0x2d, 0x49,
	0x4c, 0xce, 0x48, 0xcc, 0xcc, 0xd3, 0x07, 0xb3, 0xf2, 0x8b, 0x52, 0xf5, 0x13, 0x4b, 0x4b, 0x32,
	0xf2, 0x8b, 0x32, 0x4b, 0x2a, 0x61, 0xac, 0xaa, 0xc4, 0x92, 0xcc, 0xfc, 0x3c, 0xbd, 0x82, 0xa2,
	0xfc, 0x92, 0x7c, 0x21, 0x19, 0xb8, 0x0e, 0x3d, 0x98, 0x0e, 0x3d, 0xb8, 0x0e, 0x29, 0x91, 0xf4,
	0xfc, 0xf4, 0x7c, 0xb0, 0x42, 0x7d, 0x10, 0x0b, 0xa2, 0x47, 0x4a, 0x1b, 0xaf, 0x2d, 0x05, 0xf9,
	0x39, 0x99, 0xc9, 0x99, 0xa9, 0xc5, 0x10, 0xc5, 0x4a, 0xf5, 0x5c, 0xbc, 0x8e, 0xc8, 0xf6, 0x0a,
	0x89, 0x73, 0xb1, 0xe7, 0x16, 0xa7, 0xc7, 0x97, 0x16, 0xe5, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70,
	0x06, 0xb1, 0xe5, 0x16, 0xa7, 0x87, 0x16, 0xe5, 0x08, 0x85, 0x72, 0x09, 0xc2, 0x5c, 0x98, 0x9a,
	0x12, 0x0f, 0x36, 0xa6, 0x52, 0x82, 0x49, 0x81, 0x51, 0x83, 0xcf, 0x48, 0x43, 0x0f, 0x9f, 0x33,
	0xf5, 0x02, 0xc0, 0x6a, 0x43, 0x2a, 0x0b, 0x52, 0x83, 0x04, 0x10, 0x46, 0x40, 0x44, 0x95, 0xf2,
	0xb8, 0x04, 0x51, 0x1c, 0xe0, 0x93, 0x59, 0x5c, 0x22, 0x14, 0xc9, 0xc5, 0x87, 0x12, 0x1a, 0xc5,
	0x12, 0x8c, 0x0a, 0xcc, 0x1a, 0xdc, 0x46, 0xda, 0xf8, 0x2d, 0x42, 0x31, 0xc8, 0x89, 0xe5, 0xc4,
	0x3d, 0x79, 0x86, 0x20, 0x34, 0x83, 0x9c, 0xbc, 0x4e, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e,
	0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58,
	0x8e, 0x21, 0xca, 0x20, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x17, 0x1c, 0x70,
	0xba, 0x68, 0x61, 0x58, 0x81, 0x14, 0x8a, 0x25, 0x95, 0x05, 0xa9, 0xc5, 0x49, 0x6c, 0xe0, 0x30,
	0x34, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x88, 0x76, 0x6b, 0x6c, 0xd8, 0x01, 0x00, 0x00,
}

func (m *Authorization) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Authorization) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Authorization) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.AuthorizedPolicy != 0 {
		i = encodeVarintAuthorization(dAtA, i, uint64(m.AuthorizedPolicy))
		i--
		dAtA[i] = 0x10
	}
	if len(m.MsgUrl) > 0 {
		i -= len(m.MsgUrl)
		copy(dAtA[i:], m.MsgUrl)
		i = encodeVarintAuthorization(dAtA, i, uint64(len(m.MsgUrl)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *AuthorizationList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AuthorizationList) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AuthorizationList) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Authorizations) > 0 {
		for iNdEx := len(m.Authorizations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Authorizations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintAuthorization(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintAuthorization(dAtA []byte, offset int, v uint64) int {
	offset -= sovAuthorization(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Authorization) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.MsgUrl)
	if l > 0 {
		n += 1 + l + sovAuthorization(uint64(l))
	}
	if m.AuthorizedPolicy != 0 {
		n += 1 + sovAuthorization(uint64(m.AuthorizedPolicy))
	}
	return n
}

func (m *AuthorizationList) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Authorizations) > 0 {
		for _, e := range m.Authorizations {
			l = e.Size()
			n += 1 + l + sovAuthorization(uint64(l))
		}
	}
	return n
}

func sovAuthorization(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAuthorization(x uint64) (n int) {
	return sovAuthorization(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Authorization) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAuthorization
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
			return fmt.Errorf("proto: Authorization: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Authorization: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MsgUrl", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthorization
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
				return ErrInvalidLengthAuthorization
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthorization
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MsgUrl = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AuthorizedPolicy", wireType)
			}
			m.AuthorizedPolicy = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthorization
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AuthorizedPolicy |= PolicyType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipAuthorization(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAuthorization
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
func (m *AuthorizationList) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAuthorization
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
			return fmt.Errorf("proto: AuthorizationList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AuthorizationList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authorizations", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthorization
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
				return ErrInvalidLengthAuthorization
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAuthorization
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authorizations = append(m.Authorizations, Authorization{})
			if err := m.Authorizations[len(m.Authorizations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAuthorization(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAuthorization
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
func skipAuthorization(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAuthorization
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
					return 0, ErrIntOverflowAuthorization
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
					return 0, ErrIntOverflowAuthorization
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
				return 0, ErrInvalidLengthAuthorization
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAuthorization
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAuthorization
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAuthorization        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAuthorization          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAuthorization = fmt.Errorf("proto: unexpected end of group")
)
