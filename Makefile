# ───────────────────────────────────────────────────────────────────────
#  PromptOps Engine — Makefile
# ───────────────────────────────────────────────────────────────────────

.PHONY: help install-deps backend-test backend-run frontend-dev clean

help:
	@echo "PromptOps Engine CLI"
	@echo "Usage:"
	@echo "  make install-deps   - Install backend and frontend dependencies"
	@echo "  make backend-test   - Run backend tests (Ginkgo)"
	@echo "  make backend-run    - Run backend API server"
	@echo "  make frontend-dev   - Start frontend development server"
	@echo "  make clean          - Remove temporary files and build artifacts"

install-deps:
	@echo "Installing backend dependencies..."
	cd backend && go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

backend-test:
	@echo "Running backend tests..."
	cd backend && go test -v ./handlers/...

backend-run:
	@echo "Starting backend..."
	cd backend && go run main.go

frontend-dev:
	@echo "Starting frontend..."
	cd frontend && npm run dev

clean:
	@echo "Cleaning artifacts..."
	rm -rf backend/bin
	rm -rf frontend/.next
	rm -rf frontend/node_modules
