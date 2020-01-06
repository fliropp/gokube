package grpc

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 10000, "The server port")
)

type gokrpcServer struct {
	msg string
}

func (ps *gokrpcServer) GetPing(ctx context.Context, p *PingReq) (*PingResp, error) {
	return &PingResp{Response: "pong"}, nil
}

func (s *gokrpcServer) StreamPing(p *PingReq, stream Ping_StreamPingServer) error {
	kubez := [5]string{"gokube1", "gokube2", "gokube3", "gokube4", "gokube5"}
	for _, k := range kubez {
		err := stream.Send(&PingResp{Response: k})
		if err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

func newServer() PingServer {
	s := &gokrpcServer{msg: "how many gokubez?"}
	return s
}

func RunGrpcServer() {
	fmt.Println("gokRPC server is running")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = testdata.Path("server1.pem")
		}
		if *keyFile == "" {
			*keyFile = testdata.Path("server1.key")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	RegisterPingServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
