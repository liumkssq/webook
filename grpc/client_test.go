package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	server := grpc.NewServer()
	userServer := &Server{}
	RegisterUserServiceServer(server, userServer)
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		t.Fatal(err)
	}
	server.Serve(l)
	//server.GracefulStop()
}

func TestClient(t *testing.T) {
	cc, err := grpc.Dial("127.0.0.1:8090",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	client := NewUserServiceClient(cc)
	resp, err := client.GetById(context.Background(), &GetByIdReq{Id: 1})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.User)
}
