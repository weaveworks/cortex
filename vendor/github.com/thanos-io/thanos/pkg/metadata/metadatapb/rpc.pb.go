// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: metadata/metadatapb/rpc.proto

package metadatapb

import (
	context "context"
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	storepb "github.com/thanos-io/thanos/pkg/store/storepb"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type MetadataRequest struct {
	Metric                  string                          `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
	Limit                   int32                           `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	PartialResponseStrategy storepb.PartialResponseStrategy `protobuf:"varint,3,opt,name=partial_response_strategy,json=partialResponseStrategy,proto3,enum=thanos.PartialResponseStrategy" json:"partial_response_strategy,omitempty"`
}

func (m *MetadataRequest) Reset()         { *m = MetadataRequest{} }
func (m *MetadataRequest) String() string { return proto.CompactTextString(m) }
func (*MetadataRequest) ProtoMessage()    {}
func (*MetadataRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1d9ae5661e0dc3fc, []int{0}
}
func (m *MetadataRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MetadataRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MetadataRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MetadataRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetadataRequest.Merge(m, src)
}
func (m *MetadataRequest) XXX_Size() int {
	return m.Size()
}
func (m *MetadataRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MetadataRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MetadataRequest proto.InternalMessageInfo

type MetadataResponse struct {
	// Types that are valid to be assigned to Result:
	//	*MetadataResponse_Metadata
	//	*MetadataResponse_Warning
	Result isMetadataResponse_Result `protobuf_oneof:"result"`
}

func (m *MetadataResponse) Reset()         { *m = MetadataResponse{} }
func (m *MetadataResponse) String() string { return proto.CompactTextString(m) }
func (*MetadataResponse) ProtoMessage()    {}
func (*MetadataResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1d9ae5661e0dc3fc, []int{1}
}
func (m *MetadataResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MetadataResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MetadataResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MetadataResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetadataResponse.Merge(m, src)
}
func (m *MetadataResponse) XXX_Size() int {
	return m.Size()
}
func (m *MetadataResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MetadataResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MetadataResponse proto.InternalMessageInfo

type isMetadataResponse_Result interface {
	isMetadataResponse_Result()
	MarshalTo([]byte) (int, error)
	Size() int
}

type MetadataResponse_Metadata struct {
	Metadata *MetricMetadata `protobuf:"bytes,1,opt,name=metadata,proto3,oneof" json:"metadata,omitempty"`
}
type MetadataResponse_Warning struct {
	Warning string `protobuf:"bytes,2,opt,name=warning,proto3,oneof" json:"warning,omitempty"`
}

func (*MetadataResponse_Metadata) isMetadataResponse_Result() {}
func (*MetadataResponse_Warning) isMetadataResponse_Result()  {}

func (m *MetadataResponse) GetResult() isMetadataResponse_Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func (m *MetadataResponse) GetMetadata() *MetricMetadata {
	if x, ok := m.GetResult().(*MetadataResponse_Metadata); ok {
		return x.Metadata
	}
	return nil
}

func (m *MetadataResponse) GetWarning() string {
	if x, ok := m.GetResult().(*MetadataResponse_Warning); ok {
		return x.Warning
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*MetadataResponse) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*MetadataResponse_Metadata)(nil),
		(*MetadataResponse_Warning)(nil),
	}
}

type MetricMetadata struct {
	Metadata map[string]MetricMetadataEntry `protobuf:"bytes,1,rep,name=metadata,proto3" json:"metadata" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *MetricMetadata) Reset()         { *m = MetricMetadata{} }
func (m *MetricMetadata) String() string { return proto.CompactTextString(m) }
func (*MetricMetadata) ProtoMessage()    {}
func (*MetricMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_1d9ae5661e0dc3fc, []int{2}
}
func (m *MetricMetadata) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MetricMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MetricMetadata.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MetricMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetricMetadata.Merge(m, src)
}
func (m *MetricMetadata) XXX_Size() int {
	return m.Size()
}
func (m *MetricMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_MetricMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_MetricMetadata proto.InternalMessageInfo

type MetricMetadataEntry struct {
	Metas []Meta `protobuf:"bytes,2,rep,name=metas,proto3" json:"metas"`
}

func (m *MetricMetadataEntry) Reset()         { *m = MetricMetadataEntry{} }
func (m *MetricMetadataEntry) String() string { return proto.CompactTextString(m) }
func (*MetricMetadataEntry) ProtoMessage()    {}
func (*MetricMetadataEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_1d9ae5661e0dc3fc, []int{3}
}
func (m *MetricMetadataEntry) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MetricMetadataEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MetricMetadataEntry.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MetricMetadataEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetricMetadataEntry.Merge(m, src)
}
func (m *MetricMetadataEntry) XXX_Size() int {
	return m.Size()
}
func (m *MetricMetadataEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_MetricMetadataEntry.DiscardUnknown(m)
}

var xxx_messageInfo_MetricMetadataEntry proto.InternalMessageInfo

type Meta struct {
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type"`
	Help string `protobuf:"bytes,2,opt,name=help,proto3" json:"help"`
	Unit string `protobuf:"bytes,3,opt,name=unit,proto3" json:"unit"`
}

func (m *Meta) Reset()         { *m = Meta{} }
func (m *Meta) String() string { return proto.CompactTextString(m) }
func (*Meta) ProtoMessage()    {}
func (*Meta) Descriptor() ([]byte, []int) {
	return fileDescriptor_1d9ae5661e0dc3fc, []int{4}
}
func (m *Meta) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Meta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Meta.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Meta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Meta.Merge(m, src)
}
func (m *Meta) XXX_Size() int {
	return m.Size()
}
func (m *Meta) XXX_DiscardUnknown() {
	xxx_messageInfo_Meta.DiscardUnknown(m)
}

