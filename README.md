# Music Browser

This will be an overengineered backend in Go to show that I can work with Go,
SQL, Redis, Docker, Kubernetes, Terraform, AWS, JSON, GraphQL and gRPC.

The idea is to create a client-server web app that can browse music from 
Youtube Music, Spotify and Bandcamp in one place, as well as download it using
yt-dlp, convert into a different format with ffmpeg, and then browse and listen
to the downloaded music.

What skills are currently demostrated by the application:
- :white_check_mark: Go
- :white_check_mark: SQL
- :white_check_mark: Redis
- :white_check_mark: Docker
- :x: Kubernetes
- :x: Terraform
- :x: AWS
- :white_check_mark: JSON
- :white_check_mark: GraphQL
- :white_check_mark: gRPC

Additionally:
- :white_check_mark: NodeJS scripting

What doesn't work yet:
- :x: Exploring music locally

You can build this app by running 
```
> docker build -t music-browser .
```

When running the built docker image, set SPOTIFY_CLIENT_ID and
SPOTIFY_CLIENT_SECRET environment variables and map the port 3333:
```
> docker run -p 3333:3333 \
-e SPOTIFY_CLIENT_ID=... \
-e SPOTIFY_CLIENT_SECRET=... \
music-browser
```
You can get them at https://developer.spotify.com/dashboard.

Then you can send requests to the application like so:
```
> curl --request GET -G \
--url 'localhost:3333' \
--data-urlencode backend=yt-music \
--data-urlencode type=artist \
--data-urlencode 'query=Ne Obliviscaris'
```
I recommend piping the response to some json formatting tool, as the response
is minified by default.

Possible backend values: spotify, bandcamp, yt-music, local.
Possible type values: artists, albums.

Or a POST request to download a song:
```
> curl --request POST -G \
--url 'localhost:3333' \
--data-urlencode url='https://dreamcatalogue.bandcamp.com/track/--436'
```

There is also a working GraphQL endpoint at localhost:3333/graphql.
The schema is defined in the file [graphql.go](music_browser/graphql/graphql.go)
