# Vigil

Push-based web monitoring API. User defines a query, system watches the web continuously and POSTs new results to a webhook.

Inspired by: parallel.ai/blog/monitor-api

## Stack

- **Language:** Go
- **Router:** chi
- **DB:** Postgres (pgx driver)
- **Job queue:** asynq + Redis
- **Dedup:** Redis (URL hashes)
- **Search:** Serper API
- **Relevance filter:** Claude Haiku via Anthropic Go SDK

## Architecture

Producer/consumer pattern:
- **API layer** — user creates/manages monitors, reads results
- **Worker layer** — scheduler enqueues jobs, workers run search → dedup → Claude → webhook

## Status

- [ ] go.mod (blocked: need GitHub username for module path)
- [ ] Data design / Postgres schema
- [ ] API design (endpoints)
- [ ] Worker design
- [ ] Implementation

## Context

User is learning software design as we build. Teach concepts inline. Keep explanations grounded in what we're actually building.
