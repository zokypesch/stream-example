package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/zokypesch/streaming/srv"

	"google.golang.org/grpc"
)

// Tambahan buat web
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Stream streaming function
func Stream(cli srv.SimpleServiceClient, ws *websocket.Conn) error {
	// Read initial request from websocket
	var req srv.SimpleData
	err := ws.ReadJSON(&req)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("woiii")
	// req := srv.SimpleData{Msg: "cihuyyyy"}

	// Even if we aren't expecting further requests from the websocket, we still need to read from it to ensure we
	// get close signals
	go func() {
		for {
			if _, _, err := ws.NextReader(); err != nil {
				break
			}
		}
	}()

	log.Printf("Received Request: %v", req)

	// Send request to stream server
	stream, err := cli.SimpleRPC(context.Background())

	stream.Send(&req)
	stream.CloseSend()

	if err != nil {
		log.Println(err)
		return err
	}
	// defer stream.Close()

	// Read from the stream server and pass responses on to websocket
	for {
		// Read from stream, end request once the stream is closed
		rsp, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				return err
			}

			break
		}

		// Write server response to the websocket
		err = ws.WriteJSON(rsp)
		if err != nil {
			// End request if socket is closed
			if isExpectedClose(err) {
				log.Println("Expected Close on socket", err)
				break
			} else {
				return err
			}
		}
	}

	return nil
}

func isExpectedClose(err error) bool {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
		log.Println("Unexpected websocket close: ", err)
		return false
	}

	return true
}

func main() {
	conn, err := OpenConnectionGRPC(6000)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	client := srv.NewSimpleServiceClient(conn)
	// stream, err := client.SimpleRPC(context.Background())
	// waitc := make(chan struct{})

	// bidirect

	// msg := &srv.SimpleData{Msg: "cuyyyy"}
	// go func() {
	// 	for {
	// 		log.Println("Sleeping...")
	// 		time.Sleep(2 * time.Second)
	// 		log.Println("Sending msg...")
	// 		stream.Send(msg)
	// 		ab, _ := stream.Recv()
	// 		fmt.Println(ab.Msg)
	// 	}
	// }()

	// serverside
	// stream.Send(msg)
	// for i := 0; i < 5; i++ {
	// 	ab, _ := stream.Recv()
	// 	fmt.Println(ab.Msg)
	// }
	// stream.CloseSend()

	// use http
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal("Upgrade: ", err)
			return
		}
		defer conn.Close()

		// Handle websocket request
		if err := Stream(client, conn); err != nil {
			log.Fatal("Echo: ", err)
			return
		}
		log.Println("Stream complete")

	})

	http.ListenAndServe(":9070", nil)

	// <-waitc

}

// OpenConnectionGRPC for open connection GRPC client
func OpenConnectionGRPC(port int) (*grpc.ClientConn, error) {
	if port == 0 {
		return nil, fmt.Errorf("failed to open connection")
	}

	conn, err := grpc.Dial("localhost:"+strconv.Itoa(port), grpc.WithInsecure())

	return conn, err
}
