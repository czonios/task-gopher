#!/bin/bash
sudo docker run --name task-gopher-server --env-file ./.env --restart unless-stopped go-containerized:latest