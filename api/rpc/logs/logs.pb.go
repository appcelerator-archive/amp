// Code generated by protoc-gen-go.
// source: github.com/appcelerator/amp/api/rpc/logs/logs.proto
// DO NOT EDIT!

/*
Package logs is a generated protocol buffer package.

It is generated from these files:
	github.com/appcelerator/amp/api/rpc/logs/logs.proto

It has these top-level messages:
	LogEntry
	GetRequest
	GetReply
*/
package logs

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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

type LogEntry struct {
	Timestamp          string `protobuf:"bytes,1,opt,name=timestamp" json:"timestamp,omitempty"`
	ContainerId        string `protobuf:"bytes,2,opt,name=container_id,json=containerId" json:"container_id,omitempty"`
	ContainerName      string `protobuf:"bytes,3,opt,name=container_name,json=containerName" json:"container_name,omitempty"`
	ContainerShortName string `protobuf:"bytes,4,opt,name=container_short_name,json=containerShortName" json:"container_short_name,omitempty"`
	ContainerState     string `protobuf:"bytes,5,opt,name=container_state,json=containerState" json:"container_state,omitempty"`
	ServiceName        string `protobuf:"bytes,6,opt,name=service_name,json=serviceName" json:"service_name,omitempty"`
	ServiceId          string `protobuf:"bytes,7,opt,name=service_id,json=serviceId" json:"service_id,omitempty"`
	TaskId             string `protobuf:"bytes,8,opt,name=task_id,json=taskId" json:"task_id,omitempty"`
	StackName          string `protobuf:"bytes,9,opt,name=stack_name,json=stackName" json:"stack_name,omitempty"`
	NodeId             string `protobuf:"bytes,10,opt,name=node_id,json=nodeId" json:"node_id,omitempty"`
	Role               string `protobuf:"bytes,11,opt,name=role" json:"role,omitempty"`
	Message            string `protobuf:"bytes,12,opt,name=message" json:"message,omitempty"`
}

func (m *LogEntry) Reset()                    { *m = LogEntry{} }
func (m *LogEntry) String() string            { return proto.CompactTextString(m) }
func (*LogEntry) ProtoMessage()               {}
func (*LogEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *LogEntry) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

func (m *LogEntry) GetContainerId() string {
	if m != nil {
		return m.ContainerId
	}
	return ""
}

func (m *LogEntry) GetContainerName() string {
	if m != nil {
		return m.ContainerName
	}
	return ""
}

func (m *LogEntry) GetContainerShortName() string {
	if m != nil {
		return m.ContainerShortName
	}
	return ""
}

func (m *LogEntry) GetContainerState() string {
	if m != nil {
		return m.ContainerState
	}
	return ""
}

func (m *LogEntry) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

func (m *LogEntry) GetServiceId() string {
	if m != nil {
		return m.ServiceId
	}
	return ""
}

func (m *LogEntry) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

func (m *LogEntry) GetStackName() string {
	if m != nil {
		return m.StackName
	}
	return ""
}

func (m *LogEntry) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func (m *LogEntry) GetRole() string {
	if m != nil {
		return m.Role
	}
	return ""
}

func (m *LogEntry) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type GetRequest struct {
	Container string `protobuf:"bytes,1,opt,name=container" json:"container,omitempty"`
	Message   string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	Node      string `protobuf:"bytes,3,opt,name=node" json:"node,omitempty"`
	Size      int64  `protobuf:"varint,4,opt,name=size" json:"size,omitempty"`
	Service   string `protobuf:"bytes,5,opt,name=service" json:"service,omitempty"`
	Stack     string `protobuf:"bytes,6,opt,name=stack" json:"stack,omitempty"`
	Infra     bool   `protobuf:"varint,7,opt,name=infra" json:"infra,omitempty"`
}

func (m *GetRequest) Reset()                    { *m = GetRequest{} }
func (m *GetRequest) String() string            { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()               {}
func (*GetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *GetRequest) GetContainer() string {
	if m != nil {
		return m.Container
	}
	return ""
}

func (m *GetRequest) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *GetRequest) GetNode() string {
	if m != nil {
		return m.Node
	}
	return ""
}

func (m *GetRequest) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *GetRequest) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *GetRequest) GetStack() string {
	if m != nil {
		return m.Stack
	}
	return ""
}

