## GOre RAT
GOre is a personal project to make a RAT/C2 in Golang with some C scripts if I want. Basically, it will be an attempt to learn go through a super fun
project. GOre is a portmanteu of "Go rev" or "go reverse". It is also based on the RAT from Blackhat GO, and indeed in it's most basic form it is a 1:1 copy of that program's source code. It is not a FULL copy, as there
are some dependencies that I needed to add to get this thing to compile properly. 

The RAT cosists of three distinct pieces: a server for communicating with a (singular for the moment) implant, an Operator client to issue commands to the server (to rely them to an implant) and the implant whihc operates on the
victims device. The implant is the device that needs to be loaded onto a victim, the server exists on some external service, and the client can be run locally on a machine, through a VPN if you'd like.

## client
the client communicates with the server, which acts as a VIA to the implant. This binary is for the malware operators.
## CLIENT TODO
- encrypt client to server communications
- add better command support, IE run the binary and input commands like a shell
- allow for multiple operators to communicate without cross polination of output/input

## implant 
the implant is "implanted" onto the device to facilitate c2 commuication to the server. The implant only talks to the server, not to ay client operators directly. It runs OS commads directly on the victim as the victim. Does
not support pipes, as it is a bit raw on how it hadles command parsing.
## IMPLANT TODO

## server

## Reccomended build methods
use the given shell script with `./shell_scripts/implant_build.sh <domain/ip> <port> <executable name>`.
Shell scripts have help text which can be invoked with the -h flag, or by not supplying the correct number of arguments.
There is also a makefile for the client and server with settings for running on localhost, or on custom networking settings, depending on how you wnat to set up the framework. Implants are not generate by the Makefile, it is recomended to use the shell scripts.

## Running
The client and server executable is statically compiled, and can be run by first starting the server on the server machine (or localhost) with  `./name_of_executable` when built. Then the implant can be run on the victim machine with `./implant_name`. Then, commands can be executed with `./client_name <command> <UUID>`, where command is the OS command you want to run on linux (`ls` for example), and UUID is the UUID of the target implant. Multiple implants can be run simultaneously. UUID's are printed when the code is copmiled with the shell script or makefile. They are also printed as debug information on the server, so if you don't have it just establish a connection with the server to see the UUID. You will need Go to build the executables, and you will need to install the dependencies with `go mod download` in the project directory.
