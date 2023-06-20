#!/bin/sh

PATH="$(go env GOPATH)/bin:$PATH" protoc -I../ \
  --go_out=music_api/ --go_opt=paths=source_relative \
  --go-grpc_out=music_api/ --go-grpc_opt=paths=source_relative \
  music_api.proto
