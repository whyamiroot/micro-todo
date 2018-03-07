#!/bin/sh

prot=$1
protoc -I $prot -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/options \
-I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--grpc-gateway_out=logtostderr=true:$prot/ \
--swagger_out=logtostderr=true:$prot/ \
--go_out=plugins=grpc:$prot $prot/*.proto