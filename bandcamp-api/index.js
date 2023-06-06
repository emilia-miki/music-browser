// @ts-nocheck
const bandcamp = require('bandcamp-scraper');
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');

/**
 * @typedef {Object} ArtistLink
 * @property {string} name
 * @property {string} url
 */

/**
 * @typedef {Object} Track
 * @property {string} name
 * @property {string} url
 * @property {string} image_url
 * @property {number} duration_seconds
 * @property {string} album
 * @property {string} album_url
 * @property {Array<ArtistLink>} artists
 */


/**
 * @typedef {Object} Album
 * @property {string} name
 * @property {string} url
 * @property {string} image_url
 * @property {number} year
 * @property {number} duration_seconds
 * @property {Array<Track>} tracks
 * @property {Array<ArtistLink>} artists
 */

/** 
 * @param {string} duration 
 * @returns {number}
 */
function durationToSeconds(duration) {
  const splits = duration.split(':');
  const minutes = parseInt(splits[0]);
  const seconds = parseInt(splits[1]);
  return minutes * 60 + seconds;
}

/**
 * @param {string} url
 * @returns {Promise<Album>}
*/
function getAlbum(url) {
  return new Promise(resolve => {
    bandcamp.getAlbumInfo(url, (_, data) => {
      const releaseDateString = data.raw.album_release_date;
      const year = parseInt(releaseDateString.split(' ')[2]);
      const artists = [{
        name: data.artist,
        url: url.split('/').slice(0, -2).join('/'),
      }];

      /** @type {Array<Track>} */
      const tracks = data.tracks.map(track => ({
        name: track.name,
        url: track.url,
        image_url: data.imageUrl,
        duration_seconds: durationToSeconds(track.duration),
        album: data.title,
        album_url: url,
        artists,
      }));

      const duration_seconds = tracks.reduce(
        (prev, curr) => prev + curr.duration_seconds, 0);

      resolve({
        name: data.title,
        url: url,
        image_url: data.imageUrl,
        year,
        duration_seconds,
        tracks,
        artists,
      });
    });
  }, () => {});
}

/**
 * @typedef {Object} Artist
 * @property {string} name
 * @property {string} url
 * @property {string} image_url
 * @property {Array<Album>} albums
 */

/**
 * @param {Object} item
 * @param {string} url
 * @returns {Promise<Artist>}
 */
function getArtist(url) {
  return new Promise(resolve => {
    bandcamp.getArtistInfo(url, (_, data) => {
      const albumPromises = data.albums.map(album => getAlbum(album.url));
      Promise.all(albumPromises).then(albums => resolve(
        {
          name: data.name,
          url: url,
          image_url: data.coverImage,
          albums,
        }
      ));
    });
  }, () => {});
}

/**
 * @template {Artist | Album} T
 * @param {string} query
 * @param {string} type
 * @returns {Promise<Array<T>>}
 */
function search(query, type) {
  const params = {
    query: query,
    age: 1,
  };

  return new Promise(resolve => {
    bandcamp.search(params, (_, results) => {
      if (type === "artist") {
        mapFunction = getArtist;
      } else if (type === "album") {
        mapFunction = getAlbum;
      }

      const items = results
        .filter(item => item.type === type)
        .map(item => item.url)
        .map(item => mapFunction(item));

     Promise.all(items).then(items => resolve(items));
    });
  }, () => {});
}

if (!require.main) {
  module.exports = search;
  return;
}

port = process.env.BANDCAMP_API_PORT
if (!port) {
  console.error("environment variable BANDCAMP_API_PORT not set")
  process.exit(1)
}

const PROTO_PATH = __dirname + "/../music_api.proto";
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
});
const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
const musicApi = protoDescriptor.music_api;

var server = new grpc.Server();
server.addService(musicApi.MusicApi.service, {
  searchArtists: (call, callback) => {
    search(call.request.query, 'artist')
    .then(items => {
      callback(null, {items});
    });
  },
  searchAlbums: (call, callback) =>
    search(call.request.query, 'album')
    .then(items => callback(null, {items})),
});

server.bindAsync('0.0.0.0:' + port, grpc.ServerCredentials.createInsecure(),
  () => server.start()
);
