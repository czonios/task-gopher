#!/bin/bash
sudo docker run --name task-gopher-server --env-file ../.env --restart unless-stopped -p 8334:8334 -d -t go-containerized:latest
mkdir -p ../logs
sudo docker logs --follow task-gopher-server >> ../logs/server_logs.log