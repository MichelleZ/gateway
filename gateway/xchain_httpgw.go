/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */

package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"github.com/xuperchain/xuperunion/pb"
)

var (
	rpcEndpoint = flag.String("gateway_endpoint", "localhost:50089", "endpoint of grpc service forward to")
	// http port
	httpEndpoint = flag.String("http_endpoint", ":8098", "endpoint of http service")
	// InitialWindowSize window size
	InitialWindowSize int32 = 128 << 10
	// InitialConnWindowSize connetion window size
	InitialConnWindowSize int32 = 64 << 10
	// ReadBufferSize buffer size
	ReadBufferSize = 32 << 10
	// WriteBufferSize write buffer size
	WriteBufferSize = 32 << 10
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithInitialWindowSize(InitialWindowSize), grpc.WithWriteBufferSize(WriteBufferSize), grpc.WithInitialConnWindowSize(InitialConnWindowSize), grpc.WithReadBufferSize(ReadBufferSize)}
	err := pb.RegisterXchainHandlerFromEndpoint(ctx, mux, *rpcEndpoint, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(*httpEndpoint, mux)
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		panic(err)
	}
}
