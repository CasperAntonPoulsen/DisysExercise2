package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/CasperAntonPoulsen/DisysExercise2/proto"
	"google.golang.org/grpc"
)

const port = ":8080"

var mutex sync.Mutex

type Request struct {
	user   *pb.User
	stream pb.MutualExclusion_RequestTokenServer
}

type Server struct {
	pb.UnimplementedMutualExclusionServer
	RequestQueue (chan Request)
	Release      (chan *pb.Release)
	error        chan error
}

func (s *Server) RequestToken(rqst *pb.Request, stream pb.MutualExclusion_RequestTokenServer) error {
	mutex.Lock()

	go func() { s.RequestQueue <- Request{user: rqst.User, stream: stream} }()
	log.Printf("Request token recieved from: %v", rqst.User.Userid)
	mutex.Unlock()

	return <-s.error
}

func (s *Server) ReleaseToken(ctx context.Context, release *pb.Release) (*pb.Empty, error) {

	go func() { s.Release <- &pb.Release{User: release.User} }()

	return &pb.Empty{}, nil
}

func (s *Server) AccesCritical(ctx context.Context, user *pb.User) (*pb.Empty, error) {
	log.Printf("User: %v is accesing the critcal section", user.Userid)

	return &pb.Empty{}, nil
}

func GrantToken(rqst Request) error {
	err := rqst.stream.Send(&pb.Grant{User: rqst.user})

	return err
}

func main() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Error, couldn't create the server %v", err)
	}

	server := Server{}

	go func() {
		for {
			fmt.Print("Checking requests \n")
			rqst := <-server.RequestQueue

			fmt.Print("Recieved request \n")
			err := GrantToken(rqst)
			if err != nil {
				log.Fatalf("Failed to send grant token: %v", err)
			}

			release := <-server.Release
			log.Printf("Release token recieved from: %v", release.User.Userid)

		}
	}()

	log.Println("Starting server at port ", port)

	pb.RegisterMutualExclusionServer(grpcServer, &server)
	grpcServer.Serve(listener)

}
