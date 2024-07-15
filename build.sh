#!/bin/bash

set -e

if [ "$1" == "build" ]; then
    echo "Checking Dependencies..."
    go mod tidy

    echo "Building server..."
    go build -o server ./cmd/main.go

    echo "Building client..."
    go build -o client ./client/main.go
    echo "Build complete."
elif [ "$1" == "run-server" ]; then
    echo "Running the server..."
    ./server
elif [ "$1" == "run-client" ]; then
    if [ "$#" -lt 3 ]; then
        echo "Usage: ./build.sh run-client [-auth=<true|false>] [-tls=<true|false>] [-ca-cert=<path/to/client.pem>] <address> <port>"
        exit 1
    fi
    echo "Running the client..."
    cd client

    client_cmd="./client"

    for arg in "${@:2}"; do
        client_cmd+=" $arg"
    done

    echo "Executing command: $client_cmd"
    eval "$client_cmd"
else
    echo "Invalid command. Usage: $0 [build|run-server|run-client]"
    exit 1
fi
