package main

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"

	pb "github.com/rvmiller89/gRPC-sample"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedKVStoreServer
	m     map[string]int
	mutex sync.Mutex
}

func (s *server) Get(ctx context.Context,
	req *pb.GetRequest) (*pb.GetResponse, error) {
	key := req.GetKey()
	s.mutex.Lock()
	val, exists := s.m[key]
	defer s.mutex.Unlock()
	if !exists {
		return &pb.GetResponse{}, errors.New("Missing key")
	}
	return &pb.GetResponse{Key: key, Value: int32(val)}, nil
}

func (s *server) Set(ctx context.Context,
	req *pb.SetRequest) (*pb.SetResponse, error) {
	s.mutex.Lock()
	s.m[req.GetKey()] = int(req.GetValue())
	s.mutex.Unlock()
	return &pb.SetResponse{}, nil
}

func main() {
	log.Println("Starting server on port", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Unable to listen to on port", port)
	}
	server := server{m: make(map[string]int)}
	grpcServer := grpc.NewServer()
	pb.RegisterKVStoreServer(grpcServer, &server)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Unable to start server on port", port)
	}
}
