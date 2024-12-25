##GOre RAT##
GOre is a personal project to make a RAT/C2 in Golang with some C scripts if I want. Basically, it will be an attempt to learn go through a super fun
project. GOre is a portmanteu of "Go rev" or "go reverse". It is also based on the RAT from Blackhat GO, and indeed in it's most basic form it is a 1:1 copy of that program's source code. It is not a FULL copy, as there
are some dependencies that I needed to add to get this thing to compile properly. 

The RAT cosists of three distinct pieces: a server for communicating with a (singular for the moment) implant, an Operator client to issue commands to the server (to rely them to an implant) and the implant whihc operates on the
victims device. The implant is the device that needs to be loaded onto a victim, the server exists on some external service, and the client can be run locally on a machine, through a VPN if you'd like.

##client##
the client communicates with the server, which acts as a VIA to the implant. This binary is for the malware operators.
##CLIENT TODO##
- encrypt client to server communications
- add better command support, IE run the binary and input commands like a shell
- allow for multiple operators to communicate without cross polination of output/input

##implant## 
the implant is "implanted" onto the device to facilitate c2 commuication to the server. The implant only talks to the server, not to ay client operators directly. It runs OS commads directly on the victim as the victim. Does
not support pipes, as it is a bit raw on how it hadles command parsing.
##IMPLANT TODO##

##server##
##SERVER TODO##

##building an executable from source:##
go build -o <exec-name> -ldflags="-w -s -buildid=" -trimpath <path/to/source_file>
-maybe add "-H=windowsgui" after trimpath to make the implant a windowless application
