#!/bin/bash
set -e

REGISTRY_IMAGE="localhost:5000/tenant:latest"

function pause() {
  read -n 1 -s -r -p "Press any key to exit..."
}

trap pause EXIT

docker build -t $REGISTRY_IMAGE .
docker push $REGISTRY_IMAGE
echo "Build, test, and push complete."
curl -s http://localhost:5000/v2/_catalog