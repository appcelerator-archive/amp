// Code generated by protoc-gen-go.
// source: github.com/appcelerator/amp/api/rpc/resource/resource.proto
// DO NOT EDIT!

/*
Package resource is a generated protocol buffer package.

It is generated from these files:
	github.com/appcelerator/amp/api/rpc/resource/resource.proto

It has these top-level messages:
	ResourceEntry
	ListResourcesRequest
	ListResourcesReply
*/
package resource

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ResourceType int32

const (
	ResourceType_RESOURCE_STACK ResourceType = 0
)

var ResourceType_name = map[int32]string{
	0: "RESOURCE_STACK",
}
var ResourceType_value = map[string]int32{
	"RESOURCE_STACK": 0,
}

func (x ResourceType) String() string {
	return proto.EnumName(ResourceType_name, int32(x))
}
func (ResourceType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ResourceEntry struct {
	Id   string       `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Type ResourceType `protobuf:"varint,2,opt,name=type,enum=resource.ResourceType" json:"type,omitempty"`
	Name string       `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
}

func (m *ResourceEntry) Reset()                    { *m = ResourceEntry{} }
func (m *ResourceEntry) String() string            { return proto.CompactTextString(m) }
func (*ResourceEntry) ProtoMessage()               {}
func (*ResourceEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ResourceEntry) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ResourceEntry) GetType() ResourceType {
	if m != nil {
		return m.Type
	}
	return ResourceType_RESOURCE_STACK
}

func (m *ResourceEntry) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type ListResourcesRequest struct {
}

