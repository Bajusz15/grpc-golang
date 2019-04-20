package  main

import (
	"fmt"
	"github.com/bajusz15/grpc-go/greet/greetpb"
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

	c := greetpb.NewGreetServiceClient(cc)
	log.Printf("created client: %f", c)
}
//go run greet/greet_client/client.go