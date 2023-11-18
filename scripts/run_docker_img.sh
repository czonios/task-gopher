#!/bin/bash
TAG="czonios/task-gopher/server"
PORT=8334
docker run --name task-gopher-server --env-file ../.env -p $PORT:$PORT -t $TAG:latest