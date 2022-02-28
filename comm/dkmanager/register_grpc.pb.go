// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package dkmanager

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RegistryClient is the client API for Registry service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegistryClient interface {
	Register(ctx context.Context, in *HostRegisterInfo, opts ...grpc.CallOption) (*RegisterResult, error)
	ReportNodeStatus(ctx context.Context, in *HostReport, opts ...grpc.CallOption) (*ReportStatus, error)
	ScheduleTask(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ScheduleResult, error)
	ReleaseResource(ctx context.Context, in *ReleaseRequest, opts ...grpc.CallOption) (*ReleaseResult, error)
}

type registryClient struct {
	cc grpc.ClientConnInterface
}

func NewRegistryClient(cc grpc.ClientConnInterface) RegistryClient {
	return &registryClient{cc}
}

func (c *registryClient) Register(ctx context.Context, in *HostRegisterInfo, opts ...grpc.CallOption) (*RegisterResult, error) {
	out := new(RegisterResult)
	err := c.cc.Invoke(ctx, "/dkmanager.Registry/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryClient) ReportNodeStatus(ctx context.Context, in *HostReport, opts ...grpc.CallOption) (*ReportStatus, error) {
	out := new(ReportStatus)
	err := c.cc.Invoke(ctx, "/dkmanager.Registry/ReportNodeStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryClient) ScheduleTask(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ScheduleResult, error) {
	out := new(ScheduleResult)
	err := c.cc.Invoke(ctx, "/dkmanager.Registry/ScheduleTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryClient) ReleaseResource(ctx context.Context, in *ReleaseRequest, opts ...grpc.CallOption) (*ReleaseResult, error) {
	out := new(ReleaseResult)
	err := c.cc.Invoke(ctx, "/dkmanager.Registry/ReleaseResource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegistryServer is the server API for Registry service.
// All implementations must embed UnimplementedRegistryServer
// for forward compatibility
type RegistryServer interface {
	Register(context.Context, *HostRegisterInfo) (*RegisterResult, error)
	ReportNodeStatus(context.Context, *HostReport) (*ReportStatus, error)
	ScheduleTask(context.Context, *Empty) (*ScheduleResult, error)
	ReleaseResource(context.Context, *ReleaseRequest) (*ReleaseResult, error)
	//mustEmbedUnimplementedRegistryServer()
}

// UnimplementedRegistryServer must be embedded to have forward compatible implementations.
type UnimplementedRegistryServer struct {
}

func (UnimplementedRegistryServer) Register(context.Context, *HostRegisterInfo) (*RegisterResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedRegistryServer) ReportNodeStatus(context.Context, *HostReport) (*ReportStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportNodeStatus not implemented")
}
func (UnimplementedRegistryServer) ScheduleTask(context.Context, *Empty) (*ScheduleResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ScheduleTask not implemented")
}
func (UnimplementedRegistryServer) ReleaseResource(context.Context, *ReleaseRequest) (*ReleaseResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseResource not implemented")
}
func (UnimplementedRegistryServer) mustEmbedUnimplementedRegistryServer() {}

// UnsafeRegistryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegistryServer will
// result in compilation errors.
type UnsafeRegistryServer interface {
	mustEmbedUnimplementedRegistryServer()
}

func RegisterRegistryServer(s grpc.ServiceRegistrar, srv RegistryServer) {
	s.RegisterService(&Registry_ServiceDesc, srv)
}

func _Registry_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostRegisterInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dkmanager.Registry/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServer).Register(ctx, req.(*HostRegisterInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Registry_ReportNodeStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostReport)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServer).ReportNodeStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dkmanager.Registry/ReportNodeStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServer).ReportNodeStatus(ctx, req.(*HostReport))
	}
	return interceptor(ctx, in, info, handler)
}

func _Registry_ScheduleTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServer).ScheduleTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dkmanager.Registry/ScheduleTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServer).ScheduleTask(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Registry_ReleaseResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReleaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServer).ReleaseResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dkmanager.Registry/ReleaseResource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServer).ReleaseResource(ctx, req.(*ReleaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Registry_ServiceDesc is the grpc.ServiceDesc for Registry service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Registry_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dkmanager.Registry",
	HandlerType: (*RegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Registry_Register_Handler,
		},
		{
			MethodName: "ReportNodeStatus",
			Handler:    _Registry_ReportNodeStatus_Handler,
		},
		{
			MethodName: "ScheduleTask",
			Handler:    _Registry_ScheduleTask_Handler,
		},
		{
			MethodName: "ReleaseResource",
			Handler:    _Registry_ReleaseResource_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "register.proto",
}