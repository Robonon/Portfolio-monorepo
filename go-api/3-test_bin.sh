#!/bin/bash

go test -c ./... -o bin/test/
./bin/test/* -test.v