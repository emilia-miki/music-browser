import os
import sys
from concurrent import futures

import grpc
from ytmusicapi import YTMusic

import music_api_pb2_grpc
from music_api_pb2 import Query, Artists, Albums, Artist, Album, Track, ArtistLink

YT_MUSIC_BROWSE_BASE_URL = "https://music.youtube.com/browse/"
YT_MUSIC_VIDEO_BASE_URL = "https://music.youtube.com/watch?v="

def get_best_thumbnail_url(item: dict) -> str:
    thumbnails = item["thumbnails"]

    url = ""
    max_height = 0
    for thumbnail in thumbnails:
        if thumbnail["height"] > max_height:
            max_height = thumbnail["height"]
            url = thumbnail["url"]

    return url

def parse_artist_link(item: dict) -> ArtistLink:
    name = item["name"]
    url = None
    if "id" in item and item["id"] is not None:
        url = YT_MUSIC_BROWSE_BASE_URL + item["id"]

    return ArtistLink(
        name=name,
        url=url,
    )

def parse_track(item: dict, album_url: str, album_image_url: str) -> Track:
    name = item["title"]
    url = YT_MUSIC_VIDEO_BASE_URL + item["videoId"]
    image_url = album_image_url
    duration_seconds = item["duration_seconds"]
    album = item["album"]
    artists = [parse_artist_link(artist) for artist in item["artists"]]

    return Track(
        name=name,
        url=url,
        image_url=image_url,
        duration_seconds=duration_seconds,
        album=album,
        album_url=album_url,
        artists=artists,
    )

def get_album(yt_music: YTMusic, id: str,
    artist_link: ArtistLink | None = None) -> Album:
    album = yt_music.get_album(id)

    name = album["title"]
    url = YT_MUSIC_BROWSE_BASE_URL + id
    image_url = get_best_thumbnail_url(album)
    year = int(album["year"])
    tracks = [parse_track(track, url, image_url) 
        for track in album["tracks"] if "videoId" in track]
    duration_seconds = sum([track.duration_seconds for track in tracks])
    artists = [parse_artist_link(artist) for artist in album["artists"]]
    if artist_link is not None:
        if (len(artists) == 0 or (len(artists) == 1 
            and artists[0].url is None)):
            artists = [artist_link]

    return Album(
        name=name,
        url=url,
        image_url=image_url,
        year=year,
        duration_seconds=duration_seconds,
        tracks=tracks,
        artists=artists,
    )

def get_artist(yt_music: YTMusic, id: str) -> Artist:
    artist = yt_music.get_artist(id)

    name = artist["name"]
    url = YT_MUSIC_BROWSE_BASE_URL + artist["channelId"]
    image_url = get_best_thumbnail_url(artist)

    album_ids = set()
    if "albums" in artist and "results" in artist["albums"]:
        for album in artist["albums"]["results"]:
            album_ids.add(album["browseId"])
    elif "songs" in artist and "results" in artist["songs"]:
        for song in artist["songs"]["results"]:
            album_ids.add(song["album"]["id"])

    albums = [get_album(yt_music, id, ArtistLink(name=name, url=url)) 
        for id in album_ids]

    return Artist(name=name, url=url, image_url=image_url, albums=albums)

class MusicApiServicer(music_api_pb2_grpc.MusicApiServicer):
    def __init__(self):
        self.yt_music = YTMusic()

    def SearchArtists(self, request: Query, context) -> Artists:
        items = self.yt_music.search(request.query, "artists")
        artists = [get_artist(self.yt_music, item["browseId"]) 
            for item in items]
        return Artists(items=artists)

    def SearchAlbums(self, request: Query, context) -> Albums:
        items = self.yt_music.search(request.query, "albums")
        albums = [get_album(self.yt_music, item["browseId"]) for item in items]
        return Albums(items=albums)

if __name__ == "__main__":
    port = os.getenv("YT_MUSIC_API_PORT")
    if port is None:
        sys.stderr.write("YT_MUSIC_API_PORT environment variable not set\n")
        sys.exit(1)

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    music_api_pb2_grpc.add_MusicApiServicer_to_server(
        MusicApiServicer(), server)
    server.add_insecure_port("[::]:" + port)
    server.start()
    server.wait_for_termination()
