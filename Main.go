package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/zokypesch/streaming/srv"

	"google.golang.org/grpc"
)

func main() {
	conn, err := OpenConnectionGRPC(6000)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	client := srv.NewSimpleServiceClient(conn)
	stream, err := client.SimpleRPC(context.Background())
	waitc := make(chan struct{})

	msg := &srv.SimpleData{Msg: "cuyyyy"}
	go func() {
		for {
			log.Println("Sleeping...")
			time.Sleep(2 * time.Second)
			log.Println("Sending msg...")
			stream.Send(msg)
			ab, _ := stream.Recv()
			fmt.Println(ab.Msg)
		}
	}()
	<-waitc
	stream.CloseSend()

}

// OpenConnectionGRPC for open connection GRPC client
func OpenConnectionGRPC(port int) (*grpc.ClientConn, error) {
	if port == 0 {
		return nil, fmt.Errorf("failed to open connection")
	}

	conn, err := grpc.Dial("localhost:"+strconv.Itoa(port), grpc.WithInsecure())

	return conn, err
}
