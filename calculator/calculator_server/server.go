package main

import (
	"context"
	"fmt"
	"github.com/bajusz15/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (*server) Calculator(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	firstNumber := req.First
	secondNumber := req.Second
	sum := firstNumber + secondNumber
	res := &calculatorpb.CalculatorResponse{
		Result: sum,
	}
	return res, nil
}
func main() {
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//go run greet/greet_server/server.go
