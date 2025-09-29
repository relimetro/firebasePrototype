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

	message2 := pb.UserRequest { UserId: "name", SessionToken:response.Temp, }
	response2, err := c.GetDetails(context.Background(), &message2)
	if err != nil { log.Fatalf("Err: send msg2: %s", err) }
	log.Printf("Response from server: %s", response2.Details)

	message3 := pb.UserRequest { UserId: "nam", SessionToken:response.Temp, }
	response3, err := c.GetDetails(context.Background(), &message3)
	if err != nil { log.Fatalf("Err: send msg3: %s", err) }
	log.Printf("Response from server: %s", response3.Details)

}
