#!/bin/bash

docker run --rm --net=tp0_testing_net\
 -v ./config/netcat:/script \
 --env-file ./config/netcat/env.txt\
 netcat_test