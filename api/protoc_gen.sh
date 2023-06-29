#!/bin/sh

mkdir -p proto

PATH="$(go env GOPATH)/bin:$PATH" protoc -I../proto \
  --go_out=proto --go_opt=paths=source_relative \
  --go-grpc_out=proto --go-grpc_opt=paths=source_relative \
  music_browser.proto
