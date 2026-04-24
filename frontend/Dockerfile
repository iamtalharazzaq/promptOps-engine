# ═══════════════════════════════════════════════════════════════
#  PromptOps Engine — Next.js Frontend Dockerfile
# ═══════════════════════════════════════════════════════════════
#
#  Multi-stage build:
#    Stage 1 (deps):    installs node_modules
#    Stage 2 (builder): runs `next build` to produce .next output
#    Stage 3 (runner):  serves the built app with `next start`
#
#  Build:   docker build -t promptops-frontend ./frontend
#  Run:     docker run -p 3000:3000 promptops-frontend
# ═══════════════════════════════════════════════════════════════

# ── Stage 1: Dependencies ──────────────────────────────────────
FROM node:20-alpine AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci

# ── Stage 2: Build ─────────────────────────────────────────────
FROM node:20-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Backend API URL is baked into the build
ENV NEXT_PUBLIC_API_URL=http://localhost:8080
RUN npm run build

# ── Stage 3: Runtime ───────────────────────────────────────────
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs
EXPOSE 3000
ENV PORT=3000

CMD ["node", "server.js"]
