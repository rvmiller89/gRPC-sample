package main

import (
	"context"
	"errors"
	"log"
	"net"

	pb "github.com/rvmiller89/gRPC-sample"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedKVStoreServer
}

const (
	port = ":50051"
)

var (
	m map[string]int
)

func (s *server) Get(ctx context.Context,
	req *pb.GetRequest) (*pb.GetResponse, error) {
	key := req.GetKey()
	val, exists := m[key]
	if !exists {
		return &pb.GetResponse{}, errors.New("Missing key")
	}
	return &pb.GetResponse{Key: key, Value: int32(val)}, nil
}

func (s *server) Set(ctx context.Context,
	req *pb.SetRequest) (*pb.SetResponse, error) {
	m[req.GetKey()] = int(req.GetValue())
	return &pb.SetResponse{}, nil
}

func main() {
	log.Println("Starting server on port", port)
	m = make(map[string]int)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Unable to listen to on port", port)
	}

	s := grpc.NewServer()
	pb.RegisterKVStoreServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalln("Unable to start server on port", port)
	}
}
