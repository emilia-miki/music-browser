import * as grpc from "@grpc/grpc-js";
import { Innertube, YTNodes, YTMusic, Misc } from "youtubei.js";
import { MusicApiService } from "../ts-proto/music_api_grpc_pb";
import {
  Url,
  Query,
  ArtistWithAlbums,
  AlbumWithTracks,
  TrackWithAlbumAndArtist,
  Artist,
  Artists,
  Album,
  Albums,
  Track,
  Tracks,
} from "../ts-proto/music_api_pb";

const YT_MUSIC_BROWSE_BASE_URL = "https://music.youtube.com/browse/";
const YT_MUSIC_VIDEO_BASE_URL = "https://music.youtube.com/watch?v=";

function extractId(url: string): string {
  if (url.startsWith(YT_MUSIC_BROWSE_BASE_URL)) {
    return url.substring(YT_MUSIC_BROWSE_BASE_URL.length);
  } else if (url.startsWith(YT_MUSIC_VIDEO_BASE_URL)) {
    return url.substring(YT_MUSIC_VIDEO_BASE_URL.length);
  } else {
    return "";
  }
}

function getBestThumbnailUrl(arr: Array<Misc.Thumbnail>): string | null {
  if (arr.length === 0) {
    return null;
  }

  let bestUrl = "";
  let bestWidth = 0;
  arr.forEach((obj) => {
    if (obj.width > bestWidth) {
      bestUrl = obj.url;
      bestWidth = obj.width;
    }
  });

  return bestUrl;
}

function parseArtistWithAlbums(
  artistData: YTMusic.Artist,
  artistUrl: string
): ArtistWithAlbums {
  const artistWithAlbums = new ArtistWithAlbums();

  if (!artistData.header) {
    return artistWithAlbums;
  }

  const artist = new Artist();

  artist.setUrl(artistUrl);

  if (!artistData.header.is(YTNodes.MusicHeader)) {
    const thumbnail = artistData.header.thumbnail as YTNodes.MusicThumbnail;
    const imageUrl = getBestThumbnailUrl(thumbnail.contents);
    if (imageUrl) {
      artist.setImageUrl(imageUrl);
    }
  }

  if (artistData.header.title) {
    artist.setName(artistData.header.title.toString());
  }

  artistWithAlbums.setArtist(artist);

  const albums = new Albums();

  let albumsData;
  for (let i = 0; i < artistData.sections.length; i++) {
    const section = artistData.sections[i];
    if (section.type !== YTNodes.MusicCarouselShelf.type) {
      continue;
    }

    const mcs = section as YTNodes.MusicCarouselShelf;
    if (mcs.header?.title.toString() !== "Albums") {
      continue;
    }

    albumsData = mcs.contents.filterType(YTNodes.MusicTwoRowItem);
    break;
  }

  if (albumsData) {
    const albumsList = albumsData.map((data) => {
      const album = new Album();

      if (data.id) {
        album.setUrl(YT_MUSIC_BROWSE_BASE_URL + data.id);
      }

      album.setArtistUrl(artistUrl);

      const imageUrl = getBestThumbnailUrl(data.thumbnail);
      if (imageUrl) {
        album.setImageUrl(imageUrl);
      }

      album.setName(data.title.toString());
      if (data.year) {
        const year = parseInt(data.year);
        album.setYear(year);
      }

      return album;
    });

    albums.setAlbumsList(albumsList);
  } else {
    albums.setAlbumsList([]);
  }

  artistWithAlbums.setAlbums(albums);

  return artistWithAlbums;
}

function parseAlbumWithTracks(
  albumData: YTMusic.Album,
  albumUrl: string
): AlbumWithTracks {
  const albumWithTracks = new AlbumWithTracks();

  const album = new Album();

  if (albumData.url) {
    album.setUrl(albumUrl);
  }

  if (albumData.header) {
    const imageUrl = getBestThumbnailUrl(albumData.header.thumbnails);
    if (imageUrl) {
      album.setImageUrl(imageUrl);
    }

    if (albumData.header.author) {
      const artistId = albumData.header.author.channel_id;
      const artistUrl = YT_MUSIC_BROWSE_BASE_URL + artistId;
      album.setArtistUrl(artistUrl);
    }

    album.setName(albumData.header.title.toString());

    const year = parseInt(albumData.header.year);
    album.setYear(year);
  }

  albumWithTracks.setAlbum(album);

  const tracks = new Tracks();

  const tracksList = albumData.contents.map((trackData) => {
    const track = new Track();

    if (trackData.id) {
      track.setUrl(YT_MUSIC_VIDEO_BASE_URL + trackData.id);
    }

    let imageUrl;
    if (trackData.thumbnail) {
      imageUrl = getBestThumbnailUrl(trackData.thumbnail.contents);
    }
    if (!imageUrl) {
      imageUrl = album.getImageUrl();
    }
    if (imageUrl) {
      track.setImageUrl(imageUrl);
    }

    const artistUrl = album.getArtistUrl();
    if (artistUrl) {
      track.setArtistUrl(artistUrl);
    }

    if (albumUrl) {
      track.setAlbumUrl(albumUrl);
    }

    const name = trackData.title?.toString();
    if (name) {
      track.setName(name);
    }

    if (trackData.duration) {
      track.setDurationSeconds(trackData.duration.seconds);
    }

    return track;
  });

  tracks.setTracksList(tracksList);

  albumWithTracks.setTracks(tracks);

  return albumWithTracks;
}

