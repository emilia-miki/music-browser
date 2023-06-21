.DEFAULT_GOAL := test

ts-proto:
	cd ts-proto && npm run build
.PHONY: ts-proto

bandcamp-api: ts-proto
	cd bandcamp-api && npm run build
.PHONY: bandcamp-api

yt-music-api: ts-proto
	cd yt-music-api && npm run build
.PHONY: yt-music-api

music-browser: bandcamp-api yt-music-api
	cd music_browser && go build
.PHONY: music-browser

build: music-browser
.PHONY: build

test: build
.PHONY: test
