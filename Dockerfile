FROM golang AS builder
WORKDIR /app
COPY music_browser music_browser
WORKDIR /app/music_browser
RUN CGO_ENABLED=0 GOOS=linux go build

FROM fedora
RUN dnf update -y \
    && dnf install -y \
    redis \
    postgresql-server \
    nodejs \
    python3 \
    python3-pip \
    && pip3 install ytmusicapi grpcio-tools

WORKDIR /app
COPY music_api.proto .
COPY bandcamp-api bandcamp-api
COPY yt_music_api yt_music_api
COPY --from=builder /app/music_browser/music_browser .

RUN cd bandcamp-api && npm install && cd .. \
&& mkdir /usr/local/pgsql \
&& chown postgres /usr/local/pgsql \
&& sudo -u postgres initdb -D /usr/local/pgsql/data -A trust

EXPOSE 3333
ENV MUSIC_BROWSER_PORT=3333
ENV BANDCAMP_API_PORT=3334
ENV YT_MUSIC_API_PORT=3335
ENV REDIS_PORT=3336
ENV POSTGRES_PORT=3337

CMD node bandcamp-api > bandcamp-api/bandcamp-api.log 2>&1 \
& python3 yt_music_api > yt_music_api/yt_music_api.log 2>&1 \
& redis-server --port $REDIS_PORT --requirepass "" \
& sudo -u postgres pg_ctl start -D /usr/local/pgsql/data -o "-p ${POSTGRES_PORT}" \
&& ./music_browser

