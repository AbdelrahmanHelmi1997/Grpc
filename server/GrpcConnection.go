package server

import (
	"net"
	"test/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	proto.UnimplementedUserServer
}

func GrpcServerConnection() {

	listener, err := net.Listen("tcp", ":4040")

	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterUserServer(srv, &Server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}
