#!/usr/bin/bash

protoc \
	--proto_path=./proto \
	--go_out=./protoOut --go_opt=paths=source_relative \
	--go-grpc_out=./protoOut --go-grpc_opt=paths=source_relative \
	proto.proto

protoc \
	--proto_path=./proto \
	--go_out=./protoOut --go_opt=paths=source_relative \
	--go-grpc_out=./protoOut --go-grpc_opt=paths=source_relative \
	aiProompt.proto

# protoc --include_imports --include_source_info -o proto/firebase.pb proto/proto.proto
GOOGLEAPIS_DIR="../../googleapis"
protoc -I${GOOGLEAPIS_DIR} -I. --include_imports --include_source_info --descriptor_set_out=proto.pb proto/proto.proto ../../proto/dummy.proto

# sudo dnf install protoc protoc-gen-go protoc-gen-go-grpc

# tryed : grpc-plugins (for android)

