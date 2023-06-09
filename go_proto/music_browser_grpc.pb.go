// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: music_browser.proto

package go_proto

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
	Browser_GetArtist_FullMethodName     = "/music_api.Browser/GetArtist"
	Browser_GetAlbum_FullMethodName      = "/music_api.Browser/GetAlbum"
	Browser_GetTrack_FullMethodName      = "/music_api.Browser/GetTrack"
	Browser_SearchArtists_FullMethodName = "/music_api.Browser/SearchArtists"
	Browser_SearchAlbums_FullMethodName  = "/music_api.Browser/SearchAlbums"
	Browser_SearchTracks_FullMethodName  = "/music_api.Browser/SearchTracks"
)

// BrowserClient is the client API for Browser service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BrowserClient interface {
	GetArtist(ctx context.Context, in *Id, opts ...grpc.CallOption) (*ArtistWithAlbums, error)
	GetAlbum(ctx context.Context, in *Id, opts ...grpc.CallOption) (*AlbumWithTracks, error)
	GetTrack(ctx context.Context, in *Id, opts ...grpc.CallOption) (*TrackWithAlbumAndArtist, error)
	SearchArtists(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Artists, error)
	SearchAlbums(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Albums, error)
	SearchTracks(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Tracks, error)
}

type browserClient struct {
	cc grpc.ClientConnInterface
}

func NewBrowserClient(cc grpc.ClientConnInterface) BrowserClient {
	return &browserClient{cc}
}

func (c *browserClient) GetArtist(ctx context.Context, in *Id, opts ...grpc.CallOption) (*ArtistWithAlbums, error) {
	out := new(ArtistWithAlbums)
	err := c.cc.Invoke(ctx, Browser_GetArtist_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *browserClient) GetAlbum(ctx context.Context, in *Id, opts ...grpc.CallOption) (*AlbumWithTracks, error) {
	out := new(AlbumWithTracks)
	err := c.cc.Invoke(ctx, Browser_GetAlbum_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *browserClient) GetTrack(ctx context.Context, in *Id, opts ...grpc.CallOption) (*TrackWithAlbumAndArtist, error) {
	out := new(TrackWithAlbumAndArtist)
	err := c.cc.Invoke(ctx, Browser_GetTrack_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *browserClient) SearchArtists(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Artists, error) {
	out := new(Artists)
	err := c.cc.Invoke(ctx, Browser_SearchArtists_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *browserClient) SearchAlbums(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Albums, error) {
	out := new(Albums)
	err := c.cc.Invoke(ctx, Browser_SearchAlbums_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *browserClient) SearchTracks(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Tracks, error) {
	out := new(Tracks)
	err := c.cc.Invoke(ctx, Browser_SearchTracks_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BrowserServer is the server API for Browser service.
// All implementations must embed UnimplementedBrowserServer
// for forward compatibility
type BrowserServer interface {
	GetArtist(context.Context, *Id) (*ArtistWithAlbums, error)
	GetAlbum(context.Context, *Id) (*AlbumWithTracks, error)
	GetTrack(context.Context, *Id) (*TrackWithAlbumAndArtist, error)
	SearchArtists(context.Context, *Query) (*Artists, error)
	SearchAlbums(context.Context, *Query) (*Albums, error)
	SearchTracks(context.Context, *Query) (*Tracks, error)
	mustEmbedUnimplementedBrowserServer()
}

// UnimplementedBrowserServer must be embedded to have forward compatible implementations.
type UnimplementedBrowserServer struct {
}

func (UnimplementedBrowserServer) GetArtist(context.Context, *Id) (*ArtistWithAlbums, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArtist not implemented")
}
func (UnimplementedBrowserServer) GetAlbum(context.Context, *Id) (*AlbumWithTracks, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAlbum not implemented")
}
func (UnimplementedBrowserServer) GetTrack(context.Context, *Id) (*TrackWithAlbumAndArtist, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTrack not implemented")
}
func (UnimplementedBrowserServer) SearchArtists(context.Context, *Query) (*Artists, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchArtists not implemented")
}
func (UnimplementedBrowserServer) SearchAlbums(context.Context, *Query) (*Albums, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchAlbums not implemented")
}
func (UnimplementedBrowserServer) SearchTracks(context.Context, *Query) (*Tracks, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchTracks not implemented")
}
func (UnimplementedBrowserServer) mustEmbedUnimplementedBrowserServer() {}

// UnsafeBrowserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BrowserServer will
// result in compilation errors.
type UnsafeBrowserServer interface {
	mustEmbedUnimplementedBrowserServer()
}

func RegisterBrowserServer(s grpc.ServiceRegistrar, srv BrowserServer) {
	s.RegisterService(&Browser_ServiceDesc, srv)
}

func _Browser_GetArtist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Id)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrowserServer).GetArtist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Browser_GetArtist_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrowserServer).GetArtist(ctx, req.(*Id))
	}
	return interceptor(ctx, in, info, handler)
}

func _Browser_GetAlbum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Id)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrowserServer).GetAlbum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Browser_GetAlbum_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrowserServer).GetAlbum(ctx, req.(*Id))
	}
	return interceptor(ctx, in, info, handler)
}

func _Browser_GetTrack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Id)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrowserServer).GetTrack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Browser_GetTrack_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrowserServer).GetTrack(ctx, req.(*Id))
	}
	return interceptor(ctx, in, info, handler)
}

func _Browser_SearchArtists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Query)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrowserServer).SearchArtists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Browser_SearchArtists_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrowserServer).SearchArtists(ctx, req.(*Query))
	}
	return interceptor(ctx, in, info, handler)
}

func _Browser_SearchAlbums_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Query)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrowserServer).SearchAlbums(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Browser_SearchAlbums_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrowserServer).SearchAlbums(ctx, req.(*Query))
	}
	return interceptor(ctx, in, info, handler)
}

