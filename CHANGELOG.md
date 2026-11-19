# Changelog

All notable changes to this project are documented here.

## [Unreleased] — RelayOps API MVP

### Added
- Go 1.23 module and standard-library HTTP server.
- Health endpoint and configurable `PORT`.
- JSON request decoding with a 1 MiB body limit and unknown-field rejection.
- Contact normalization endpoint with name and email validation.
- Webhook ingestion endpoint with provider/event ID idempotency.
- Event listing endpoint with optional provider filtering.
- Unit and HTTP-handler tests using `testing` and `httptest`.
- GitHub Actions workflow for tests and build.
- Render free-plan Blueprint with health check.

### Pending verification
- Latest GitHub Actions test and build run.
- Render service creation and live endpoint verification.
- Durable event storage and provider signature verification.
