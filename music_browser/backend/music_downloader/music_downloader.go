package music_downloader

import (
	"database/sql"

	"google.golang.org/grpc"
)

type MusicDownloader struct {
	YtDlpConn    *grpc.ClientConn
	PostgresConn *sql.DB
}

func (*MusicDownloader) DownloadTrack(url string) {
	// TODO
}
