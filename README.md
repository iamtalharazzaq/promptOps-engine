# ⚙️ PromptOps Engine
**LLM Orchestration Platform with Schema Validation, Metrics & CI/CD**

---

![PromptOps Banner](https://img.shields.io/badge/Status-Development-emerald?style=for-the-badge)
![Go](https://img.shields.io/badge/Backend-Go-00ADD8?style=for-the-badge&logo=go)
![Next.js](https://img.shields.io/badge/Frontend-Next.js-black?style=for-the-badge&logo=next.js)
![Docker](https://img.shields.io/badge/Deploy-Docker-2496ED?style=for-the-badge&logo=docker)

---

## 🚀 Overview

**PromptOps Engine** is a high-performance, developer-centric platform designed to bridge the gap between local LLM inference and production-ready applications. It provides a robust orchestration layer with a focus on **Type Safety**, **Real-time Observability**, and **CI/CD Best Practices**.

> [!IMPORTANT]
> This project has completed the **Week 02** milestone (**Schema Guard**). Core orchestration, JSON validation, and automated retries are operational.

---

## ✨ Features

| Feature | Description | Status |
| :--- | :--- | :---: |
| **Real-time Streaming** | Multi-token SSE streaming for ultra-low latency responses. | ✅ |
| **Monorepo Workflow** | Unified Go backend and Next.js frontend management. | ✅ |
| **Local Inference** | Native Ollama integration for zero-cost, private AI usage. | ✅ |
| **Premium UI** | Stunning Emerald/Black "hacker" aesthetic with glassmorphism. | ✅ |
| **Schema Validation** | Type-safe JSON output enforcement for industrial use. | ✅ |
| **Metrics & Monitoring** | Native Prometheus metrics for token usage and latency. | ⏳ |
| **CI/CD Driven** | Automated deployment pipelines and container security. | ⏳ |

---

## 🛠️ Tech Stack

### Backend
- **Language**: Go (v1.22+)
- **Routing**: `chi` router (lightweight & fast)
- **Streaming**: Server-Sent Events (SSE)
- **Configuration**: Environment-driven with `.env` support

### Frontend
- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript
- **Styling**: Premium Vanilla CSS (Zero-utility bloat)
- **Design**: Emerald accent, Dark Mode, Glassmorphism

---

## 📦 Getting Started

### Prerequisites
- [Docker](https://www.docker.com/) & Docker Compose
- [Ollama](https://ollama.com/) (running locally)

### Quick Start
```bash
# Clone the repository
git clone https://github.com/promptops/engine.git

# Install dependencies and start (using Makefile)
make install-deps
make backend-run & make frontend-dev

# Alternatively, using Docker
docker compose up --build
```

### Useful Commands (Makefile)
| Command | Result |
| :--- | :--- |
| `make backend-test` | Run backend tests (Ginkgo) |
| `make backend-run` | Start the Go API server |
| `make frontend-dev` | Start the Next.js dev server |
| `make install-deps` | Install all dependencies |

### Pulling a Model (Required for first run)
```bash
docker compose exec ollama ollama pull tinyllama
```

---

## 📄 Documentation
- 📖 **[Developer Guide](docs/DEV_GUIDE.md)**
- 🗺️ **[10-Week Roadmap](docs/ROADMAP.md)**
- 📝 **[Weekly Progress Log](docs/WEEKLY_LOG.md)**

---

## 🛡️ License
Distributed under the MIT License. See `LICENSE` for more information.
