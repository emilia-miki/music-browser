"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const bandcamp_scraper_1 = __importDefault(require("bandcamp-scraper"));
const grpc_js_1 = require("@grpc/grpc-js");
const music_api_grpc_pb_1 = require("../ts-proto/music_api_grpc_pb");
const music_api_pb_1 = require("../ts-proto/music_api_pb");
function getArtist(url) {
    return new Promise(resolve => {
        bandcamp_scraper_1.default.getArtistInfo(url, (_, data) => {
            const imageUrl = data.coverImage;
            const name = data.name;
            const artist = new music_api_pb_1.Artist();
            artist.setUrl(url);
            artist.setImageUrl(imageUrl);
            artist.setName(name);
            const albums = new music_api_pb_1.Albums();
            Promise.all(data.albums.map((album) => getAlbumWithoutTracks(album.url))).then(albumsList => {
                albums.setAlbumsList(albumsList);
                const result = new music_api_pb_1.ArtistWithAlbums();
                result.setArtist(artist);
                result.setAlbums(albums);
                resolve(result);
            });
        });
    });
}
function getAlbum(url) {
    return new Promise((resolve) => {
        bandcamp_scraper_1.default.getAlbumInfo(url, (_, data) => {
            const imageUrl = data.imageUrl;
            const artistUrl = extractArtistUrl(url);
            const name = data.title;
            const dateString = data.raw.album_release_date;
            const year = dateToYear(dateString);
            const album = new music_api_pb_1.Album();
            album.setUrl(url);
            album.setImageUrl(imageUrl);
            album.setName(name);
            album.setYear(year);
            album.setArtistUrl(artistUrl);
            const albumUrl = url;
            const tracks = new music_api_pb_1.Tracks();
            tracks.setTracksList(data.tracks.map((data) => {
                const url = data.url;
                const name = data.name;
                const duration = durationToSeconds(data.duration);
                const track = new music_api_pb_1.Track();
                track.setUrl(url);
                track.setImageUrl(imageUrl);
                track.setArtistUrl(artistUrl);
                track.setAlbumUrl(albumUrl);
                track.setName(name);
                track.setDurationSeconds(duration);
                return track;
            }));
            const result = new music_api_pb_1.AlbumWithTracks();
            result.setAlbum(album);
            result.setTracks(tracks);
            resolve(result);
        });
    });
}
function getTrack(url) {
    return new Promise(resolve => {
        bandcamp_scraper_1.default.getTrackInfo(url, (_, data) => {
            const artistUrl = extractArtistUrl(url);
            const albumUrl = data.raw.album_url;
            const name = data.title;
            const durationFloat = data.raw.trackinfo[0].duration;
            const duration = Math.round(durationFloat);
            const promise = function () {
                return __awaiter(this, void 0, void 0, function* () {
                    const artist = yield getArtistWithoutAlbums(artistUrl);
                    const album = yield getAlbumWithoutTracks(albumUrl);
                    const track = new music_api_pb_1.Track();
                    track.setUrl(url);
                    const imageUrl = album.getImageUrl();
                    if (imageUrl) {
                        track.setImageUrl(imageUrl);
                    }
                    track.setArtistUrl(artistUrl);
                    track.setAlbumUrl(albumUrl);
                    track.setName(name);
                    track.setDurationSeconds(duration);
                    const result = new music_api_pb_1.TrackWithAlbumAndArtist();
                    result.setArtist(artist);
                    result.setAlbum(album);
                    result.setTrack(track);
                    return result;
                });
            };
            promise().then(result => resolve(result));
        });
    });
}
function searchArtists(query) {
    const params = {
        query: query,
        age: 1,
    };
    return new Promise(resolve => {
        bandcamp_scraper_1.default.search(params, (_, data) => {
            const items = data
                .filter((item) => item.type === "artist")
                .map((data) => getArtist(data.url));
            const result = new music_api_pb_1.Artists();
            Promise.all(items).then(artistsList => {
                result.setArtistsList(artistsList);
                resolve(result);
            });
        });
    });
}
function searchAlbums(query) {
    const params = {
        query: query,
        age: 1,
    };
    return new Promise(resolve => {
        bandcamp_scraper_1.default.search(params, (_, data) => {
            const items = data
                .filter((item) => item.type === "album")
                .map((data) => getAlbum(data.url));
            const result = new music_api_pb_1.Albums();
            Promise.all(items).then(albumsList => {
                result.setAlbumsList(albumsList);
                resolve(result);
            });
        });
    });
}
function getArtistWithoutAlbums(url) {
    return new Promise(resolve => {
        bandcamp_scraper_1.default.getArtistInfo(url, (_, data) => {
            const imageUrl = data.coverImage;
            const name = data.name;
            const result = new music_api_pb_1.Artist();
            result.setUrl(url);
            result.setImageUrl(imageUrl);
            result.setName(name);
            resolve(result);
        });
    });
}
function getAlbumWithoutTracks(url) {
    return new Promise(resolve => {
        bandcamp_scraper_1.default.getAlbumInfo(url, (_, data) => {
            const imageUrl = data.imageUrl;
            const artistUrl = extractArtistUrl(url);
            const name = data.title;
            const dateString = data.raw.album_release_date;
            const year = dateToYear(dateString);
            const result = new music_api_pb_1.Album();
            result.setUrl(url);
            result.setImageUrl(imageUrl);
            result.setArtistUrl(artistUrl);
            result.setName(name);
            result.setYear(year);
            resolve(result);
        });
    });
}
function durationToSeconds(duration) {
    const splits = duration.split(':');
    const minutes = parseInt(splits[0]);
    const seconds = parseInt(splits[1]);
    return minutes * 60 + seconds;
}
function dateToYear(date) {
    return parseInt(date.split(" ")[2]);
}
function extractArtistUrl(url) {
    return url.split('/').slice(0, -2).join('/');
}
const BandcampApiServer = {
    getArtist(call, callback) {
        getArtist(call.request.getUrl())
            .then(artist => callback(null, artist));
    },
    getAlbum(call, callback) {
        getAlbum(call.request.getUrl())
            .then(album => callback(null, album));
    },
    getTrack(call, callback) {
        getTrack(call.request.getUrl())
            .then(track => callback(null, track));
    },
    searchArtists(call, callback) {
        searchArtists(call.request.getQuery())
            .then(artists => callback(null, artists));
    },
    searchAlbums(call, callback) {
        searchAlbums(call.request.getQuery())
            .then(albums => callback(null, albums));
    },
};
if (require.main == module) {
    const port = process.env.BANDCAMP_API_PORT;
    if (!port) {
        console.error("environment variable BANDCAMP_API_PORT not set");
        process.exit(1);
    }
    const uri = `0.0.0.0:${port}`;
    var server = new grpc_js_1.Server();
    server.addService(music_api_grpc_pb_1.MusicApiService, BandcampApiServer);
    server.bindAsync(uri, grpc_js_1.ServerCredentials.createInsecure(), () => server.start());
}
module.exports = BandcampApiServer;
