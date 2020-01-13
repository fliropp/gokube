package grpcclient

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/fliropp/gokube/pkg/protokube"
	"google.golang.org/grpc"
)

func streamBullshit(client protokube.StreamerClient, b *protokube.BullshitIn) {
	log.Printf("Stream bullshit from pykube to gokube over gRPC (%s)", b)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.StreamBullshit(ctx, b)
	if err != nil {
		log.Fatalf("%v.StreamBullshit(_) = _, %v: ", client, err)
	}
	for {
		bs, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.StreamBullshit(_) = _, %v", client, err)
		}
		log.Println(bs)
	}
}

func RunGrpcClient() {
	//serverAddr := "192.168.64.2:50051"
	serverAddr := "pykube-service:50051"

	fmt.Println("gRPC client is running...")
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := protokube.NewStreamerClient(conn)

	streamBullshit(client, &protokube.BullshitIn{Bi: "unleash the bullshit....!"})

}
