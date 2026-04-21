# 🛠️ PromptOps Engine — Developer Guide

> A production-grade orchestration platform for local LLMs, focusing on reliability, observability, and premium aesthetics.

---

## 🏗️ Architecture Overview

The PromptOps Engine is a monorepo consisting of a Go backend, a Next.js frontend, and a local Ollama inference server.

> [!TIP]
> For a more detailed technical explanation including sequence diagrams and data flow, see the **[Architectural Deep Dive](ARCHITECTURE.md)**.

```mermaid
graph TD
    User([User]) <--> Frontend[Next.js Frontend :3000]
    Frontend <--> Backend[Go Backend :8080]
    Backend <--> Ollama[Ollama Server :11434]
    
    subgraph Observability
        Prometheus[Prometheus :9090] -- scrapes --> Backend
        Grafana[Grafana :3001] -- visualizes --> Prometheus
    end
```

### Core Technologies
- **Backend (Go 1.24)**: High-performance API using `chi` for routing and `slog` for structured logging.
- **Frontend (Next.js 14)**: Premium interface with TypeScript and real-time SSE streaming.
- **Observability**: Native Prometheus instrumentation (`client_golang`) and Request ID tracing (`uuid`).
- **Reliability (Schema Guard)**: JSON Schema enforcement using `gojsonschema`.
- **Inference**: Orchestration of local LLMs via the Ollama REST API.

> [!NOTE]
> For a detailed list of packages and the rationale behind their selection, see the **[Dependency Rationale](ARCHITECTURE.md#dependency--service-rationale)** in the architecture guide.

---

## 📂 Project Structure

```bash
.
├── backend/                # Go API Server
│   ├── cmd/                # Entry points (if refactored)
│   ├── config/             # Environment configuration
│   ├── handlers/           # HTTP handlers (chat, health)
│   ├── middleware/         # CORS, Logging, Metrics, UUID
│   ├── pkg/                # Internal packages
│   │   ├── metrics/        # Prometheus metrics definitions
│   │   └── utils/          # Generic utilities (SSE helper)
│   └── services/           # Service layer (Ollama, Validator)
├── frontend/               # Next.js Web App
│   ├── app/                # App Router (page.tsx, layout.tsx)
│   ├── components/         # React components (ChatMessage, ChatInput)
│   ├── lib/                # API client and logic
│   └── public/             # Static assets
├── monitoring/             # Monitoring config (Prometheus/Grafana)
└── docs/                   # Documentation and Roadmap
```

---

## 🛡️ Schema Guard (Week 2)

PromptOps Engine ensures LLM reliability through **Structured Output Validation**.

1. **Schema Definition**: The frontend provides a JSON Schema.
2. **Strict Validation**: The backend validates the LLM's response using `gojsonschema`.
3. **Automated Retries**: If validation fails, the engine automatically retries (up to 3 times) with a correction prompt.
4. **Real-time Feedback**: UI badges show "Validating", "Valid", "Invalid", or "Retrying".

---

## 📉 Observability & Tracing (Week 3)

### Structured Logging
Every request is logged as JSON via `slog` and tagged with a unique `request_id` (UUIDv4).
```json
{
  "time": "2026-04-21T18:00:00Z",
  "level": "INFO",
  "msg": "request completed",
  "request_id": "8da45474-97c2-469e-91c1-1087aa2a4139",
  "method": "POST",
  "path": "/chat",
  "status": 200,
  "latency": "12.34ms"
}
```

### Metrics Endpoint
Exposed at `:8080/metrics` for Prometheus.
- `promptops_http_requests_total`: Request counts by status/method.
- `promptops_ollama_token_usage_total`: Input/Output token counts.
- `promptops_ollama_request_duration_seconds`: LLM response time latency.

---

## 🚀 Development Workflow

### Prerequisites
- **Docker & Docker Compose**
- **Go 1.24** (optional for local run)
- **Node.js 20** (optional for local run)

### Quick Start (Docker)
```bash
# 1. Start the full stack
docker compose up --build

# 2. Pull the default model
docker compose exec ollama ollama pull tinyllama

# 3. Access the UI
# http://localhost:3000
```

### Makefile Commands
- `make install-deps`: Install dev dependencies.
- `make backend-run`: Run backend locally.
- `make frontend-dev`: Run frontend locally.
- `make backend-test`: Run BDD test suite (Ginkgo).

---

## ✨ Design Principles
- **Aesthetics First**: Every component must use glassmorphism and emerald/hacker theme.
- **Zero Placeholder**: No generic placeholders; use generated assets or meaningful defaults.
- **Type Safety**: End-to-end TypeScript and Go struct validation.
