package grpc

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func streamPings(client PingClient, p *PingReq) {
	log.Printf("Sping pykube over gRPC (%s)", p)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.StreamPing(ctx, p)
	if err != nil {
		log.Fatalf("%v.GetPingReq(_) = _, %v: ", client, err)
	}
	for {
		png, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.StreamPings(_) = _, %v", client, err)
		}
		log.Println(png)
	}
}

func getSinglePing(client PingClient, p *PingReq) {
	log.Printf("Single ping pykube over gRPC (%s)", p)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	png, err := client.GetPing(ctx, p)
	if err != nil {
		log.Fatalf("%v.GetPingReq(_) = _, %v: ", client, err)
	}
	log.Println(png)
}

func RunGrpcClient() {
	serverAddr := "localhost:50051"
	fmt.Println("gRPC ping pykube client is running")
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := NewPingClient(conn)

	getSinglePing(client, &PingReq{Request: "ping"})
	fmt.Println("-----------")
	streamPings(client, &PingReq{Request: "pings"})

}
