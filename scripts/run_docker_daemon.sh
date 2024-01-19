#!/bin/bash
TAG="czonios/task-gopher/server"
PORT=8334
docker run --name task-gopher-server --env-file ../.env --restart unless-stopped -p $PORT:$PORT -d -t $TAG:latest
mkdir -p ../logs
docker logs --follow task-gopher-server >> ../logs/server_logs.log