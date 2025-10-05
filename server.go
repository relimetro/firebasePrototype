package main

import(
	"log" // for loggin
	"sync" // for mutex
	"google.golang.org/api/iterator"

	// Grpc
	"context"
	"net"
	pb "example/proto_example/protoOut"
	"google.golang.org/grpc"

	// firebase
	firebase "firebase.google.com/go"
	firestore "cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)



///////////////////////////////////////////////////////////////
/// stuff and things


// Firebase
type UserRecord struct {
	Username string
	RiskFactor int32
	// Password string
}



// Auth/Tokens
type Session_Tokens_Type struct {
	data [65535]string
	mu sync.RWMutex
	idx int64
	// todo free list, more info stored about token not just username
}
var Session_Tokens = Session_Tokens_Type{ idx: 0}


func ValidateLogin(username string, password string) bool {
	return true }



///////////////////////////////////////////////////////////////
/// GRPC
type server struct{
	pb.UnimplementedServServer
}



// Login (UserLogin -> SessionToken)
func (s *server) Login(ctx context.Context, x *pb.UserLogin) (*pb.SessionToken, error) {
	log.Printf("login: %s, %s\n\n",x.UserName, x.PlaintextPassword)

	// Mutex Write Lock
	Session_Tokens.mu.Lock()
	defer Session_Tokens.mu.Unlock()

	// Validate Login and assign token
	if ValidateLogin(x.UserName, x.PlaintextPassword) {
		Session_Tokens.data[Session_Tokens.idx] = x.UserName
	} else { Session_Tokens.data[Session_Tokens.idx] = "__invalid__"}

	returnVal := Session_Tokens.idx
	Session_Tokens.idx = Session_Tokens.idx +1

	// return session token
	return &pb.SessionToken{
		Temp: returnVal,
	}, nil
}



// GetDetails (UserRequest -> UserDetails)
func (s *server) GetDetails(ctx context.Context, x *pb.UserRequest) (*pb.UserDetails, error) {
	log.Printf("GetDetails: %s, %s", x.SessionToken, x.UserId)

	// Mutex Read Lock
	Session_Tokens.mu.RLock()
	defer Session_Tokens.mu.RUnlock()

	idx := x.SessionToken

	// check if user can access required data
	log.Printf("Idx: %d, UserId: %s, Session: %s\n\n",idx,x.UserId, Session_Tokens.data[idx])
	if Session_Tokens.data[idx] == x.UserId {
		return &pb.UserDetails{
			Details:"some details" }, nil
	}
	// unauthorized access response
	return &pb.UserDetails{
		Details:"invalid" }, nil
}



// GetRisk (SessionToken -> RiskScore)
func (s *server) GetRisk(ctx context.Context, x *pb.SessionToken) (*pb.RiskScore, error) {

	// Mutex Read Lock
	Session_Tokens.mu.RLock()
	username := Session_Tokens.data[x.Temp] // todo, validate valid session Token (not out of bounds etc)
	Session_Tokens.mu.RUnlock()

	log.Printf("GetRisk: Session: %d username: %s", x.Temp, username)

	// find
	iter := client.Collection("users").Documents(context.Background())
	for { // todo: probably a way to do this on server
		// iterate
		doc, err := iter.Next()
		if err == iterator.Done { break }
		if err != nil { log.Fatalf("failed to iterate:\n%v",err)}

		// get data of record
		var docData UserRecord
		if err := doc.DataTo(&docData); err != nil {
			log.Fatalf("err2") }

		// check if target user
		if docData.Username == username {
			log.Printf("%d",docData.RiskFactor)
			return &pb.RiskScore{ Score: docData.RiskFactor, }, nil
		}
	}

	// Dummy Response
	return &pb.RiskScore{ Score: 0, }, nil
}






///////////////////////////////////////////////////////////////
/// global firebase client, initialized at startup

var client *firestore.Client
func firebaseInit(){
	FBctx := context.Background()
	sa := option.WithCredentialsFile("./firebase.json")
	app, err := firebase.NewApp(FBctx,nil,sa)
	if err != nil { log.Fatalf("Firebase: failed to create app:\n%v",err)}
	var err2 error
	client, err2 = app.Firestore(FBctx)
	if err2 != nil { log.Fatalf("Firebase: failed to access store:\n%v",err)}
}





///////////////////////////////////////////////////////////////
/// Main

func main() {

	// firebase settup
	firebaseInit()
	defer client.Close()

	// FBctx := context.Background()
	// test firebase add
	// _, _, err2 := client.Collection("users").Add(FBctx, map[string]interface{}{
	// 	"username":"ada",
	// 	"password":"12345",
	// 	"riskFactor":19,
	// })
	// if err2 != nil { log.Fatalf("Failed adding\n%v", err2)}
	





	// grpc connection
	lis, err := net.Listen("tcp", ":9000");
	if err != nil { log.Fatalf("GRPC: failed to listen:\n%v", err) }

	// serv GRPC
	grpcServer := grpc.NewServer()
	pb.RegisterServServer(grpcServer, &server{})
	log.Printf("Ready!! >:0")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("GRPC: Failed to serve:\n%v",err) }

}
