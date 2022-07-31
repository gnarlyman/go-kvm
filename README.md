# go-kvm

openssl genrsa -out ca.key 2048

openssl req -new -key ca.key -x509 -days 3650 -out ca.crt -config config 

openssl genrsa -out server.key 2048

openssl req -new -nodes -key server.key -out server.csr -config config

openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt

openssl genrsa -out client.key 2048

openssl req -new -nodes -key client.key -out client.csr -config config

openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt