from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Album(_message.Message):
    __slots__ = ["artists", "duration_seconds", "image_url", "name", "tracks", "url", "year"]
    ARTISTS_FIELD_NUMBER: _ClassVar[int]
    DURATION_SECONDS_FIELD_NUMBER: _ClassVar[int]
    IMAGE_URL_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    TRACKS_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    YEAR_FIELD_NUMBER: _ClassVar[int]
    artists: _containers.RepeatedCompositeFieldContainer[ArtistLink]
    duration_seconds: int
    image_url: str
    name: str
    tracks: _containers.RepeatedCompositeFieldContainer[Track]
    url: str
    year: int
    def __init__(self, name: _Optional[str] = ..., url: _Optional[str] = ..., image_url: _Optional[str] = ..., year: _Optional[int] = ..., duration_seconds: _Optional[int] = ..., tracks: _Optional[_Iterable[_Union[Track, _Mapping]]] = ..., artists: _Optional[_Iterable[_Union[ArtistLink, _Mapping]]] = ...) -> None: ...

class Albums(_message.Message):
    __slots__ = ["items"]
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[Album]
    def __init__(self, items: _Optional[_Iterable[_Union[Album, _Mapping]]] = ...) -> None: ...

class Artist(_message.Message):
    __slots__ = ["albums", "image_url", "name", "url"]
    ALBUMS_FIELD_NUMBER: _ClassVar[int]
    IMAGE_URL_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    albums: _containers.RepeatedCompositeFieldContainer[Album]
    image_url: str
    name: str
    url: str
    def __init__(self, name: _Optional[str] = ..., url: _Optional[str] = ..., image_url: _Optional[str] = ..., albums: _Optional[_Iterable[_Union[Album, _Mapping]]] = ...) -> None: ...

class ArtistLink(_message.Message):
    __slots__ = ["name", "url"]
    NAME_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    name: str
    url: str
    def __init__(self, name: _Optional[str] = ..., url: _Optional[str] = ...) -> None: ...

class Artists(_message.Message):
    __slots__ = ["items"]
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[Artist]
    def __init__(self, items: _Optional[_Iterable[_Union[Artist, _Mapping]]] = ...) -> None: ...

class Query(_message.Message):
    __slots__ = ["query"]
    QUERY_FIELD_NUMBER: _ClassVar[int]
    query: str
    def __init__(self, query: _Optional[str] = ...) -> None: ...

class Track(_message.Message):
    __slots__ = ["album", "album_url", "artists", "duration_seconds", "image_url", "name", "url"]
    ALBUM_FIELD_NUMBER: _ClassVar[int]
    ALBUM_URL_FIELD_NUMBER: _ClassVar[int]
    ARTISTS_FIELD_NUMBER: _ClassVar[int]
    DURATION_SECONDS_FIELD_NUMBER: _ClassVar[int]
    IMAGE_URL_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    album: str
    album_url: str
    artists: _containers.RepeatedCompositeFieldContainer[ArtistLink]
    duration_seconds: int
    image_url: str
    name: str
    url: str
    def __init__(self, name: _Optional[str] = ..., url: _Optional[str] = ..., image_url: _Optional[str] = ..., duration_seconds: _Optional[int] = ..., album: _Optional[str] = ..., album_url: _Optional[str] = ..., artists: _Optional[_Iterable[_Union[ArtistLink, _Mapping]]] = ...) -> None: ...
