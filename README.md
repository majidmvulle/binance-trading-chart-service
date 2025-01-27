
# Binance Trading Chart Service

## Overview

This project implements a real-time trading chart service in Go. The service connects to the Binance API WebSocket stream, 
fetches real-time tick data for specified cryptocurrency symbols, aggregates this data into 1-minute OHLC (Open-High-Low-Close) candlesticks, 
broadcasts the candlestick data via a gRPC streaming API, and persists completed candlesticks to a PostgreSQL database. 
The service is designed to be deployed to a Kubernetes cluster using Terraform for infrastructure management.

**Key Features:**

*   **Real-time Data Ingestion:** Fetches tick data from the Binance WebSocket API for BTCUSDT, ETHUSDT, and PEPEUSDT symbols (configurable through environment variables).
*   **OHLC Candlestick Aggregation:** Aggregates tick data into 1-minute OHLC candlesticks.
*   **gRPC Streaming API:** Provides a gRPC streaming service to broadcast real-time candlestick data to clients.
*   **Data Persistence:** Persists completed 1-minute candlesticks to a PostgreSQL database for historical data storage.
*   **Kubernetes Deployment:** Deployed to a local Kubernetes cluster (using kind) and managed with Terraform for Infrastructure as Code (IaC).
*   **Unit Tests:** Includes unit tests for the core OHLC aggregation logic.

**Components:**

*   **Binance WebSocket API:** Provides real-time tick data for cryptocurrency symbols.
*   **Ingestor Service (Go):**
    *   Connects to Binance WebSocket API.
    *   Fetches and parses tick data.
    *   Aggregates tick data into 1-minute OHLC candlesticks.
    *   Broadcasts current candlestick data via gRPC streaming API.
    *   Written in Go (`ingestor` directory).
*   **Aggregator (Go Package):**  Internal Go package within `ingestor` (`internal/aggregator`) responsible for the OHLC aggregation logic.
*   **gRPC Streaming API:** Implemented using gRPC and protocol buffers (`internal/grpc`). Allows clients to subscribe to a real-time stream of candlestick data.
*   **Persistor Service (Go):**
    *   gRPC client that subscribes to the `ingestor`'s candlestick stream.
    *   Persists completed 1-minute candlestick data to a PostgreSQL database.
    *   Written in Go (`persistor` directory).
*   **PostgreSQL Database:** Relational database used to store historical candlestick data.
*   **Kubernetes Cluster:** Target deployment environment for the services (provisioned locally using kind and Terraform).
*   **Terraform:** Infrastructure-as-Code tool used to provision the local Kubernetes cluster and deploy the services.

## Setup and Prerequisites

**Software Requirements:**

