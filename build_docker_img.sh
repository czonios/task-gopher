#!/bin/bash
docker build . -t go-containerized:latest
echo "---------------"
echo " "
echo "Docker images found:"
docker image ls | grep go-containerized