package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/whyamiroot/micro-todo/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Logger struct {
	Lock          *sync.Mutex
	Config        *EnvironmentConfig
	LogCountCache map[string]map[uint32]uint64
}

func NewLogger(config *EnvironmentConfig) {

}

func (l *Logger) StartLoggerServiceAndListen() {
	if l.Config == nil {
		panic("No configuration")
	}
	envConf := l.Config

	if envConf.RPCPort == 0 {
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
	proto.RegisterLoggerServiceServer(server, l)
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

func (l *Logger) GetHealth(c context.Context, _ *proto.Empty) (*proto.Health, error) {
	panic("implement me")
}

func (l *Logger) AddLog(c context.Context, entry *proto.LogEntry) (*proto.LoggerResponse, error) {
	panic("implement me")
}

func (l *Logger) GetLogStream(info *proto.SignedServiceInfo, streamServer proto.LoggerService_GetLogStreamServer) error {
	panic("implement me")
}

func (l *Logger) GetLogStreamWithConstraint(r *proto.ConstraintedLogRequest, streamServer proto.LoggerService_GetLogStreamWithConstraintServer) error {
	panic("implement me")
}

func (l *Logger) GetLogEntriesCount(c context.Context, info *proto.SignedServiceInfo) (*proto.LogEntriesCount, error) {
	panic("implement me")
}
