package main

import(
	"net"
	"context"
	pb "example/proto_example/protoOut"
	"google.golang.org/grpc"
	"log"
)

type server struct{
	pb.UnimplementedServServer
}



func (s *server) Login(ctx context.Context, x *pb.UserLogin) (*pb.SessionToken, error) {
	log.Printf("login: %s, %s",x.UserName, x.PlaintextPassword)
	return &pb.SessionToken{
		Temp: x.UserName,
	}, nil
}

func (s *server) GetDetails(ctx context.Context, x *pb.UserRequest) (*pb.UserDetails, error) {
	return &pb.UserDetails{
		Details:"none" }, nil
}

func (s *server) GetRisk(ctx context.Context, x *pb.SessionToken) (*pb.RiskScore, error) {
	return &pb.RiskScore{
		Score: 0, }, nil
}




func main() {
	lis, err := net.Listen("tcp", ":9000");
	if err != nil { log.Fatalf("failed to listen") }

	grpcServer := grpc.NewServer()
	pb.RegisterServServer(grpcServer, &server{})
	log.Printf("Ready!!")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve") }
	log.Printf("Ready!!")
}
