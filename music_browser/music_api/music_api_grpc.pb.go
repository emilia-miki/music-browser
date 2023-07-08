// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: music_api.proto

package music_api

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
	MusicApi_GetArtist_FullMethodName     = "/music_api.MusicApi/GetArtist"
	MusicApi_GetAlbum_FullMethodName      = "/music_api.MusicApi/GetAlbum"
	MusicApi_GetTrack_FullMethodName      = "/music_api.MusicApi/GetTrack"
	MusicApi_SearchArtists_FullMethodName = "/music_api.MusicApi/SearchArtists"
	MusicApi_SearchAlbums_FullMethodName  = "/music_api.MusicApi/SearchAlbums"
)

// MusicApiClient is the client API for MusicApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MusicApiClient interface {
	GetArtist(ctx context.Context, in *Url, opts ...grpc.CallOption) (*ArtistWithAlbums, error)
	GetAlbum(ctx context.Context, in *Url, opts ...grpc.CallOption) (*AlbumWithTracks, error)
	GetTrack(ctx context.Context, in *Url, opts ...grpc.CallOption) (*TrackWithAlbumAndArtist, error)
	SearchArtists(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Artists, error)
	SearchAlbums(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Albums, error)
}

type musicApiClient struct {
	cc grpc.ClientConnInterface
}

func NewMusicApiClient(cc grpc.ClientConnInterface) MusicApiClient {
	return &musicApiClient{cc}
}

func (c *musicApiClient) GetArtist(ctx context.Context, in *Url, opts ...grpc.CallOption) (*ArtistWithAlbums, error) {
	out := new(ArtistWithAlbums)
	err := c.cc.Invoke(ctx, MusicApi_GetArtist_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *musicApiClient) GetAlbum(ctx context.Context, in *Url, opts ...grpc.CallOption) (*AlbumWithTracks, error) {
	out := new(AlbumWithTracks)
	err := c.cc.Invoke(ctx, MusicApi_GetAlbum_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *musicApiClient) GetTrack(ctx context.Context, in *Url, opts ...grpc.CallOption) (*TrackWithAlbumAndArtist, error) {
	out := new(TrackWithAlbumAndArtist)
	err := c.cc.Invoke(ctx, MusicApi_GetTrack_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *musicApiClient) SearchArtists(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Artists, error) {
	out := new(Artists)
	err := c.cc.Invoke(ctx, MusicApi_SearchArtists_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *musicApiClient) SearchAlbums(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Albums, error) {
	out := new(Albums)
	err := c.cc.Invoke(ctx, MusicApi_SearchAlbums_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MusicApiServer is the server API for MusicApi service.
// All implementations must embed UnimplementedMusicApiServer
// for forward compatibility
type MusicApiServer interface {
	GetArtist(context.Context, *Url) (*ArtistWithAlbums, error)
	GetAlbum(context.Context, *Url) (*AlbumWithTracks, error)
	GetTrack(context.Context, *Url) (*TrackWithAlbumAndArtist, error)
	SearchArtists(context.Context, *Query) (*Artists, error)
	SearchAlbums(context.Context, *Query) (*Albums, error)
	mustEmbedUnimplementedMusicApiServer()
}

// UnimplementedMusicApiServer must be embedded to have forward compatible implementations.
type UnimplementedMusicApiServer struct {
}

func (UnimplementedMusicApiServer) GetArtist(context.Context, *Url) (*ArtistWithAlbums, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArtist not implemented")
}
func (UnimplementedMusicApiServer) GetAlbum(context.Context, *Url) (*AlbumWithTracks, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAlbum not implemented")
}
func (UnimplementedMusicApiServer) GetTrack(context.Context, *Url) (*TrackWithAlbumAndArtist, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTrack not implemented")
}
func (UnimplementedMusicApiServer) SearchArtists(context.Context, *Query) (*Artists, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchArtists not implemented")
}
func (UnimplementedMusicApiServer) SearchAlbums(context.Context, *Query) (*Albums, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchAlbums not implemented")
}
func (UnimplementedMusicApiServer) mustEmbedUnimplementedMusicApiServer() {}

// UnsafeMusicApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MusicApiServer will
// result in compilation errors.
type UnsafeMusicApiServer interface {
	mustEmbedUnimplementedMusicApiServer()
}

func RegisterMusicApiServer(s grpc.ServiceRegistrar, srv MusicApiServer) {
	s.RegisterService(&MusicApi_ServiceDesc, srv)
}

func _MusicApi_GetArtist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Url)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MusicApiServer).GetArtist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MusicApi_GetArtist_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MusicApiServer).GetArtist(ctx, req.(*Url))
	}
	return interceptor(ctx, in, info, handler)
}

func _MusicApi_GetAlbum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Url)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MusicApiServer).GetAlbum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MusicApi_GetAlbum_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MusicApiServer).GetAlbum(ctx, req.(*Url))
	}
	return interceptor(ctx, in, info, handler)
}

func _MusicApi_GetTrack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Url)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MusicApiServer).GetTrack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MusicApi_GetTrack_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MusicApiServer).GetTrack(ctx, req.(*Url))
	}
	return interceptor(ctx, in, info, handler)
}

func _MusicApi_SearchArtists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Query)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MusicApiServer).SearchArtists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MusicApi_SearchArtists_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MusicApiServer).SearchArtists(ctx, req.(*Query))
	}
	return interceptor(ctx, in, info, handler)
}

func _MusicApi_SearchAlbums_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Query)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MusicApiServer).SearchAlbums(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MusicApi_SearchAlbums_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MusicApiServer).SearchAlbums(ctx, req.(*Query))
	}
	return interceptor(ctx, in, info, handler)
}

// MusicApi_ServiceDesc is the grpc.ServiceDesc for MusicApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MusicApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "music_api.MusicApi",
	HandlerType: (*MusicApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetArtist",
			Handler:    _MusicApi_GetArtist_Handler,
		},
		{
			MethodName: "GetAlbum",
			Handler:    _MusicApi_GetAlbum_Handler,
		},
		{
			MethodName: "GetTrack",
			Handler:    _MusicApi_GetTrack_Handler,
		},
		{
			MethodName: "SearchArtists",
			Handler:    _MusicApi_SearchArtists_Handler,
		},
		{
			MethodName: "SearchAlbums",
			Handler:    _MusicApi_SearchAlbums_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "music_api.proto",
}