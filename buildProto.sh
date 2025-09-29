#!/usr/bin/bash

protoc \
	--proto_path=./proto \
	--go_out=./protoOut --go_opt=paths=source_relative \
	--go-grpc_out=./protoOut --go-grpc_opt=paths=source_relative \
	proto.proto

# sudo dnf install protoc protoc-gen-go protoc-gen-go-grpc

# tryed : grpc-plugins (for android)

