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

# Load the Docker image into the kind cluster
kind load docker-image go-api:latest -n test-cluster

# Set kube context to the kind cluster
kubectl cluster-info --context kind-test-cluster

# Install ingress controller and metrics server
kubectl apply -f https://kind.sigs.k8s.io/examples/ingress/deploy-ingress-nginx.yaml
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
kubectl -n kube-system patch deployment metrics-server \
  --type='json' \
  -p='[{"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"--kubelet-insecure-tls"}]'


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

kubectl wait --namespace default --for=condition=Ready pod --selector=app=go-api --timeout=120s
start cmd "/k kubectl logs -f -l app=go-api --all-containers=true"
kubectl wait --namespace default --for=condition=Ready pod --selector=app=ollama --timeout=120s
start cmd "/k kubectl logs -f -l app=ollama --all-containers=true"
start cmd "/k kubectl port-forward svc/ingress-nginx-controller -n ingress-nginx 8080:80"
