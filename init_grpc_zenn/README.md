# gRPC

## バージョンなど

```bash
$ go version
go version go1.20.5 darwin/amd64
```

## ディレクトリ構成

```bash
tree .
.
├── README.md
├── api
│   └── hello.proto
├── go.mod
├── go.sum
└── pkg
    └── grpc

4 directories, 4 files
```

## CLIコマンド/Goパッケージのインストール

- `protoc`: protoファイルからコードを自動生成するためのコマンド

```bash
$ brew install protobuf

$ which protoc
/usr/local/bin/protoc
```

- [https://github.com/grpc/grpc-go](https://github.com/grpc/grpc-go): GoでgRPCを扱うためのパッケージ
- [https://github.com/grpc/grpc-go]

```bash
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

$ export PATH="$PATH:$(go env GOPATH)/bin"
```

## コードを生成

- `protoc` コマンドで `.proto` 拡張子のファイルからgPRC用のコードを生成する

```bash
protoc --go_out=../pkg/grpc --go_opt=paths=source_relative \
       --go-grpc_out=../pkg/grpc --go-grpc_opt=paths=source_relative \
       hello.proto
```

### 生成されたファイル

- `pkg/grpc/hello.pb.go`: protoファイル内で定義したメッセージ`HelloRequest`/`HelloResponse`型を、Goの構造体に定義しなおしたものが自動生成される
- `pkg/grpc/hello_grpc.pb.go`: protoファイルから自動生成されたサービス部分のコード

### `pkg/grpc/hello.pb.go` の内容

```go
// 型の定義
type HelloRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

type HelloResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}
```

### `pkg/grpc/hello_grpc.pb.go` の内容

```go
// GreetingServiceClient is the client API for GreetingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GreetingServiceClient interface {
	// サービスが持つメソッドの定義
	Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
}

// GreetingServiceServer is the server API for GreetingService service.
// All implementations must embed UnimplementedGreetingServiceServer
// for forward compatibility
type GreetingServiceServer interface {
	// サービスが持つメソッドの定義
	Hello(context.Context, *HelloRequest) (*HelloResponse, error)
	mustEmbedUnimplementedGreetingServiceServer()
}
```

## gRPCサーバーの実装

### サーバーの実装

- gRPCサーバーを `localhost:8080` で動かすためのコード

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	hellopb "github.com/k3forx/grpcpractice/pkg/grpc"
)

type myServer struct {
	hellopb.UnimplementedGreetingServiceServer
}

func (s *myServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}

func NewMyServer() *myServer {
	return &myServer{}
}

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

    // hellopb.RegisterGreetingServiceServerのシグネチャ
    // RegisterGreetingServiceServer(s grpc.ServiceRegistrar, srv GreetingServiceServer)
	hellopb.RegisterGreetingServiceServer(s, NewMyServer())

	go func() {
		log.Printf("start gRPC server port: %d", port)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
```

### gRPCurlのインストール

-  gPRCのサーバーに `curl` コマンドのようにリクエストを送るためのツール

```bash
$ brew install grpcurl

$ which grpcurl
```

### サーバーリフレクションの設定

- [サーバーリフレクションについて](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md)

```go
	reflection.Register(s)
```

### 動作確認

- サーバーを起動

```bash
$ cd cmd/server

$ go run ./main.go
2023/08/02 21:54:03 start gRPC server port: 8080
```

- リクエストを送る

```bash
$ grpcurl -plaintext -d '{"name": "k3forx"}' localhost:8080 myapp.GreetingService.Hello | jq -r
{
  "message": "Hello, k3forx!"
}
```

## gRPCクライアントの実装

- サーバーを起動

```bash
$ cd cmd/server

$ go run ./main.go
2023/08/02 22:14:53 start gRPC server port: 8080
```

- クライアントを起動

```bash
$ cd cmd/client/

$ go run main.go
start gRPC client
1: Send request
2: exit
please enter >1
Please enter your name
k3forx
Hello, k3forx!

1: Send request
2: exit
please enter >2
bye...
```

## gRPCで実現できるストリーミング処理

### gRPCで可能な4種類の通信方式

- Unary RPC
- Server streaming RPC
- Client streaming RPC
- Bidirectional streaming RPC

## gRPCにおけるステータスコード

- gRPCの場合は**メソッドの呼び出しに成功した場合には、中で何が起ころうともHTTPレスポンスステータスコードは`200 OK`を返す**
