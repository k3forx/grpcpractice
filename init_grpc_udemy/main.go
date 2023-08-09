package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang/protobuf/jsonpb"
	"github.com/k3forx/grpcpractice/init_grpc_udemy/pb"
	"google.golang.org/protobuf/proto"
)

func main() {
	employee := &pb.Employee{
		Id:          1,
		Name:        "Suzuki",
		Email:       "test@example.com",
		Occupation:  pb.Occupation_ENGINEER,
		PhoneNumber: []string{"080-1234-5678", "090-1234-5678"},
		Project: map[string]*pb.Company_Project{
			"ProjectX": {},
		},
		Profile: &pb.Employee_Text{
			Text: "My name is Suzuki",
		},
		Birthday: &pb.Date{
			Year: 2000, Month: 1, Day: 1,
		},
	}

	// バイナリデータをwriter/readする
	binData, err := proto.Marshal(employee)
	if err != nil {
		log.Fatalln("can't serialize", err)
	}

	if err := os.WriteFile("test.bin", binData, 0666); err != nil {
		log.Fatalln("can't write", err)
	}

	in, err := os.ReadFile("test.bin")
	if err != nil {
		log.Fatalln("can't read", err)
	}

	readEmployee := &pb.Employee{}
	if err := proto.Unmarshal(in, readEmployee); err != nil {
		log.Fatalln("failed to unmarshal", err)
	}

	fmt.Printf("employee: %+v\n", employee)

	// データをJSONとして扱う
	m := jsonpb.Marshaler{}
	out, err := m.MarshalToString(employee)
	if err != nil {
		log.Fatalln("failed to marshal", err)
	}

	readEmployee = &pb.Employee{}
	if err := jsonpb.UnmarshalString(out, readEmployee); err != nil {
		log.Fatalln("failed to unmarshal", err)
	}
	fmt.Printf("employee: %+v\n", employee)
}
