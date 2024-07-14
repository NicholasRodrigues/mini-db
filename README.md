# Mini DB

Mini DB is a simple key-value store that communicates through a TCP interface. It supports basic operations such as setting and looking up key-value pairs, and it persists data using Parquet files.

## Features

- **SET and LOOKUP commands**: Store and retrieve key-value pairs.
- **Data persistence**: Data is persisted using Parquet files to ensure durability across server restarts.
- **Optional TLS support**: Secure your connections with TLS.
- **Token-based authentication**: Secure your database with an authentication token.
- **Structured logging**: Logging is handled using Logrus, providing structured log outputs.

## Server

### Configuration

The server configuration is managed using a YAML file. Below is an example configuration file (`config.yaml`):

```yaml
server:
  port: "8080"
  tls: false
  tls_cert_file: "config/server.crt"
  tls_key_file: "config/server.key"

storage:
  file_path: "data/db.parquet"

logging:
  level: "info"

security:
  auth_enabled: true
  auth_token: "your-secure-auth-token"
```

### Running the Server

1. Clone the repository and navigate to the project directory:

    ```sh
    git clone <repository_url>
    cd mini-db
    ```

2. Build and run the server application:

    ```sh
    go build -o server cmd/main.go
    ./server
    ```

3. Using Docker:

    ```sh
    docker-compose up --build
    ```

### Testing the Server

To run tests:

```sh
go test ./...
```

### Stopping the Server

To stop the server gracefully, send a termination signal (e.g., `Ctrl+C` in the terminal running the server).

## Client

The client application allows you to interact with the Mini DB server via the command line.

### Running the Client

1. Navigate to the client directory:

    ```sh
    cd mini-db/client
    ```

2. Build and run the client application:

    ```sh
    go build -o client main.go
    ./client <address> <port> [--tls]
    ```

    Replace `<address>` and `<port>` with the appropriate server address and port. Use the `--tls` flag if the server is using TLS.

### Commands

- **SET**: Stores a key-value pair. Optionally, an auth token can be included.
  
    ```sh
    SET [<auth_token>] my_key my_value
    ```

- **LOOKUP**: Retrieves the value associated with a key. Optionally, an auth token can be included.
  
    ```sh
    LOOKUP [<auth_token>] my_key
    ```

- **EXIT**: Exits the client application.

### Example Usage

```sh
./client localhost 8080
> SET your-secure-auth-token my_key my_value
OK
> LOOKUP your-secure-auth-token my_key
my_value
> EXIT
Exiting...
```

### Build and Run Using build.sh

The `build.sh` script simplifies the process of building and running the Mini DB server and client applications. Here are the detailed instructions on how to use this script:

#### Requirements

Before running the script, ensure you have the following prerequisites installed:

- **Go**: The server and client are built using Go. Make sure you have Go installed on your system. You can download it from [Go's official website](https://golang.org/dl/).

#### Using the build.sh Script

The `build.sh` script provides a convenient way to build and run the server and client. Follow these steps:

1. **Grant Execution Permissions**:
   
   First, ensure that the `build.sh` script is executable. You can set the execution permission with the following command:
   
   ```sh
   chmod +x build.sh
   ```

2. **Build the Server and Client**:
   
   To compile both the server and client applications, use the `build` command:
   
   ```sh
   ./build.sh build
   ```

   This command will generate executable files named `server` and `client` in your project directory.

3. **Run the Server**:
   
   After building, you can start the server using the `run-server` command:
   
   ```sh
   ./build.sh run-server
   ```

   Ensure that the server's configuration in `config.yaml` is correct before starting the server.

4. **Run the Client**:
   
   To start the client, use the `run-client` command followed by the server's address and port. Include the `--tls` option if TLS is enabled:
   
   ```sh
   ./build.sh run-client <address> <port> [--tls]
   ```

   Replace `<address>` and `<port>` with the appropriate values. For example:
   
   ```sh
   ./build.sh run-client localhost 8080
   ```

#### Troubleshooting

If you encounter any issues while using the `build.sh` script, check the following:

- Ensure all file paths and environment variables are set correctly in the script and your system.
- Verify that all dependencies, particularly Go, are installed and configured correctly.
- Check the script output for any error messages that can help diagnose issues.
