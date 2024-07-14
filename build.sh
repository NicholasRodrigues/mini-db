#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Define build and run modes
if [ "$1" == "build" ]; then
    echo "Building server..."
    go build -o server ./cmd/main.go

    echo "Building client..."
    go build -o client ./client/main.go
    echo "Build complete."
elif [ "$1" == "run-server" ]; then
    echo "Running the server..."
    ./server
elif [ "$1" == "run-client" ]; then
    if [ "$#" -ne 3 ]; then
        echo "Usage: ./build.sh run-client <address> <port> [--tls]"
        exit 1
    fi
    echo "Running the client..."
    cd client
    ./client "$2" "$3" "$4"
else
    echo "Invalid command. Usage: $0 [build|run-server|run-client]"
    exit 1
fi
