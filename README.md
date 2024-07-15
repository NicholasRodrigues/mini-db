


![CI Workflow Badge](https://github.com/NicholasRodrigues/mini-db/actions/workflows/ci.yml/badge.svg)
![CD Workflow Badge](https://github.com/NicholasRodrigues/mini-db/actions/workflows/cd.yml/badge.svg)
---

# Mini DB

Mini DB is a simple key-value store that communicates through a TCP interface, designed to handle basic operations such as setting and looking up key-value pairs. It utilizes Parquet files for data persistence, ensuring data durability across server restarts.

## Quick Start

1. **Clone the repository:**
   ```sh
   git clone <repository_url>
   cd mini-db
   ```

2. **Build and run using Docker:**
   ```sh
   docker-compose up --build
   ```

3. **Or build and run locally:**
   ```sh
   ./build.sh build
   ./build.sh run-server
   ```

## Features

- **SET and LOOKUP Commands:** Store and retrieve key-value pairs efficiently.
- **Data Persistence:** Utilizes Parquet files for reliable data storage.
- **TLS Support:** Optional TLS to secure connections.
- **Authentication:** Supports token-based authentication for enhanced security.
- **Structured Logging:** Uses Logrus for clear and structured logging.

## Configuration

Configuration is managed via a YAML file (`config.yaml`). Here’s an example:

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

## Building and Running

### Using build.sh Script

Ensure Go is installed, and the script is executable:

```sh
chmod +x build.sh
./build.sh build
./build.sh run-server
./build.sh run-client localhost 8080 [--tls]
```

### Docker Instructions

Ensure Docker is installed and run:

```sh
docker-compose up --build
```

## Client Application

Interact with the Mini DB server using the client application:

```sh
cd client
go build -o client main.go
./client <address> <port> [--tls]
```

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

## Server Management

**Testing:** Run tests with `go test ./...`.

**Stopping the Server:** Use `Ctrl+C` in the terminal or `docker-compose down` for Docker users.

## Troubleshooting

- Check all file paths and environment variables.
- Ensure Go and Docker are properly installed and configured.
- Review script outputs and server logs for error messages.

---