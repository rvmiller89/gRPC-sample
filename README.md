# gRPC-sample

#### Purpose

This is a sample golang application to practice four types of gRPC communication:

1. Unary RPC
1. Server-to-client Streaming RPC
1. Client-to-server Streaming RPC
1. Bidirectional Streaming RPC

#### Building

##### Protobuf

To compile `kvstore.proto` you will need the `protoc-gen-go` Protobuf compiler plugin for Go. For more details see the [gRPC Quick Start](https://grpc.io/docs/languages/go/quickstart/) guide.

```sh
protoc --go_out=paths=source_relative,plugins=grpc:. kvstore.proto
```

