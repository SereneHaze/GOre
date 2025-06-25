#source: 
https://www.handracs.info/blog/grpcmtlsgo/

#generate a certificate authority with eliptical curve prime 256v1 with name cakey.key
1.) 
`openssl ecparam -name prime256v1 -genkey -noout -out cakey.key`

#create private key with common name "TestCA' and country "MY"
2.) 
`openssl req -x509 -new -nodes -key cakey.key -subj "/CN=TestCA/C=MY" -days 730 -out cacert.pem` 

#create a server key
3.) 
`openssl ecparam -name prime256v1 -genkey -noout -out server.key`

#create a config file for signing the keys
4.)
`cat > csr.conf <<EOF
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

EOF`

#create CSR(Certificate signing rqequest) 
5.)
`openssl req -new -key server.key -out server.csr -config csr.conf`

#generate SSL cert, will be self-signed.
6.)
`openssl x509 -req -in server.csr -CA cacert.pem -CAkey cakey.key -CAcreateserial -out server.pem -days 90 -extfile csr.conf -extensions req_ext`

#generate client certificate key
7.)
`openssl ecparam -name prime256v1 -genkey -noout -out client.key`

#create csr config file for client
8.)
`cat > csrclient.conf <<EOF
[ req ]
default_bits = 256
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[ dn ]
C = US
CN = client

[ req_ext ]
keyUsage = keyEncipherment
extendedKeyUsage = clientAuth

EOF`

#generate CSR for client
9.)
`openssl req -new -key client.key -out client.csr -config csrclient.conf`

#sign the CSR for the client and issue a cert for the client
10.)
`openssl x509 -req -in client.csr -CA cacert.pem -CAkey cakey.key -CAcreateserial -out client.pem -days 90 -extfile csrclient.conf -extensions req_ext`

--The certificates are now generated--

#moving certificates to the proper directories
now that the certificates are generated and singed by our self-signed authority, you need to move them to the proper directories.
The implant and server will compile with the TLS certificates as embeded files, meaning that they will need to present in each one's directory to properly
copmile. the basic directory structure should look like:

- implant/
  - implant.go
  - client.pem
  - client.key
  - cacert.pem

- server/
  - server.go
  - server.pem
  - server.key
  - cacert.pem

As of now there is no shell script to make this an easy and automated process, but one will be added in future builds.


## Some Extra Notes

Replace the CN and the IP of the `alt_names` section in the server's csr.conf file with the DNS and the IP of the server. The client CN seems to not really matter.

I have the defaults set to localhost.