*   [Go](https://go.dev/dl/) (version >= 1.22 recommended)
*   [Docker](https://docs.docker.com/get-docker/)
*   [Terraform](https://www.terraform.io/downloads) (version >= 1.0 recommended)
*   [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
*   [kind](https://kind.sigs.k8s.io/docs/user/quick-start/) (for local Kubernetes cluster provisioning - **optional**)
*   [PostgreSQL](https://www.postgresql.org/download/) (or Docker for PostgreSQL - for local database setup)
*   [Goose](https://github.com/pressly/goose) (for database migrations)
*   [GolangCI-Lint](https://golangci-lint.run/usage/install/) (for linting go code)
*   `protoc` and Go gRPC tools (installation instructions in gRPC documentation)
*   [protoc-go-inject-tag](https://pkg.go.dev/github.com/syncore/protoc-go-inject-tag)
*   `mockery` (for generating mocks - installation instructions in mockery documentation)

**Environment Variables:**

Create `.env` files in the `ingestor` and `persistor` directories (example `.env.example` files are provided in each directory). You will need to set the following environment variables:

*   **`ingestor/.env`:**
    *   `APP_GRPC_PORT`: Port for the ingestor gRPC server (e.g., `50051`).
    *   `BINANCE_WEBSOCKET_BASE_URL`: Base URL for Binance WebSocket API (e.g., `wss://stream.binance.com:9443`).
    *   `BINANCE_SYMBOLS`: Space-separated list of symbols to fetch (e.g., `BTCUSDT ETHUSDT PEPEUSDT`).
    *   `APP_DEBUG`: Set to `true` for debug logging, `false` for production.

*   **`persistor/.env`:**
    *   `SERVER_ADDRESS`: The address of the gRPC server (ingestor service) to consume the stream from.
    *   `DB_WRITE_DSN`: PostgreSQL database _write_ connection string (e.g., `postgres://user:password@host:port/dbname?sslmode=disable`). **Important: For production, use Kubernetes Secrets to manage database credentials securely instead of hardcoding in `.env` files.**
    *   `DB_READ_DSN`: PostgreSQL database _read_ connection string (e.g., `postgres://user:password@host:port/dbname?sslmode=disable`). **Important: For production, use Kubernetes Secrets to manage database credentials securely instead of hardcoding in `.env` files. (Good practice for future read replicas)**

**Infrastructure Setup (Kubernetes with Terraform):**

This project uses Terraform to deploy to a Kubernetes cluster (you can optionally provision a local Kubernetes cluster using `kind`). 

1.  **Install `kind` (Optional):** Follow the [kind installation instructions](https://kind.sigs.k8s.io/docs/user/quick-start/) if you want to provision a local Kubernetes cluster using kind with Terraform. If you are using an existing Kubernetes cluster, skip this step.
2.  **Install Terraform:** Follow the [Terraform installation instructions](https://www.terraform.io/downloads).
3.  **Navigate to the `infra/terraform/kind` directory:** `cd infra/terraform/kind` to set up the `kind` cluster.
4.  **Apply Terraform Configuration:** `terraform apply` (This will provision a kind cluster named `binance-trading-chart-service-cluster` (configurable - if you choose to provision a new cluster).

## Build Instructions

**1. Generate gRPC Protobuf Code:**

Navigate to the `ingestor` directory and run:

```bash
cd ingestor
make gen/pb
```

This will generate the Go gRPC code from the `aggregator.proto` file in `ingestor/pkg/api/aggregator`.

**2. Build Go Services:**

Build the Go binaries for `ingestor` and `persistor` services:

```bash
# From the root project directory:
cd ingestor
go build -o cmd/ingestor ./cmd 

cd ../persistor
go build -o cmd/persistor ./cmd
```

**3. Build Docker Images:**

Build the Docker images `ingestor` and `persistor` services.

```bash
# From the root project directory:
cd ingestor
make docker/build
make docker/push

cd ../persistor
make docker/build
make docker/push
```

**4. Database Migrations:**

Run them using:

```bash
# From the root project directory:
cd persistor
make migrate/up
```

**Note:** Ensure your `DB_WRITE_DSN` and `DB_READ_DSN` environment variables are correctly set before running migrations.

## Run Instructions

**1. Deploy to Kubernetes using Terraform:**

If you haven't already, deploy the infrastructure and services to your local Kubernetes cluster (kind or any other) using Terraform:
Ensure you set the Kubernetes cluster environment variables in `infra/kubernetes/terraform.tfvars`.

- Copy `infra/kubernetes/terraform.tfvars.example` to `infra/kubernetes/terraform.tfvars`

*  **`infra/kubernetes/terraform.tfvars`:**
   * `db_write_dsn`: PostgreSQL database _write_ connection string.
   * `db_read_dsn`: PostgreSQL database _read_ connection string.
   * `k8s_host`: K8s host server endpoint.
   * `k8s_client_certificate`: K8s client certificate for Kubernetes API.
   * `k8s_client_key`: K8s client key for Kubernetes API.
   * `k8s_cluster_ca_certificate`: K8s cluster client ca certificate for Kubernetes API.


```bash
cd infra/terraform/kubernetes
terraform apply
```

Terraform will deploy the Kubernetes manifests for `ingestor` and `persistor` services.

**2. Verify Deployment:**

Use `kubectl` to verify that the pods and services are running correctly in your Kubernetes cluster:

```bash
kubectl get deployments -n default 
kubectl get pods -n default
kubectl get services -n default
```

**3. Check Persisted Data (PostgreSQL):**

Connect to your PostgreSQL database (using `psql` or a database client tool) and query the `agg_trade_ticks` table in your database to verify that candlestick data is being persisted by the `persistor` service.

## Test Instructions

**Run Unit Tests:**

To run the unit tests for the `aggregator` package, navigate to the test directory and run `make test`:

```bash
cd ingestor
make test
```

This will execute the unit tests defined in `aggregator_test.go` and report the test results.

## Future Improvements

Potential areas for future improvement and enhancements:

*   **Better Error Handling and Resilience:** Implement more comprehensive error handling throughout the services (Binance WebSocket client, aggregator, gRPC server, persistence), including reconnection logic, retry mechanisms, circuit breakers, and logging/monitoring for production reliability.
*   **Scalability and High Availability:** Design the services for horizontal scalability (e.g., using Kubernetes replica sets, horizontal pod autoscaling for `ingestor` and `persistor`). Also to consider database scalability and high availability as well.
*   **Advanced Kubernetes Deployment Strategies:** Implement more advanced Kubernetes deployment patterns (e.g., blue/green deployments, canary releases, rolling updates, health checks, readiness/liveness probes, resource quotas, etc.) for improved deployment and management in production.
*   **Monitoring and Logging:** Integrate proper monitoring and logging for all services (e.g., using Prometheus, Grafana, Elasticsearch, Kibana) to track performance, health, and identify issues in production.
*   **Enhanced Security:** Implement more robust security measures, especially for secret management (using Kubernetes Secrets encryption at rest, external secret stores), network policies (more restrictive egress rules), and secure communication (TLS/SSL for gRPC and potentially for WebSocket).
*   **Charting UI:** Develop a front-end charting UI (e.g., using React, Vue.js, or similar) that consumes the gRPC streaming API and visualizes the real-time candlestick data in interactive charts.
*   **Support for More Symbols and Timeframes:** Extend the service to support a wider range of cryptocurrency symbols and different candlestick timeframes (e.g., 5-minute, 15-minute, hourly candles).

## Author

*   Majid Mvulle
*   Date: 2025-01-27 
*   [@majidmvulle](https://github.com/majidmvulle)

## License
None - Do with as you please!
