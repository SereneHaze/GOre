#!/bin/bash
NEW_UUID=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
echo $NEW_UUID
go build -o ./implant -ldflags="-X 'main.uuid=$NEW_UUID'-w -s -buildid=" -trimpath ./implant/implant.go
