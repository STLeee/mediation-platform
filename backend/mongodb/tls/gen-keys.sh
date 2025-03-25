#/bin/bash

# Generate CA, Server and Client keys
openssl genrsa -out mongodb-test-ca.key 4096
openssl req -new -x509 -days 1826 -key mongodb-test-ca.key -out mongodb-test-ca.crt -config openssl-test-ca.cnf
openssl genrsa -out mongodb-test-ia.key 4096
openssl req -new -key mongodb-test-ia.key -out mongodb-test-ia.csr -config openssl-test-ca.cnf
openssl x509 -sha256 -req -days 730 -in mongodb-test-ia.csr -CA mongodb-test-ca.crt -CAkey mongodb-test-ca.key -set_serial 01 -out mongodb-test-ia.crt -extfile openssl-test-ca.cnf -extensions v3_ca
cat mongodb-test-ia.crt mongodb-test-ca.crt > test-ca.pem

openssl genrsa -out mongodb-test-server1.key 4096
openssl req -new -key mongodb-test-server1.key -out mongodb-test-server1.csr -config openssl-test-server.cnf
openssl x509 -sha256 -req -days 365 -in mongodb-test-server1.csr -CA mongodb-test-ia.crt -CAkey mongodb-test-ia.key -CAcreateserial -out mongodb-test-server1.crt -extfile openssl-test-server.cnf -extensions v3_req
cat mongodb-test-server1.crt mongodb-test-server1.key > test-server1.pem

openssl genrsa -out mongodb-test-client.key 4096
openssl req -new -key mongodb-test-client.key -out mongodb-test-client.csr -config openssl-test-client.cnf
openssl x509 -sha256 -req -days 365 -in mongodb-test-client.csr -CA mongodb-test-ia.crt -CAkey mongodb-test-ia.key -CAcreateserial -out mongodb-test-client.crt -extfile openssl-test-client.cnf -extensions v3_req
cat mongodb-test-client.crt mongodb-test-client.key > test-client.pem

# Insert to python package `certifi`
cat test-ca.pem >> $(python -m certifi)

# Generate replica set key
openssl rand -base64 756 | tr -d '\n' > mongodb-test-keyfile
chmod 600 mongodb-test-keyfile
