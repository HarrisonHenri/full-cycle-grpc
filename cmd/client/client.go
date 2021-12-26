package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/HarrisonHenri/fullcycle-grpc/pb"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUsersBi(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{Id: "any_id", Name: "any_name", Email: "any_email"}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not add the user: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{Id: "any_id", Name: "any_name", Email: "any_email"}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not add the user: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not rcv the msg: %v", err)
		}

		fmt.Println(stream)
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id: "1",
		},
		{
			Id: "2",
		},
		{
			Id: "3",
		},
		{
			Id: "4",
		},
		{
			Id: "5",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Could not add the users: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Could not rcv the response: %v", err)
	}

	fmt.Println(res)

}

func AddUsersBi(client pb.UserServiceClient) {
	stream, err := client.AddUsersBiStream(context.Background())

	if err != nil {
		log.Fatalf("Could not init the stream: %v", err)
	}

	reqs := []*pb.User{
		{
			Id: "1",
		},
		{
			Id: "2",
		},
		{
			Id: "3",
		},
		{
			Id: "4",
		},
		{
			Id: "5",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Id)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Could not rcv the msg: %v", err)
			}

			fmt.Println(res)
		}
		close(wait)
	}()

	<-wait
}
