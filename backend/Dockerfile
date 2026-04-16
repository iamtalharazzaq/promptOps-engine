# ═══════════════════════════════════════════════════════════════
#  PromptOps Engine — Go Backend Dockerfile
# ═══════════════════════════════════════════════════════════════
#
#  Multi-stage build:
#    Stage 1 (builder): compiles the Go binary with CGO disabled
#    Stage 2 (runtime): copies the ~15MB static binary into a
#                       minimal alpine image for production
#
#  Build:   docker build -t promptops-backend ./backend
#  Run:     docker run -p 8080:8080 promptops-backend
# ═══════════════════════════════════════════════════════════════

# ── Stage 1: Build ──────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache module downloads separately from source code
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /promptops-backend main.go

# ── Stage 2: Runtime ────────────────────────────────────────────
FROM alpine:3.19

# Install CA certs for HTTPS calls to Ollama (if remote)
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /promptops-backend .

# Default env vars (overridable via docker-compose or -e flags)
ENV PORT=8080
ENV OLLAMA_HOST=http://ollama:11434
ENV OLLAMA_MODEL=tinyllama
ENV FRONTEND_URL=http://localhost:3000

EXPOSE 8080

# Health check: hit /health every 30s
HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
  CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["./promptops-backend"]
