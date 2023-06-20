import bcfetch, {
  Track as TrackInfo,
  Album as AlbumInfo,
  Artist as ArtistInfo,
  Label,
} from 'bandcamp-fetch';
import {
  Server,
  ServerCredentials,
  ServerUnaryCall,
  sendUnaryData
} from '@grpc/grpc-js';

import { MusicApiService } from '../ts-proto/music_api_grpc_pb';
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
} from '../ts-proto/music_api_pb';

function parseArtist(
  artistInfo: Omit<ArtistInfo, "type"> | Omit<Label, "type">,
): Artist {
  const artist = new Artist();

  if (artistInfo.url) {
    artist.setUrl(artistInfo.url);
  }

  if (artistInfo.imageUrl) {
    artist.setImageUrl(artistInfo.imageUrl);
  }

  artist.setName(artistInfo.name);

  return artist;
}

function parseAlbum(albumInfo: Omit<AlbumInfo, "type">): Album {
  const album = new Album();

  if (albumInfo.url) {
    album.setUrl(albumInfo.url)
  }

  if (albumInfo.imageUrl) {
    album.setImageUrl(albumInfo.imageUrl);
  }

  if (albumInfo.artist?.url) {
    album.setArtistUrl(albumInfo.artist.url);
  }

  album.setName(albumInfo.name);

  if (albumInfo.releaseDate) {
    const year = parseInt(albumInfo.releaseDate.split(" ")[2]);
    album.setYear(year);
  }

  return album;
}

function parseTrack(trackInfo: Omit<TrackInfo, "type">): Track {
  const track = new Track();

  if (trackInfo.url) {
    track.setUrl(trackInfo.url);
  }

  if (trackInfo.imageUrl) {
    track.setImageUrl(trackInfo.imageUrl);
  }

  if (trackInfo.artist?.url) {
    track.setArtistUrl(trackInfo.artist.url);
  }

  if (trackInfo.album?.url) {
    track.setAlbumUrl(trackInfo.album.url);
  }

  track.setName(trackInfo.name);

  if (trackInfo.duration) {
    track.setDurationSeconds(Math.round(trackInfo.duration));
  }

  return track;
}

async function getArtist(url: string): Promise<ArtistWithAlbums> {
  const artistInfo = await bcfetch.band.getInfo({ bandUrl: url });

  const artistWithAlbums = new ArtistWithAlbums();

  const artist = parseArtist(artistInfo);
  artistWithAlbums.setArtist(artist);

  const albumsInfo = await bcfetch.band.getDiscography({ bandUrl: url });
  const albums = new Albums();
  const albumsList = await Promise.all(albumsInfo
    .filter((info): info is AlbumInfo => !!info)
    .map(async info => {
      if (info.url) {
        const albumInfo = await bcfetch.album.getInfo({ albumUrl: info.url });
        return parseAlbum(albumInfo);
      } else {
        return parseAlbum(info);
      }
    }));
  albums.setAlbumsList(albumsList);
  artistWithAlbums.setAlbums(albums);

  return artistWithAlbums;
}

async function getAlbum(url: string): Promise<AlbumWithTracks> {
  const albumInfo = await bcfetch.album.getInfo({ albumUrl: url });

  const albumWithTracks = new AlbumWithTracks();

  const album = parseAlbum(albumInfo);
  albumWithTracks.setAlbum(album);

  const tracks = new Tracks();
  const tracksList = albumInfo.tracks
    ? albumInfo.tracks.map(trackInfo => {
      trackInfo.imageUrl = albumInfo.imageUrl;
      trackInfo.artist = { url: albumInfo.artist?.url, name: "" };
      trackInfo.album = { url, name: "" };
      return parseTrack(trackInfo);
    })
    : new Array<Track>();
  tracks.setTracksList(tracksList);
  albumWithTracks.setTracks(tracks);

  return albumWithTracks;
}

async function getTrack(url: string): Promise<TrackWithAlbumAndArtist> {
  const trackInfo = await bcfetch.track.getInfo({ trackUrl: url });

  const trackWithAlbumAndArtist = new TrackWithAlbumAndArtist();

  const track = parseTrack(trackInfo);
  trackWithAlbumAndArtist.setTrack(track);

  let album = new Album();
  const albumUrl = track.getAlbumUrl();
  if (albumUrl) {
    const albumInfo = await bcfetch.album.getInfo({ albumUrl });
    album = parseAlbum(albumInfo);

    if (albumInfo.tracks) {
      for (let i = 0; i < albumInfo.tracks.length; i++) {
        const info = albumInfo.tracks[i];
        if (trackInfo.url && info.url === trackInfo.url) {
          if (info.duration) {
            track.setDurationSeconds(Math.round(info.duration));
          }
        }
      }
    }
  }
  trackWithAlbumAndArtist.setAlbum(album);

  let artist = new Artist();
  const artistUrl = track.getArtistUrl();
  if (artistUrl) {
    const artistInfo = await bcfetch.band.getInfo({ bandUrl: artistUrl });
    artist = parseArtist(artistInfo);
  }
  trackWithAlbumAndArtist.setArtist(artist);

  return trackWithAlbumAndArtist;
}

async function searchArtists(query: string): Promise<Artists> {
  const results = await bcfetch.search.artistsAndLabels({ query });

  const artists = new Artists();

  const artistsList = results.items.map(result => parseArtist(result));
  artists.setArtistsList(artistsList);

  return artists;
}

async function searchAlbums(query: string): Promise<Albums> {
  const results = await bcfetch.search.albums({ query });

  const albums = new Albums();

  const albumsList = await Promise.all(results.items.map(async result => {
    const albumInfo = await bcfetch.album.getInfo({ albumUrl: result.url });
    return parseAlbum(albumInfo);
  }));
  albums.setAlbumsList(albumsList);

  return albums;
}

const BandcampApiServer = {
  getArtist(
    call: ServerUnaryCall<Url, ArtistWithAlbums>,
    callback: sendUnaryData<ArtistWithAlbums>
  ) {
    getArtist(call.request.getUrl())
      .then(artist => callback(null, artist));
  },

  getAlbum(
    call: ServerUnaryCall<Url, AlbumWithTracks>,
    callback: sendUnaryData<AlbumWithTracks>
  ) {
    getAlbum(call.request.getUrl())
      .then(album => callback(null, album));
  },

  getTrack(
    call: ServerUnaryCall<Url, TrackWithAlbumAndArtist>,
    callback: sendUnaryData<TrackWithAlbumAndArtist>
  ) {
    getTrack(call.request.getUrl())
      .then(track => callback(null, track));
  },

  searchArtists(
    call: ServerUnaryCall<Query, Artists>,
    callback: sendUnaryData<Artists>
  ) {
    searchArtists(call.request.getQuery())
      .then(artists => callback(null, artists));
  },

  searchAlbums(
    call: ServerUnaryCall<Query, Albums>,
    callback: sendUnaryData<Albums>
  ) {
    searchAlbums(call.request.getQuery())
      .then(albums => callback(null, albums));
  },
}

if (require.main == module) {
  const port = process.env.BANDCAMP_API_PORT;
  if (!port) {
    console.error("environment variable BANDCAMP_API_PORT not set");
    process.exit(1);
  }
  const uri = `0.0.0.0:${port}`;

  var server = new Server();
  server.addService(MusicApiService, BandcampApiServer);
  server.bindAsync(
    uri,
    ServerCredentials.createInsecure(),
    () => server.start()
  );
}

module.exports = BandcampApiServer;
