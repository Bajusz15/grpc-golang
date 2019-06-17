package main

import (
	"context"
	"fmt"
	"github.com/bajusz15/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
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
func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error{
	fmt.Printf("received PrimeNumberDecomposition RPC: %v\n", req)
	number := req.GetNumber()
	divisor := int64(2)
	//N = 210
	for number > 1{
		if number % divisor == 0{// if k evenly divides into N
			// this is a factor
			_ = stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor // divide N by k so that we have the rest of the number left.
		} else {
			divisor++
			fmt.Println("incremented divisor: %v\n", divisor)
		}
	}
	return nil
}
func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("received ComputeAverage RPC: \n")
	sum := int32(0)
	count := 0
	for  {
		req, err :=stream.Recv()
		if err == io.EOF {
			average := float64(sum)/float64(count)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			fmt.Printf("computeAverage failed, error: %v", err)
			return err
		}
		sum += req.GetNumber()
		count++
	}
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
