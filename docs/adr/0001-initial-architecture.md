# ADR 0001: Initial Architecture

- Status: Proposed
- Date: 2026-11-06
- Decision owners: Project maintainers

## Context

The project aims to demonstrate production-oriented Go backend engineering, real-time communication, concurrency, and distributed systems.

The architecture must support multiple WebSocket gateway instances while remaining manageable for a single developer to build and maintain.

## Decisions

### 1. Modular Monolith Instead of Microservices

We will start with a modular monolith.

Reasons:

- Easier local development.
- Less infrastructure complexity.
- Shared domain logic.
- Clear module boundaries.
- Can evolve into separate services when justified.

### 2. Separate API Server and WebSocket Gateway Executables

The REST API and WebSocket Gateway will run as separate processes but share application and domain packages.

Reasons:

- Different connection and scaling characteristics.
- Independent deployment and scaling later.
- Clear separation of responsibilities.

### 3. PostgreSQL as the Source of Truth

All persistent business data will be stored in PostgreSQL.

Reasons:

- Transaction support.
- Relational integrity.
- Reliable message persistence.
- Flexible querying.

### 4. Redis Pub/Sub for Initial Event Distribution

Redis Pub/Sub will distribute events between gateway instances.

Reasons:

- Simple to operate.
- Easy to simulate locally.
- Suitable for learning distributed event delivery.

Limitation:

Redis Pub/Sub does not persist messages for disconnected subscribers. The application must recover missed messages through database synchronization.

### 5. Transactional Outbox for Reliable Publishing

Message persistence and event creation will happen in the same database transaction.

A separate worker will publish pending events to Redis.

Reasons:

- Prevents a committed message from being lost merely because publishing failed immediately afterward.
- Supports retry.
- Separates persistence from event distribution.

### 6. MinIO for Local Object Storage

Attachments will be stored in MinIO using an S3-compatible API.

PostgreSQL will store attachment metadata and object keys.

### 7. Docker Compose for Local Distributed Simulation

Docker Compose will run the infrastructure and multiple gateway instances locally.

Reasons:

- Reproducible environment.
- Easier integration testing.
- Demonstrates cross-instance communication without requiring cloud infrastructure.

## Consequences

Positive:

- Clear architectural boundaries.
- Practical distributed-systems learning.
- Testable components.
- Lower initial operational complexity.

Trade-offs:

- Redis Pub/Sub requires recovery and synchronization mechanisms.
- The outbox worker introduces additional processing and retry logic.
- Separate processes require configuration and lifecycle management.
- A modular monolith still requires discipline to prevent package coupling.

## Future Considerations

Revisit these decisions when:

- Event volume requires a durable message broker such as Kafka.
- A module needs independent deployment.
- Presence or connection management becomes a bottleneck.
- Cloud deployment introduces additional reliability requirements.

Do not introduce additional infrastructure until a measurable requirement justifies it.
