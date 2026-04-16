# 👻 GhostAI Lite

> A premium, production-ready AI chat foundation built with Go and Next.js.

GhostAI Lite is an 8-week structured project designed to teach full-stack AI engineering, from SSE streaming and local inference (Ollama) to CI/CD and production deployment.

![GhostAI Lite UI](file:///home/ubuntu/.gemini/antigravity/brain/4e83cdf3-7cca-49d6-9c30-3f7f98bc2fb8/final_ghostai_ui_state_1776335039842.png)

## 🚀 Week 1 Accomplishments

- **Monorepo Architecture**: Clean separation between `backend/` (Go) and `frontend/` (Next.js).
- **Go SSE Streaming**: High-performance backend using Chi router and Server-Sent Events to stream LLM responses from Ollama.
- **Premium UI/UX**: Emerald green & black "hacker" aesthetic with glassmorphism, scanline effects, and a custom ghost logo.
- **Docker Orchestration**: Multi-stage Dockerfiles and a root `docker-compose.yml` for instant local development.
- **Configurable Inference**: Integrated token limits (`num_predict`) and model selection (defaulting to `tinyllama`).

## 🛠️ Tech Stack

- **Backend**: Go (Chi, SSE, NDJSON)
- **Frontend**: Next.js 14 (App Router, TS, Vanilla CSS)
- **AI**: Ollama (local LLM inference)
- **DevOps**: Docker, Docker Compose

## 🚦 Quick Start

### Prerequisites
- [Docker](https://docs.docker.com/get-docker/) & [Ollama](https://ollama.com/)
- `ollama pull tinyllama`

### Spin up the stack
```bash
docker compose up --build
```
The UI will be available at `http://localhost:3000` and the API at `http://localhost:8080`.

## 📖 Documentation
- [Developer Guide](docs/DEV_GUIDE.md) — Deep dive into architecture and features.
- [Week 1 Walkthrough](.system_generated/walkthrough.md) — Final verification screenshots and feature recap.

---
Built with 💚 and ghosts.
