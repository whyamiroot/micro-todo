FROM ubuntu:16.04 as protobuf

WORKDIR ~

RUN apt-get update && apt-get install -y autoconf automake libtool curl make g++ unzip git
RUN git clone https://github.com/google/protobuf.git && cd protobuf && git checkout 3.5.x && ./autogen.sh && ./configure --prefix=/usr && make -j 4 && make install && ldconfig

FROM golang:1.9.2 as todo

RUN go get -u github.com/golang/dep/cmd/dep

EXPOSE 3000
EXPOSE 80

WORKDIR /go/src/github.com/whyamiroot/micro-todo

COPY . .

RUN dep ensure
RUN ./build.sh

CMD ./registry.bin
