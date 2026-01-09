# Singularity

> [!WARNING]
> **Educational project only.**
> This project is not intended for production use.
> For production-ready ecosystems, see **SimpleCloud** or **CloudNet**.

---

## Overview

**Singularity** is a lightweight **control-plane** built as part of the **Event Horizon** study project.

Its purpose is to experiment with:

* server lifecycle management
* service discovery
* basic scaling orchestration
* inter-service communication

Target use case: **Hytale server infrastructures**, with a strong focus on **minimalism**, **explicit behavior**, and **low cognitive overhead**.

Singularity does **not** run game servers.
It coordinates metadata, state, and orchestration signals consumed by external agents and containers.

---

## Role in Event Horizon

Within the Event Horizon ecosystem, Singularity acts as the **central coordinator**.

Responsibilities:

* Track registered servers and their runtime state
* Expose an HTTP API for nodes and services
* Store ephemeral state and metadata
* Provide basic scaling and orchestration hooks

Non-responsibilities:

* No binary distribution
* No automatic high-availability

---

## Technology Stack

Chosen for simplicity, performance, and ecosystem maturity.

* **Language:** Go 1.25
* **Web Framework:** Gin
* **Database:** MongoDB
  Used for persistent metadata and historical state.
* **Key-Value Store:** Redis
  Used for ephemeral state, heartbeats, and fast lookups.
* **Containerization:** Docker & Docker Compose

---

## Project Structure

Standard Go layout with strict separation of concerns:

```
cmd/
 └─ api/            # Application entrypoint

internal/
 ├─ auth/           # JWT authentication
 ├─ config/         # Configuration loading (TOML)
 ├─ data/           # MongoDB repositories
 ├─ docker/         # Container and runtime interaction
 ├─ route/          # HTTP routes and handlers
 └─ strategy/       # Scaling and orchestration strategies

blueprints/
 └─ *.toml          # Blueprints configuration
```

---

## API Overview

Singularity exposes a RESTful HTTP API consumed by game nodes and auxiliary services.

### Server Endpoints

* `GET    /v1/servers`
  List all registered servers

* `POST   /v1/servers`
  Create a new server from a blueprint

* `GET    /v1/servers/:id`
  Retrieve server details

* `DELETE /v1/servers/:id`
  Delete a server and its container

* `PATCH  /v1/servers/:id/status`
  Update server lifecycle status

* `PATCH  /v1/servers/:id/report`
  Update runtime metrics / report

* `POST   /v1/servers/:id/restart`
  Restart a server instance

---

## Getting Started

### Prerequisites

* Go 1.25+
* Docker
* Docker Compose

---

### Installation

Clone the repository:

```bash
git clone <repository-url>
go build -o singularity-api ./cmd/api
```

Edit `config.toml` and define:

* `jwt_secret_key`
* MongoDB URI
* Redis URI

---

### Running with Docker

```bash
docker compose up --build
go build cmd\api\main.go
```

Default services:

* API: `http://localhost:8080`
* MongoDB: `localhost:27017`
* Redis: `localhost:6379`

---

## Non-Goals

This project intentionally does **not** attempt to solve:

* High availability
* Leader election
* Distributed consensus
* Auto-healing
* Production-grade security
* Multi-region deployments

If you need those, this project is the wrong tool.

---

## Contributing

This is a study project.
Contributions are welcome as long as they:

* keep complexity low
* avoid over-engineering
* align with the experimental nature of the project

---

## License

Educational use only.
