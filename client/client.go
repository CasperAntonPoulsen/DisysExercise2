package main

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/CasperAntonPoulsen/DisysExercise2/proto"
)

const (
	address = "localhost:8080"
)

var client pb.MutualExclusionClient
var wait *sync.WaitGroup
var mutex sync.Mutex

func init() {
	wait = &sync.WaitGroup{}
}

type User struct {
	id    int32
	time  int32
	state string
}

func evaluaterequest(rqst *pb.User, user *pb.User) bool {
	if user.State == "HELD" || user.State == "WANTED" && (user.Time < rqst.Time || user.Userid < rqst.Userid) {
		return false
	}
	return true
}

func reply(rply *pb.Reply) {

}

func recieverequests(user *pb.User) error {
	var streamerror error
	stream, err := client.RecieveRequests(context.Background(), user)
	if err != nil {
		return fmt.Errorf("connection has fauled: %v", err)
	}
	wait.Add(1)
	//recieve requests
	go func(str pb.MutualExclusion_RecieveRequestsClient) {
		defer wait.Done()

		request, err := str.Recv()
		if err != nil {
			streamerror = fmt.Errorf("error recieving request: %v", err)
		}

		if evaluaterequest(request, user) {
			reply(&pb.Reply{&pb.User{
				User: request,
			}})
		}

	}(stream)
	return streamerror
}

func main() {

}