func (m *GetRequest) GetInfra() bool {
	if m != nil {
		return m.Infra
	}
	return false
}

type GetReply struct {
	Entries []*LogEntry `protobuf:"bytes,1,rep,name=entries" json:"entries,omitempty"`
}

func (m *GetReply) Reset()                    { *m = GetReply{} }
func (m *GetReply) String() string            { return proto.CompactTextString(m) }
func (*GetReply) ProtoMessage()               {}
func (*GetReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GetReply) GetEntries() []*LogEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func init() {
	proto.RegisterType((*LogEntry)(nil), "logs.LogEntry")
	proto.RegisterType((*GetRequest)(nil), "logs.GetRequest")
	proto.RegisterType((*GetReply)(nil), "logs.GetReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Logs service

type LogsClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error)
	GetStream(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (Logs_GetStreamClient, error)
}

type logsClient struct {
	cc *grpc.ClientConn
}

func NewLogsClient(cc *grpc.ClientConn) LogsClient {
	return &logsClient{cc}
}

func (c *logsClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error) {
	out := new(GetReply)
	err := grpc.Invoke(ctx, "/logs.Logs/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logsClient) GetStream(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (Logs_GetStreamClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Logs_serviceDesc.Streams[0], c.cc, "/logs.Logs/GetStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &logsGetStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Logs_GetStreamClient interface {
	Recv() (*LogEntry, error)
	grpc.ClientStream
}

type logsGetStreamClient struct {
	grpc.ClientStream
}

func (x *logsGetStreamClient) Recv() (*LogEntry, error) {
	m := new(LogEntry)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Logs service

type LogsServer interface {
	Get(context.Context, *GetRequest) (*GetReply, error)
	GetStream(*GetRequest, Logs_GetStreamServer) error
}

func RegisterLogsServer(s *grpc.Server, srv LogsServer) {
	s.RegisterService(&_Logs_serviceDesc, srv)
}

func _Logs_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogsServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logs.Logs/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogsServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logs_GetStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(LogsServer).GetStream(m, &logsGetStreamServer{stream})
}

type Logs_GetStreamServer interface {
	Send(*LogEntry) error
	grpc.ServerStream
}

type logsGetStreamServer struct {
	grpc.ServerStream
}

func (x *logsGetStreamServer) Send(m *LogEntry) error {
	return x.ServerStream.SendMsg(m)
}

var _Logs_serviceDesc = grpc.ServiceDesc{
	ServiceName: "logs.Logs",
	HandlerType: (*LogsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Logs_Get_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetStream",
			Handler:       _Logs_GetStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "github.com/appcelerator/amp/api/rpc/logs/logs.proto",
}

func init() {
	proto.RegisterFile("github.com/appcelerator/amp/api/rpc/logs/logs.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 433 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0x5f, 0x6a, 0xdb, 0x40,
	0x10, 0xc6, 0x51, 0xe4, 0x58, 0xf6, 0x38, 0x75, 0xcb, 0x12, 0xc8, 0x52, 0x5a, 0x48, 0x0d, 0x25,
	0x7e, 0xb2, 0x42, 0xd2, 0x2b, 0x14, 0x63, 0x08, 0x7d, 0x90, 0x0f, 0xd0, 0x6e, 0xa4, 0xa9, 0xb2,
	0xc4, 0xd2, 0x6e, 0x77, 0x27, 0x05, 0xf7, 0x10, 0x3d, 0x4c, 0x4f, 0x58, 0x76, 0x56, 0x7f, 0xdc,
	0x87, 0xbe, 0x88, 0x9d, 0xdf, 0x7c, 0xf3, 0xcd, 0xce, 0x68, 0xe1, 0xbe, 0xd6, 0xf4, 0xf4, 0xf2,
	0xb8, 0x29, 0x4d, 0x93, 0x2b, 0x6b, 0x4b, 0x3c, 0xa0, 0x53, 0x64, 0x5c, 0xae, 0x1a, 0x9b, 0x2b,
	0xab, 0x73, 0x67, 0xcb, 0xfc, 0x60, 0x6a, 0xcf, 0x9f, 0x8d, 0x75, 0x86, 0x8c, 0x98, 0x84, 0xf3,
	0xea, 0x77, 0x0a, 0xb3, 0x07, 0x53, 0x7f, 0x6e, 0xc9, 0x1d, 0xc5, 0x3b, 0x98, 0x93, 0x6e, 0xd0,
	0x93, 0x6a, 0xac, 0x4c, 0xae, 0x93, 0xf5, 0xbc, 0x18, 0x81, 0xf8, 0x00, 0x17, 0xa5, 0x69, 0x49,
	0xe9, 0x16, 0xdd, 0x57, 0x5d, 0xc9, 0x33, 0x16, 0x2c, 0x06, 0xb6, 0xab, 0xc4, 0x47, 0x58, 0x8e,
	0x92, 0x56, 0x35, 0x28, 0x53, 0x16, 0xbd, 0x1a, 0xe8, 0x17, 0xd5, 0xa0, 0xb8, 0x85, 0xcb, 0x51,
	0xe6, 0x9f, 0x8c, 0xa3, 0x28, 0x9e, 0xb0, 0x58, 0x0c, 0xb9, 0x7d, 0x48, 0x71, 0xc5, 0x0d, 0xbc,
	0x3e, 0xa9, 0x20, 0x45, 0x28, 0xcf, 0x59, 0x3c, 0xf6, 0xdb, 0x07, 0x1a, 0x2e, 0xe9, 0xd1, 0xfd,
	0xd4, 0x25, 0x46, 0xcb, 0x69, 0xbc, 0x64, 0xc7, 0xd8, 0xeb, 0x3d, 0x40, 0x2f, 0xd1, 0x95, 0xcc,
	0xe2, 0x98, 0x1d, 0xd9, 0x55, 0xe2, 0x0a, 0x32, 0x52, 0xfe, 0x39, 0xe4, 0x66, 0x9c, 0x9b, 0x86,
	0x70, 0x57, 0x71, 0x1d, 0xa9, 0xf2, 0x39, 0x1a, 0xcf, 0xbb, 0xba, 0x40, 0xd8, 0xf6, 0x0a, 0xb2,
	0xd6, 0x54, 0xec, 0x09, 0xb1, 0x2e, 0x84, 0xbb, 0x4a, 0x08, 0x98, 0x38, 0x73, 0x40, 0xb9, 0x60,
	0xca, 0x67, 0x21, 0x21, 0x6b, 0xd0, 0x7b, 0x55, 0xa3, 0xbc, 0x60, 0xdc, 0x87, 0xab, 0x3f, 0x09,
	0xc0, 0x16, 0xa9, 0xc0, 0x1f, 0x2f, 0xe8, 0x29, 0xfc, 0x92, 0x61, 0xc2, 0xfe, 0x97, 0x0c, 0xe0,
	0xd4, 0xe6, 0xec, 0x1f, 0x9b, 0xd0, 0x34, 0xb4, 0xef, 0xf6, 0xcf, 0xe7, 0xc0, 0xbc, 0xfe, 0x15,
	0xd7, 0x9c, 0x16, 0x7c, 0x0e, 0x0e, 0xdd, 0xe8, 0xdd, 0x42, 0xfb, 0x50, 0x5c, 0xc2, 0x39, 0x0f,
	0xd7, 0xad, 0x30, 0x06, 0x81, 0xea, 0xf6, 0xbb, 0x53, 0xbc, 0xb7, 0x59, 0x11, 0x83, 0xd5, 0x27,
	0x98, 0xf1, 0x9d, 0xed, 0xe1, 0x28, 0xd6, 0x90, 0x61, 0x4b, 0x4e, 0xa3, 0x97, 0xc9, 0x75, 0xba,
	0x5e, 0xdc, 0x2d, 0x37, 0xfc, 0xea, 0xfa, 0x57, 0x56, 0xf4, 0xe9, 0xbb, 0x6f, 0x30, 0x79, 0x30,
	0xb5, 0x17, 0x37, 0x90, 0x6e, 0x91, 0xc4, 0x9b, 0xa8, 0x1b, 0x87, 0x7f, 0xbb, 0x3c, 0x21, 0xc1,
	0x3a, 0x87, 0xf9, 0x16, 0x69, 0x4f, 0x0e, 0x55, 0xf3, 0x7f, 0x79, 0xdf, 0xe8, 0x36, 0x79, 0x9c,
	0xf2, 0x53, 0xbf, 0xff, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x20, 0xc2, 0x3c, 0xbb, 0x21, 0x03, 0x00,
	0x00,
}
