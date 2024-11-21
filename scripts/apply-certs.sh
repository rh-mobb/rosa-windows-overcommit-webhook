#!/usr/bin/env bash

# Variables
CA_KEY="tmp/ca.key"
CA_CERT="tmp/ca.crt"
SERVER_KEY="tmp/server.key"
SERVER_CSR="tmp/server.csr"
SERVER_CERT="tmp/server.crt"
CONFIG_FILE="tmp/openssl.cnf"

# Create Kubernetes secrets
echo "Creating Kubernetes secrets..."

# Secret for the CA certificate
kubectl create secret generic webhook-ca \
    --namespace=windows-overcommit-webhook \
    --from-file=ca.crt=$CA_CERT \
    --dry-run=client -o yaml | kubectl apply -f -

# Secret for the server certificate and key
kubectl create secret generic webhook-certs \
    --namespace=windows-overcommit-webhook \
    --from-file=tls.crt=$SERVER_CERT \
    --from-file=tls.key=$SERVER_KEY \
    --from-file=ca.crt=$CA_CERT \
    --dry-run=client -o yaml | kubectl apply -f -

echo "Kubernetes secrets created:"
echo "- webhook-ca (contains ca.crt)"
echo "- webhook-certs (contains tls.crt, tls.key, and ca.crt)"