function parseTrackWithAlbumAndArtist(
  artistData: YTMusic.Artist | null,
  albumData: YTMusic.Album | null,
  trackData: YTMusic.TrackInfo | null,
  albumId: string | null
): TrackWithAlbumAndArtist {
  const trackWithAlbumAndArtist = new TrackWithAlbumAndArtist();

  const artist = new Artist();

  const artistId = trackData?.basic_info.channel_id;
  if (artistData?.header) {
    const header = artistData.header;

    if (artistId) {
      artist.setUrl(YT_MUSIC_BROWSE_BASE_URL + artistId);
    }

    if (
      header.is(YTNodes.MusicImmersiveHeader) ||
      header.is(YTNodes.MusicVisualHeader)
    ) {
      let imageUrl: string | null = null;
      if (header.thumbnail) {
        if (header.thumbnail instanceof YTNodes.MusicThumbnail) {
          imageUrl = getBestThumbnailUrl(header.thumbnail.contents);
        } else {
          imageUrl = getBestThumbnailUrl(header.thumbnail);
        }
      }
      if (imageUrl) {
        artist.setImageUrl(imageUrl);
      }
    }

    if (header.title) {
      artist.setName(header.title.toString());
    }
  }

  trackWithAlbumAndArtist.setArtist(artist);

  const album = new Album();

  if (albumId) {
    album.setUrl(YT_MUSIC_BROWSE_BASE_URL + albumId);
  }

  if (albumId && artistId) {
    album.setArtistUrl(YT_MUSIC_BROWSE_BASE_URL + artistId);
  }

  if (albumData?.header) {
    const header = albumData.header;

    const imageUrl = getBestThumbnailUrl(header.thumbnails);
    if (imageUrl) {
      album.setImageUrl(imageUrl);
    }

    album.setName(header.title.toString());

    const year = parseInt(header.year);
    album.setYear(year);
  }

  trackWithAlbumAndArtist.setAlbum(album);

  const track = new Track();

  const trackId = trackData?.basic_info.id;
  if (trackId) {
    track.setUrl(YT_MUSIC_VIDEO_BASE_URL + trackId);
  }

  if (trackData) {
    const imageUrl = album.getImageUrl();
    if (imageUrl) {
      track.setImageUrl(imageUrl);
    }

    const artistUrl = artist.getUrl();
    if (artistUrl) {
      track.setArtistUrl(artistUrl);
    }

    const albumUrl = album.getUrl();
    if (albumUrl) {
      track.setAlbumUrl(albumUrl);
    }

    const name = trackData.basic_info.title?.toString();
    if (name) {
      track.setName(name);
    }

    if (trackData.basic_info.duration) {
      track.setDurationSeconds(trackData.basic_info.duration);
    }
  }

  trackWithAlbumAndArtist.setTrack(track);

  return trackWithAlbumAndArtist;
}

