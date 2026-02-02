# Singularity

> [!WARNING]
> **Educational project only.**
> This project is not intended for production use.

## Overview

**Singularity** is a lightweight control-plane for managing containerized game servers.

It coordinates metadata, state, and orchestration signals consumed by containers. Singularity does **not** run game servers directly - it manages their lifecycle through Docker.

## Technology Stack

- **Language:** Go 1.25
- **Web Framework:** Gin
- **Database:** MongoDB
- **Containerization:** Docker
- **Configuration:** [Pkl](https://pkl-lang.org/)

## Project Structure

```
cmd/
  api/                  # Application entrypoint

internal/
  auth/                 # JWT authentication
  config/               # Configuration loading
  data/                 # Data models
  docker/               # Container interaction
  manager/              # Blueprint and server managers
  route/                # HTTP routes
  strategy/             # Container strategies

pkl/
  Blueprint.pkl         # Blueprint schema

gen/
  blueprint/            # Generated Go code from Pkl
```

## API

### Blueprints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/v1/blueprints/list` | List all blueprints |
| `GET` | `/v1/blueprints/:id` | Get blueprint by ID |
| `POST` | `/v1/blueprints/reload` | Reload blueprints from disk |

### Servers

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/v1/servers` | List servers |
| `POST` | `/v1/servers` | Create server |
| `GET` | `/v1/servers/:id` | Get server |
| `DELETE` | `/v1/servers/:id` | Delete server |
| `PATCH` | `/v1/servers/:id/status` | Update status |
| `PATCH` | `/v1/servers/:id/report` | Update metrics |
| `POST` | `/v1/servers/:id/restart` | Restart server |

### Metrics

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/v1/metrics/prometheus` | Prometheus service discovery |

## Development

### Prerequisites

- Go 1.25+
- Docker
- [Pkl CLI](https://pkl-lang.org/main/current/pkl-cli/index.html)

### Setup

```bash
git clone <repository-url>
go mod download
```

### Regenerating Blueprint Code

After modifying `pkl/Blueprint.pkl`:

```bash
pkl-gen-go pkl/Blueprint.pkl --base-path singularity
```

### Building

```bash
go build -o singularity ./cmd/api
```

### Running

```bash
docker compose up -d   # MongoDB
./singularity
```

## License

Educational use only.
