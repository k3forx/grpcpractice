package main

import (
	"context"
	"fmt"
	"log"

	"github.com/k3forx/grpcpractice/init_grpc_udemy/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v\n", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)
	callListFiles(client)
}

func callListFiles(client pb.FileServiceClient) {
	res, err := client.ListFiles(context.Background(), &pb.ListFilesRequest{})
	if err != nil {
		log.Fatalf("failed to list files: %v\n", err)
	}

	fmt.Printf("filenames: %s\n", res.GetFilenames())
}
