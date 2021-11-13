package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/CasperAntonPoulsen/DisysExercise2/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8080"
)

var (
	client pb.MutualExclusionClient
	wait   *sync.WaitGroup
	id     = flag.Int("id", 1, "id of the node")
)

func init() {
	wait = &sync.WaitGroup{}
}

func requestToken(rqst *pb.Request) error {
	var streamerror error
	stream, err := client.RequestToken(context.Background(), rqst)
	if err != nil {
		return fmt.Errorf("connection has fauled: %v", err)
	}

	//recieve grant token

	for {
		_, err := stream.Recv()
		if err != nil {
			streamerror = fmt.Errorf("error recieving grant token: %v", err)
			break
		} else {
			log.Print("Grant token recieved, accessing critical section")
			// access critical section
			client.AccesCritical(context.Background(), rqst.User)
			log.Print("Finished, sending release token")
			// then release
			client.ReleaseToken(context.Background(), &pb.Release{User: rqst.User})
			break
		}
	}

	return streamerror
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect to service : %v", err)
	}

	client = pb.NewMutualExclusionClient(conn)
	user := &pb.User{Userid: int32(*id)}

	for {
		time.Sleep(4 * time.Second)
		requestToken(&pb.Request{User: user})

	}

}
