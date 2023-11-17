#!/bin/bash
sudo docker run --name task-gopher-server --env-file ./.env --restart unless-stopped -p 8334:8334 -d -t go-containerized:latest
docker logs task-gopher-server >> ./server_logs.log