import bcfetch from 'bandcamp-fetch';
import { Server, ServerCredentials, ServerUnaryCall, sendUnaryData } from '@grpc/grpc-js';
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

function getArtist(url: string): Promise<ArtistWithAlbums> {
  return new Promise(resolve => {
    bandcamp.getArtistInfo(url, (_: any, data: any) => {
      const imageUrl = data.coverImage;
      const name = data.name;

      const artist = new Artist();
      artist.setUrl(url);
      artist.setImageUrl(imageUrl);
      artist.setName(name);

      const albums = new Albums();
      Promise.all(
        data.albums.map((album: any) => getAlbumWithoutTracks(album.url))
      ).then(albumsList => {
        albums.setAlbumsList(albumsList);

        const result = new ArtistWithAlbums();
        result.setArtist(artist);
        result.setAlbums(albums);

        resolve(result);
      });
    });
  });
}

function getAlbum(url: string): Promise<AlbumWithTracks> {
  return new Promise((resolve) => {
    bandcamp.getAlbumInfo(url, (_: any, data: any) => {
      const imageUrl = data.imageUrl;
      const artistUrl = extractArtistUrl(url);
      const name = data.title;
      const dateString = data.raw.album_release_date;
      const year = dateToYear(dateString);

      const album = new Album();
      album.setUrl(url);
      album.setImageUrl(imageUrl);
      album.setName(name);
      album.setYear(year);
      album.setArtistUrl(artistUrl);

      const albumUrl = url;
      const tracks = new Tracks();
      tracks.setTracksList(data.tracks.map((data: any) => {
        const url = data.url;
        const name = data.name;
        const duration = durationToSeconds(data.duration);

        const track = new Track();
        track.setUrl(url);
        track.setImageUrl(imageUrl);
        track.setArtistUrl(artistUrl);
        track.setAlbumUrl(albumUrl);
        track.setName(name);
        track.setDurationSeconds(duration);
        return track;
      }));

      const result = new AlbumWithTracks();
      result.setAlbum(album);
      result.setTracks(tracks);

      resolve(result);
    });
  });
}

function getTrack(url: string): Promise<TrackWithAlbumAndArtist> {
  return new Promise(resolve => {
    bandcamp.getTrackInfo(url, (_: any, data: any) => {
      const artistUrl = extractArtistUrl(url);
      const albumUrl = data.raw.album_url;
      const name = data.title;
      const durationFloat = data.raw.trackinfo[0].duration;
      const duration = Math.round(durationFloat);

      const promise = async function() {
        const artist = await getArtistWithoutAlbums(artistUrl);

        const album = await getAlbumWithoutTracks(albumUrl);

        const track = new Track();
        track.setUrl(url);

        const imageUrl = album.getImageUrl();
        if (imageUrl) {
          track.setImageUrl(imageUrl);
        }

        track.setArtistUrl(artistUrl);
        track.setAlbumUrl(albumUrl);
        track.setName(name);
        track.setDurationSeconds(duration);


        const result = new TrackWithAlbumAndArtist();
        result.setArtist(artist);
        result.setAlbum(album);
        result.setTrack(track);
        
        return result;
      };

      promise().then(result => resolve(result));
    });
  });
}

function searchArtists(query: string): Promise<Artists> {
  const params = {
    query: query,
    age: 1,
  };

  return new Promise(resolve => {
    bandcamp.search(params, (_: any, data: any) => {
      const items = data
        .filter((item: any) => item.type === "artist")
        .map((data: any) => getArtist(data.url));

      const result = new Artists();
      Promise.all(items).then(artistsList => {
        result.setArtistsList(artistsList);

        resolve(result);
      })
    });
  });
}

function searchAlbums(query: string): Promise<Albums> {
  const params = {
    query: query,
    age: 1,
  };

  return new Promise(resolve => {
    bandcamp.search(params, (_: any, data: any) => {
      const items = data
        .filter((item: any) => item.type === "album")
        .map((data: any) => getAlbum(data.url));

      const result = new Albums();
      Promise.all(items).then(albumsList => {
        result.setAlbumsList(albumsList);

        resolve(result);
      });
    });
  });
}

function getArtistWithoutAlbums(url: string): Promise<Artist> {
  return new Promise(resolve => {
    bandcamp.getArtistInfo(url, (_: any, data: any) => {
      const imageUrl = data.coverImage;
      const name = data.name;

      const result = new Artist();
      result.setUrl(url);
      result.setImageUrl(imageUrl);
      result.setName(name);

      resolve(result);
    });
  });
}

function getAlbumWithoutTracks(url: string): Promise<Album> {
  return new Promise(resolve => {
    bandcamp.getAlbumInfo(url, (_: any, data: any) => {
      const imageUrl = data.imageUrl;
      const artistUrl = extractArtistUrl(url);
      const name = data.title;
      const dateString = data.raw.album_release_date;
      const year = dateToYear(dateString);

      const result = new Album();
      result.setUrl(url);
      result.setImageUrl(imageUrl);
      result.setArtistUrl(artistUrl);
      result.setName(name);
      result.setYear(year);

      resolve(result);
    });
  });
}

function durationToSeconds(duration: string): number {
  const splits = duration.split(':');
  const minutes = parseInt(splits[0]);
  const seconds = parseInt(splits[1]);
  return minutes * 60 + seconds;
}

function dateToYear(date: string): number {
  return parseInt(date.split(" ")[2]);
}

function extractArtistUrl(url: string): string {
  return url.split('/').slice(0, -2).join('/');
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
