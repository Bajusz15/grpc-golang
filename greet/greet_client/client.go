package main

import (
	"context"
	"fmt"
	"github.com/bajusz15/grpc-go/greet/greetpb"
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

	c := greetpb.NewGreetServiceClient(cc)

	//log.Printf("created client: %f", c)
	//doUnary(c)
	//interop.DoServerStreaming(c)
	doServerStreaming(c)
}
func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do Unary RPC")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Mate",
			LastName:  "Bajusz",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	fmt.Println("response from greet: %v ", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do a server streaming RPC...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Mate",
			LastName:  "Bajusz",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			fmt.Println("end of stream")
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		fmt.Println("response from GreetManyTimes: %v", msg.GetResult())
	}

}

//go run greet/greet_client/client.go
