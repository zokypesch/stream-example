package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	proto "github.com/zokypesch/streaming/srv"
	"google.golang.org/grpc"
)

// SimpleStreamming for struct stream
type SimpleStreamming struct{}

// SimpleRPC for register service
func (s *SimpleStreamming) SimpleRPC(stream proto.SimpleService_SimpleRPCServer) error {
	for {
		in, err := stream.Recv()
		log.Println("Received value")
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Println("Got " + in.Msg)
		stream.Send(&proto.SimpleData{Msg: "cihoyyyy"})
	}
}

// SimpleHandler for register service
func (s *SimpleStreamming) SimpleHandler(ctx context.Context, in *proto.SimpleRequestHandler) (*proto.SimpleResponseHandler, error) {
	return nil, fmt.Errorf("failed")
}

func main() {
	fmt.Println("Welcome !")
	grpcServer := grpc.NewServer()
	proto.RegisterSimpleServiceServer(grpcServer, &SimpleStreamming{})

	l, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Listening on tcp://localhost:6000")
	grpcServer.Serve(l)
}
