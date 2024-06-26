// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: iam/v1/iam.proto

package iam

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

const (
	IamService_UpsertUser_FullMethodName         = "/iam.v1.IamService/UpsertUser"
	IamService_GetUser_FullMethodName            = "/iam.v1.IamService/GetUser"
	IamService_Login_FullMethodName              = "/iam.v1.IamService/Login"
	IamService_RefreshToken_FullMethodName       = "/iam.v1.IamService/RefreshToken"
	IamService_GetListUser_FullMethodName        = "/iam.v1.IamService/GetListUser"
	IamService_DeleteUser_FullMethodName         = "/iam.v1.IamService/DeleteUser"
	IamService_GetCurrentUserInfo_FullMethodName = "/iam.v1.IamService/GetCurrentUserInfo"
)

// IamServiceClient is the client API for IamService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IamServiceClient interface {
	// UpsertUser ...
	UpsertUser(ctx context.Context, in *UpsertUserRequest, opts ...grpc.CallOption) (*UpsertUserResponse, error)
	// GetUser ...
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
	// Login ...
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	// Refresh token ...
	RefreshToken(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*RefreshTokenResponse, error)
	// GetListUser ...
	GetListUser(ctx context.Context, in *GetListUserRequest, opts ...grpc.CallOption) (*GetListUserResponse, error)
	// DeleteUser ...
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error)
	// GetCurrentUserInfo ...
	GetCurrentUserInfo(ctx context.Context, in *GetCurrentUserInfoRequest, opts ...grpc.CallOption) (*GetCurrentUserInfoResponse, error)
}

type iamServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIamServiceClient(cc grpc.ClientConnInterface) IamServiceClient {
	return &iamServiceClient{cc}
}

func (c *iamServiceClient) UpsertUser(ctx context.Context, in *UpsertUserRequest, opts ...grpc.CallOption) (*UpsertUserResponse, error) {
	out := new(UpsertUserResponse)
	err := c.cc.Invoke(ctx, IamService_UpsertUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamServiceClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	out := new(GetUserResponse)
	err := c.cc.Invoke(ctx, IamService_GetUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, IamService_Login_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamServiceClient) RefreshToken(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*RefreshTokenResponse, error) {
	out := new(RefreshTokenResponse)
	err := c.cc.Invoke(ctx, IamService_RefreshToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamServiceClient) GetListUser(ctx context.Context, in *GetListUserRequest, opts ...grpc.CallOption) (*GetListUserResponse, error) {
	out := new(GetListUserResponse)
	err := c.cc.Invoke(ctx, IamService_GetListUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamServiceClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error) {
	out := new(DeleteUserResponse)
	err := c.cc.Invoke(ctx, IamService_DeleteUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iamServiceClient) GetCurrentUserInfo(ctx context.Context, in *GetCurrentUserInfoRequest, opts ...grpc.CallOption) (*GetCurrentUserInfoResponse, error) {
	out := new(GetCurrentUserInfoResponse)
	err := c.cc.Invoke(ctx, IamService_GetCurrentUserInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IamServiceServer is the server API for IamService service.
// All implementations must embed UnimplementedIamServiceServer
// for forward compatibility
type IamServiceServer interface {
	// UpsertUser ...
	UpsertUser(context.Context, *UpsertUserRequest) (*UpsertUserResponse, error)
	// GetUser ...
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	// Login ...
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	// Refresh token ...
	RefreshToken(context.Context, *RefreshTokenRequest) (*RefreshTokenResponse, error)
	// GetListUser ...
	GetListUser(context.Context, *GetListUserRequest) (*GetListUserResponse, error)
	// DeleteUser ...
	DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error)
	// GetCurrentUserInfo ...
	GetCurrentUserInfo(context.Context, *GetCurrentUserInfoRequest) (*GetCurrentUserInfoResponse, error)
	mustEmbedUnimplementedIamServiceServer()
}

// UnimplementedIamServiceServer must be embedded to have forward compatible implementations.
type UnimplementedIamServiceServer struct {
}

func (UnimplementedIamServiceServer) UpsertUser(context.Context, *UpsertUserRequest) (*UpsertUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertUser not implemented")
}
func (UnimplementedIamServiceServer) GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedIamServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedIamServiceServer) RefreshToken(context.Context, *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshToken not implemented")
}
func (UnimplementedIamServiceServer) GetListUser(context.Context, *GetListUserRequest) (*GetListUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListUser not implemented")
}
func (UnimplementedIamServiceServer) DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedIamServiceServer) GetCurrentUserInfo(context.Context, *GetCurrentUserInfoRequest) (*GetCurrentUserInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCurrentUserInfo not implemented")
}
func (UnimplementedIamServiceServer) mustEmbedUnimplementedIamServiceServer() {}

// UnsafeIamServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IamServiceServer will
// result in compilation errors.
type UnsafeIamServiceServer interface {
	mustEmbedUnimplementedIamServiceServer()
}

func RegisterIamServiceServer(s grpc.ServiceRegistrar, srv IamServiceServer) {
	s.RegisterService(&IamService_ServiceDesc, srv)
}

func _IamService_UpsertUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpsertUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamServiceServer).UpsertUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IamService_UpsertUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamServiceServer).UpsertUser(ctx, req.(*UpsertUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IamService_GetUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamServiceServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IamService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamService_RefreshToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamServiceServer).RefreshToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IamService_RefreshToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamServiceServer).RefreshToken(ctx, req.(*RefreshTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamService_GetListUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamServiceServer).GetListUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IamService_GetListUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamServiceServer).GetListUser(ctx, req.(*GetListUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamService_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamServiceServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IamService_DeleteUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamServiceServer).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IamService_GetCurrentUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCurrentUserInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IamServiceServer).GetCurrentUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IamService_GetCurrentUserInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IamServiceServer).GetCurrentUserInfo(ctx, req.(*GetCurrentUserInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// IamService_ServiceDesc is the grpc.ServiceDesc for IamService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IamService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "iam.v1.IamService",
	HandlerType: (*IamServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpsertUser",
			Handler:    _IamService_UpsertUser_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _IamService_GetUser_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _IamService_Login_Handler,
		},
		{
			MethodName: "RefreshToken",
			Handler:    _IamService_RefreshToken_Handler,
		},
		{
			MethodName: "GetListUser",
			Handler:    _IamService_GetListUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _IamService_DeleteUser_Handler,
		},
		{
			MethodName: "GetCurrentUserInfo",
			Handler:    _IamService_GetCurrentUserInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "iam/v1/iam.proto",
}
