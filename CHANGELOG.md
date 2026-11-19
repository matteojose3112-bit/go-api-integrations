# Changelog

All notable changes to this project are documented here.

## [0.1.0] — 2026-11-19

### Added
- Go 1.23 module and standard-library HTTP server.
- Health endpoint and configurable `PORT`.
- JSON request decoding with a 1 MiB body limit and unknown-field rejection.
- Contact normalization endpoint with name and email validation.
- Webhook ingestion endpoint with provider/event ID idempotency.
- Event listing endpoint with optional provider filtering.
- Unit and HTTP-handler tests using `testing` and `httptest`.
- GitHub Actions workflow for tests and build.
- Render deployment configuration and free-plan web service.

### Deployment
- Deployed service: [relayops-go-api](https://dashboard.render.com/web/srv-dbvmqt60tbcc7390b2q0)
- Live API: https://relayops-go-api.onrender.com
- Health endpoint verified successfully: `{"service":"relayops-api","status":"ok"}`.

### Known limitations
- Events are held in memory and are not durable across restarts.
- Persistent storage, authentication, rate limiting, and provider-specific webhook signature verification remain future work.
