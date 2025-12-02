#!/bin/bash

set -e

./1-test.sh 
./2-build.sh
./3-test_bin.sh
./4-build_image.sh
./5-integration_test.sh
./6-deploy_local.sh