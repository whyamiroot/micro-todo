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
)

func main() {
	envConf := GetConfig()
	registry := NewRegistry()
	registry.StartHealthChecks()
	if envConf.RPCPort == 0 {
		//TODO add logging to the logging service
		panic("No RPC port is specified, unable to start")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", envConf.RPCPort))
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
	var server *grpc.Server
	if !envConf.TLSEnabled {
		server = grpc.NewServer()
	} else {
		server = grpc.NewServer(grpc.Creds(cred))
	}
	proto.RegisterRegistryServiceServer(server, registry)
	go server.Serve(lis)

	mux := runtime.NewServeMux()

	proto.RegisterRegistryServiceHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf(":%d", envConf.RPCPort), opts)
	var httpErr error
	if envConf.TLSEnabled {
		fmt.Println("Starting HTTPS gateway...") //TODO add logging to the logging service
		httpErr = http.ListenAndServeTLS(fmt.Sprintf(":%d", envConf.HTTPSPort), envConf.CertFile, envConf.KeyFile, mux)
	} else {
		fmt.Println("Starting HTTP gateway...")
		httpErr = http.ListenAndServe(fmt.Sprintf(":%d", envConf.HTTPPort), mux)
	}
	if httpErr != nil {
		panic(httpErr.Error())
	}
}
