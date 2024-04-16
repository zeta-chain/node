// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: authority/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// QueryGetPoliciesRequest is the request type for the Query/Policies RPC method.
type QueryGetPoliciesRequest struct {
}

func (m *QueryGetPoliciesRequest) Reset()         { *m = QueryGetPoliciesRequest{} }
func (m *QueryGetPoliciesRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGetPoliciesRequest) ProtoMessage()    {}
func (*QueryGetPoliciesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b64d9e7f9da035b5, []int{0}
}
func (m *QueryGetPoliciesRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetPoliciesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetPoliciesRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetPoliciesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetPoliciesRequest.Merge(m, src)
}
func (m *QueryGetPoliciesRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetPoliciesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetPoliciesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetPoliciesRequest proto.InternalMessageInfo

// QueryGetPoliciesResponse is the response type for the Query/Policies RPC method.
type QueryGetPoliciesResponse struct {
	Policies Policies `protobuf:"bytes,1,opt,name=policies,proto3" json:"policies"`
}

func (m *QueryGetPoliciesResponse) Reset()         { *m = QueryGetPoliciesResponse{} }
func (m *QueryGetPoliciesResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGetPoliciesResponse) ProtoMessage()    {}
func (*QueryGetPoliciesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b64d9e7f9da035b5, []int{1}
}
func (m *QueryGetPoliciesResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetPoliciesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetPoliciesResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetPoliciesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetPoliciesResponse.Merge(m, src)
}
func (m *QueryGetPoliciesResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetPoliciesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetPoliciesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetPoliciesResponse proto.InternalMessageInfo

func (m *QueryGetPoliciesResponse) GetPolicies() Policies {
	if m != nil {
		return m.Policies
	}
	return Policies{}
}

func init() {
	proto.RegisterType((*QueryGetPoliciesRequest)(nil), "authority.QueryGetPoliciesRequest")
	proto.RegisterType((*QueryGetPoliciesResponse)(nil), "authority.QueryGetPoliciesResponse")
}

func init() { proto.RegisterFile("authority/query.proto", fileDescriptor_b64d9e7f9da035b5) }

var fileDescriptor_b64d9e7f9da035b5 = []byte{
	// 306 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0xb3, 0xa2, 0x52, 0xd7, 0x5b, 0x54, 0xac, 0x45, 0xd6, 0x12, 0x41, 0x44, 0x30, 0x6b,
	0x2b, 0xbe, 0x40, 0x2f, 0x82, 0x27, 0xdb, 0xa3, 0xb7, 0x4d, 0x18, 0xb6, 0x0b, 0xed, 0xce, 0x36,
	0xbb, 0x11, 0xab, 0x37, 0xf1, 0x01, 0x04, 0x5f, 0xaa, 0xc7, 0x82, 0x17, 0x4f, 0x22, 0xad, 0x0f,
	0x22, 0x4d, 0x93, 0x18, 0x28, 0x7a, 0xfb, 0x99, 0xfd, 0xfe, 0x7f, 0xff, 0x19, 0xba, 0x27, 0x52,
	0xd7, 0xc7, 0x44, 0xb9, 0x31, 0x1f, 0xa5, 0x90, 0x8c, 0x43, 0x93, 0xa0, 0x43, 0x7f, 0xab, 0x1c,
	0x37, 0xea, 0xbf, 0x84, 0xc1, 0x81, 0x8a, 0x15, 0xd8, 0x25, 0xd4, 0x38, 0x8b, 0xd1, 0x0e, 0xd1,
	0xf2, 0x48, 0x58, 0x58, 0xba, 0xf9, 0x7d, 0x2b, 0x02, 0x27, 0x5a, 0xdc, 0x08, 0xa9, 0xb4, 0x70,
	0x0a, 0x75, 0xce, 0xee, 0x4a, 0x94, 0x98, 0x49, 0xbe, 0x50, 0xf9, 0xf4, 0x50, 0x22, 0xca, 0x01,
	0x70, 0x61, 0x14, 0x17, 0x5a, 0xa3, 0xcb, 0x2c, 0x79, 0x7e, 0x70, 0x40, 0xf7, 0xbb, 0x8b, 0xd4,
	0x6b, 0x70, 0xb7, 0xf9, 0xcf, 0x3d, 0x18, 0xa5, 0x60, 0x5d, 0xd0, 0xa5, 0xf5, 0xd5, 0x27, 0x6b,
	0x50, 0x5b, 0xf0, 0xaf, 0x68, 0xad, 0x28, 0x5a, 0x27, 0x4d, 0x72, 0xba, 0xdd, 0xde, 0x09, 0xcb,
	0x1d, 0xc2, 0x02, 0xef, 0xac, 0x4f, 0x3e, 0x8f, 0xbc, 0x5e, 0x89, 0xb6, 0x5f, 0x08, 0xdd, 0xc8,
	0x32, 0xfd, 0x27, 0x5a, 0x2b, 0x28, 0x3f, 0xa8, 0x58, 0xff, 0x28, 0xd3, 0x38, 0xfe, 0x97, 0x59,
	0xb6, 0x0a, 0x4e, 0x9e, 0xdf, 0xbf, 0xdf, 0xd6, 0x9a, 0x3e, 0xe3, 0x8f, 0xe0, 0xc4, 0x79, 0xdc,
	0x17, 0x4a, 0xf3, 0xd5, 0xd3, 0x76, 0x6e, 0x26, 0x33, 0x46, 0xa6, 0x33, 0x46, 0xbe, 0x66, 0x8c,
	0xbc, 0xce, 0x99, 0x37, 0x9d, 0x33, 0xef, 0x63, 0xce, 0xbc, 0xbb, 0x0b, 0xa9, 0x5c, 0x3f, 0x8d,
	0xc2, 0x18, 0x87, 0xd5, 0x8c, 0x85, 0x8c, 0x31, 0x01, 0xfe, 0x50, 0x89, 0x73, 0x63, 0x03, 0x36,
	0xda, 0xcc, 0xee, 0x78, 0xf9, 0x13, 0x00, 0x00, 0xff, 0xff, 0x50, 0xcc, 0xe2, 0x70, 0xe5, 0x01,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Queries Policies
	Policies(ctx context.Context, in *QueryGetPoliciesRequest, opts ...grpc.CallOption) (*QueryGetPoliciesResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Policies(ctx context.Context, in *QueryGetPoliciesRequest, opts ...grpc.CallOption) (*QueryGetPoliciesResponse, error) {
	out := new(QueryGetPoliciesResponse)
	err := c.cc.Invoke(ctx, "/authority.Query/Policies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Queries Policies
	Policies(context.Context, *QueryGetPoliciesRequest) (*QueryGetPoliciesResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Policies(ctx context.Context, req *QueryGetPoliciesRequest) (*QueryGetPoliciesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Policies not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Policies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetPoliciesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Policies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authority.Query/Policies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Policies(ctx, req.(*QueryGetPoliciesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "authority.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Policies",
			Handler:    _Query_Policies_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "authority/query.proto",
}

func (m *QueryGetPoliciesRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetPoliciesRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetPoliciesRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryGetPoliciesResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetPoliciesResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetPoliciesResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Policies.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryGetPoliciesRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryGetPoliciesResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Policies.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryGetPoliciesRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryGetPoliciesRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetPoliciesRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryGetPoliciesResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryGetPoliciesResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetPoliciesResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Policies", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Policies.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
