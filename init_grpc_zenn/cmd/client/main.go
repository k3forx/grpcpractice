package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	hellopb "github.com/k3forx/grpcpractice/init_grpc_zenn/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	scanner *bufio.Scanner
	client  hellopb.GreetingServiceClient
)

func main() {
	fmt.Println("start gRPC client")

	scanner = bufio.NewScanner(os.Stdin)

	address := "localhost:8080"
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("connection failed")
		return
	}
	defer conn.Close()

	client = hellopb.NewGreetingServiceClient(conn)

	for {
		fmt.Println("1: Send request")
		fmt.Println("2: exit")
		fmt.Print("please enter >")

		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			hello()
		case "2":
			fmt.Println("bye...")
			goto M
		}

		fmt.Println("")
	}
M:
}

func hello() {
	fmt.Println("Please enter your name")
	scanner.Scan()

	req := &hellopb.HelloRequest{
		Name: scanner.Text(),
	}

	res, err := client.Hello(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.GetMessage())
}
