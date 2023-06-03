FROM golang AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build

FROM fedora
WORKDIR /app
RUN dnf update -y \
    && dnf install -y \
    redis \
    postgresql-server \
    nodejs \
    python3 \
    python3-pip \
    && npm install bandcamp-scraper \
    && pip3 install ytmusicapi
COPY --from=builder /app/music-browser .
COPY bandcamp_api.js .
COPY yt_music_api.py .
RUN mkdir /usr/local/pgsql \
&& chown postgres /usr/local/pgsql \
&& sudo -u postgres initdb -D /usr/local/pgsql/data -A trust

ENV REDIS_PORT=3336
ENV POSTGRES_PORT=3337
ENV BANDCAMP_API_PORT=3334
ENV YT_MUSIC_API_PORT=3335
ENV MUSIC_BROWSER_PORT=3333
EXPOSE 3333

# Start all required servers sequentially
CMD node bandcamp_api.js & python3 yt_music_api.py & redis-server --port $REDIS_PORT --requirepass "" & sudo -u postgres pg_ctl start -D /usr/local/pgsql/data -o "-p ${POSTGRES_PORT}" && ./music-browser

