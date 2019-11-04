#!/bin/bash

kill -9 $(pgrep localstack)
docker rm -f localstack_main
