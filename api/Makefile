.DEFAULT_GOAL := build

music_api:
	./protoc-gen.sh

build: music_api
	go build
.PHONY: build
