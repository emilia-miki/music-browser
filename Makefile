.DEFAULT_GOAL := build

ts-proto:
	make -C ts-proto
.PHONY: ts-proto

bandcamp-api: ts-proto
	make -C bandcamp-api
.PHONY: bandcamp-api

yt-music-api: ts-proto
	make -C yt-music-api
.PHONY: yt-music-api

music_browser:
	make -C music_browser
.PHONY: music_browser

build: bandcamp-api yt-music-api music_browser
.PHONY: build
