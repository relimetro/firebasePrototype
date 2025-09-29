package main

import (
	"log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "example/proto_example/protoOut"
)

func main() {
	// connection
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure() )
	if err != nil { log.Fatalf("could not connect, %s", err)}
	defer conn.Close()
	c := pb.NewServClient(conn)

	// send msg
	message := pb.UserLogin { UserName: "name", PlaintextPassword:"pass", }
	response, err := c.Login(context.Background(), &message)
	if err != nil { log.Fatalf("Err: send msg: %s", err) }
	log.Printf("Response from server: %d", response.Temp)

}
