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

kubectl cluster-info --context kind-test-cluster