async function createYtMusicApiServer() {
  const innerTube = await Innertube.create();
  return {
    getArtist(
      call: grpc.ServerUnaryCall<Url, ArtistWithAlbums>,
      callback: grpc.sendUnaryData<ArtistWithAlbums>
    ): void {
      const artistUrl = call.request.getUrl();
      const artistId = extractId(artistUrl);

      innerTube.music.getArtist(artistId).then((data) => {
        const artistWithAlbums = parseArtistWithAlbums(data, artistUrl);
        callback(null, artistWithAlbums);
      });
    },

    getAlbum(
      call: grpc.ServerUnaryCall<Url, AlbumWithTracks>,
      callback: grpc.sendUnaryData<AlbumWithTracks>
    ): void {
      const albumUrl = call.request.getUrl();
      const albumId = extractId(albumUrl);

      innerTube.music.getAlbum(albumId).then((data) => {
        const albumWithTracks = parseAlbumWithTracks(data, albumUrl);
        callback(null, albumWithTracks);
      });
    },

    getTrack(
      call: grpc.ServerUnaryCall<Url, TrackWithAlbumAndArtist>,
      callback: grpc.sendUnaryData<TrackWithAlbumAndArtist>
    ): void {
      const trackUrl = call.request.getUrl();
      const trackId = extractId(trackUrl);

      const getPromise = async function () {
        let artistData: YTMusic.Artist | null = null;
        let albumData: YTMusic.Album | null = null;
        let albumId: string | null = null;
        const trackData = await innerTube.music.getInfo(trackId);

        const artistId = trackData.basic_info.channel_id ?? null;
        if (artistId) {
          artistData = await innerTube.music.getArtist(artistId);
        }

        if (artistData) {
          const songsData = await artistData.getAllSongs();
          const songsList = songsData?.contents;
          if (songsList) {
            for (let i = 0; i < songsList.length; i++) {
              if (songsList[i].id !== trackId) {
                continue;
              }

              albumId = songsList[i].album?.id ?? null;
              if (!albumId) {
                break;
              }

              albumData = await innerTube.music.getAlbum(albumId);
              break;
            }
          }
        }

        const trackWithAlbumAndArtist = parseTrackWithAlbumAndArtist(
          artistData,
          albumData,
          trackData,
          albumId
        );

        return trackWithAlbumAndArtist;
      };

      getPromise().then((track) => callback(null, track));
    },

    searchArtists(
      call: grpc.ServerUnaryCall<Query, Artists>,
      callback: grpc.sendUnaryData<Artists>
    ): void {
      const query = call.request.getQuery();

      innerTube.music.search(query, { type: "artist" }).then((search) => {
        const artists = new Artists();
        const artistsList: Artist[] = [];
        if (search.artists) {
          search.artists.contents.forEach((item) => {
            if (item.id && item.name) {
              const artist = new Artist();

              artist.setUrl(YT_MUSIC_BROWSE_BASE_URL + item.id);

              if (item.thumbnail) {
                const imageUrl = getBestThumbnailUrl(item.thumbnail.contents);
                if (imageUrl) {
                  artist.setImageUrl(imageUrl);
                }
              }

              artist.setName(item.name);

              artistsList.push(artist);
            } else {
              return null;
            }
          });
        }

        artists.setArtistsList(artistsList);

        callback(null, artists);
      });
    },

    searchAlbums(
      call: grpc.ServerUnaryCall<Query, Albums>,
      callback: grpc.sendUnaryData<Albums>
    ): void {
      const query = call.request.getQuery();

      innerTube.music.search(query, { type: "album" }).then((search) => {
        const albums = new Albums();
        const albumsList: Album[] = [];
        if (search.contents && search.contents[0].contents) {
          search.contents[0].contents
            .filterType(YTNodes.MusicResponsiveListItem)
            .forEach((item) => {
              if (item.id && item.title) {
                const album = new Album();

                album.setUrl(YT_MUSIC_BROWSE_BASE_URL + item.id);

                if (item.thumbnail) {
                  const imageUrl = getBestThumbnailUrl(item.thumbnail.contents);
                  if (imageUrl) {
                    album.setImageUrl(imageUrl);
                  }
                }

                if (item.author?.channel_id) {
                  album.setArtistUrl(
                    YT_MUSIC_BROWSE_BASE_URL + item.author?.channel_id
                  );
                }

                album.setName(item.title);

                if (item.year) {
                  const year = parseInt(item.year);
                  album.setYear(year);
                }

                albumsList.push(album);
              }
            });
        }

        albums.setAlbumsList(albumsList);

        callback(null, albums);
      });
    },
  };
}

if (require.main === module) {
  const port = process.env.YT_MUSIC_API_PORT;
  if (!port) {
    console.error("environment variable YT_MUSIC_API_PORT not set");
    process.exit(1);
  }
  const uri = `0.0.0.0:${port}`;

  createYtMusicApiServer().then((YtMusicApiServer) => {
    const server = new grpc.Server();
    server.addService(MusicApiService, YtMusicApiServer as any);
    server.bindAsync(uri, grpc.ServerCredentials.createInsecure(), () => {
      server.start();
    });
  });
}

module.exports = {
  create: createYtMusicApiServer,
};
