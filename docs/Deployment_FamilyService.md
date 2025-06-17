# Deployment Document

## Family Service GraphQL

### 1. Introduction

#### 1.1 Purpose
This document provides instructions for deploying the Family Service GraphQL application to a Docker environment using docker-compose. It covers environment setup, configuration, deployment steps, and monitoring.

#### 1.2 Scope
This deployment document covers:
- Prerequisites for deployment
- Environment configuration
- Docker and docker-compose setup
- Application deployment
- Database setup
- Monitoring and maintenance

#### 1.3 References
- Software Requirements Specification (SRS)
- Software Design Document (SDD)
- Software Test Plan (STP)

### 2. Deployment Prerequisites

#### 2.1 Hardware Requirements
- Minimum 2 CPU cores
- Minimum 4GB RAM
- Minimum 10GB disk space

#### 2.2 Software Requirements
- Docker Engine 20.10.0+
- Docker Compose 2.0.0+
- Git (for source code management)

#### 2.3 Network Requirements
- Outbound internet access for pulling Docker images
- Ports 8089, 27017, and 5432 available on the host

### 3. Environment Setup

#### 3.1 Directory Structure
Create a deployment directory with the following structure:
```
family-service/
├── .env
├── docker-compose.yml
├── logs/
├── secrets/          # Contains sensitive configuration files
│   ├── gf_security_admin_user.txt
│   ├── gf_security_admin_password.txt
│   ├── mongo_root_username.txt
│   ├── mongo_root_password.txt
│   ├── postgresql_username.txt
│   ├── postgresql_password.txt
│   └── redis_password.txt
└── data/
    ├── mongo/
    └── postgres/
```

See the [Secrets Setup Guide](Secrets_Setup_Guide.md) for details on creating the required files in the secrets folder.

#### 3.2 Environment Configuration
Create a `.env` file with the following variables:

```
# Application configuration
PORT=8089
DB_DRIVER=mongo  # or postgres
LOG_LEVEL=info

# MongoDB configuration
MONGO_URI=mongodb://mongo:27017/family_service
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=password

# PostgreSQL configuration
DB_URI=postgres://postgres:postgres@postgres:5432/familydb
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=familydb
```

### 4. Docker Configuration

#### 4.1 Docker Compose File
Create a `docker-compose.yml` file with the following content:

```yaml
version: '3.8'
services:
  family-service:
    image: family-service-graphql:latest
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8089:8089"
    environment:
      - PORT=${PORT}
      - DB_DRIVER=${DB_DRIVER}
      - MONGO_URI=${MONGO_URI}
      - DB_URI=${DB_URI}
      - LOG_LEVEL=${LOG_LEVEL}
    depends_on:
      - mongo
      - postgres
    restart: unless-stopped
    volumes:
      - ./logs:/app/logs
      - ./secrets:/app/secrets  # Mount secrets folder
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8089/healthz"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s

  mongo:
    image: mongo:6
    ports:
      - "27017:27017"
    environment:
      - MONGO_INITDB_ROOT_USERNAME=${MONGO_INITDB_ROOT_USERNAME}
      - MONGO_INITDB_ROOT_PASSWORD=${MONGO_INITDB_ROOT_PASSWORD}
    volumes:
      - ./data/mongo:/data/db
    restart: unless-stopped
    healthcheck:
      test: echo 'db.runCommand("ping").ok' | mongosh localhost:27017/test --quiet
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=${POSTGRES_DB}
    ports:
      - "5432:5432"
    volumes:
      - ./data/postgres:/var/lib/postgresql/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
```

#### 4.2 Dockerfile
Create a `Dockerfile` in the project root with the following content:

```
# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o family-service ./cmd

# Final stage
FROM alpine:3.17

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates curl

# Copy the binary from the builder stage
COPY --from=builder /app/family-service .

# Expose the application port
EXPOSE 8089

# Set the entry point
ENTRYPOINT ["/app/family-service"]
```

### 5. Database Setup

#### 5.1 MongoDB Setup
The MongoDB container will be initialized with the root username and password specified in the `.env` file. The database will be created automatically when the application first connects to it.

#### 5.2 PostgreSQL Setup
The PostgreSQL container will be initialized with the username, password, and database name specified in the `.env` file. The database schema will be created automatically when the application first connects to it.

#### 5.3 Database Initialization
For a production environment, you may want to initialize the databases with schema and initial data. Create initialization scripts in the following directories:
- `./data/mongo/init/` for MongoDB
- `./data/postgres/init/` for PostgreSQL

### 6. Deployment Steps

#### 6.1 Building the Application
To build the application Docker image:

```bash
docker-compose build
```

#### 6.2 Starting the Services
To start all services:

```bash
docker-compose up -d
```

#### 6.3 Verifying Deployment
To verify that all services are running:

```bash
docker-compose ps
```

To check the application logs:

```bash
docker-compose logs -f family-service
```

To verify the application is responding:

```bash
curl http://localhost:8089/healthz
```

#### 6.4 Stopping the Services
To stop all services:

```bash
docker-compose down
```

### 7. Scaling and High Availability

#### 7.1 Horizontal Scaling
For horizontal scaling, you can run multiple instances of the application:

```bash
docker-compose up -d --scale family-service=3
```

Note: This requires additional configuration for load balancing.

#### 7.2 Database High Availability
For production environments, consider using:
- MongoDB replica sets
- PostgreSQL with replication

#### 7.3 Load Balancing
For production environments, consider using:
- Nginx as a reverse proxy
- Kubernetes for container orchestration

### 8. Monitoring and Maintenance

#### 8.1 Logging
Application logs are stored in the `./logs` directory. You can view them with:

```bash
docker-compose logs -f family-service
```

#### 8.2 Health Checks
The application provides a health check endpoint at `/healthz`. You can monitor it with:

```bash
curl http://localhost:8089/healthz
```

#### 8.3 Monitoring with Prometheus and Grafana

The Family Service includes built-in support for monitoring and observability using Prometheus and Grafana. This allows you to collect and visualize metrics about the application's performance and behavior.

##### 8.3.1 Prometheus Integration

Prometheus is used to scrape and store metrics from the Family Service. The service exposes metrics at the `/metrics` endpoint, which Prometheus scrapes at regular intervals.

**Configuration**:

The Prometheus configuration is defined in `prometheus.yml`:

```yaml
global:
  scrape_interval:     15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: "family_service"
    metrics_path: "/metrics"
    scrape_interval: 15s
    static_configs:
      - targets: ["family_service:8089"]
```

**Accessing Prometheus**:

Once the application is deployed, you can access the Prometheus UI at:

```
http://localhost:9090
```

From here, you can query metrics and create graphs. For example, you can view information about the Go runtime:

![Prometheus Go Info Query](assets/prometheus_go_info_query.png)

You can also check that Prometheus is successfully scraping metrics from the Family Service:

![Prometheus Target Query](assets/prometheus_target_query.png)

##### 8.3.2 Grafana Integration

Grafana is used to create dashboards that visualize the metrics collected by Prometheus. The Family Service includes a custom Grafana dashboard that provides insights into the application's performance.

**Configuration**:

The Grafana configuration is defined in `grafana.ini`. The custom dashboard for the Family Service is provided in `docs/assets/grafana_dashboard_for_family_service.json`.

**Accessing Grafana**:

Once the application is deployed, you can access the Grafana UI at:

```
http://localhost:3000
```

Default credentials are stored in the `secrets` folder:
- Username: Contents of `secrets/gf_security_admin_user.txt`
- Password: Contents of `secrets/gf_security_admin_password.txt`

**Custom Dashboard**:

The Family Service includes a custom Grafana dashboard that visualizes key metrics:

![Grafana Dashboard](assets/grafana_dashboard_for_family_service.png)

This dashboard shows:
- Heap allocations over time
- (Additional metrics can be added as needed)

**Importing the Dashboard**:

To import the custom dashboard:

1. Log in to Grafana
2. Click on "+" > "Import"
3. Click "Upload JSON file"
4. Select `docs/assets/grafana_dashboard_for_family_service.json`
5. Click "Import"

##### 8.3.3 Available Metrics

The Family Service exposes the following metrics:

- **Go Runtime Metrics**:
  - `go_memstats_heap_alloc_bytes`: Current heap allocation
  - `go_memstats_heap_inuse_bytes`: Heap memory in use
  - `go_goroutines`: Number of goroutines
  - `go_threads`: Number of OS threads created

- **HTTP Metrics**:
  - `http_requests_total`: Total number of HTTP requests
  - `http_request_duration_seconds`: HTTP request duration
  - `http_requests_in_flight`: Current number of HTTP requests being processed

