#!/bin/bash
#this is not printing as nicely as I want it to.
echo "
  ________              _________                __   
 /  _____/  ____   ____ \_   ___ \  ____________/  |_ 
/   \  ____/ __ \ /    \/    \  \/_/ __ \_  __ \   __\\
\    \_\  \  ___/|   |  \     \___\  ___/|  | \/|  |  
 \______  /\___  >___|  /\______  /\___  >__|   |__|  
        \/     \/     \/        \/     \/             

"

echo "[?] This script is intedned to be used in the home directory of the GOre project, and assumes existing folder names."

#Source: https://www.handracs.info/blog/grpcmtlsgo/
#generate certificate authority, used to have -noout.
openssl ecparam -name prime256v1 -genkey -noout -out ./tls/CAkey.key
echo "[+] Generated CAkey"
#generate private key
openssl req -x509 -new -nodes -key ./tls/CAkey.key -subj "/CN=Synapse/C=US" -days 730 -out ./server/CAcert.pem
#copy CAcert to implant
cp ./server/CAcert.pem ./implant/CAcert.pem
echo "[+] CA cert and key generated; CAcer.pem copied to implant directory."
#create server key
openssl ecparam -name prime256v1 -genkey -noout -out ./server/server.key
echo "[+] Generated server .key file."
#create config file for key singing
#I believe the fields under [ dn ] are the actually usefull and important ones in both scripts. 
cat > ./server/csr.conf <<EOF
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
echo "[+] Generated config file."
#create new Certificate Signing Request
openssl req -new -key ./server/server.key -out ./server/server.csr -config ./server/csr.conf
openssl x509 -req -in ./server/server.csr -CA ./server/CAcert.pem -CAkey ./tls/CAkey.key -CAcreateserial -out ./server/server.pem -days 90 -extfile ./server/csr.conf -extensions req_ext
echo "[+] Key has been signed."
