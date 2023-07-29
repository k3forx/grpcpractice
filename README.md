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