func (m *ListResourcesRequest) Reset()                    { *m = ListResourcesRequest{} }
func (m *ListResourcesRequest) String() string            { return proto.CompactTextString(m) }
func (*ListResourcesRequest) ProtoMessage()               {}
func (*ListResourcesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type ListResourcesReply struct {
	Resources []*ResourceEntry `protobuf:"bytes,1,rep,name=resources" json:"resources,omitempty"`
}

func (m *ListResourcesReply) Reset()                    { *m = ListResourcesReply{} }
func (m *ListResourcesReply) String() string            { return proto.CompactTextString(m) }
func (*ListResourcesReply) ProtoMessage()               {}
func (*ListResourcesReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ListResourcesReply) GetResources() []*ResourceEntry {
	if m != nil {
		return m.Resources
	}
	return nil
}

func init() {
	proto.RegisterType((*ResourceEntry)(nil), "resource.ResourceEntry")
	proto.RegisterType((*ListResourcesRequest)(nil), "resource.ListResourcesRequest")
	proto.RegisterType((*ListResourcesReply)(nil), "resource.ListResourcesReply")
	proto.RegisterEnum("resource.ResourceType", ResourceType_name, ResourceType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Resource service

type ResourceClient interface {
	ListResources(ctx context.Context, in *ListResourcesRequest, opts ...grpc.CallOption) (*ListResourcesReply, error)
}

type resourceClient struct {
	cc *grpc.ClientConn
}

func NewResourceClient(cc *grpc.ClientConn) ResourceClient {
	return &resourceClient{cc}
}

func (c *resourceClient) ListResources(ctx context.Context, in *ListResourcesRequest, opts ...grpc.CallOption) (*ListResourcesReply, error) {
	out := new(ListResourcesReply)
	err := grpc.Invoke(ctx, "/resource.Resource/ListResources", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Resource service

type ResourceServer interface {
	ListResources(context.Context, *ListResourcesRequest) (*ListResourcesReply, error)
}

func RegisterResourceServer(s *grpc.Server, srv ResourceServer) {
	s.RegisterService(&_Resource_serviceDesc, srv)
}

func _Resource_ListResources_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListResourcesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).ListResources(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/resource.Resource/ListResources",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).ListResources(ctx, req.(*ListResourcesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Resource_serviceDesc = grpc.ServiceDesc{
	ServiceName: "resource.Resource",
	HandlerType: (*ResourceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListResources",
			Handler:    _Resource_ListResources_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/appcelerator/amp/api/rpc/resource/resource.proto",
}

func init() {
	proto.RegisterFile("github.com/appcelerator/amp/api/rpc/resource/resource.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 293 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0x4f, 0x4b, 0xc3, 0x30,
	0x18, 0xc6, 0x6d, 0x37, 0x64, 0x7b, 0x75, 0x53, 0x5e, 0x74, 0x96, 0x31, 0x64, 0xf4, 0x34, 0x76,
	0x58, 0x70, 0xe2, 0xc9, 0x93, 0x8c, 0x9e, 0x26, 0x08, 0xd9, 0x3c, 0x8f, 0xac, 0x0d, 0x35, 0xd0,
	0x36, 0x31, 0x49, 0x85, 0x5e, 0xfd, 0x0a, 0x7e, 0x34, 0xbf, 0x82, 0x1f, 0x44, 0x88, 0xb6, 0xf3,
	0xef, 0xed, 0x4d, 0xf2, 0x7b, 0xf8, 0x3d, 0xe4, 0x85, 0xeb, 0x54, 0xd8, 0x87, 0x72, 0x3b, 0x8b,
	0x65, 0x4e, 0x98, 0x52, 0x31, 0xcf, 0xb8, 0x66, 0x56, 0x6a, 0xc2, 0x72, 0x45, 0x98, 0x12, 0x44,
	0xab, 0x98, 0x68, 0x6e, 0x64, 0xa9, 0x63, 0xde, 0x0c, 0x33, 0xa5, 0xa5, 0x95, 0xd8, 0xa9, 0xcf,
	0xc3, 0x51, 0x2a, 0x65, 0x9a, 0x71, 0x97, 0x60, 0x45, 0x21, 0x2d, 0xb3, 0x42, 0x16, 0xe6, 0x83,
	0x0b, 0x37, 0xd0, 0xa3, 0x9f, 0x64, 0x54, 0x58, 0x5d, 0x61, 0x1f, 0x7c, 0x91, 0x04, 0xde, 0xd8,
	0x9b, 0x74, 0xa9, 0x2f, 0x12, 0x9c, 0x42, 0xdb, 0x56, 0x8a, 0x07, 0xfe, 0xd8, 0x9b, 0xf4, 0xe7,
	0x83, 0x59, 0xe3, 0xa9, 0x63, 0xeb, 0x4a, 0x71, 0xea, 0x18, 0x44, 0x68, 0x17, 0x2c, 0xe7, 0x41,
	0xcb, 0xa5, 0xdd, 0x1c, 0x0e, 0xe0, 0xe4, 0x56, 0x18, 0x5b, 0xd3, 0x86, 0xf2, 0xc7, 0x92, 0x1b,
	0x1b, 0x2e, 0x01, 0x7f, 0xdc, 0xab, 0xac, 0xc2, 0x2b, 0xe8, 0xd6, 0x02, 0x13, 0x78, 0xe3, 0xd6,
	0xe4, 0x60, 0x7e, 0xf6, 0x5b, 0xe9, 0x9a, 0xd2, 0x1d, 0x39, 0x0d, 0xe1, 0xf0, 0x6b, 0x1d, 0x44,
	0xe8, 0xd3, 0x68, 0x75, 0x77, 0x4f, 0x17, 0xd1, 0x66, 0xb5, 0xbe, 0x59, 0x2c, 0x8f, 0xf7, 0xe6,
	0x0a, 0x3a, 0x35, 0x83, 0x09, 0xf4, 0xbe, 0xc9, 0xf1, 0x7c, 0x27, 0xf9, 0xab, 0xed, 0x70, 0xf4,
	0xef, 0xbb, 0xca, 0xaa, 0xf0, 0xf4, 0xf9, 0xf5, 0xed, 0xc5, 0x3f, 0xc2, 0x1e, 0x79, 0xba, 0x68,
	0x16, 0x61, 0xb6, 0xfb, 0xee, 0x8b, 0x2f, 0xdf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xdc, 0x55, 0xb7,
	0x00, 0xc9, 0x01, 0x00, 0x00,
}
