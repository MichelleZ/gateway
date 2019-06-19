/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */

package etcd

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	etcd3 "go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"

	"github.com/xuperchain/xuperunion/common/log"
	pb "github.com/xuperchain/xuperunion/gateway/consistent_hash/helloworld"
	loggw "github.com/xuperchain/xuperunion/gateway/log"
)

func (ts *testServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

type testServer struct {
	server *grpc.Server
	addr   string
}

func NewtestServer(addr string) *testServer {
	server := grpc.NewServer()
	return &testServer{
		server: server,
		addr:   addr,
	}
}

func (ts *testServer) Run() {
	lis, err := net.Listen("tcp", ts.addr)
	if err != nil {
		fmt.Errorf("failed to listen %v.", err)
	}

	pb.RegisterGreeterServer(ts.server, ts)
	ts.server.Serve(lis)
}

func (ts *testServer) Stop() {
	ts.server.GracefulStop()
}

func TestResolver(t *testing.T) {
	etcdConfig := etcd3.Config{
		Endpoints: []string{"http://10.100.78.75:2379"},
	}

	xlog, _ := log.OpenLog(loggw.CreateLog())

	registry, err := NewRegistry(
		Input{
			EtcdConfig: etcdConfig,
			Prefix:     "gateway/testResolver", // real: etchxchain/gateway/testResolver/
			Addr:       "10.100.78.75:11000",
			TTL:        5,
			Xlog:       xlog,
		})
	if err != nil {
		t.Errorf("Register failed: %v.", err)
	}

	server := NewtestServer(fmt.Sprintf("10.100.78.75:11000"))

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		registry.Register()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Run()
	}()

	r := NewResolver(etcdConfig, xlog)
	resolver.Register(r)

	cc, err := grpc.Dial(r.Scheme()+":///gateway/testResolver", grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		t.Fatalf("failed to dial: %v.", err)
	}
	defer cc.Close()

	testc := pb.NewGreeterClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	fmt.Println("call SayHello")
	res, err := testc.SayHello(ctx, &pb.HelloRequest{Name: "world"})
	if err != nil {
		t.Fatalf("SayHello() err: %v.", err)
	}
	t.Logf("res: %v", res.Message)

	go func() {
		time.Sleep(10 * time.Second)
		t.Log("Unregister && server stop")
		server.Stop()
		registry.UnRegister()
	}()

	wg.Wait()
}
