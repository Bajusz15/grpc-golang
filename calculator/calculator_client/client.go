package main

import (
	"context"
	"fmt"
	"github.com/bajusz15/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) //not for production
	if err != nil {
		log.Fatalf("could nto connect %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//log.Printf("created client: %f", c)
	//doUnary(c)
	//doServerStreaming(c)
	 doClientStreaming(c)
}
func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("starting to sum calculator  Unary RPC")
	req := &calculatorpb.CalculatorRequest{
		First:  5,
		Second: 40,
	}
	res, err := c.Calculator(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Calculator RPC: %v", err)
	}
	log.Printf("response from calculator: %v", res.Result)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("starting prime number decompostition streaming RPC")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 12,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("stg bad happened")
		}
		fmt.Println(res.GetPrimeFactor())
	}

}

func doClientStreaming(c calculatorpb.CalculatorServiceClient){
	fmt.Println("starting compute average client streaming RPC")
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		fmt.Println("compute average client streaming failed")
	}
	numbers := []int32{3,5,9,54,23}
	for _, number := range numbers{
		fmt.Printf("sending number: %v", number)
		_ = stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println("error while receiving response")
	}
	fmt.Printf("The average is: %v", res.GetAverage())
}


//go run greet/greet_client/client.go
