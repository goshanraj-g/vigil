# Vigil

An AI web monitoring tool that searches the web continuously and delivers relevant, summarized results to your webhook.

<img alt="vigildemo" src="https://github.com/user-attachments/assets/cbc13079-256d-4491-9e17-ecf9f96ab222" />

## How it works
1) You create a monitor with a query (e.g. `"Anthropic Claude Mythos News"`)
2) The scheduler enqueues a search job every 30 minutes via Redis
3) The worker searches the web using Serper (Google Search API), filtered to the last 24 hours
4) Each result is checked against a Redis dedup cache
5) Claude filters remaining results for relevance for query
6) Relevant results are scraped with Parallel Extract and summarized by Haiku
7) Results are stored in Postgres and delivered to your webhook


## Technologies
<img width="750" alt="technologies" src="https://github.com/user-attachments/assets/a16cf98b-7310-416a-8aee-29958539d3ff" />

## System Architecture
<img width="750" alt="vigilarchitecture" src="https://github.com/user-attachments/assets/e6f12c91-9dd2-44ad-9868-19998ffe86d6" />

## Inspiration
I came across an article about [Parallel's Monitor API](https://parallel.ai/blog/monitor-api), and was pretty curious on how it actually works, so I build my own version to learn more.

## Getting Started

### Prerequisties 
- Go 1.22+
- Docker (for Postgres and Redis)

### Setup

1. Clone the repo and install dependencies:
   ```bash
   git clone https://github.com/goshanraj-g/vigil
   cd vigil
   go mod tidy

2. Copy the example env file and fill in your API keys:
cp .env.example .env
3. Start Postgres and Redis:
docker compose up -d
4. Run the API server and worker:
go run ./cmd/api
go run ./cmd/worker

Slack Commands
```
┌───────────────────┬──────────────────────────┐
│      Command      │       Description        │
├───────────────────┼──────────────────────────┤
│ /monitor <query>  │ Start monitoring a query │
├───────────────────┼──────────────────────────┤
│ /monitor list     │ List active monitors     │
├───────────────────┼──────────────────────────┤
│ /monitor stop <n> │ Stop monitor #n          │
└───────────────────┴──────────────────────────┘
```
