package main

import (
	"context"
	"fmt"
	"github.com/bajusz15/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
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
	doUnary(c)
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

//go run greet/greet_client/client.go
