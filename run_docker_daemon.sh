#!/bin/bash
sudo docker run --name task-gopher-server --env-file ./.env --restart unless-stopped -p 8334:8334 -d go-containerized:latest
docker logs task-gopher-server 2>> ./server_logs.log