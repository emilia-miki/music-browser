CREATE TABLE IF NOT EXISTS image (
    url TEXT PRIMARY KEY,

    path TEXT
);

CREATE TABLE IF NOT EXISTS artist (
    url TEXT PRIMARY KEY,
    image_url TEXT,

    name TEXT,

    FOREIGN KEY (image_url) REFERENCES image (url)
);

CREATE TABLE IF NOT EXISTS album (
    url TEXT PRIMARY KEY,
    image_url TEXT,

    name TEXT,
    year INTEGER,

    FOREIGN KEY (image_url) REFERENCES image (url)
);

CREATE TABLE IF NOT EXISTS track (
    url TEXT PRIMARY KEY,
    image_url TEXT,
    album_url TEXT,

    path TEXT,
    name TEXT,
    duration_seconds INTEGER,

    FOREIGN KEY (image_url) REFERENCES image (url),
    FOREIGN KEY (album_url) REFERENCES album (url)
);

CREATE TABLE IF NOT EXISTS artist_album (
    artist_url TEXT,
    album_url TEXT,

    FOREIGN KEY (artist_url) REFERENCES artist (url),
    FOREIGN KEY (album_url) REFERENCES album (url),
    PRIMARY KEY (artist_url, album_url)
);

CREATE TABLE IF NOT EXISTS artist_track (
    artist_url TEXT,
    track_url TEXT,

    FOREIGN KEY (artist_url) REFERENCES artist (url),
    FOREIGN KEY (track_url) REFERENCES track (url),
    PRIMARY KEY (artist_url, track_url)
);

CREATE TABLE IF NOT EXISTS album_track (
    album_url TEXT,
    track_url TEXT,

    FOREIGN KEY (album_url) REFERENCES album (url),
    FOREIGN KEY (track_url) REFERENCES track (url),
    PRIMARY KEY (album_url, track_url)
);
