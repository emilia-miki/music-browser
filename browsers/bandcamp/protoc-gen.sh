#!/bin/sh

PROTO_DIR=../

npx grpc_tools_node_protoc \
  --js_out=import_style=commonjs,binary:./ \
  --grpc_out=./ \
  --plugin=protoc-gen-grpc=./node_modules/.bin/grpc_tools_node_protoc_plugin \
  -I../../proto \
  music_browser.proto

npx grpc_tools_node_protoc \
  --ts_out=./ \
  --grpc_out=./ \
  --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts \
  -I../../proto \
  music_browser.proto

sed -ie 's/^import \* as grpc from "grpc";$/import * as grpc from "@grpc\/grpc-js";/' music_browser_grpc_pb.d.ts
sed -ie "s/^var grpc = require('grpc');\$/var grpc = require('@grpc\/grpc-js');/" music_browser_grpc_pb.js
