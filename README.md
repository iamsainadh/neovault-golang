# Neovault

## Multi-Tenant Secret Management Platform

Neovault is a secure, multi-tenant secret management platform built using Go microservices. It provides encrypted storage, organization-level isolation, secret versioning, audit tracking, and a foundation for automated secret rotation.

Neovault uses envelope encryption to protect sensitive data and separates cryptographic responsibilities between Root Master Key (RMK), Key Encryption Keys (KEK), and Data Encryption Keys (DEK).

# Why Neovault?

Traditional applications often store sensitive credentials directly in configuration files or environment variables, making secret management difficult at scale.

Neovault provides a centralized platform for securely storing, encrypting, versioning, and managing secrets across multiple organizations while maintaining strict tenant isolation.

# Features

## Current Features

- Multi-tenant secret management
- Organization-based isolation
- Secure encrypted secret storage
- Envelope encryption architecture
- Secret version management
- PostgreSQL credential password update support
- Manual secret rotation support
- Secret lifecycle management
  - Active
  - Archived
  - Deleted
- Audit logging through Kafka
- gRPC-based internal service communication
- REST API gateway
- Authentication and authorization service
- Kubernetes deployment support
- Prometheus metrics
- Grafana dashboards

---

# Planned Features

- Automatic secret rotation engine
- Amazon RDS credential rotation
- MySQL credential rotation
- MongoDB credential rotation
- Payment gateway key rotation scripts
- Notification integration for rotation success/failure
- Additional secret providers

---

# Technology Stack

| Component                | Technology           |
| ------------------------ | -------------------- |
| Language                 | Go 1.26              |
| API Gateway              | Gin                  |
| Internal Communication   | gRPC                 |
| Database                 | PostgreSQL           |
| Password Hashing         | Argon2               |
| Scheduler                | robfig/cron/v3       |
| Event Streaming          | Kafka                |
| Kafka Client             | franz-go/pkg/kgo     |
| Development              | Air (Live Reload)    |
| Containerization         | Docker               |
| Local Orchestration      | Docker Compose       |
| Production Orchestration | Kubernetes           |
| Monitoring               | Prometheus + Grafana |

---

# Architecture Overview

Neovault follows a microservice architecture.

```
                    External Clients
                          |
                          |
                   API Gateway
                       Gin
                          |
    ------------------------------------------------
    |              |              |                |
   Auth          Secret          Vault          Notification
    |              |              |                |
    ------------------------------------------------
                          |
                        gRPC
                          |
                    PostgreSQL
                          |
                       Kafka
                          |
                     Audit Service
```

---

# Multi-Tenant Architecture

Every organization is logically isolated.

Each tenant owns:

- Users
- Organization Members
- KEKs
- Vault Keys
- Secrets
- Secret Versions

Secret access is always scoped using Organization ID, ensuring complete tenant isolation.

# Microservices

## API Gateway

Technology:

- Gin

Responsibilities:

- Public REST API exposure
- Request validation
- Authentication middleware
- Organization context injection
- Routing requests to internal services through gRPC

---

## Auth Service

Responsibilities:

- User management
- Organization management
- Organization membership
- Authentication
- Password hashing using Argon2
- Authorization management

---

## Secret Service

Responsibilities:

- Secret creation
- Secret retrieval
- Secret versioning
- Encryption workflow
- Secret lifecycle management

---

## Vault Service

Responsibilities:

- Organization key management
- KEK storage
- Key versioning
- Encryption metadata management

---

## Scheduler Service

Responsibilities:

- Executes scheduled jobs
- Checks rotation policies
- Triggers secret rotation workflows

Uses:

- robfig/cron/v3

---

## Audit Service

Responsibilities:

- Stores security audit events
- Consumes Kafka events
- Tracks secret operations

---

## Notification Service

Responsibilities:

- Sends rotation status notifications
- Handles success and failure events

---

# Encryption Architecture

Neovault uses envelope encryption.

The encryption hierarchy:

                RMK
      (Root Master Key - ENV)
                   |
                   |
          Encrypts Organization KEK
                   |
                   |
          Organization KEK
                   |
                   |
          Encrypts Secret DEK
                   |
                   |
             Secret DEK
                   |
                   |
          Encrypts Secret Payload

---

# Key Management Model

## Root Master Key (RMK)

- Stored securely as an environment variable
- Not stored in PostgreSQL
- Used only for protecting organization-level keys
- No AWS KMS dependency

---

## Key Encryption Key (KEK)

Each organization receives an independent KEK.

Properties:

- Generated per organization
- Stored encrypted
- Managed through Vault Service
- Supports future key rotation

Database mapping:
VAULT_KEYS

contains:

- Organization ID
- Key version
- Encrypted KEK
- Algorithm
- Key status

Example:

```
Organization A

RMK
|
└── KEK Version 1
|
└── Secret DEKs

Organization B

RMK
|
└── KEK Version 1
|
└── Secret DEKs
```

---

## Data Encryption Key (DEK)

DEKs are generated for secret encryption.

Flow:

```
Secret Value
        |
        |
Generate DEK
        |
        |
Encrypt Secret Payload
        |
        |
Encrypt DEK using Organization KEK
```

---

# Secret Storage Model

## Secret Table

Stores secret metadata.

Supported secret types:

```
generic
api_key
password
token
```

Secret lifecycle:

```
active
archived
deleted
```

---

# Database Architecture

PostgreSQL is used as the primary database.

Database diagram:
<img src=docs/db_architecture.png>

Main domains:

```
Authentication
|
├── users
├── organizations
└── organization_members

Secrets
|
├── secrets
└── secret_versions

Vault
|
└── vault_keys

Audit
|
└── audit_logs
```

---

# Local Development

## Requirements

Install:

- Go 1.26
- Docker
- Docker Compose
- PostgreSQL
- Kafka

---

# Running With Docker Compose

Build services:

```
docker compose build
```

Start environment:

```
docker compose up
```

Stop environment:

```
docker compose down
```

---

# Kubernetes Deployment

Neovault provides Kubernetes manifests.

Includes:

- Services
- Deployments
- ConfigMaps
- RBAC configuration

Deployment example:

```
kubectl apply -f k8s/
```

# Monitoring

Neovault includes built-in Prometheus metrics for HTTP traffic, application performance, and Go runtime statistics. Metrics are exposed for scraping by Prometheus and can be visualized using Grafana dashboards.

## Current Metrics

### HTTP Metrics

The API Gateway exposes request-level metrics through middleware.

| Metric                              | Type      | Description                                                          |
| ----------------------------------- | --------- | -------------------------------------------------------------------- |
| `neovault_http_requests_total`      | Counter   | Total HTTP requests grouped by method, endpoint, and response status |
| `neovault_request_duration_seconds` | Histogram | Request latency grouped by method and endpoint                       |
| `neovault_http_errors_total`        | Counter   | Total HTTP error responses grouped by method, endpoint, and status   |

Example labels:

```text
method="POST"
endpoint="/api/v1/secrets"
status="200"
```

---

### Go Runtime Metrics

Neovault registers the standard Prometheus Go collector.

Collected metrics include:

- Goroutines
- Memory allocation
- Heap usage
- Garbage collection statistics
- Thread count
- CPU scheduling statistics
- Go version information

Examples:

```text
go_goroutines
go_memstats_alloc_bytes
go_memstats_heap_alloc_bytes
go_gc_duration_seconds
```

---

### Process Metrics

Neovault also exposes process-level metrics.

Examples include:

- CPU usage
- Memory usage
- Open file descriptors
- Virtual memory
- Resident memory
- Process start time

Examples:

```text
process_cpu_seconds_total
process_resident_memory_bytes
process_virtual_memory_bytes
process_open_fds
process_start_time_seconds
```

---

## Prometheus Integration

Prometheus scrapes the metrics endpoint exposed by the API Gateway.

```
                API Gateway
                     │
                     │
                /metrics
                     │
                     │
               Prometheus
                     │
                     │
                Grafana
```

---

## Future Metrics

Neovault will continue expanding observability with service-specific metrics.

### Secret Metrics

- Secrets created
- Secrets updated
- Secrets deleted
- Secret retrieval count
- Secret version count
- Secret access latency

### Vault Metrics

- KEKs created
- KEK rotations
- Active key versions
- Encryption operations
- Decryption operations

### Rotation Metrics

- Rotation jobs executed
- Successful rotations
- Failed rotations
- Rotation duration
- Secrets awaiting rotation

### Authentication Metrics

- Successful logins
- Failed logins
- JWT validation failures
- Organization creation count

### Database Metrics

- Query latency
- Active connections
- Connection pool utilization
- Failed queries
- Transaction duration

### gRPC Metrics

- RPC requests
- RPC latency
- Error rate
- Active connections

## Future SDKs

Official SDKs will simplify secure secret retrieval for backend applications.

Planned SDKs

- Go SDK (Go modules)
- Node.js SDK (npm)
- Java SDK (Maven Central)

The SDKs will provide:

- Authentication
- Organization-aware clients
- Secret retrieval
- Secret version access
- Secret creation
- Automatic retries
- Connection pooling

# Project Outcomes

Neovault demonstrates the design and implementation of a secure, scalable, and cloud-native secret management platform built using modern Go microservices.

Key outcomes include:

- Designed a multi-tenant architecture with strict organization-level data isolation.
- Implemented envelope encryption using RMK, KEK, and DEK without relying on external KMS providers.
- Built independent microservices communicating through gRPC for improved scalability and maintainability.
- Exposed REST APIs through a Gin-based API Gateway while keeping internal services private.
- Implemented versioned secret storage to support secure updates and future key rotation.
- Established an event-driven audit logging pipeline using Apache Kafka.
- Integrated Prometheus metrics and Grafana dashboards for observability and operational monitoring.
- Containerized services using Docker and Docker Compose for local development.
- Prepared Kubernetes deployment manifests with RBAC configuration for production environments.
- Designed the platform to support future secret providers and automated credential rotation workflows.
