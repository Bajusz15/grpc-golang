package main

import (
	"context"
	"fmt"
	"github.com/bajusz15/grpc-go/blog/blogpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
fmt.Println("Blog client")
cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) //not for production
if err != nil {
log.Fatalf("could nto connect %v", err)
}

defer cc.Close()

c := blogpb.NewBlogServiceClient(cc)
blog := &blogpb.Blog{
	AuthorId: "MAte",
	Title: "valami",
	Content: "tartalom",
}
	response, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("blog has been created: %v", response)
	blogId := response.GetBlog().GetId()

	//read Blog
	fmt.Println("reading the blog")

	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "5d076708ec6b9b2270db3657"})
	if err2 != nil {
		fmt.Printf("error happened whiel reading: %v\n", err2)
	}
	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogId }
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("error happened whiel reading: %v", readBlogErr)
	}
	fmt.Printf("blog was read: %v", readBlogRes)

	//update the blog
	newBlog := &blogpb.Blog{
		Id: blogId,
		AuthorId: "Bajusz",
		Title: "valami 2",
		Content: "tartalom updated",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: newBlog,
	})
	if updateErr != nil {
		fmt.Printf("error while updating: %v", updateErr)
	}
	fmt.Printf("Blog was read: %v", updateRes)

	//delete Blog
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogId})
	if deleteErr != nil {
		fmt.Printf("error while deleting: %v\n", deleteErr)
	}
	fmt.Printf("Blog was deleted: %v", deleteRes)

	// list blogs
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})

	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stg bad happened; %v", err)
		}
		fmt.Println(res.GetBlog())
	}

}