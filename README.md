# gRPC-sample

## Purpose

This is a sample golang application to practice four types of gRPC communication:

1. Unary RPC
1. Server-to-client Streaming RPC
1. Client-to-server Streaming RPC
1. Bidirectional Streaming RPC

There are two components, a client and a server, which emulate a [Memcached](https://memcached.org/)-style in-memory key-value store.

After starting the server, you can start the client which serves as a REPL for Memcached-like commands:

```sh
Connected. Commands are:

SET <key> <value>
GET <key>
EXIT

> GET abc
Error getting value rpc error: code = Unknown desc = Missing key
> SET abc 7
> GET abc
Got value 7
```

## Building

### Protobuf

To compile `kvstore.proto` you will need the `protoc-gen-go` Protobuf compiler plugin for Go. For more details see the [gRPC Quick Start](https://grpc.io/docs/languages/go/quickstart/) guide.

```sh
$ protoc --go_out=paths=source_relative,plugins=grpc:. kvstore.proto
```

### Server

```sh
$ go run server/server.go
2020/06/10 19:58:33 Starting server on port :50051
```

### Client

```sh
$ go run client/client.go 
Connecting client on localhost:50051
```