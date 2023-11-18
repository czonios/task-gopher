#!/bin/bash
docker run --name task-gopher-server --env-file ../.env -p 8334:8334 -t go-containerized:latest