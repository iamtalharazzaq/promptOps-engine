# ═══════════════════════════════════════════════════════════════
#  GhostAI Lite — Next.js Frontend Dockerfile
# ═══════════════════════════════════════════════════════════════
#
#  Multi-stage build:
#    Stage 1 (deps):    installs node_modules
#    Stage 2 (builder): runs `next build` to produce .next output
#    Stage 3 (runner):  serves the built app with `next start`
#
#  Build:   docker build -t ghostai-frontend ./frontend
#  Run:     docker run -p 3000:3000 ghostai-frontend
# ═══════════════════════════════════════════════════════════════

# ── Stage 1: Dependencies ──────────────────────────────────────
FROM node:20-alpine AS deps

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --only=production=false

# ── Stage 2: Build ─────────────────────────────────────────────
FROM node:20-alpine AS builder

WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Backend API URL is baked into the build (Next.js inlines NEXT_PUBLIC_ vars)
ENV NEXT_PUBLIC_API_URL=http://localhost:8080

RUN npm run build

# ── Stage 3: Runtime ───────────────────────────────────────────
FROM node:20-alpine AS runner

WORKDIR /app
ENV NODE_ENV=production

# Copy only what's needed to run
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

EXPOSE 3000

CMD ["npm", "start"]