- **Database Metrics**:
  - `db_operations_total`: Total number of database operations
  - `db_operation_duration_seconds`: Database operation duration
  - `db_connections_open`: Current number of open database connections

- **Application Metrics**:
  - `app_errors_total`: Total number of application errors

#### 8.4 Backup and Restore

##### 8.4.1 MongoDB Backup
To backup MongoDB data:

```bash
docker exec -it family-service_mongo_1 mongodump --out /data/db/backup
```

##### 8.4.2 PostgreSQL Backup
To backup PostgreSQL data:

```bash
docker exec -it family-service_postgres_1 pg_dump -U postgres familydb > backup.sql
```

##### 8.4.3 Restore from Backup
To restore MongoDB data:

```bash
docker exec -it family-service_mongo_1 mongorestore /data/db/backup
```

To restore PostgreSQL data:

```bash
cat backup.sql | docker exec -i family-service_postgres_1 psql -U postgres familydb
```

#### 8.4 Updating the Application
To update the application to a new version:

1. Pull the latest code:
   ```bash
   git pull
   ```

2. Rebuild and restart the services:
   ```bash
   docker-compose down
   docker-compose build
   docker-compose up -d
   ```

### 9. Troubleshooting

#### 9.1 Common Issues

##### 9.1.1 Application Fails to Start
- Check the application logs: `docker-compose logs family-service`
- Verify environment variables in `.env` file
- Ensure databases are running: `docker-compose ps`

##### 9.1.2 Database Connection Issues
- Check database logs: `docker-compose logs mongo` or `docker-compose logs postgres`
- Verify database credentials in `.env` file
- Ensure database ports are not in use by other services

##### 9.1.3 API Not Responding
- Check if the application container is running: `docker-compose ps`
- Check application logs for errors: `docker-compose logs family-service`
- Verify the port mapping in `docker-compose.yml`

#### 9.2 Debugging
To access a shell in the application container for debugging:

```bash
docker-compose exec family-service sh
```

### 10. Security Considerations

#### 10.1 Environment Variables and Secrets
- Use strong passwords for database credentials
- Do not commit `.env` file to version control
- Set up the required secrets files as described in the [Secrets Setup Guide](Secrets_Setup_Guide.md)
- Ensure the `secrets` folder is not committed to version control (it's already in .gitignore)
- Consider using a secrets management solution for production

#### 10.2 Network Security
- Use a reverse proxy (e.g., Nginx) to handle TLS termination
- Restrict access to database ports
- Use a private Docker network for inter-service communication

#### 10.3 Container Security
- Use the latest Docker images
- Run containers with non-root users
- Regularly update dependencies

### 11. Production Considerations

#### 11.1 Production Environment
For a production environment, consider:
- Using a container orchestration platform like Kubernetes
- Implementing a CI/CD pipeline for automated deployment
- Setting up monitoring and alerting
- Configuring backups and disaster recovery

#### 11.2 Performance Tuning
- Adjust database connection pool sizes
- Configure appropriate JVM memory settings
- Implement caching where appropriate

### 12. Appendices

#### 12.1 Environment Variable Reference
| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Application port | 8089 |
| DB_DRIVER | Database driver (mongo or postgres) | mongo |
| LOG_LEVEL | Logging level | info |
| MONGO_URI | MongoDB connection URI | mongodb://mongo:27017/family_service |
| DB_URI | PostgreSQL connection URI | postgres://postgres:postgres@postgres:5432/familydb |

#### 12.2 Docker Compose Commands Reference
| Command | Description |
|---------|-------------|
| `docker-compose build` | Build the application image |
| `docker-compose up -d` | Start all services in detached mode |
| `docker-compose down` | Stop all services |
| `docker-compose logs -f` | View logs from all services |
| `docker-compose ps` | List running services |

#### 12.3 Deployment Diagrams

##### 12.3.1 Physical Deployment Diagram
![Deployment Diagram](diagrams/Deployment%20Diagram.svg)

This deployment diagram illustrates the physical architecture of the system, showing the Docker host, containers, volumes, and connections between components. It includes details about exposed ports and environment variables.

##### 12.3.2 Container Diagram
![Container Diagram](diagrams/Deployment%20Container%20Diagram.svg)

This container diagram focuses on the Docker Compose configuration, showing the relationships between containers, networks, and volumes. It includes details about service dependencies, port forwarding, and deployment commands.

#### 12.4 Database Schema Reference
Refer to the Software Design Document (SDD) for database schema details.
