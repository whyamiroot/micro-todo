package main

import (
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/whyamiroot/micro-todo/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	envConf := LoadConfigFromEnv()
	registry := NewRegistry()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%D", envConf.RPCPort))
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
		os.Exit(1)
	}

	var opts []grpc.DialOption
	var cred credentials.TransportCredentials
	if !envConf.TLSEnabled {
		opts = append(opts, grpc.WithInsecure())
	} else {
		cred, err := credentials.NewServerTLSFromFile(envConf.CertFile, envConf.KeyFile)
		if err != nil {
			fmt.Println("Failed to load TLS credentials")
		}
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	fmt.Println("Starting gRPC server...")
	server := grpc.NewServer(grpc.Creds(cred))
	proto.RegisterRegistryServiceServer(server, registry)
	go server.Serve(lis)

	mux := runtime.NewServeMux()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(envConf.Timeout)*time.Second)
	defer cancel()
	proto.RegisterRegistryServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", envConf.RPCPort), opts)
	var httpErr error
	if envConf.TLSEnabled {
		fmt.Println("Starting HTTPS gateway...")
		httpErr = http.ListenAndServeTLS(fmt.Sprintf(":%d", envConf.HTTPSPort), envConf.CertFile, envConf.KeyFile, mux)
	} else {
		fmt.Println("Starting HTTP gateway...")
		httpErr = http.ListenAndServe(fmt.Sprintf(":%d", envConf.HTTPPort), mux)
	}
	if httpErr != nil {
		panic(httpErr.Error())
	}
}
