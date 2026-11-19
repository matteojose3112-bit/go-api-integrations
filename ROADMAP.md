# Roadmap

Items are complete only when implemented and verified.

## Phase 1 — API foundation

- [x] Establish Go module and standard-library HTTP server.
- [x] Add health endpoint and configurable port.
- [x] Add JSON request decoding with a body-size limit.
- [x] Add structured JSON error responses.
- [x] Add contact normalization endpoint with basic validation.
- [x] Add webhook ingestion endpoint.
- [x] Add in-memory idempotency by provider and event ID.
- [x] Add event listing with provider filter.
- [x] Pass GitHub Actions tests and build on the latest verified commit.

## Phase 2 — Deployment

- [x] Prepare Render deployment configuration.
- [x] Create the Render free-plan web service in My Workspace.
- [x] Deploy API and verify the public health endpoint.
- [x] Publish live API URL and deployment details.
- [ ] Add API examples or a small interactive request client.

**Live API:** https://relayops-go-api.onrender.com  
**Health:** https://relayops-go-api.onrender.com/health  
**Render dashboard:** https://dashboard.render.com/web/srv-dbvmqt60tbcc7390b2q0

## Phase 3 — Reliability and integrations

- [ ] Persist events in PostgreSQL or another durable store.
- [ ] Add provider-specific webhook signature verification.
- [ ] Add structured logging, request IDs, and metrics.
- [ ] Add retry-safe background processing for asynchronous integrations.
- [ ] Add authentication and rate limiting for protected endpoints.
- [ ] Expand tests for malformed, oversized, and duplicate requests.

## Definition of done

- API endpoints are documented and reproducible locally.
- Automated tests and build pass.
- Deployed health endpoint is publicly reachable.
- Webhook persistence and security limitations are explicit.
- No credentials or private business data are committed.
