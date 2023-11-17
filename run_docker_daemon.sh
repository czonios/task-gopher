#!/bin/bash
sudo docker run --name task-gopher-server --env-file ./.env --restart unless-stopped -p 8334:8334 go-containerized:latest