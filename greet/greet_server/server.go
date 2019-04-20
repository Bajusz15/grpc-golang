package main

import (
	"fmt"
	"github.com/bajusz15/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {

}
func main(){
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	s:= grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err:= s.Serve(lis); err != nil{
		log.Fatalf("failed to serve: %v", err)
	}
}
//go run greet/greet_server/server.go
