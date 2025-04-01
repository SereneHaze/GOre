#variables for the certificates
CA_COUNTRY=US
CN=Synapse

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
custom:
	#create SSL certificates for the server
	echo "[:] generating SSL certificates for GOre server"
	openssl req -x509 -name prime256v1 -genkey -out tls/cakey.key
	echo "[+] Generated Certificate Authority key"
	openssl req -x509 -new -nodes -key cakey.key -subj "/CN=Synapse/C=US" -days 730 -out tls/cacert.pem
	echo "[+] Generated Certificate Authority private key"
	openssl ecparam -name prime256v1 -genkey -out tls/server.key
	echo "[+] Generated server key"
	cat > tls/csr.conf <<EOF                                                                                                                                                                                                                      
	[ req ]                                                                                                                                                                                                                                    
	default_bits = 256                                                                                                                                                                                                                         
	prompt = no
	default_md = sha256
	req_extensions = req_ext
	distinguished_name = dn
	[ dn ]
	C = US
	CN = localhost
	[ req_ext ]
	keyUsage = keyEncipherment, dataEncipherment
	extendedKeyUsage = serverAuth
	subjectAltName = @alt_names
	[ alt_names ]
	DNS.1 = localhost
	IP.1 = 127.0.0.1
	EOF
	echo "[+] built csr config file"
	openssl req -new -key server.key -out server.csr -config tls/csr.conf		
	echo "[+] Generated certificate signing request"
	openssl x509 -req -in tls/server.csr -CA cacert.pem -CAkey tls/cakey.key -CAcreateserial -out tls/server.pem -days 90 -extfile tls/csr.conf -extensions req_ext
	echo "[+] Generated server certificate"
	echo "[:] building server and client"
	go build -o ${SERVER_BINARY} -ldflags=" -X 'main.server_ip=${SERVER_IP}' -X 'main.server_port=${SERVER_PORT}' -X 'main.operator_port=${OPERATOR_PORT}'" server/server.go
	go build -o ${CLIENT_BINARY} -ldflags="-X 'main.server_ip=${SERVER_IP}' -X 'main.server_port=${OPERATOR_PORT}'" client/client.go
#removal of compiled server/operator binaries. Removal of implants is more manual.
clean:
	rm ${SERVER_BINARY}
	rm ${CLIENT_BINARY}
