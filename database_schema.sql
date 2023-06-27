CREATE TABLE IF NOT EXISTS link_map (
    sp_url TEXT PRIMARY KEY,
    yt_url TEXT
);

CREATE TABLE IF NOT EXISTS image (
    url TEXT PRIMARY KEY,
    path TEXT
);

CREATE TABLE IF NOT EXISTS artist (
    url TEXT PRIMARY KEY,
    image_url TEXT REFERENCES image (url),
    name TEXT
);

CREATE TABLE IF NOT EXISTS album (
    url TEXT PRIMARY KEY,
    image_url TEXT REFERENCES image (url),

    name TEXT,
    year INTEGER
);

CREATE TABLE IF NOT EXISTS track (
    url TEXT PRIMARY KEY,
    image_url TEXT REFERENCES image (url),
    album_url TEXT REFERENCES album (url),

    path TEXT,
    name TEXT,
    duration_seconds INTEGER
);

CREATE TABLE IF NOT EXISTS artist_album (
    artist_url TEXT REFERENCES artist (url),
    album_url TEXT REFERENCES album (url),

    PRIMARY KEY (artist_url, album_url)
);

CREATE TABLE IF NOT EXISTS artist_track (
    artist_url TEXT REFERENCES artist (url),
    track_url TEXT REFERENCES track (url),

    PRIMARY KEY (artist_url, track_url)
);

CREATE TABLE IF NOT EXISTS album_track (
    album_url TEXT REFERENCES album (url),
    track_url TEXT REFERENCES track (url),

    PRIMARY KEY (album_url, track_url)
);
