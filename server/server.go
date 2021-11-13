package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/CasperAntonPoulsen/DisysExercise2/proto"
	"google.golang.org/grpc"
)

const port = ":8080"

type Request struct {
	user   *pb.User
	stream pb.MutualExclusion_RequestTokenServer
}

type Server struct {
	pb.UnimplementedMutualExclusionServer
	RequestQueue chan Request
	Release      chan *pb.Release
	error        chan error
}

func (s *Server) RequestToken(rqst *pb.Request, stream pb.MutualExclusion_RequestTokenServer) error {

	request := Request{user: rqst.User, stream: stream}
	log.Printf("Request token recieved from: %v", rqst.User.Userid)

	s.RequestQueue <- request

	log.Printf("Request token added to queue from: %v", rqst.User.Userid)
	return <-s.error
}

func (s *Server) ReleaseToken(ctx context.Context, release *pb.Release) (*pb.Empty, error) {
	log.Printf("Recieved release token from: %v", release.User.Userid)
	s.Release <- &pb.Release{User: release.User}

	return &pb.Empty{}, nil
}

func (s *Server) AccesCritical(ctx context.Context, user *pb.User) (*pb.Empty, error) {
	log.Printf("User: %v is accesing the critcal section", user.Userid)
	time.Sleep(2 * time.Second)
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
	requestqueue := make(chan Request)
	releasequeue := make(chan *pb.Release)
	server := Server{
		RequestQueue: requestqueue,
		Release:      releasequeue,
	}

	log.Println("Starting server at port ", port)

	pb.RegisterMutualExclusionServer(grpcServer, &server)
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
	grpcServer.Serve(listener)

}
