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

func streamToPykube(client protokube.BiStreamerClient, b *protokube.Vessel) {
	log.Printf("Send reuquest to Py-Kube over gRPC (%s)", b)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.BidirectionalStream(ctx)
	if err != nil {
		log.Fatalf("%v.BidirectionalStream_) = _, %v: ", client, err)
	}
	time.Sleep(1 * time.Second)
	stream.Send(&protokube.Vessel{Val: 1})
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.BidirectionalStream(_) = _, %v", client, err)
		}
		stream.Send(&protokube.Vessel{Val: response.Val + 1})
		log.Println(response)
	}
}

func RunGrpcClient() {
	serverAddr := "localhost:50051"
	//serverAddr := "pykube-service:50051"

	fmt.Println("Go-Kube RPC client is running...")
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := protokube.NewBiStreamerClient(conn)

	streamToPykube(client, &protokube.Vessel{Val: 1})

}
