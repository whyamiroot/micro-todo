FROM ubuntu:16.04 as protobuf

WORKDIR ~

RUN apt-get update && apt-get install -y autoconf automake libtool curl make g++ unzip git
RUN git clone https://github.com/google/protobuf.git && cd protobuf && git checkout 3.5.x && ./autogen.sh && ./configure --prefix=/usr && make -j 4 && make install && ldconfig

FROM golang:1.9.2

MAINTAINER Lev Sizov levsizov2015@outlook.com

WORKDIR /go/src/github.com/whyamiroot/micro-todo

RUN go get -u google.golang.org/grpc
ARG CACHEBUST=1
ONBUILD RUN go build -o ./build/registry.bin github.com/whyamiroot/micro-todo/registry

ENTRYPOINT ["./build/registry.bin"]
