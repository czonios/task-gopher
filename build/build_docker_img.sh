#!/bin/bash
cd ..
docker build . -t go-containerized:latest -f ./build/dockerfile
echo "---------------"
echo " "
echo "Docker images found:"
docker image ls | grep go-containerized