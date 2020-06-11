# gRPC-sample

## Purpose

This is a sample golang application to practice gRPC communication.

There are two components, a client and a server, which emulate a [Memcached](https://memcached.org/)-style in-memory key-value store.

Although we only store integer values in this example, for the purpose of this exercise we can assume the value has a much larger memory footprint than the key, necessitating streaming values for bulk operations.

This way we can exercise different styles of gRPC communication:

1. Unary RPC (`GET`, `SET`)
1. Server-to-client Streaming RPC (`GETBULK`)
1. Client-to-server Streaming RPC (`SETBULK`)

The following style is not yet exercised in this example:

1. Bidirectional Streaming RPC

## Usage

After starting the server (see [Building](#building) section below), you can start the client which serves as a REPL for Memcached-like commands:

```sh
Connected. Commands are:

        SET <key> <value>
        GET <key>
        GETBULK <key1>,<key2>,...
        SETBULK <key1> <value1>,<key2> <value2>,...
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