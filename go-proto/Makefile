.DEFAULT_GOAL := proto-gen

proto-gen:
	PATH="$(go env GOPATH)/bin:${PATH}" protoc -I../proto \
	--go_out=. \
	--go_opt=Mmusic_browser.proto=github.com/emilia-miki/music-browser/go_proto \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=Mmusic_browser.proto=github.com/emilia-miki/music-browser/go_proto \
	--go-grpc_opt=paths=source_relative \
	music_browser.proto
.PHONY: proto-gen
