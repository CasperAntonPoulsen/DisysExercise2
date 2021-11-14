# DisysExercise2

To compile the proto file use the command:

make compile-proto

To run the server:

go run server/server.go

To run a client/node:

go run client/client.go -id {int}


Currently clients are just set to send a request token with a two second delay if they're not currently in the queue. Accesing critical just logs it on the server and sleeps for a set amount of time.