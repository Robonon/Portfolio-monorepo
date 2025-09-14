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

kind load docker-image go-api:latest -n test-cluster
kubectl cluster-info --context kind-test-cluster
kubectl apply -f https://kind.sigs.k8s.io/examples/ingress/deploy-ingress-nginx.yaml

# Wait for ingress-nginx controller to be ready
echo "Waiting for ingress-nginx-controller to be ready..."
until kubectl wait --namespace ingress-nginx \
  --for=condition=Ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s; do
  echo "Ingress-nginx-controller is not ready yet..."
  sleep 5
done
echo "Ingress-nginx-controller is ready."
helm upgrade --install go-api ./chart --kube-context kind-test-cluster

# If everything succeeds, you can optionally unset cleanup
trap - ERR

kubectl port-forward svc/ingress-nginx-controller -n ingress-nginx 8080:80
