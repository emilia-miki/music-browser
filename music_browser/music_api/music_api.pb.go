// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: music_api.proto

package music_api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Query struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Query string `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
}

func (x *Query) Reset() {
	*x = Query{}
	if protoimpl.UnsafeEnabled {
		mi := &file_music_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Query) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Query) ProtoMessage() {}

func (x *Query) ProtoReflect() protoreflect.Message {
	mi := &file_music_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Query.ProtoReflect.Descriptor instead.
func (*Query) Descriptor() ([]byte, []int) {
	return file_music_api_proto_rawDescGZIP(), []int{0}
}

func (x *Query) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

type Url struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *Url) Reset() {
	*x = Url{}
	if protoimpl.UnsafeEnabled {
		mi := &file_music_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Url) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Url) ProtoMessage() {}

func (x *Url) ProtoReflect() protoreflect.Message {
	mi := &file_music_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Url.ProtoReflect.Descriptor instead.
func (*Url) Descriptor() ([]byte, []int) {
	return file_music_api_proto_rawDescGZIP(), []int{1}
}

func (x *Url) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type Artist struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url      *string `protobuf:"bytes,1,opt,name=url,proto3,oneof" json:"url,omitempty"`
	ImageUrl *string `protobuf:"bytes,2,opt,name=image_url,json=imageUrl,proto3,oneof" json:"image_url,omitempty"`
	Name     *string `protobuf:"bytes,3,opt,name=name,proto3,oneof" json:"name,omitempty"`
}

func (x *Artist) Reset() {
	*x = Artist{}
	if protoimpl.UnsafeEnabled {
		mi := &file_music_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Artist) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Artist) ProtoMessage() {}

func (x *Artist) ProtoReflect() protoreflect.Message {
	mi := &file_music_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Artist.ProtoReflect.Descriptor instead.
func (*Artist) Descriptor() ([]byte, []int) {
	return file_music_api_proto_rawDescGZIP(), []int{2}
}

func (x *Artist) GetUrl() string {
	if x != nil && x.Url != nil {
		return *x.Url
	}
	return ""
}

func (x *Artist) GetImageUrl() string {
	if x != nil && x.ImageUrl != nil {
		return *x.ImageUrl
	}
	return ""
}

func (x *Artist) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

type Artists struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Artists []*Artist `protobuf:"bytes,1,rep,name=artists,proto3" json:"artists,omitempty"`
}

func (x *Artists) Reset() {
	*x = Artists{}
	if protoimpl.UnsafeEnabled {
		mi := &file_music_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Artists) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Artists) ProtoMessage() {}

func (x *Artists) ProtoReflect() protoreflect.Message {
	mi := &file_music_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Artists.ProtoReflect.Descriptor instead.
func (*Artists) Descriptor() ([]byte, []int) {
	return file_music_api_proto_rawDescGZIP(), []int{3}
}

func (x *Artists) GetArtists() []*Artist {
	if x != nil {
		return x.Artists
	}
	return nil
}

type Album struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url       *string `protobuf:"bytes,1,opt,name=url,proto3,oneof" json:"url,omitempty"`
	ImageUrl  *string `protobuf:"bytes,2,opt,name=image_url,json=imageUrl,proto3,oneof" json:"image_url,omitempty"`
	ArtistUrl *string `protobuf:"bytes,3,opt,name=artist_url,json=artistUrl,proto3,oneof" json:"artist_url,omitempty"`
	Name      *string `protobuf:"bytes,4,opt,name=name,proto3,oneof" json:"name,omitempty"`
	Year      *uint32 `protobuf:"varint,5,opt,name=year,proto3,oneof" json:"year,omitempty"`
}

func (x *Album) Reset() {
	*x = Album{}
	if protoimpl.UnsafeEnabled {
		mi := &file_music_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Album) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Album) ProtoMessage() {}

func (x *Album) ProtoReflect() protoreflect.Message {
	mi := &file_music_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Album.ProtoReflect.Descriptor instead.
func (*Album) Descriptor() ([]byte, []int) {
	return file_music_api_proto_rawDescGZIP(), []int{4}
}

func (x *Album) GetUrl() string {
	if x != nil && x.Url != nil {
		return *x.Url
	}
	return ""
}

func (x *Album) GetImageUrl() string {
	if x != nil && x.ImageUrl != nil {
		return *x.ImageUrl
	}
	return ""
}

func (x *Album) GetArtistUrl() string {
	if x != nil && x.ArtistUrl != nil {
		return *x.ArtistUrl
	}
	return ""
}

func (x *Album) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *Album) GetYear() uint32 {
	if x != nil && x.Year != nil {
		return *x.Year
	}
	return 0
}

type Albums struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Albums []*Album `protobuf:"bytes,1,rep,name=albums,proto3" json:"albums,omitempty"`
}

func (x *Albums) Reset() {
	*x = Albums{}
	if protoimpl.UnsafeEnabled {
		mi := &file_music_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Albums) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Albums) ProtoMessage() {}

func (x *Albums) ProtoReflect() protoreflect.Message {
	mi := &file_music_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Albums.ProtoReflect.Descriptor instead.
func (*Albums) Descriptor() ([]byte, []int) {
	return file_music_api_proto_rawDescGZIP(), []int{5}
}

func (x *Albums) GetAlbums() []*Album {
	if x != nil {
		return x.Albums
	}
	return nil
}

type Track struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url             *string `protobuf:"bytes,1,opt,name=url,proto3,oneof" json:"url,omitempty"`
	ImageUrl        *string `protobuf:"bytes,2,opt,name=image_url,json=imageUrl,proto3,oneof" json:"image_url,omitempty"`
	ArtistUrl       *string `protobuf:"bytes,3,opt,name=artist_url,json=artistUrl,proto3,oneof" json:"artist_url,omitempty"`
	AlbumUrl        *string `protobuf:"bytes,4,opt,name=album_url,json=albumUrl,proto3,oneof" json:"album_url,omitempty"`
	Name            *string `protobuf:"bytes,5,opt,name=name,proto3,oneof" json:"name,omitempty"`
	DurationSeconds *uint32 `protobuf:"varint,6,opt,name=duration_seconds,json=durationSeconds,proto3,oneof" json:"duration_seconds,omitempty"`
}

func (x *Track) Reset() {
	*x = Track{}
	if protoimpl.UnsafeEnabled {
		mi := &file_music_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Track) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Track) ProtoMessage() {}

func (x *Track) ProtoReflect() protoreflect.Message {
	mi := &file_music_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Track.ProtoReflect.Descriptor instead.
func (*Track) Descriptor() ([]byte, []int) {
	return file_music_api_proto_rawDescGZIP(), []int{6}
}

func (x *Track) GetUrl() string {
	if x != nil && x.Url != nil {
		return *x.Url
	}
	return ""
}

func (x *Track) GetImageUrl() string {
	if x != nil && x.ImageUrl != nil {
		return *x.ImageUrl
	}
	return ""
}

func (x *Track) GetArtistUrl() string {
	if x != nil && x.ArtistUrl != nil {
		return *x.ArtistUrl
	}
	return ""
}

func (x *Track) GetAlbumUrl() string {
	if x != nil && x.AlbumUrl != nil {
		return *x.AlbumUrl
	}
	return ""
}

func (x *Track) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *Track) GetDurationSeconds() uint32 {
	if x != nil && x.DurationSeconds != nil {
		return *x.DurationSeconds
	}
	return 0
}

type Tracks struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tracks []*Track `protobuf:"bytes,1,rep,name=tracks,proto3" json:"tracks,omitempty"`
}

func (x *Tracks) Reset() {
	*x = Tracks{}
	if protoimpl.UnsafeEnabled {
		mi := &file_music_api_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tracks) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tracks) ProtoMessage() {}

func (x *Tracks) ProtoReflect() protoreflect.Message {
	mi := &file_music_api_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tracks.ProtoReflect.Descriptor instead.
func (*Tracks) Descriptor() ([]byte, []int) {
	return file_music_api_proto_rawDescGZIP(), []int{7}
}

func (x *Tracks) GetTracks() []*Track {
	if x != nil {
		return x.Tracks
	}
	return nil
}

type ArtistWithAlbums struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Artist *Artist `protobuf:"bytes,1,opt,name=artist,proto3" json:"artist,omitempty"`
	Albums *Albums `protobuf:"bytes,2,opt,name=albums,proto3" json:"albums,omitempty"`
}

func (x *ArtistWithAlbums) Reset() {
	*x = ArtistWithAlbums{}
	if protoimpl.UnsafeEnabled {
		mi := &file_music_api_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ArtistWithAlbums) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArtistWithAlbums) ProtoMessage() {}

func (x *ArtistWithAlbums) ProtoReflect() protoreflect.Message {
	mi := &file_music_api_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArtistWithAlbums.ProtoReflect.Descriptor instead.
func (*ArtistWithAlbums) Descriptor() ([]byte, []int) {
	return file_music_api_proto_rawDescGZIP(), []int{8}
}

func (x *ArtistWithAlbums) GetArtist() *Artist {
	if x != nil {
		return x.Artist
	}
	return nil
}

func (x *ArtistWithAlbums) GetAlbums() *Albums {
	if x != nil {
		return x.Albums
	}
	return nil
}

type AlbumWithTracks struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Album  *Album  `protobuf:"bytes,1,opt,name=album,proto3" json:"album,omitempty"`
	Tracks *Tracks `protobuf:"bytes,2,opt,name=tracks,proto3" json:"tracks,omitempty"`
}

func (x *AlbumWithTracks) Reset() {
	*x = AlbumWithTracks{}
	if protoimpl.UnsafeEnabled {
		mi := &file_music_api_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AlbumWithTracks) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AlbumWithTracks) ProtoMessage() {}

func (x *AlbumWithTracks) ProtoReflect() protoreflect.Message {
	mi := &file_music_api_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AlbumWithTracks.ProtoReflect.Descriptor instead.
func (*AlbumWithTracks) Descriptor() ([]byte, []int) {
	return file_music_api_proto_rawDescGZIP(), []int{9}
}

func (x *AlbumWithTracks) GetAlbum() *Album {
	if x != nil {
		return x.Album
	}
	return nil
}

func (x *AlbumWithTracks) GetTracks() *Tracks {
	if x != nil {
		return x.Tracks
	}
	return nil
}

type TrackWithAlbumAndArtist struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Artist *Artist `protobuf:"bytes,1,opt,name=artist,proto3" json:"artist,omitempty"`
	Album  *Album  `protobuf:"bytes,2,opt,name=album,proto3" json:"album,omitempty"`
	Track  *Track  `protobuf:"bytes,3,opt,name=track,proto3" json:"track,omitempty"`
}

func (x *TrackWithAlbumAndArtist) Reset() {
	*x = TrackWithAlbumAndArtist{}
	if protoimpl.UnsafeEnabled {
		mi := &file_music_api_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TrackWithAlbumAndArtist) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrackWithAlbumAndArtist) ProtoMessage() {}

func (x *TrackWithAlbumAndArtist) ProtoReflect() protoreflect.Message {
	mi := &file_music_api_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrackWithAlbumAndArtist.ProtoReflect.Descriptor instead.
func (*TrackWithAlbumAndArtist) Descriptor() ([]byte, []int) {
	return file_music_api_proto_rawDescGZIP(), []int{10}
}

func (x *TrackWithAlbumAndArtist) GetArtist() *Artist {
	if x != nil {
		return x.Artist
	}
	return nil
}

func (x *TrackWithAlbumAndArtist) GetAlbum() *Album {
	if x != nil {
		return x.Album
	}
	return nil
}

func (x *TrackWithAlbumAndArtist) GetTrack() *Track {
	if x != nil {
		return x.Track
	}
	return nil
}

var File_music_api_proto protoreflect.FileDescriptor

var file_music_api_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x22, 0x1d, 0x0a, 0x05,
	0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x22, 0x17, 0x0a, 0x03, 0x55,
	0x72, 0x6c, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x75, 0x72, 0x6c, 0x22, 0x79, 0x0a, 0x06, 0x41, 0x72, 0x74, 0x69, 0x73, 0x74, 0x12, 0x15,
	0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x03, 0x75,
	0x72, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x20, 0x0a, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x75,
	0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x08, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x55, 0x72, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01,
	0x42, 0x06, 0x0a, 0x04, 0x5f, 0x75, 0x72, 0x6c, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x22,
	0x36, 0x0a, 0x07, 0x41, 0x72, 0x74, 0x69, 0x73, 0x74, 0x73, 0x12, 0x2b, 0x0a, 0x07, 0x61, 0x72,
	0x74, 0x69, 0x73, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6d, 0x75,
	0x73, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x72, 0x74, 0x69, 0x73, 0x74, 0x52, 0x07,
	0x61, 0x72, 0x74, 0x69, 0x73, 0x74, 0x73, 0x22, 0xcd, 0x01, 0x0a, 0x05, 0x41, 0x6c, 0x62, 0x75,
	0x6d, 0x12, 0x15, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x03, 0x75, 0x72, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x20, 0x0a, 0x09, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x08, 0x69,
	0x6d, 0x61, 0x67, 0x65, 0x55, 0x72, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x22, 0x0a, 0x0a, 0x61, 0x72,
	0x74, 0x69, 0x73, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02,
	0x52, 0x09, 0x61, 0x72, 0x74, 0x69, 0x73, 0x74, 0x55, 0x72, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x17,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x03, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x79, 0x65, 0x61, 0x72, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x04, 0x52, 0x04, 0x79, 0x65, 0x61, 0x72, 0x88, 0x01, 0x01,
	0x42, 0x06, 0x0a, 0x04, 0x5f, 0x75, 0x72, 0x6c, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x61, 0x72, 0x74, 0x69, 0x73,
	0x74, 0x5f, 0x75, 0x72, 0x6c, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x07,
	0x0a, 0x05, 0x5f, 0x79, 0x65, 0x61, 0x72, 0x22, 0x32, 0x0a, 0x06, 0x41, 0x6c, 0x62, 0x75, 0x6d,
	0x73, 0x12, 0x28, 0x0a, 0x06, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x6c,
	0x62, 0x75, 0x6d, 0x52, 0x06, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x73, 0x22, 0xa0, 0x02, 0x0a, 0x05,
	0x54, 0x72, 0x61, 0x63, 0x6b, 0x12, 0x15, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x20, 0x0a, 0x09,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x01, 0x52, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x72, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x22,
	0x0a, 0x0a, 0x61, 0x72, 0x74, 0x69, 0x73, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x02, 0x52, 0x09, 0x61, 0x72, 0x74, 0x69, 0x73, 0x74, 0x55, 0x72, 0x6c, 0x88,
	0x01, 0x01, 0x12, 0x20, 0x0a, 0x09, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x5f, 0x75, 0x72, 0x6c, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x03, 0x52, 0x08, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x55, 0x72,
	0x6c, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x04, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x2e, 0x0a,
	0x10, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64,
	0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x05, 0x52, 0x0f, 0x64, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x88, 0x01, 0x01, 0x42, 0x06, 0x0a,
	0x04, 0x5f, 0x75, 0x72, 0x6c, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f,
	0x75, 0x72, 0x6c, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x61, 0x72, 0x74, 0x69, 0x73, 0x74, 0x5f, 0x75,
	0x72, 0x6c, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x5f, 0x75, 0x72, 0x6c,
	0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x64, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x22, 0x32,
	0x0a, 0x06, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x73, 0x12, 0x28, 0x0a, 0x06, 0x74, 0x72, 0x61, 0x63,
	0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63,
	0x5f, 0x61, 0x70, 0x69, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x52, 0x06, 0x74, 0x72, 0x61, 0x63,
	0x6b, 0x73, 0x22, 0x68, 0x0a, 0x10, 0x41, 0x72, 0x74, 0x69, 0x73, 0x74, 0x57, 0x69, 0x74, 0x68,
	0x41, 0x6c, 0x62, 0x75, 0x6d, 0x73, 0x12, 0x29, 0x0a, 0x06, 0x61, 0x72, 0x74, 0x69, 0x73, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x61,
	0x70, 0x69, 0x2e, 0x41, 0x72, 0x74, 0x69, 0x73, 0x74, 0x52, 0x06, 0x61, 0x72, 0x74, 0x69, 0x73,
	0x74, 0x12, 0x29, 0x0a, 0x06, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x11, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x6c,
	0x62, 0x75, 0x6d, 0x73, 0x52, 0x06, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x73, 0x22, 0x64, 0x0a, 0x0f,
	0x41, 0x6c, 0x62, 0x75, 0x6d, 0x57, 0x69, 0x74, 0x68, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x73, 0x12,
	0x26, 0x0a, 0x05, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10,
	0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x6c, 0x62, 0x75, 0x6d,
	0x52, 0x05, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x12, 0x29, 0x0a, 0x06, 0x74, 0x72, 0x61, 0x63, 0x6b,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f,
	0x61, 0x70, 0x69, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x73, 0x52, 0x06, 0x74, 0x72, 0x61, 0x63,
	0x6b, 0x73, 0x22, 0x94, 0x01, 0x0a, 0x17, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x57, 0x69, 0x74, 0x68,
	0x41, 0x6c, 0x62, 0x75, 0x6d, 0x41, 0x6e, 0x64, 0x41, 0x72, 0x74, 0x69, 0x73, 0x74, 0x12, 0x29,
	0x0a, 0x06, 0x61, 0x72, 0x74, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x72, 0x74, 0x69, 0x73,
	0x74, 0x52, 0x06, 0x61, 0x72, 0x74, 0x69, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x05, 0x61, 0x6c, 0x62,
	0x75, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63,
	0x5f, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x6c, 0x62, 0x75, 0x6d, 0x52, 0x05, 0x61, 0x6c, 0x62, 0x75,
	0x6d, 0x12, 0x26, 0x0a, 0x05, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x54, 0x72, 0x61,
	0x63, 0x6b, 0x52, 0x05, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x32, 0xb2, 0x02, 0x0a, 0x08, 0x4d, 0x75,
	0x73, 0x69, 0x63, 0x41, 0x70, 0x69, 0x12, 0x3a, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x41, 0x72, 0x74,
	0x69, 0x73, 0x74, 0x12, 0x0e, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2e,
	0x55, 0x72, 0x6c, 0x1a, 0x1b, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2e,
	0x41, 0x72, 0x74, 0x69, 0x73, 0x74, 0x57, 0x69, 0x74, 0x68, 0x41, 0x6c, 0x62, 0x75, 0x6d, 0x73,
	0x22, 0x00, 0x12, 0x38, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x62, 0x75, 0x6d, 0x12, 0x0e,
	0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x72, 0x6c, 0x1a, 0x1a,
	0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x6c, 0x62, 0x75, 0x6d,
	0x57, 0x69, 0x74, 0x68, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x73, 0x22, 0x00, 0x12, 0x40, 0x0a, 0x08,
	0x47, 0x65, 0x74, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x12, 0x0e, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63,
	0x5f, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x72, 0x6c, 0x1a, 0x22, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63,
	0x5f, 0x61, 0x70, 0x69, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x57, 0x69, 0x74, 0x68, 0x41, 0x6c,
	0x62, 0x75, 0x6d, 0x41, 0x6e, 0x64, 0x41, 0x72, 0x74, 0x69, 0x73, 0x74, 0x22, 0x00, 0x12, 0x37,
	0x0a, 0x0d, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x41, 0x72, 0x74, 0x69, 0x73, 0x74, 0x73, 0x12,
	0x10, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x51, 0x75, 0x65, 0x72,
	0x79, 0x1a, 0x12, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x72,
	0x74, 0x69, 0x73, 0x74, 0x73, 0x22, 0x00, 0x12, 0x35, 0x0a, 0x0c, 0x53, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x41, 0x6c, 0x62, 0x75, 0x6d, 0x73, 0x12, 0x10, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f,
	0x61, 0x70, 0x69, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x1a, 0x11, 0x2e, 0x6d, 0x75, 0x73, 0x69,
	0x63, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x6c, 0x62, 0x75, 0x6d, 0x73, 0x22, 0x00, 0x42, 0x30,
	0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6d, 0x69,
	0x6c, 0x69, 0x61, 0x2d, 0x6d, 0x69, 0x6b, 0x69, 0x2f, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x2d, 0x62,
	0x72, 0x6f, 0x77, 0x73, 0x65, 0x72, 0x2f, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x69,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_music_api_proto_rawDescOnce sync.Once
	file_music_api_proto_rawDescData = file_music_api_proto_rawDesc
)

func file_music_api_proto_rawDescGZIP() []byte {
	file_music_api_proto_rawDescOnce.Do(func() {
		file_music_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_music_api_proto_rawDescData)
	})
	return file_music_api_proto_rawDescData
}

var file_music_api_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_music_api_proto_goTypes = []interface{}{
	(*Query)(nil),                   // 0: music_api.Query
	(*Url)(nil),                     // 1: music_api.Url
	(*Artist)(nil),                  // 2: music_api.Artist
	(*Artists)(nil),                 // 3: music_api.Artists
	(*Album)(nil),                   // 4: music_api.Album
	(*Albums)(nil),                  // 5: music_api.Albums
	(*Track)(nil),                   // 6: music_api.Track
	(*Tracks)(nil),                  // 7: music_api.Tracks
	(*ArtistWithAlbums)(nil),        // 8: music_api.ArtistWithAlbums
	(*AlbumWithTracks)(nil),         // 9: music_api.AlbumWithTracks
	(*TrackWithAlbumAndArtist)(nil), // 10: music_api.TrackWithAlbumAndArtist
}
var file_music_api_proto_depIdxs = []int32{
	2,  // 0: music_api.Artists.artists:type_name -> music_api.Artist
	4,  // 1: music_api.Albums.albums:type_name -> music_api.Album
	6,  // 2: music_api.Tracks.tracks:type_name -> music_api.Track
	2,  // 3: music_api.ArtistWithAlbums.artist:type_name -> music_api.Artist
	5,  // 4: music_api.ArtistWithAlbums.albums:type_name -> music_api.Albums
	4,  // 5: music_api.AlbumWithTracks.album:type_name -> music_api.Album
	7,  // 6: music_api.AlbumWithTracks.tracks:type_name -> music_api.Tracks
	2,  // 7: music_api.TrackWithAlbumAndArtist.artist:type_name -> music_api.Artist
	4,  // 8: music_api.TrackWithAlbumAndArtist.album:type_name -> music_api.Album
	6,  // 9: music_api.TrackWithAlbumAndArtist.track:type_name -> music_api.Track
	1,  // 10: music_api.MusicApi.GetArtist:input_type -> music_api.Url
	1,  // 11: music_api.MusicApi.GetAlbum:input_type -> music_api.Url
	1,  // 12: music_api.MusicApi.GetTrack:input_type -> music_api.Url
	0,  // 13: music_api.MusicApi.SearchArtists:input_type -> music_api.Query
	0,  // 14: music_api.MusicApi.SearchAlbums:input_type -> music_api.Query
	8,  // 15: music_api.MusicApi.GetArtist:output_type -> music_api.ArtistWithAlbums
	9,  // 16: music_api.MusicApi.GetAlbum:output_type -> music_api.AlbumWithTracks
	10, // 17: music_api.MusicApi.GetTrack:output_type -> music_api.TrackWithAlbumAndArtist
	3,  // 18: music_api.MusicApi.SearchArtists:output_type -> music_api.Artists
	5,  // 19: music_api.MusicApi.SearchAlbums:output_type -> music_api.Albums
	15, // [15:20] is the sub-list for method output_type
	10, // [10:15] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_music_api_proto_init() }
func file_music_api_proto_init() {
	if File_music_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_music_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Query); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_music_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Url); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_music_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Artist); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_music_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Artists); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_music_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Album); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_music_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Albums); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_music_api_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Track); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_music_api_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tracks); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_music_api_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ArtistWithAlbums); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_music_api_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AlbumWithTracks); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_music_api_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TrackWithAlbumAndArtist); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_music_api_proto_msgTypes[2].OneofWrappers = []interface{}{}
	file_music_api_proto_msgTypes[4].OneofWrappers = []interface{}{}
	file_music_api_proto_msgTypes[6].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_music_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_music_api_proto_goTypes,
		DependencyIndexes: file_music_api_proto_depIdxs,
		MessageInfos:      file_music_api_proto_msgTypes,
	}.Build()
	File_music_api_proto = out.File
	file_music_api_proto_rawDesc = nil
	file_music_api_proto_goTypes = nil
	file_music_api_proto_depIdxs = nil
}