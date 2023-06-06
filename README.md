# Music Browser

This will be an overengineered backend in Go to show that I can work with Go,
SQL, Redis, Docker, Kubernetes, Terraform, AWS. JSON, GraphQL and gRPC.

The idea is to create a client-server web app that can browse music from 
Youtube Music, Spotify and Bandcamp in one place, as well as download it using
yt-dlp, convert into a different format with ffmpeg, and then browse and listen
to the downloaded music.

What skills are currently demostrated by the application:
- :white_check_mark: Go
- :x: SQL
- :x: Redis
- :white_check_mark: Docker
- :x: Kubernetes
- :x: Terraform
- :x: AWS
- :white_check_mark: JSON
- :x: GraphQL
- :white_check_mark: gRPC

Additionally:
- :white_check_mark: Python and NodeJS scripting

What else doesn't work yet:
- Downloading tracks
- Exploring local backend

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
--data-urlencode search-type=artist \
--data-urlencode 'query=Ne Obliviscaris' \
| python3 -m json.tool
```
Piping to python3 -m json.tool just formats the JSON, because the application
returns it minified by default. 

Possible backend values: spotify, bandcamp, yt-music, local.
Possible search-type values: album, artist.

Or a POST request to download a song:
```
> curl --request POST -G \
--url 'localhost:3333' \
--data-urlencode url='https://dreamcatalogue.bandcamp.com/track/--436' \
| python3 -m json.tool
```
