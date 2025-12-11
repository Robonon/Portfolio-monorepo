#!/bin/bash

set -e

CLUSTER_NAME="local-cluster"
CLEANUP_CLUSTER=true

cleanup() {
    if $CLEANUP_CLUSTER; then
        echo "Cleaning up: Deleting kind cluster..."
        kind delete cluster -n $CLUSTER_NAME || true
    fi
}
trap cleanup ERR

# Check if kind cluster already exists
if kind get clusters | grep -q "^$CLUSTER_NAME$"; then
    echo "Kind cluster '$CLUSTER_NAME' already exists. Continuing..."
    CLEANUP_CLUSTER=false
else
    kind create cluster -n $CLUSTER_NAME --config kind-config.yaml
    CLEANUP_CLUSTER=true
fi

# Install Gateway API CRDs
kubectl apply --server-side -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.4.1/standard-install.yaml

kubectl cluster-info --context kind-$CLUSTER_NAME

docker compose -f ./docker-compose.yaml up -d
