# Go API Integrations

> A small REST integration hub for normalizing CRM contact data and accepting idempotent webhook events.

**Status:** Go API MVP implemented; CI verification pending  
**Language:** Go 1.23  
**Framework:** Go standard library `net/http`  
**Deployment configuration:** Render free web service blueprint prepared  
**Live API:** Pending deployment authorization and verification

## Overview

RelayOps API demonstrates a lightweight integration layer for business systems. It exposes a health endpoint, a contact-normalization endpoint, and a webhook ingestion endpoint with duplicate-event detection.

The project intentionally uses Go's standard library to keep the service small, portable, and easy to audit.

## Endpoints

| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/health` | Service health check |
| POST | `/api/v1/contacts/normalize` | Trim and normalize contact name, email, and company |
| POST | `/api/v1/webhooks/ingest` | Accept and deduplicate webhook events |
| GET | `/api/v1/events` | List events, optionally filtered by provider |

### Normalize a contact

```bash
curl -X POST http://localhost:8080/api/v1/contacts/normalize \
  -H "Content-Type: application/json" \
  -d '{"name":" Ada  Lovelace ","email":" ADA@example.com ","company":" Analytical Engines "}'
```

### Ingest a webhook

```bash
curl -X POST http://localhost:8080/api/v1/webhooks/ingest \
  -H "Content-Type: application/json" \
  -d '{"provider":"crm","event_id":"evt-001","event_type":"contact.created","payload":{"contact_id":42}}'
```

A first event returns HTTP 202. Replaying the same provider/event ID returns HTTP 200 with `duplicate: true`.

## Run locally

Requires Go 1.23 or newer.

```bash
go test ./...
go run ./cmd/api
```

The API listens on port 8080 by default. Set `PORT` to override it.

## Project structure

```text
go-api-integrations/
├── cmd/api/main.go
├── internal/api/
│   ├── router.go
│   └── router_test.go
├── .github/workflows/ci.yml
├── go.mod
├── render.yaml
├── README.md
├── ROADMAP.md
└── CHANGELOG.md
```

## Quality and deployment

- Unit and HTTP-handler tests use Go's built-in `testing` and `httptest` packages.
- GitHub Actions runs `go test ./...` and `go build ./...`.
- `render.yaml` defines a free-plan Go web service with a `/health` health check.

The Render blueprint is prepared, but no additional Render service has been created yet. The live API URL will be published after deployment is explicitly authorized and verified.

## Current limitations

- Webhook events are stored in memory and are lost when the process restarts.
- No persistent database, queue, authentication, or provider-specific signature verification is included.
- The API is a demonstration integration layer, not a production webhook receiver.
- Request bodies are size-limited and unknown JSON fields are rejected.

## Roadmap and changelog

- [Roadmap](ROADMAP.md)
- [Changelog](CHANGELOG.md)
