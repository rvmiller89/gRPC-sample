package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/rvmiller89/gRPC-sample"
	"google.golang.org/grpc"
)

const (
	address  = "localhost:50051"
	commands = `
SET <key> <value>
GET <key>
EXIT`
)

var client pb.KVStoreClient

func set(key string, value int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := client.Set(ctx, &pb.SetRequest{Key: key, Value: int32(value)})
	if err != nil {
		fmt.Println("Error setting key", err)
	}
}

func get(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.Get(ctx, &pb.GetRequest{Key: key})
	if err != nil {
		fmt.Println("Error getting value", err)
		return
	}
	fmt.Println("Got value", resp.GetValue())
}

func main() {
	fmt.Println("Connecting client on", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Println("Could not connect: %v", err)
	}
	defer conn.Close()

	client = pb.NewKVStoreClient(conn)
	fmt.Printf("Connected. Commands are:\n%v\n\n", commands)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		line := scanner.Text()
		if line == "EXIT" {
			break
		} else if strings.TrimSpace(line) == "" {
			continue
		} else {
			cmd := strings.Fields(line)
			switch cmd[0] {
			case "SET":
				val, _ := strconv.Atoi(cmd[2])
				set(cmd[1], val)
			case "GET":
				get(cmd[1])
			}
		}
	}
}
