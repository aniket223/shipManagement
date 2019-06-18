package main
import(
	"context"
	"sync"
	"log"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "github.com/aniket223/shipManagement/consignment-service/proto/consignment"
)

const(
	port = ":50051"
)

type repository interface{
	Create(*pb.Consignment)(*pb.Consignment,error) 
}

type Repository struct{
	mu sync.RWMutex
	consignments []*pb.Consignment
}

//function to create a new consignment
func (repo *Repository)Create(consignment *pb.Consignment)(*pb.Consignment,error){
	//no one else can create a consignment so we use mutex locks
	repo.mu.Lock()
	//append the consignment to already existing consignments
	updated:=append(repo.consignments,consignment)
	//update the consignments
	repo.consignments = updated
	//unlock
	repo.mu.Unlock()
	return consignment,nil
}

//service should implement all of the methods to satisfy the service
//we defined in our protobuf definition.
type service struct{
	repo repository
}

func(s *service)createConsignment(ctx context.Context,req *pb.Consignment)(*pb.Response,error){
	consignment,err:=s.repo.Create(req)
	if err!=nil{
		return nil,err
	}
	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func main(){
	repo := &Repository{}

	// Set-up our gRPC server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	pb.RegisterShippingServiceServer(s, &service{repo})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

