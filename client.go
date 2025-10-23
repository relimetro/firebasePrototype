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
	// conn, err := grpc.Dial(":9000", grpc.WithInsecure() )
	conn, err := grpc.Dial(":9000", grpc.WithInsecure() )
	if err != nil { log.Fatalf("GRPC: could not connect,\n%s", err)}
	defer conn.Close()
	c := pb.NewServClient(conn)

	// Login
	message := pb.UserLogin { UserName: "name", PlaintextPassword:"pass", }
	mySession, err := c.Login(context.Background(), &message)
	if err != nil { log.Fatalf("Err: send msg: %s", err) }
	log.Printf("Response from server: %d", mySession.Temp)

	// Get Details
	message2 := pb.UserRequest { UserId: "name", SessionToken:mySession.Temp, }
	response2, err := c.GetDetails(context.Background(), &message2)
	if err != nil { log.Fatalf("Err: send msg2: %s", err) }
	log.Printf("Response from server: %s", response2.Details)

	// Invalid GetDetails
	message3 := pb.UserRequest { UserId: "nam", SessionToken:mySession.Temp, }
	response3, err := c.GetDetails(context.Background(), &message3)
	if err != nil { log.Fatalf("Err: send msg3: %s", err) }
	log.Printf("Response from server: %s", response3.Details)

	// Get Risk
	message4 := *mySession
	response4, err := c.GetRisk(context.Background(), &message4)
	if err != nil { log.Fatalf("Err: send msg4: %s", err) }
	log.Printf("Response from server: %d", response4.Score)

	// Send lifestyle
	message5 := pb.LifestyleRequest { Message:"123"}
	response5, err := c.SendLifestyle(context.Background(), &message5)
	if err != nil { log.Fatalf("Err: send msg5: %s", err) }
	log.Printf("Response from server: %d", response5.Success)

}
