package main

import (
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/whyamiroot/micro-todo/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
)

func main() {
	registry := NewRegistry()

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Starting gRPC server...")
	server := grpc.NewServer()
	proto.RegisterRegistryServiceServer(server, registry)
	go server.Serve(lis)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Println("Starting HTTP gateway...")
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	proto.RegisterRegistryServiceHandlerFromEndpoint(ctx, mux, ":3000", opts)
	http.ListenAndServe(":8080", mux)
}
