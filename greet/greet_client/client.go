package main

import (
	"context"
	"fmt"
	"github.com/bajusz15/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
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
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDirectionalStreaming(c)
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

func doClientStreaming(c greetpb.GreetServiceClient){
	fmt.Println("starting to do a client streaming RPC...")
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mate",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "someone",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "yournamehere",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "another name",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bajusz",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		fmt.Println("doClientStreaming failed")
	}
	//iterate over our slice
	for _, req := range requests{
		fmt.Println("sending request: ", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println("stream.CloseAndRecv() failed")
	}
	fmt.Printf("LongGreet Response: %v\n", res)
}

func doBiDirectionalStreaming(c greetpb.GreetServiceClient){
	fmt.Println("starting to do a Client Streaming RPC")

	// we create a stream by invoking client
	stream, err := c.GreetEveryOne(context.Background())
	if err != nil {
		log.Fatalf("error while crteating stream: %v", err)
		return
	}
	requests := []*greetpb.GreetEveryOneRequest{
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mate",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "someone",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "yournamehere",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "another name",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bajusz",
			},
		},
	}
	waitc := make(chan struct{})
	//we send bunch of messages to the server (go routine)
	go func() {
		// send messages
		for _, req := range requests{
			fmt.Printf("sending message: %v", req)
			err :=stream.Send(req)
			if err != nil {
				log.Printf("error sending messages to stream: %v",err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// we receivew  bunch of messages from the server (go routine)
	go func() {
		//receive messages
		for  {
			res, err :=stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving: %v", err)
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	//block until everythign is done
	<-waitc
}

//go run greet/greet_client/client.go
