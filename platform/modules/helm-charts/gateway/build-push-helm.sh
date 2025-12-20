#!/bin/bash
set -e

CHART_DIR="./"
CHART_NAME="gateway"
CHART_VERSION="0.1.1"
REGISTRY="localhost:5000"
helm template $CHART_DIR
helm package $CHART_DIR --version $CHART_VERSION
helm push ${CHART_NAME}-${CHART_VERSION}.tgz oci://$REGISTRY
echo "Helm chart packaged and pushed to $REGISTRY/$CHART_NAME"
rm ${CHART_NAME}-${CHART_VERSION}.tgz
read -n 1 -s -r -p "Press any key to exit..."