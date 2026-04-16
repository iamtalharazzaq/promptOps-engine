/**
 * @fileoverview API client for communicating with the GhostAI Lite backend.
 *
 * This module provides a streaming-aware fetch wrapper that consumes
 * Server-Sent Events (SSE) from the Go backend's /chat endpoint and
 * delivers tokens to the UI one at a time for the "typing" effect.
 *
 * @module lib/api
 */

/** Base URL of the Go backend API. Falls back to localhost:8080 in dev. */
const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

/**
 * Shape of a single SSE chat event from the backend.
 *
 * @property content - A text fragment (usually a single token)
 * @property done    - `true` on the final event, signalling end of generation
 */
export interface ChatEvent {
  content: string;
  done: boolean;
}

/**
 * Shape of the /health endpoint response.
 *
 * @property status    - Always "ok" when the backend is reachable
 * @property timestamp - ISO-8601 server time
 * @property service   - Backend service name
 * @property version   - Semantic version string
 */
export interface HealthResponse {
  status: string;
  timestamp: string;
  service: string;
  version: string;
}

/**
 * Check backend health.
 *
 * @returns The parsed HealthResponse from GET /health
 * @throws  If the backend is unreachable or returns non-200
 */
export async function checkHealth(): Promise<HealthResponse> {
  const res = await fetch(`${API_BASE}/health`);
  if (!res.ok) throw new Error(`Health check failed: ${res.status}`);
  return res.json();
}

/**
 * Stream a chat response from the backend using Server-Sent Events.
 *
 * This function sends the user's message to POST /chat and reads the
 * streaming SSE response using the Fetch API's ReadableStream. Each
 * `data:` line is parsed as JSON and delivered via the `onChunk` callback.
 *
 * ### SSE Protocol
 * The backend sends lines in this format:
 * ```
 * data: {"content":"Hello","done":false}\n\n
 * data: {"content":"","done":true}\n\n
 * ```
 *
 * ### Usage
 * ```ts
 * await streamChat("Hello!", {
 *   onChunk: (ev) => appendToMessage(ev.content),
 *   onDone:  ()   => setStreaming(false),
 *   onError: (e)  => showError(e.message),
 * });
 * ```
 *
 * @param message  - The user's message to send
 * @param handlers - Callback functions for stream events
 * @param model    - Optional model override (defaults to backend config)
 */
export async function streamChat(
  message: string,
  handlers: {
    onChunk: (event: ChatEvent) => void;
    onDone: () => void;
    onError: (error: Error) => void;
  },
  model?: string
): Promise<void> {
  try {
    const res = await fetch(`${API_BASE}/chat`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ message, ...(model ? { model } : {}) }),
    });

    if (!res.ok) {
      throw new Error(`Chat request failed: ${res.status}`);
    }

    const reader = res.body?.getReader();
    if (!reader) {
      throw new Error("ReadableStream not supported");
    }

    const decoder = new TextDecoder();
    let buffer = "";

    // Read chunks from the stream until done
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });

      // Process complete SSE lines from the buffer
      const lines = buffer.split("\n");
      buffer = lines.pop() || ""; // Keep incomplete line in buffer

      for (const line of lines) {
        const trimmed = line.trim();
        if (!trimmed || !trimmed.startsWith("data: ")) continue;

        try {
          const json = trimmed.slice(6); // Remove "data: " prefix
          const event: ChatEvent = JSON.parse(json);
          handlers.onChunk(event);

          if (event.done) {
            handlers.onDone();
            return;
          }
        } catch {
          // Skip malformed lines (keep-alive comments, etc.)
        }
      }
    }

    handlers.onDone();
  } catch (error) {
    handlers.onError(
      error instanceof Error ? error : new Error(String(error))
    );
  }
}
