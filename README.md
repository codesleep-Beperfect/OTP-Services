# OTP Verification System

A distributed microservices-based OTP (One-Time Password) system built in Go. It allows tenants to register, generate OTPs for their users, and deliver them via email using an event-driven architecture.

---

## Architecture Overview

```
                    +------------------------------------------+
                    |                 CLIENT                   |
                    +------------------------------------------+
                         |                        |
                 register/validate          send/resend/verify
                         |                        |
                         v                        v
              +--------------------+    +----------------------+
              |   TENANT SERVICE   |    |     OTP SERVICE      |
              |     Port 8080      |<---|     Port 8081        |
              |                    |    |                      |
              | POST /register     |    | POST /otp/send       |
              | GET  /validate     |    | POST /otp/resend     |
              +--------+-----------+    | POST /otp/verify     |
                       |               +----------+-----------+
                       v                          |
                  +---------+            +--------+--------+
                  |  MySQL  |            |                 |
                  | tenants |          Redis            Kafka
                  +---------+        (OTP store)     (otp-email)
                                      10-min TTL          |
                                                          v
                                             +-----------------------+
                                             |     EMAIL SERVICE     |
                                             |   (Kafka Consumer)    |
                                             |                       |
                                             | Group:                |
                                             | email-service-group   |
                                             +-----------+-----------+
                                                         |
                                                         v
                                                   +-----------+
                                                   |   SMTP    |
                                                   |  (Gmail)  |
                                                   |           |
                                                   | Sends OTP |
                                                   |  to user  |
                                                   +-----------+
```

---

## How It Works

### Step 1 — Register Tenant
```
Client --> POST /v1/tenant/register --> Tenant Service --> MySQL
                                              |
                                        Returns api_key (sk_...)
```

### Step 2 — Send OTP
```
Client + api_key --> OTP Service
        |
  Validate api_key --> Tenant Service (HTTP)
        |
  Generate 4-digit OTP
        |
  Store hash in Redis (TTL: 10 min)
        |
  Publish OTPEvent to Kafka topic "otp-email"
        |
  Email Service consumes event --> Send email via SMTP
```

### Step 3 — Verify OTP
```
Client + api_key + otp --> OTP Service
        |
  Fetch from Redis --> Hash input OTP
        |
  Compare hashes
        |
  Match    --> Delete from Redis --> return { valid: true }
  No match --> return { valid: false }
```

---

## Services

### 1. Tenant Service — Port `8080`
Handles tenant registration and API key management.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v1/tenant/register` | Register a new tenant, returns API key |
| GET | `/v1/tenant/validate` | Validate API key, returns tenant ID |

```json
// POST /v1/tenant/register
Request:  { "name": "my-app", "email": "admin@myapp.com" }
Response: { "id": "uuid", "name": "my-app", "email": "admin@myapp.com", "api_key": "sk_..." }
```

---

### 2. OTP Service — Port `8081`
Generates, validates, and resends OTPs. Requires `x-api-key` header on every request.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v1/otp/send` | Generate OTP and trigger email |
| POST | `/v1/otp/resend` | Resend OTP (max 3 times) |
| POST | `/v1/otp/verify` | Verify OTP code |

```json
// POST /v1/otp/send  or  /v1/otp/resend
Request:  { "identifier": "user@example.com" }
Response: { "Message": "OTP Sent Successfully" }

// POST /v1/otp/verify
Request:  { "identifier": "user@example.com", "otp": "4821" }
Response: { "valid": true }
```

---

### 3. Email Service — Background Worker
Consumes OTP events from Kafka and sends OTP emails via SMTP. No HTTP endpoints.

- Consumer Group: `email-service-group`
- Skips expired OTPs before sending (checks `ExpiresAt` from Kafka event)

---

## Security

| Feature | Detail |
|---------|--------|
| API Key Auth | Every OTP request validates `x-api-key` against Tenant Service |
| OTP Hashing | Stored as `SHA256(tenantID:identifier:otp)` — never plain text |
| Expiry | Redis TTL of 10 minutes auto-deletes the OTP |
| Resend Limit | Max 3 resends per OTP session |
| Auto-delete | OTP deleted from Redis immediately on successful verify |

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.23+ |
| HTTP Framework | Gin |
| Configuration | Viper (YAML) |
| Tenant Storage | MySQL |
| OTP Cache | Redis |
| Event Streaming | Apache Kafka |
| Email | SMTP (Gmail) |

---

## Project Structure

```
OTP_system/
|
+-- tenant-service/              # Tenant registration & API key management
|   +-- cmd/main.go
|   +-- config/
|   +-- internal/
|       +-- handler/
|       +-- service/
|       +-- repository/
|       +-- model/
|       +-- db/
|       +-- router/
|       +-- utils/
|
+-- otp-service/                 # OTP generation, verification, resend
|   +-- cmd/main.go
|   +-- config/
|   +-- internal/
|       +-- handler/
|       +-- service/
|       +-- repository/
|       +-- model/
|       +-- kafka/
|       +-- client/
|       +-- router/
|       +-- utils/
|
+-- email-service/               # Kafka consumer & email delivery
    +-- cmd/main.go
    +-- config/
    +-- internal/
        +-- kafka/
        +-- email/
        +-- model/
```

---

## Running Locally

### Prerequisites
- Go 1.23+
- Docker

### 1. Start Infrastructure

```bash
# Redis
docker run -d --name redis -p 6379:6379 redis

# Kafka (single broker setup)
docker run -d --name kafka -p 9092:9092 -e KAFKA_CFG_NODE_ID=1 -e KAFKA_CFG_PROCESS_ROLES=broker,controller -e KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093 -e KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 -e KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT -e KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@localhost:9093 -e KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER -e KAFKA_CFG_OFFSETS_TOPIC_REPLICATION_FACTOR=1 -e KAFKA_CFG_TRANSACTION_STATE_LOG_REPLICATION_FACTOR=1 -e KAFKA_CFG_TRANSACTION_STATE_LOG_MIN_ISR=1 bitnami/kafka:3.7

# Create Kafka topic (wait ~10 sec after Kafka starts)
docker exec -it kafka kafka-topics.sh --create --topic otp-email --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
```

### 2. Configure Each Service
Update `config/config.yml` in each service with your credentials (MySQL DSN, Gmail password, etc).

### 3. Run Services

```bash
# Terminal 1
cd tenant-service && go run cmd/main.go

# Terminal 2
cd otp-service && go run cmd/main.go

# Terminal 3
cd email-service && go run cmd/main.go
```

---

## Dependencies

```
github.com/gin-gonic/gin
github.com/redis/go-redis/v9
github.com/segmentio/kafka-go
github.com/spf13/viper
github.com/go-sql-driver/mysql
github.com/google/uuid
```
