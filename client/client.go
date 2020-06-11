package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
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
	GETBULK <key1>,<key2>,...
	SETBULK <key1> <value1>,<key2> <value2>,...
	EXIT`
	timeout = 10 * time.Second
)

var client pb.KVStoreClient

func set(key string, value int) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := client.Set(ctx, &pb.SetRequest{Key: key, Value: int32(value)})
	if err != nil {
		fmt.Println("Error setting key", err)
	}
}

func get(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.Get(ctx, &pb.GetRequest{Key: key})
	if err != nil {
		fmt.Println("Error getting value", err)
		return
	}
	fmt.Println("Got value", resp.GetValue())
}

func getBulk(keys []string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	stream, err := client.GetBulk(ctx, &pb.GetBulkRequest{Keys: keys})
	if err != nil {
		fmt.Println("Error getting stream", err)
		return
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error getting value", err)
			return
		}
		fmt.Printf("Key: %v has value: %v\n", resp.GetKey(), resp.GetValue())
	}
}

func setBulk(keyValues []*pb.SetRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	stream, err := client.SetBulk(ctx)
	if err != nil {
		fmt.Println("Error getting stream", err)
		return
	}
	for _, kv := range keyValues {
		if err := stream.Send(kv); err != nil {
			fmt.Println("Error setting value", err)
			return
		}
	}
	_, err = stream.CloseAndRecv()
	if err != nil {
		fmt.Println("Error closing stream", err)
		return
	}
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
			case "GETBULK":
				keys := strings.Split(cmd[1], ",")
				getBulk(keys)
			case "SETBULK":
				rest := strings.TrimPrefix(line, "SETBULK ")
				keyValues := strings.Split(rest, ",")
				keyValueArr := []*pb.SetRequest{}
				for _, kvStr := range keyValues {
					kv := strings.Fields(kvStr)
					key := kv[0]
					val, _ := strconv.Atoi(kv[1])
					keyValueArr = append(keyValueArr,
						&pb.SetRequest{Key: key, Value: int32(val)})
				}
				setBulk(keyValueArr)
			default:
				fmt.Println("Unknown command:", cmd[0])
			}
		}
	}
}
