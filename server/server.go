package main

import (
	"context"
	"errors"
	"io"
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

func (s *server) GetBulk(req *pb.GetBulkRequest,
	stream pb.KVStore_GetBulkServer) error {
	keys := req.GetKeys()
	for _, key := range keys {
		s.mutex.Lock()
		val, exists := s.m[key]
		s.mutex.Unlock()
		if !exists {
			return errors.New("Missing key: " + key)
		}
		stream.Send(&pb.GetResponse{Key: key, Value: int32(val)})
	}
	return nil
}

func (s *server) SetBulk(stream pb.KVStore_SetBulkServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.SetResponse{})
		}
		if err != nil {
			return err
		}
		s.mutex.Lock()
		s.m[req.GetKey()] = int(req.GetValue())
		s.mutex.Unlock()
	}
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
