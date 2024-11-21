#!/usr/bin/env bash

set -e

# Variables
CA_KEY="tmp/ca.key"
CA_CERT="tmp/ca.crt"
SERVER_KEY="tmp/server.key"
SERVER_CSR="tmp/server.csr"
SERVER_CERT="tmp/server.crt"
CONFIG_FILE="tmp/openssl.cnf"

mkdir -p tmp

# Generate the CA private key and certificate
echo "Generating Certificate Authority (CA)..."
openssl genrsa -out $CA_KEY 4096
openssl req -x509 -new -nodes -key $CA_KEY -sha256 -days 3650 -out $CA_CERT -subj "/CN=Kubernetes CA"

# Generate the server private key
echo "Generating server private key..."
openssl genrsa -out $SERVER_KEY 4096

# Create the OpenSSL config file for the certificate request
cat <<EOF > $CONFIG_FILE
[ req ]
default_bits       = 4096
distinguished_name = req_distinguished_name
x509_extensions    = v3_ca
req_extensions     = req_ext
prompt             = no

[ req_distinguished_name ]
C  = US
ST = Nebraska
L  = Lincoln
O  = Dustin Scott
OU = Webhook
CN = windows-overcommit-webhook

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = localhost
DNS.2 = windows-overcommit-webhook
DNS.3 = windows-overcommit-webhook.windows-overcommit-webhook
DNS.4 = windows-overcommit-webhook.windows-overcommit-webhook.svc
DNS.5 = windows-overcommit-webhook.windows-overcommit-webhook.svc.cluster.local
IP.1  = 127.0.0.1
EOF

# Generate the server certificate signing request
echo "Generating server certificate signing request..."
openssl req -new -key $SERVER_KEY -out $SERVER_CSR -config $CONFIG_FILE

# Generate the server certificate signed by the CA
echo "Generating server certificate signed by the CA..."
openssl x509 -req -in $SERVER_CSR -CA $CA_CERT -CAkey $CA_KEY -CAcreateserial \
    -out $SERVER_CERT -days 3650 -extensions req_ext -extfile $CONFIG_FILE

echo "Certificates generated:"
echo "CA Key: $CA_KEY"
echo "CA Certificate: $CA_CERT"
echo "Server Key: $SERVER_KEY"
echo "Server Certificate: $SERVER_CERT"
