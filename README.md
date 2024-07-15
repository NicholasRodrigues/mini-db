![CI Workflow Badge](https://github.com/NicholasRodrigues/mini-db/actions/workflows/ci.yml/badge.svg)
![CD Workflow Badge](https://github.com/NicholasRodrigues/mini-db/actions/workflows/cd.yml/badge.svg)

# Mini DB

Mini DB is a simple key-value store that communicates through a TCP interface, designed to handle basic operations such as setting and looking up key-value pairs. It utilizes Parquet files for data persistence, ensuring data durability across server restarts.

[Jump to Safe Usage with TLS and Authentication](#safe-usage-with-tls-and-authentication)

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

Configuration is managed via a YAML file on (`config/config.yaml`). Here’s an example:

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
  auth_enabled: false
  auth_token: "your-secure-auth-token"
```

## Building and Running

### Using build.sh Script

Ensure Go is installed, and the script is executable:

```sh
chmod +x build.sh
./build.sh build
./build.sh run-server
./build.sh run-client -auth=true -tls=true -ca-cert=path/to/client.pem localhost 8080
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
./client -auth=true -tls=true -ca-cert=path/to/client.pem <address> <port>
```

### Example Usage

```sh
./client -auth=true -tls=true -ca-cert=path/to/client.pem localhost 8080
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

## Safe Usage with TLS and Authentication

For enhanced security, Mini DB supports both TLS for encrypted connections and token-based authentication to control access. Below are guidelines on how to configure and use these features securely.

### Enabling TLS

1. **Generate TLS Certificates:** Follow the instructions above to generate `server.crt` and `server.key`.
2. **Configure the Client:** Ensure the client is set up to use the CA certificate through `-ca-cert` flag pointing to `server.crt`.

### Generating TLS Certificates

1. **Install OpenSSL**:
   If you don't have OpenSSL installed on your system, you can install it using your package manager. For example, on Ubuntu, you can install it with:
   ```bash
   sudo apt-get install openssl
   ```

2. **Generate a Private Key and a Self-Signed Certificate**:
   Use the following commands to generate your private key and a self-signed certificate:

   ```bash
   openssl genrsa -out server.key 2048
   openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365
   ```
   During the certificate creation process (`req -new -x509`), you will be prompted to enter details such as country, state, and organization. These details are used to fill the subject field of the certificate. For local testing, you can fill these with any values.

3. **Place the Generated Files Appropriately**:
   Once generated, place `server.crt` and `server.key` in a directory accessible by your server application, such as within a `config` folder or directly in the directory from which you run your server.

### Updating Your Application Configuration

Make sure your server application is configured to use these TLS files. Update the configuration settings in your `config.yaml` or wherever your server application expects them:

   ```yaml
    tls: true
    tls_cert_file: "path/to/server.crt"
    tls_key_file: "path/to/server.key"
   ```

### Running the Server with TLS

Start your server normally using your script or directly. If your server setup script or application is configured to output logs, watch for any TLS-related errors.

### Testing TLS

To test if TLS is working:

1. **Using `curl` or a similar tool**: You can test the HTTPS connection using `curl`. Replace `localhost` and `8080` with your actual server address and port.

   ```bash
   curl https://localhost:8080 --cacert path/to/server.crt
   ```

   This command tells `curl` to use your self-signed certificate as the CA certificate. If everything is set up correctly, you should not see any SSL errors, and you should be able to communicate securely with your server.

2. **Using a Web Browser**: Access your server using a web browser. You might need to import your `server.crt` into the browser or accept a security exception because the browser will warn you about the self-signed certificate.

### Handling Errors

If you encounter TLS-related errors during testing, check the following:

- Ensure the paths to `server.crt` and `server.key` are correctly specified in your server's configuration.
- Verify that the server is configured to listen over HTTPS and not HTTP.
- Check for detailed error messages in your server's logs to diagnose issues related to certificate loading or TLS handshake failures.

### Using Authentication

1. **Configure an Auth Token:** Set `auth_enabled` to `true` in the server configuration and specify an `auth_token`.
2. **Client Usage:** When starting the client, use `-auth=true` and provide the token in your commands as shown in the examples.

Follow these guidelines to ensure that your interactions with Mini DB are secure and your data is protected against unauthorized access.

---