var xxx_messageInfo_Meta proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MetadataRequest)(nil), "thanos.MetadataRequest")
	proto.RegisterType((*MetadataResponse)(nil), "thanos.MetadataResponse")
	proto.RegisterType((*MetricMetadata)(nil), "thanos.MetricMetadata")
	proto.RegisterMapType((map[string]MetricMetadataEntry)(nil), "thanos.MetricMetadata.MetadataEntry")
	proto.RegisterType((*MetricMetadataEntry)(nil), "thanos.MetricMetadataEntry")
	proto.RegisterType((*Meta)(nil), "thanos.Meta")
}

func init() { proto.RegisterFile("metadata/metadatapb/rpc.proto", fileDescriptor_1d9ae5661e0dc3fc) }

var fileDescriptor_1d9ae5661e0dc3fc = []byte{
	// 463 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x53, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xf6, 0xe6, 0x8f, 0x74, 0x02, 0xa5, 0x5a, 0xaa, 0xd6, 0x35, 0xe0, 0x44, 0x16, 0x07, 0x9f,
	0x62, 0x30, 0x1c, 0x10, 0x97, 0x4a, 0x91, 0x40, 0x95, 0x50, 0x25, 0xb4, 0x5c, 0x10, 0x1c, 0xca,
	0xa6, 0xac, 0x52, 0x0b, 0xc7, 0x5e, 0x76, 0xd7, 0x20, 0xbf, 0x05, 0x0f, 0xc0, 0x53, 0xf0, 0x14,
	0x39, 0xf6, 0xc8, 0xa9, 0x82, 0xe4, 0xc6, 0x53, 0xa0, 0xdd, 0xb5, 0x5b, 0x47, 0xf8, 0x32, 0xfa,
	0x66, 0xbe, 0xcf, 0x33, 0x9f, 0x66, 0xc7, 0xf0, 0x70, 0xc9, 0x14, 0xfd, 0x44, 0x15, 0x8d, 0x6a,
	0xc0, 0xe7, 0x91, 0xe0, 0xe7, 0x53, 0x2e, 0x72, 0x95, 0xe3, 0x81, 0xba, 0xa0, 0x59, 0x2e, 0xbd,
	0x23, 0xa9, 0x72, 0xc1, 0x22, 0x13, 0xf9, 0x3c, 0x52, 0x25, 0x67, 0xd2, 0x4a, 0xbc, 0xfd, 0x45,
	0xbe, 0xc8, 0x0d, 0x8c, 0x34, 0xb2, 0xd5, 0xe0, 0x07, 0x82, 0xbb, 0xa7, 0x55, 0x47, 0xc2, 0xbe,
	0x14, 0x4c, 0x2a, 0x7c, 0x00, 0x83, 0x25, 0x53, 0x22, 0x39, 0x77, 0xd1, 0x04, 0x85, 0x3b, 0xa4,
	0xca, 0xf0, 0x3e, 0xf4, 0xd3, 0x64, 0x99, 0x28, 0xb7, 0x33, 0x41, 0x61, 0x9f, 0xd8, 0x04, 0x7f,
	0x80, 0x23, 0x4e, 0x85, 0x4a, 0x68, 0x7a, 0x26, 0x98, 0xe4, 0x79, 0x26, 0xd9, 0x99, 0x54, 0x82,
	0x2a, 0xb6, 0x28, 0xdd, 0xee, 0x04, 0x85, 0xbb, 0xf1, 0x78, 0x6a, 0xed, 0x4d, 0xdf, 0x58, 0x21,
	0xa9, 0x74, 0x6f, 0x2b, 0x19, 0x39, 0xe4, 0xed, 0x44, 0x90, 0xc1, 0xde, 0x8d, 0x3b, 0xcb, 0xe1,
	0x67, 0x30, 0xac, 0x77, 0x60, 0x0c, 0x8e, 0xe2, 0x83, 0xba, 0xff, 0xa9, 0x31, 0x5a, 0x7f, 0x71,
	0xe2, 0x90, 0x6b, 0x25, 0xf6, 0xe0, 0xd6, 0x37, 0x2a, 0xb2, 0x24, 0x5b, 0x18, 0xfb, 0x3b, 0x27,
	0x0e, 0xa9, 0x0b, 0xb3, 0x21, 0x0c, 0x04, 0x93, 0x45, 0xaa, 0x82, 0x9f, 0x08, 0x76, 0xb7, 0x9b,
	0xe0, 0x57, 0x5b, 0xe3, 0xba, 0xe1, 0x28, 0x7e, 0xd4, 0x3e, 0x6e, 0x5a, 0x83, 0x97, 0x99, 0x12,
	0xe5, 0xac, 0xb7, 0xba, 0x1a, 0x37, 0x0c, 0x78, 0xef, 0xe0, 0xce, 0x96, 0x00, 0xef, 0x41, 0xf7,
	0x33, 0x2b, 0xab, 0x1d, 0x6b, 0x88, 0x9f, 0x40, 0xff, 0x2b, 0x4d, 0x0b, 0x66, 0x1c, 0x8e, 0xe2,
	0xfb, 0xed, 0x73, 0xcc, 0xd7, 0xc4, 0x2a, 0x5f, 0x74, 0x9e, 0xa3, 0xe0, 0x18, 0xee, 0xb5, 0x28,
	0x70, 0x08, 0x7d, 0x3d, 0x5c, 0xba, 0x1d, 0xe3, 0xfa, 0x76, 0xa3, 0x1b, 0xad, 0xdc, 0x59, 0x41,
	0xf0, 0x11, 0x7a, 0xba, 0x88, 0x1f, 0x40, 0x4f, 0x5f, 0x8c, 0xb5, 0x34, 0x1b, 0xfe, 0xbd, 0x1a,
	0x9b, 0x9c, 0x98, 0xa8, 0xd9, 0x0b, 0x96, 0x72, 0xbb, 0x3e, 0xcb, 0xea, 0x9c, 0x98, 0xa8, 0xd9,
	0x22, 0x4b, 0x94, 0x79, 0xf1, 0x8a, 0xd5, 0x39, 0x31, 0x31, 0x7e, 0x0d, 0xc3, 0xeb, 0x85, 0x1e,
	0x37, 0xf0, 0x61, 0xd3, 0x54, 0xe3, 0x06, 0x3d, 0xf7, 0x7f, 0xc2, 0x3e, 0xff, 0x63, 0x34, 0x0b,
	0x57, 0x7f, 0x7c, 0x67, 0xb5, 0xf6, 0xd1, 0xe5, 0xda, 0x47, 0xbf, 0xd7, 0x3e, 0xfa, 0xbe, 0xf1,
	0x9d, 0xcb, 0x8d, 0xef, 0xfc, 0xda, 0xf8, 0xce, 0x7b, 0xb8, 0xf9, 0x41, 0xe6, 0x03, 0x73, 0xe4,
	0x4f, 0xff, 0x05, 0x00, 0x00, 0xff, 0xff, 0x08, 0x77, 0xe4, 0x56, 0x3e, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MetadataClient is the client API for Metadata service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MetadataClient interface {
	Metadata(ctx context.Context, in *MetadataRequest, opts ...grpc.CallOption) (Metadata_MetadataClient, error)
}

type metadataClient struct {
	cc *grpc.ClientConn
}

func NewMetadataClient(cc *grpc.ClientConn) MetadataClient {
	return &metadataClient{cc}
}

func (c *metadataClient) Metadata(ctx context.Context, in *MetadataRequest, opts ...grpc.CallOption) (Metadata_MetadataClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Metadata_serviceDesc.Streams[0], "/thanos.Metadata/Metadata", opts...)
	if err != nil {
		return nil, err
	}
	x := &metadataMetadataClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Metadata_MetadataClient interface {
	Recv() (*MetadataResponse, error)
	grpc.ClientStream
}

type metadataMetadataClient struct {
	grpc.ClientStream
}

func (x *metadataMetadataClient) Recv() (*MetadataResponse, error) {
	m := new(MetadataResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MetadataServer is the server API for Metadata service.
type MetadataServer interface {
	Metadata(*MetadataRequest, Metadata_MetadataServer) error
}

// UnimplementedMetadataServer can be embedded to have forward compatible implementations.
type UnimplementedMetadataServer struct {
}

func (*UnimplementedMetadataServer) Metadata(req *MetadataRequest, srv Metadata_MetadataServer) error {
	return status.Errorf(codes.Unimplemented, "method Metadata not implemented")
}

func RegisterMetadataServer(s *grpc.Server, srv MetadataServer) {
	s.RegisterService(&_Metadata_serviceDesc, srv)
}

func _Metadata_Metadata_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(MetadataRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MetadataServer).Metadata(m, &metadataMetadataServer{stream})
}

type Metadata_MetadataServer interface {
	Send(*MetadataResponse) error
	grpc.ServerStream
}

type metadataMetadataServer struct {
	grpc.ServerStream
}

func (x *metadataMetadataServer) Send(m *MetadataResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _Metadata_serviceDesc = grpc.ServiceDesc{
	ServiceName: "thanos.Metadata",
	HandlerType: (*MetadataServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Metadata",
			Handler:       _Metadata_Metadata_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "metadata/metadatapb/rpc.proto",
}

func (m *MetadataRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MetadataRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MetadataRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.PartialResponseStrategy != 0 {
		i = encodeVarintRpc(dAtA, i, uint64(m.PartialResponseStrategy))
		i--
		dAtA[i] = 0x18
	}
	if m.Limit != 0 {
		i = encodeVarintRpc(dAtA, i, uint64(m.Limit))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Metric) > 0 {
		i -= len(m.Metric)
		copy(dAtA[i:], m.Metric)
		i = encodeVarintRpc(dAtA, i, uint64(len(m.Metric)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MetadataResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MetadataResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MetadataResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Result != nil {
		{
			size := m.Result.Size()
			i -= size
			if _, err := m.Result.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	return len(dAtA) - i, nil
}

func (m *MetadataResponse_Metadata) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MetadataResponse_Metadata) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Metadata != nil {
		{
			size, err := m.Metadata.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintRpc(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func (m *MetadataResponse_Warning) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MetadataResponse_Warning) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	i -= len(m.Warning)
	copy(dAtA[i:], m.Warning)
	i = encodeVarintRpc(dAtA, i, uint64(len(m.Warning)))
	i--
	dAtA[i] = 0x12
	return len(dAtA) - i, nil
}
func (m *MetricMetadata) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MetricMetadata) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MetricMetadata) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Metadata) > 0 {
		for k := range m.Metadata {
			v := m.Metadata[k]
			baseI := i
			{
				size, err := (&v).MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintRpc(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintRpc(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintRpc(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *MetricMetadataEntry) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MetricMetadataEntry) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MetricMetadataEntry) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Metas) > 0 {
		for iNdEx := len(m.Metas) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Metas[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintRpc(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	return len(dAtA) - i, nil
}

func (m *Meta) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Meta) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Meta) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Unit) > 0 {
		i -= len(m.Unit)
		copy(dAtA[i:], m.Unit)
		i = encodeVarintRpc(dAtA, i, uint64(len(m.Unit)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Help) > 0 {
		i -= len(m.Help)
		copy(dAtA[i:], m.Help)
		i = encodeVarintRpc(dAtA, i, uint64(len(m.Help)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Type) > 0 {
		i -= len(m.Type)
		copy(dAtA[i:], m.Type)
		i = encodeVarintRpc(dAtA, i, uint64(len(m.Type)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintRpc(dAtA []byte, offset int, v uint64) int {
	offset -= sovRpc(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MetadataRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Metric)
	if l > 0 {
		n += 1 + l + sovRpc(uint64(l))
	}
	if m.Limit != 0 {
		n += 1 + sovRpc(uint64(m.Limit))
	}
	if m.PartialResponseStrategy != 0 {
		n += 1 + sovRpc(uint64(m.PartialResponseStrategy))
	}
	return n
}

func (m *MetadataResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Result != nil {
		n += m.Result.Size()
	}
	return n
}

func (m *MetadataResponse_Metadata) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Metadata != nil {
		l = m.Metadata.Size()
		n += 1 + l + sovRpc(uint64(l))
	}
	return n
}
func (m *MetadataResponse_Warning) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Warning)
	n += 1 + l + sovRpc(uint64(l))
	return n
}
func (m *MetricMetadata) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Metadata) > 0 {
		for k, v := range m.Metadata {
			_ = k
			_ = v
			l = v.Size()
			mapEntrySize := 1 + len(k) + sovRpc(uint64(len(k))) + 1 + l + sovRpc(uint64(l))
			n += mapEntrySize + 1 + sovRpc(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *MetricMetadataEntry) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Metas) > 0 {
		for _, e := range m.Metas {
			l = e.Size()
			n += 1 + l + sovRpc(uint64(l))
		}
	}
	return n
}

func (m *Meta) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Type)
	if l > 0 {
		n += 1 + l + sovRpc(uint64(l))
	}
	l = len(m.Help)
	if l > 0 {
		n += 1 + l + sovRpc(uint64(l))
	}
	l = len(m.Unit)
	if l > 0 {
		n += 1 + l + sovRpc(uint64(l))
	}
	return n
}

func sovRpc(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRpc(x uint64) (n int) {
	return sovRpc(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MetadataRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRpc
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
			return fmt.Errorf("proto: MetadataRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MetadataRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Metric", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpc
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
				return ErrInvalidLengthRpc
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRpc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Metric = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Limit", wireType)
			}
			m.Limit = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpc
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Limit |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PartialResponseStrategy", wireType)
			}
			m.PartialResponseStrategy = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpc
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PartialResponseStrategy |= storepb.PartialResponseStrategy(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipRpc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRpc
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
func (m *MetadataResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRpc
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
			return fmt.Errorf("proto: MetadataResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MetadataResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Metadata", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpc
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
				return ErrInvalidLengthRpc
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRpc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &MetricMetadata{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Result = &MetadataResponse_Metadata{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Warning", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpc
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
				return ErrInvalidLengthRpc
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRpc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Result = &MetadataResponse_Warning{string(dAtA[iNdEx:postIndex])}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRpc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRpc
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
func (m *MetricMetadata) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRpc
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
			return fmt.Errorf("proto: MetricMetadata: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MetricMetadata: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Metadata", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpc
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
				return ErrInvalidLengthRpc
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRpc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Metadata == nil {
				m.Metadata = make(map[string]MetricMetadataEntry)
			}
			var mapkey string
			mapvalue := &MetricMetadataEntry{}
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowRpc
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowRpc
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthRpc
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthRpc
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowRpc
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthRpc
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return ErrInvalidLengthRpc
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &MetricMetadataEntry{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipRpc(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthRpc
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Metadata[mapkey] = *mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRpc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRpc
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
func (m *MetricMetadataEntry) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRpc
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
			return fmt.Errorf("proto: MetricMetadataEntry: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MetricMetadataEntry: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Metas", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpc
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
				return ErrInvalidLengthRpc
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRpc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Metas = append(m.Metas, Meta{})
			if err := m.Metas[len(m.Metas)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRpc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRpc
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
func (m *Meta) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRpc
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
			return fmt.Errorf("proto: Meta: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Meta: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpc
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
				return ErrInvalidLengthRpc
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRpc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Type = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Help", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpc
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
				return ErrInvalidLengthRpc
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRpc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Help = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Unit", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRpc
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
				return ErrInvalidLengthRpc
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRpc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Unit = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRpc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRpc
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
func skipRpc(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRpc
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
					return 0, ErrIntOverflowRpc
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
					return 0, ErrIntOverflowRpc
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
				return 0, ErrInvalidLengthRpc
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupRpc
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthRpc
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthRpc        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRpc          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupRpc = fmt.Errorf("proto: unexpected end of group")
)
