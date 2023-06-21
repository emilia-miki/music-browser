FROM alpine
RUN apk add \
    redis \
    postgresql \
    protoc \
    go \
    nodejs \
    npm \
    ffmpeg \
    yt-dlp \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

WORKDIR /app
COPY music_api.proto .
COPY database_schema.sql .
COPY music_browser music_browser
COPY ts-proto ts-proto
COPY bandcamp-api bandcamp-api
COPY yt-music-api yt-music-api

EXPOSE 3333
ENV MUSIC_BROWSER_PORT=3333
ENV BANDCAMP_API_PORT=3334
ENV YT_MUSIC_API_PORT=3335
ENV REDIS_PORT=3336
ENV POSTGRES_PORT=3337

RUN cd music_browser && ./build.sh && CGO_ENABLED=0 GOOS=linux go build && cd .. \
&& cd ts-proto && npm install && ./build.sh && cd .. \
&& cd bandcamp-api && npm install && npx tsc && cd .. \
&& cd yt-music-api && npm install && npx tsc && cd .. \
&& mkdir /usr/local/pgsql \
&& chown postgres /usr/local/pgsql \
&& su -m postgres -c 'initdb -D /usr/local/pgsql/data -A trust' \
&& mkdir -p /run/postgresql \
&& chown postgres /run/postgresql

CMD node bandcamp-api > bandcamp-api/bandcamp-api.log 2>&1 \
& node yt-music-api > yt-music-api/yt-music-api.log 2>&1 \
& redis-server --port $REDIS_PORT --requirepass "" \
--maxmemory 1gb --maxmemory-policy allkeys-lfu --save "" \
& su -m postgres -c 'pg_ctl start -D /usr/local/pgsql/data -o "-p ${POSTGRES_PORT}"' \
&& psql -p ${POSTGRES_PORT} -U postgres -d postgres -f database_schema.sql \
&& ./music_browser/music_browser
