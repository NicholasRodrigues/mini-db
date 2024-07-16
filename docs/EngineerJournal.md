
---

# Engineer's Journal: Building Mini DB

Welcome to the detailed engineer's journal for the development of Mini DB, a minimalistic key-value store designed to handle TCP-based command interactions. This document aims to provide a comprehensive view of the technical decisions, considerations, and implementation details that shaped Mini DB. It will serve as an insightful resource for understanding the entire development process from inception to deployment.

## Table of Contents

1. [Project Overview](#project-overview)
2. [Design Considerations](#design-considerations)
3. [Implementation Details](#implementation-details)
4. [Security Enhancements](#security-enhancements)
5. [Testing and Validation](#testing-and-validation)
6. [Monitoring and Alerting](#monitoring-and-alerting)
7. [Future Enhancements](#future-enhancements)
8. [Conclusion](#conclusion)

## Project Overview

**Mini DB** is crafted to operate as a simple yet robust key-value store that accepts commands over a TCP connection. The primary functionalities include setting values to keys and retrieving these values through lookup commands. A significant feature of this database is its ability to persist data, ensuring durability across server restarts.

## Design Considerations

### Choosing the Technology Stack

- **Go (Golang)**: Selected for its excellent support for concurrent operations and built-in TCP server capabilities. Its performance in networked applications and binary efficiency were key factors.
- **File System for Persistence**: Simple file-based storage was chosen for durability to ensure ease of implementation and reliability.

### Protocol Design

- **TCP Interface**: TCP was chosen for reliable, ordered, and error-checked delivery of streams of bytes. This protocol is well-suited for command-response patterns.
- **Command Parsing**: Commands are delimited by line breaks (`\n`), with simple parsing logic to interpret `SET` and `LOOKUP` operations.

## Implementation Details

### Core Functionality

- **Command Handling**: Implemented a TCP server that listens for incoming connections. Each connection is handled in a separate goroutine to manage multiple clients simultaneously.
- **Data Storage**:
  - In-memory dictionary for quick access.
  - Periodic and event-driven serialization of data to a file to ensure durability.

### Error Handling

- Robust error handling was implemented to log and manage commands that do not adhere to the expected syntax, enhancing the system's reliability and usability.

### Persistence

- **Parquet Files**: Chosen for its efficient data compression and encoding, making it suitable for large-scale data processing.

## Security Enhancements

### Secure Connections

- **TLS Integration**: Added optional TLS support to encrypt TCP connections, protecting data integrity and privacy over the network.
- **Authentication**: Implemented token-based authentication to control access, allowing only authorized clients to execute commands.
  - Quick note: the auth token-based system is only experimental and not secure for production use.

## Testing and Validation

### Unit and Integration Testing

- Developed a suite of tests to cover functional and edge cases, ensuring the database handles expected inputs and gracefully manages unexpected or erroneous conditions.
- Utilized Go's testing framework and `testify` library for assertions.

### Performance Testing

- Conducted load testing to validate the system's performance under stress and to benchmark its response times and memory usage. (Go script can be found on `monitoring/load.go`)
- Monitored resource utilization and latency to ensure the system meets performance requirements.

## Monitoring and Alerting

### Monitoring Setup with Prometheus and Grafana

1. **Prometheus Setup**

   - **Prometheus Configuration**: Created a `prometheus.yml` file with appropriate scrape configurations.
   - **Docker Compose for Monitoring**: Set up a Docker Compose configuration to run Prometheus and Grafana.

2. **Grafana Setup**

   - **Access Grafana**: Logged into Grafana and added Prometheus as a data source.
   - **Create Dashboards**: Created dashboards to visualize metrics like total requests, request rates, CPU usage, and memory usage.
   - **Customize Dashboard Colors**: Customized colors for better visual distinction and readability.

### Example Queries

1. **Total Requests**:
   ```prometheus
   sum(requests_total) by (method)
   ```

2. **Request Rate**:
   ```prometheus
   rate(requests_total[1m])
   ```

3. **CPU Usage**:
   ```prometheus
   sum(rate(process_cpu_seconds_total[1m])) by (instance)
   ```

4. **Memory Usage**:
   ```prometheus
   process_resident_memory_bytes
   ```

### Structured Logging and Health Checks

- Implemented structured logging using Logrus for better visibility into system behavior.
- Incorporated health checks to ensure system availability and reliability.

## Future Enhancements

### Scalability

- Plans to introduce distributed processing and storage capabilities to scale beyond a single server instance.

### API Extensions

- Considering adding more sophisticated query capabilities and support for different data types beyond simple strings.

### Advanced Security

- Implement more advanced security features such as OAuth or JWT-based authentication.

### User Interface

- Develop a web-based user interface for easier interaction with Mini DB.

## Conclusion

Thank you for your interest in Mini DB. This project aims to provide a simple yet robust solution for managing key-value pairs with durability and security in mind. With features like TLS support, authentication, and comprehensive monitoring using Prometheus and Grafana, Mini DB is well-suited for various applications. We encourage you to explore, use, and contribute to the project, ensuring its continuous improvement and adaptability to meet evolving needs.

For a detailed guide on how to install, configure, and use Mini DB, please refer to the [README.md](../README.md) document. This document includes all necessary commands and configurations, along with examples of how to interact with the database securely using TLS and authentication.

---
