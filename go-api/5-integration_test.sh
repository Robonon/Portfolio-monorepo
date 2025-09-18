#!/bin/bash

PORT=9999
PASS=true
echo "Starting integration tests on port $PORT"
echo "Creating and running Docker container..."
docker run -e PORT=$PORT -d -p $PORT:$PORT go-api:latest

curl -s -X GET http://localhost:$PORT -w "\n"

r=$(curl -s -X GET http://localhost:$PORT/calculations/max \
    -H "Content-Type: application/json" \
    -d '{"numbers": [1, 2, 3, 4, 5]}')

if [[ $r != *"5"* ]]; then
    echo "Test failed: Expected max to be 5, got $r"
    PASS=false
fi
echo "Test passed: Max calculation $r"

r=$(curl -s -X GET http://localhost:$PORT/calculations/reverse \
    -H "Content-Type: application/json" \
    -d '{"numbers": [1, 2, 3, 4, 5]}')
if [[ $r != *"[5,4,3,2,1]"* ]]; then
    echo "Test failed: Expected reverse to be [5,4,3,2,1], got $r"
    PASS=false
fi
echo "Test passed: Reverse calculation $r"

r=$(curl -s -X GET http://localhost:$PORT/calculations/countUnique \
    -H "Content-Type: application/json" \
    -d '{"numbers": [1, 2, 2, 5, 5]}')
if [[ $r != *"3"* ]]; then
    echo "Test failed: Expected countUnique to be 3, got $r"
    PASS=false
fi
echo "Test passed: CountUnique calculation $r"

r=$(curl -s -X GET http://localhost:$PORT/calculations/sum \
    -H "Content-Type: application/json" \
    -d '{"numbers": [1, 2, 3, 4, 5]}')

if [[ $r != *"15"* ]]; then
    echo "Test failed: Expected sum to be 15, got $r"
    PASS=false
fi
echo "Test passed: Sum calculation $r"

echo "Stopping and removing Docker container..."
docker ps | grep go-api | awk '{print $1}' | xargs docker stop | xargs docker rm

if [ "$PASS" = true ]; then
    echo "All tests passed!"
    exit 0
else
    echo "Some tests failed."
    exit 1
fi