func _Browser_SearchTracks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Query)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrowserServer).SearchTracks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Browser_SearchTracks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrowserServer).SearchTracks(ctx, req.(*Query))
	}
	return interceptor(ctx, in, info, handler)
}

// Browser_ServiceDesc is the grpc.ServiceDesc for Browser service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Browser_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "music_api.Browser",
	HandlerType: (*BrowserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetArtist",
			Handler:    _Browser_GetArtist_Handler,
		},
		{
			MethodName: "GetAlbum",
			Handler:    _Browser_GetAlbum_Handler,
		},
		{
			MethodName: "GetTrack",
			Handler:    _Browser_GetTrack_Handler,
		},
		{
			MethodName: "SearchArtists",
			Handler:    _Browser_SearchArtists_Handler,
		},
		{
			MethodName: "SearchAlbums",
			Handler:    _Browser_SearchAlbums_Handler,
		},
		{
			MethodName: "SearchTracks",
			Handler:    _Browser_SearchTracks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "music_browser.proto",
}

const (
	GetUpNext_GetUpNext_FullMethodName = "/music_api.GetUpNext/GetUpNext"
)

// GetUpNextClient is the client API for GetUpNext service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GetUpNextClient interface {
	GetUpNext(ctx context.Context, in *Id, opts ...grpc.CallOption) (*Tracks, error)
}

type getUpNextClient struct {
	cc grpc.ClientConnInterface
}

func NewGetUpNextClient(cc grpc.ClientConnInterface) GetUpNextClient {
	return &getUpNextClient{cc}
}

func (c *getUpNextClient) GetUpNext(ctx context.Context, in *Id, opts ...grpc.CallOption) (*Tracks, error) {
	out := new(Tracks)
	err := c.cc.Invoke(ctx, GetUpNext_GetUpNext_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetUpNextServer is the server API for GetUpNext service.
// All implementations must embed UnimplementedGetUpNextServer
// for forward compatibility
type GetUpNextServer interface {
	GetUpNext(context.Context, *Id) (*Tracks, error)
	mustEmbedUnimplementedGetUpNextServer()
}

// UnimplementedGetUpNextServer must be embedded to have forward compatible implementations.
type UnimplementedGetUpNextServer struct {
}

func (UnimplementedGetUpNextServer) GetUpNext(context.Context, *Id) (*Tracks, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUpNext not implemented")
}
func (UnimplementedGetUpNextServer) mustEmbedUnimplementedGetUpNextServer() {}

// UnsafeGetUpNextServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GetUpNextServer will
// result in compilation errors.
type UnsafeGetUpNextServer interface {
	mustEmbedUnimplementedGetUpNextServer()
}

func RegisterGetUpNextServer(s grpc.ServiceRegistrar, srv GetUpNextServer) {
	s.RegisterService(&GetUpNext_ServiceDesc, srv)
}

func _GetUpNext_GetUpNext_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Id)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GetUpNextServer).GetUpNext(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GetUpNext_GetUpNext_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GetUpNextServer).GetUpNext(ctx, req.(*Id))
	}
	return interceptor(ctx, in, info, handler)
}

// GetUpNext_ServiceDesc is the grpc.ServiceDesc for GetUpNext service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GetUpNext_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "music_api.GetUpNext",
	HandlerType: (*GetUpNextServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUpNext",
			Handler:    _GetUpNext_GetUpNext_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "music_browser.proto",
}

const (
	GetByTag_GetByTag_FullMethodName = "/music_api.GetByTag/GetByTag"
)

// GetByTagClient is the client API for GetByTag service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GetByTagClient interface {
	GetByTag(ctx context.Context, in *Tag, opts ...grpc.CallOption) (*Tracks, error)
}

type getByTagClient struct {
	cc grpc.ClientConnInterface
}

func NewGetByTagClient(cc grpc.ClientConnInterface) GetByTagClient {
	return &getByTagClient{cc}
}

func (c *getByTagClient) GetByTag(ctx context.Context, in *Tag, opts ...grpc.CallOption) (*Tracks, error) {
	out := new(Tracks)
	err := c.cc.Invoke(ctx, GetByTag_GetByTag_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetByTagServer is the server API for GetByTag service.
// All implementations must embed UnimplementedGetByTagServer
// for forward compatibility
type GetByTagServer interface {
	GetByTag(context.Context, *Tag) (*Tracks, error)
	mustEmbedUnimplementedGetByTagServer()
}

// UnimplementedGetByTagServer must be embedded to have forward compatible implementations.
type UnimplementedGetByTagServer struct {
}

func (UnimplementedGetByTagServer) GetByTag(context.Context, *Tag) (*Tracks, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByTag not implemented")
}
func (UnimplementedGetByTagServer) mustEmbedUnimplementedGetByTagServer() {}

// UnsafeGetByTagServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GetByTagServer will
// result in compilation errors.
type UnsafeGetByTagServer interface {
	mustEmbedUnimplementedGetByTagServer()
}

func RegisterGetByTagServer(s grpc.ServiceRegistrar, srv GetByTagServer) {
	s.RegisterService(&GetByTag_ServiceDesc, srv)
}

func _GetByTag_GetByTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Tag)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GetByTagServer).GetByTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GetByTag_GetByTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GetByTagServer).GetByTag(ctx, req.(*Tag))
	}
	return interceptor(ctx, in, info, handler)
}

// GetByTag_ServiceDesc is the grpc.ServiceDesc for GetByTag service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GetByTag_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "music_api.GetByTag",
	HandlerType: (*GetByTagServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetByTag",
			Handler:    _GetByTag_GetByTag_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "music_browser.proto",
}
