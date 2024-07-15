#!/bin/bash

set -e

echo "Running go mod tidy..."
go mod tidy

DOCKER_NETWORK="mini-db-network"
if ! docker network ls | grep -q $DOCKER_NETWORK; then
  echo "Creating Docker network $DOCKER_NETWORK..."
  docker network create $DOCKER_NETWORK
else
  echo "Docker network $DOCKER_NETWORK already exists."
fi

echo "Running Docker Compose for the server..."
docker-compose up --build -d

echo "Running Docker Compose for monitoring..."
docker-compose -f monitoring/docker-compose.monitoring.yml up --build -d

echo "Mini DB server and monitoring are up and running."

echo "Running containers:"
docker ps
