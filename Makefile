#names of the generated files
SERVER_BINARY=gore_server
CLIENT_BINARY=gore_client
#for values for the injectable variables. The only thing that really changes is the server IP, which can be a DNS name as well. 
SERVER_IP=ix-dev.cs.uoregon.edu
SERVER_PORT=5000
OPERATOR_PORT=9090
#debug build with localhost values; doesn't connect to the internet.
debug:
	go build -o ${SERVER_BINARY} -ldflags=" -X 'main.server_ip=localhost' -X 'main.server_port=5000' -X 'main.operator_port=9090'" server/server.go
	go build -o ${CLIENT_BINARY} -ldflags=" -X 'main.server_ip=localhost' -X 'main.server_port=9090'" client/client.go

#custom build with specified values; this will typically connect to the internet so be EXTRA sure you know what you're doing 
build:
	go build -o ${SERVER_BINARY} -ldflags=" -X 'main.server_ip=${SERVER_IP}' -X 'main.server_port=${SERVER_PORT}' -X 'main.operator_port=${OPERATOR_PORT}'" server/server.go
	go build -o ${CLIENT_BINARY} -ldflags="-X 'main.server_ip=${SERVER_IP}' -X 'main.server_port=${OPERATOR_PORT}'" client/client.go
#removal of compiled server/operator binaries. Removal of implants is more manual.
clean:
	rm ${SERVER_BINARY}
	rm ${CLIENT_BINARY}
