## GOre RAT
GOre is a remote access trojan that uses build-time obfuscation do avoid detection and hinder dissasembly. There are 3 parts to the C2 framework: Malicious implant,
server, and admin operator client. The client forwards commands to the server, which is added to a "work" channel, that each implant reaches out for to recieve said
commands and execute them on a victim machine.

## Implant
Each implant has a compile-time UUID, as well as a self-signed certificate that it uses to secure communications to the server. If you don't use the garble wrapper, 
the chances of detection become almost certain. The certificate and UUID are present in the binary, meaning a determined security analyst could reverse engineer the 
network traffic and the GRPC connections.

## Server
The server acts as a middle-man to the client and implants. It adds commands to a queue that is accessed by all implnats in a "reverse-backdoor" fashion; each implant
reaches out to the server for work, instead of listening on a port passivly for commands.

## Client
The client is a simple and fast implimentation to send commands to the server and recieve an output. This is not built with security in mind, as it is assumed that
you will not be using this on the same network that the implant is on. This is something to fix in the future,

## Build guide
The first thing todo is to use the `cert_gen.sh` shell script to automatically generate self-signed certs used by the framework. You are free to use your own certificates,
you must place them in the correct folders prior to generating the server and implants. The next step is to use the command `make custom` to generate server and client
binaries. You can go into the makefile and manually change the data for this, as it is somewhat assumed that this is a "build-once-run-many" program. After this is done,
you are able to produce implants eitehr manually, with the `implant_build.sh` shell script, or automatically with your own custom shell script or wrapper for my shell script.


## Disclaimer
This tool is a passion project of mine, as well as an educational tool. It is not intended for use to access computers or systems that you do not have prior permission
to access. This program is presented AS IS, with no warranty. 

Basically, don't be an idiot. This is a hacktool, if you use it illegally you will very likley face leagal reprucssions. Use the power for good, and never for evil.
