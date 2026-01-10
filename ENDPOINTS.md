# API Endpoints Documentation

This document describes the API endpoints exposed by `internal/route/server` and `internal/route/blueprint`.

## Authentication

All endpoints described below are protected by an authentication middleware and require a valid Bearer token in the `Authorization` header.

## Blueprint API

**Base Path:** `/v1/blueprints`

### List All Blueprints

Retrieves a list of all P blueprints.

- **URL:** `/v1/blueprints/list`
- **Method:** `GET`
- **Handler:** `ListAllBlueprintsHandler`
- **Response:**
    - `200 OK`: JSON array of blueprint objects.

### Get Blueprint

Retrieves the details of a specific blueprint by its ID.

- **URL:** `/v1/blueprints/:id`
- **Method:** `GET`
- **Handler:** `GetBlueprintHandler`
- **URL Parameters:**
    - `id`: The unique identifier of the blueprint.
- **Response:**
    - `200 OK`: JSON object representing the blueprint.
    - `404 Not Found`: If the blueprint does not exist.

---

## Server API

**Base Path:** `/v1/servers`

### List All Servers

Retrieves a list of all managed servers.

- **URL:** `/v1/servers`
- **Method:** `GET`
- **Handler:** `ListAllServersHandler`
- **Response:**
    - `200 OK`: JSON array of server objects.

### Get Server

Retrieves the details of a specific server by its ID.

- **URL:** `/v1/servers/:id`
- **Method:** `GET`
- **Handler:** `GetServerHandler`
- **URL Parameters:**
    - `id`: The unique identifier of the server.
- **Response:**
    - `200 OK`: JSON object representing the server.
    - `404 Not Found`: If the server does not exist.

### Create Server

Creates a new server instance based on a blueprint. This process initializes the server record and attempts to create the underlying container.

- **URL:** `/v1/servers`
- **Method:** `POST`
- **Handler:** `CreateServerHandler`
- **Query Parameters:**
    - `blueprintId` (required): The ID of the blueprint to use for creating the server.
- **Response:**
    - `201 Created`: JSON object of the newly created server.
    - `400 Bad Request`: If `blueprintId` is missing.
    - `404 Not Found`: If the specified blueprint does not exist.
    - `500 Internal Server Error`: If the underlying container creation fails.

### Delete Server

Deletes an existing server and its associated container.

- **URL:** `/v1/servers/:id`
- **Method:** `DELETE`
- **Handler:** `DeleteServerHandler`
- **URL Parameters:**
    - `id`: The unique identifier of the server to delete.
- **Response:**
    - `200 OK`: JSON object of the deleted server.
    - `404 Not Found`: If the server does not exist.
    - `500 Internal Server Error`: If the underlying container deletion fails.

### Update Server Report

Updates the internal report/statistics for a specific server.

- **URL:** `/v1/servers/:id/report`
- **Method:** `PATCH`
- **Handler:** `UpdateServerReportHandler`
- **URL Parameters:**
    - `id`: The unique identifier of the server.
- **Body:** JSON object conforming to `data.ServerReport`.
- **Response:**
    - `200 OK`: JSON object of the updated server.
    - `400 Bad Request`: If the server ID is invalid or the request body is malformed.
    - `404 Not Found`: If the server does not exist.

### Update Server Status

Updates the status of a specific server.

- **URL:** `/v1/servers/:id/status`
- **Method:** `PATCH`
- **Handler:** `UpdateServerStatusHandler`
- **URL Parameters:**
    - `id`: The unique identifier of the server.
- **Body:**
    ```json
    {
      "status": "string"
    }
    ```
- **Response:**
    - `200 OK`: JSON object `{"status": "new_status"}`.
    - `400 Bad Request`: If the server ID is invalid, the payload is malformed, or the status value is invalid.
    - `404 Not Found`: If the server does not exist (or update failed).

### Restart Server

Restarts a server by deleting its existing container and creating a new one.

- **URL:** `/v1/servers/:id/restart`
- **Method:** `POST`
- **Handler:** `RestartServerHandler`
- **URL Parameters:**
    - `id`: The unique identifier of the server to restart.
- **Response:**
    - `200 OK`: JSON object `{"status": "restarted"}`.
    - `400 Bad Request`: If the server ID is invalid.
    - `404 Not Found`: If the server does not exist.
    - `500 Internal Server Error`: If container deletion or creation fails.
