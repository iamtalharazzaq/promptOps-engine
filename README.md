# 👻 GhostAI Lite

> AI assistant with structured outputs, observability & full CI/CD

A production-ready AI chat application built with **Go**, **Next.js**, **Ollama**, and **Docker**. Run LLMs locally with a premium ChatGPT-style interface.

---

## Tech Stack

| Layer      | Technology           | Purpose                        |
| ---------- | -------------------- | ------------------------------ |
| Frontend   | Next.js 14 + TypeScript | Chat UI with streaming render |
| Backend    | Go + Chi router      | REST API with SSE streaming    |
| AI Engine  | Ollama               | Local LLM inference            |
| Database   | Supabase *(Week 3+)* | Conversation persistence       |
| DevOps     | Docker + CircleCI    | Containerisation & CI/CD       |

---

## Quick Start

### Prerequisites

- **Go** 1.22+
- **Node.js** 18+
- **Ollama** ([install guide](https://ollama.ai))
- **Docker** & Docker Compose *(optional, for containerised setup)*

### Option 1: Docker Compose (recommended)

```bash
# Clone and start all services
git clone <your-repo-url> ghostAI-lite
cd ghostAI-lite
docker compose up --build

# Pull a model (first time only)
docker compose exec ollama ollama pull tinyllama

# Open http://localhost:3000
```

### Option 2: Run locally

```bash
# 1. Start Ollama
ollama serve
ollama pull tinyllama

# 2. Start backend
cd backend
cp .env.example .env
go run main.go

# 3. Start frontend
cd frontend
npm install
npm run dev

# Open http://localhost:3000
```

---

## Project Structure

```
ghostAI-lite/
├── backend/               # Go REST API
│   ├── main.go           # Entry point & route wiring
│   ├── config/           # Environment configuration
│   ├── handlers/         # HTTP handlers (health, chat)
│   ├── middleware/        # CORS, request logging
│   ├── services/         # Ollama client
│   ├── Dockerfile        # Multi-stage build
│   └── .env.example      # Env var template
├── frontend/              # Next.js chat UI
│   ├── app/              # App Router pages & layout
│   ├── components/       # React components
│   ├── lib/              # API client & utilities
│   └── Dockerfile        # Multi-stage build
├── docs/                  # Project documentation
│   └── DEV_GUIDE.md      # Full developer guide
├── docker-compose.yml     # Local dev orchestration
└── .env.example           # Root env template
```

---

## API Endpoints

| Method | Path      | Description                          |
| ------ | --------- | ------------------------------------ |
| GET    | `/health` | Liveness check — returns `{"status":"ok"}` |
| POST   | `/chat`   | Stream an LLM response via SSE      |

📖 See [`docs/DEV_GUIDE.md`](docs/DEV_GUIDE.md) for full API reference.

---

## Roadmap

- [x] **Week 1** — Monorepo, Go API, Ollama, Chat UI
- [ ] **Week 2** — Supabase integration, conversation history
- [ ] **Week 3** — Structured JSON output mode
- [ ] **Week 4** — Observability (latency, tokens/sec, error tracking)
- [ ] **Week 5** — CI/CD pipeline with CircleCI
- [ ] **Week 6** — Deployment (Vercel + Railway)
- [ ] **Week 7** — Polish & advanced features
- [ ] **Week 8** — Terraform infrastructure

---

## License

MIT
