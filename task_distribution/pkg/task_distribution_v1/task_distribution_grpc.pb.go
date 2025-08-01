// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.29.3
// source: task_distribution.proto

package task_distribution_v1

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

// TaskServiceClient is the client API for TaskService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TaskServiceClient interface {
	SubmitFileProcessingTask(ctx context.Context, in *SubmitFileProcessingTaskRequest, opts ...grpc.CallOption) (*SubmitFileProcessingTaskResponse, error)
	GetTaskStatus(ctx context.Context, in *GetTaskStatusRequest, opts ...grpc.CallOption) (*GetTaskStatusResponse, error)
}

type taskServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTaskServiceClient(cc grpc.ClientConnInterface) TaskServiceClient {
	return &taskServiceClient{cc}
}

func (c *taskServiceClient) SubmitFileProcessingTask(ctx context.Context, in *SubmitFileProcessingTaskRequest, opts ...grpc.CallOption) (*SubmitFileProcessingTaskResponse, error) {
	out := new(SubmitFileProcessingTaskResponse)
	err := c.cc.Invoke(ctx, "/taskdistribution.v1.TaskService/SubmitFileProcessingTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) GetTaskStatus(ctx context.Context, in *GetTaskStatusRequest, opts ...grpc.CallOption) (*GetTaskStatusResponse, error) {
	out := new(GetTaskStatusResponse)
	err := c.cc.Invoke(ctx, "/taskdistribution.v1.TaskService/GetTaskStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TaskServiceServer is the server API for TaskService service.
// All implementations must embed UnimplementedTaskServiceServer
// for forward compatibility
type TaskServiceServer interface {
	SubmitFileProcessingTask(context.Context, *SubmitFileProcessingTaskRequest) (*SubmitFileProcessingTaskResponse, error)
	GetTaskStatus(context.Context, *GetTaskStatusRequest) (*GetTaskStatusResponse, error)
	mustEmbedUnimplementedTaskServiceServer()
}

// UnimplementedTaskServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTaskServiceServer struct {
}

func (UnimplementedTaskServiceServer) SubmitFileProcessingTask(context.Context, *SubmitFileProcessingTaskRequest) (*SubmitFileProcessingTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitFileProcessingTask not implemented")
}
func (UnimplementedTaskServiceServer) GetTaskStatus(context.Context, *GetTaskStatusRequest) (*GetTaskStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTaskStatus not implemented")
}
func (UnimplementedTaskServiceServer) mustEmbedUnimplementedTaskServiceServer() {}

// UnsafeTaskServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TaskServiceServer will
// result in compilation errors.
type UnsafeTaskServiceServer interface {
	mustEmbedUnimplementedTaskServiceServer()
}

func RegisterTaskServiceServer(s grpc.ServiceRegistrar, srv TaskServiceServer) {
	s.RegisterService(&TaskService_ServiceDesc, srv)
}

func _TaskService_SubmitFileProcessingTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitFileProcessingTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).SubmitFileProcessingTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/taskdistribution.v1.TaskService/SubmitFileProcessingTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).SubmitFileProcessingTask(ctx, req.(*SubmitFileProcessingTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_GetTaskStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTaskStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).GetTaskStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/taskdistribution.v1.TaskService/GetTaskStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).GetTaskStatus(ctx, req.(*GetTaskStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TaskService_ServiceDesc is the grpc.ServiceDesc for TaskService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TaskService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "taskdistribution.v1.TaskService",
	HandlerType: (*TaskServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitFileProcessingTask",
			Handler:    _TaskService_SubmitFileProcessingTask_Handler,
		},
		{
			MethodName: "GetTaskStatus",
			Handler:    _TaskService_GetTaskStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "task_distribution.proto",
}
