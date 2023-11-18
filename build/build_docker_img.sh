#!/bin/bash
TAG="czonios/task-gopher/server"
cd ..
docker build . -t $TAG:latest -f ./build/Dockerfile
echo "---------------"
echo " "
echo "Docker images found:"
docker image ls | grep $TAG