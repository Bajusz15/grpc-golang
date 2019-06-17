package main

import (
	"context"
	"fmt"
	"github.com/bajusz15/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello" + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + "number" + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		_ = stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet function was invoked with streaming request")
	result := "hello"
	for {
		req, err := stream.Recv()
		if err == io.EOF{
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
			fmt.Println("finished reading the client")
		}
		if err != nil {
			fmt.Println("valami nem jo a client streamnel: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		//lastName := req.GetGreeting().GetLastName()
		result += " " + firstName + "!"
	}
}

func (*server) GreetEveryOne(stream greetpb.GreetService_GreetEveryOneServer) error{
	fmt.Println("GreetEveryOne function was invoked with streaming request")

	for {
		req, err := stream.Recv()
		if err == io.EOF{
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "
		err = stream.Send(&greetpb.GreetEveryOneResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("Error while sending data to cloent: %v", err)
			return err
		}
	}
}
func main() {
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//go run greet/greet_server/server.go
