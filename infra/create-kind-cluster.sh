#!/bin/bash

set -e

CLEANUP_CLUSTER=true

cleanup() {
    if $CLEANUP_CLUSTER; then
        echo "Cleaning up: Deleting kind cluster..."
        kind delete cluster -n test-cluster || true
    fi
}
trap cleanup ERR

# Check if kind cluster already exists
if kind get clusters | grep -q "^test-cluster$"; then
    echo "Kind cluster 'test-cluster' already exists. Continuing..."
    CLEANUP_CLUSTER=false
else
    kind create cluster -n test-cluster
    CLEANUP_CLUSTER=true
fi

# Install Gateway API CRDs
kubectl apply --server-side -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.4.1/standard-install.yaml

kubectl cluster-info --context kind-test-cluster

docker compose -f ./docker-compose.yaml up -d
