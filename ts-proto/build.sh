#!/bin/sh

PROTO_DIR=../

npx grpc_tools_node_protoc \
  --js_out=import_style=commonjs,binary:./ \
  --grpc_out=./ \
  --plugin=protoc-gen-grpc=./node_modules/.bin/grpc_tools_node_protoc_plugin \
  -I../ \
  music_api.proto

npx grpc_tools_node_protoc \
  --ts_out=./ \
  --grpc_out=./ \
  --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts \
  -I../ \
  music_api.proto

sed -i '' 's#^import \* as grpc from "grpc";$#import * as grpc from "@grpc/grpc-js";#' music_api_grpc_pb.d.ts
