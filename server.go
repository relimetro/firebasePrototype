package main

import(
	"log" // for loggin
	"sync" // for mutex

	// Grpc
	"context"
	"net"
	pb "example/proto_example/protoOut"
	"google.golang.org/grpc"

	// firebase - firestore
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)



///////////////////////////////////////////////////////////////
/// stuff and things

type Session_Tokens_Type struct {
	data [65535]string
	mu sync.RWMutex
	idx int64
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

	// Mutex
	Session_Tokens.mu.Lock()
	defer Session_Tokens.mu.Unlock()

	if ValidateLogin(x.UserName, x.PlaintextPassword) {
		Session_Tokens.data[Session_Tokens.idx] = x.UserName
	} else { Session_Tokens.data[Session_Tokens.idx] = "__invalid__"}

	returnVal := Session_Tokens.idx
	Session_Tokens.idx = Session_Tokens.idx +1

	return &pb.SessionToken{
		Temp: returnVal,
	}, nil
}



// GetDetails (UserRequest -> UserDetails)
func (s *server) GetDetails(ctx context.Context, x *pb.UserRequest) (*pb.UserDetails, error) {
	log.Printf("GetDetails: %s, %s", x.SessionToken, x.UserId)

	// Mutex
	Session_Tokens.mu.RLock()
	defer Session_Tokens.mu.RUnlock()

	idx := x.SessionToken

	log.Printf("Idx: %d, UserId: %s, Session: %s\n\n",idx,x.UserId, Session_Tokens.data[idx])
	if Session_Tokens.data[idx] == x.UserId {
		return &pb.UserDetails{
			Details:"some details" }, nil
	}
	return &pb.UserDetails{
		Details:"invalid" }, nil

}



// GetRisk (SessionToken -> RiskScore)
func (s *server) GetRisk(ctx context.Context, x *pb.SessionToken) (*pb.RiskScore, error) {
	log.Printf("GetRisk: %d", x.Temp)
	return &pb.RiskScore{
		Score: 0, }, nil
}






///////////////////////////////////////////////////////////////
/// firebase

type FireDB struct {
	*db.Client }
var fireDB FireDB

func (db *FireDB) Connect() error {
	home := "/home/cathal/notes/MTU/project/goTest/"
	ctx := context.Background()
	opt := option.WithCredentialsFile(home+"firebase.json")
	dbUrl := "https://test-fd3ea.firebaseio.com"
	config := &firebase.Config{DatabaseURL: dbUrl}
	app, err := firebase.NewApp(ctx,config,opt)
	if err != nil { log.Fatalf("error init app: %v", err ); return err }
	client, err := app.Database(ctx)
	if err != nil { log.Fatalf("error init db: %v", err ); return err }
	db.Client = client
	return nil
}

func FirebaseDB() *FireDB {
	return &fireDB
}

type UserDB struct {
	BIN string
	username string
	password string
	riskFactor int

}

type Store struct { *FireDB }
func NewStore() *Store {
	d := FirebaseDB()
	return &Store{ FireDB:d, }
}

// Create a new user
func (s *Store) Create(b *UserDB) error {
	println("20")
	if err := s.NewRef("users/"+b.BIN).Set(context.Background(), b); err != nil {
		log.Fatalf("Create: %v", err); return err }
	println("21")
	return nil
}

func (s *Store) Delete(b *UserDB) error {
	return s.NewRef("users/" + b.BIN).Delete(context.Background())
}

func (s *Store) GetByBin(b string) (*UserDB, error) {
	println("10")
	bin := &UserDB{}
	println("11")
	if err := s.NewRef("bins/"+b).Get(context.Background(), bin); err != nil {
		log.Fatalf("error getBin db: %v", err ); return nil, err }
	println("12")
	if bin.BIN == "" {
		return nil, nil }
	println("13")
	return bin, nil
}

func (s *Store) Update(b string, m map[string]interface{}) error {
	return s.NewRef("users/"+b).Update(context.Background(), m)
}





///////////////////////////////////////////////////////////////
/// Main

func main() {
	lis, err := net.Listen("tcp", ":9000");
	if err != nil { log.Fatalf("failed to listen: %v", err) }


	// connect
	dbErr := FirebaseDB().Connect()
	if dbErr != nil { log.Fatalf("error init db: %v", dbErr ) }
	// store := NewStore()

	// A new BIN creation
	// err = store.Create(&UserDB{
	// 	BIN: "1234",
	// 	username: "cathal",
	// 	password: "bob",
	// 	riskFactor: 1,
	// })


	// store := NewStore()
	// bin, getErr := store.GetByBin("users")
	// if getErr != nil { log.Fatalf("error store: %v", dbErr ) }
	// print("username: %s",bin)


	ctx := context.Background()
	home := "/home/cathal/notes/MTU/project/goTest/"
	sa := option.WithCredentialsFile(home+"firebase.json")
	app, err := firebase.NewApp(ctx,nil,sa)
	if err != nil { log.Fatalln(err)}
	client, err := app.Firestore(ctx)
	if err != nil { log.Fatalln(err)}
	defer client.Close()

	_, _, err2 := client.Collection("users").Add(ctx, map[string]interface{}{
		"username":"ada",
		"password":"12345",
		"riskFactor":19,
	})
	if err2 != nil { log.Fatalf("Failed adding %v", err2)}
	




	grpcServer := grpc.NewServer()
	pb.RegisterServServer(grpcServer, &server{})
	log.Printf("Ready!!")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve") }
}
