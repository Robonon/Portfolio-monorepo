#!/bin/bash

MODULE_NAME=$1

if [ -z "$MODULE_NAME" ]; then
  echo "Usage: $1 <module-name>"
  exit 1
fi

mkdir -p "./$MODULE_NAME"
touch "./$MODULE_NAME/main.tf"
touch "./$MODULE_NAME/variables.tf"
touch "./$MODULE_NAME/outputs.tf"
touch "./$MODULE_NAME/README.md"

echo "Module '$MODULE_NAME' created successfully."