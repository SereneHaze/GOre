SERVER_BINARY=gore_server
CLIENT_BINARY=gore_client

build:
	go build -o ${SERVER_BINARY} server/server.go
	go build -o ${CLIENT_BINARY} client/client.go

clean:
	rm ${SERVER_BINARY}
	rm ${CLIENT_BINARY}

