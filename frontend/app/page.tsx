/**
 * @fileoverview Main chat page for GhostAI Lite.
 *
 * This is the home page ("/") and provides the full ChatGPT-style
 * conversation interface. It manages:
 *
 * - **Message state**: Array of { role, content } message objects
 * - **Streaming state**: Tracks whether a response is currently streaming
 * - **Auto-scroll**: Keeps the latest message visible as tokens arrive
 * - **Session management**: Messages persist in component state (future: localStorage)
 *
 * ### Data flow
 * 1. User types a message in ChatInput
 * 2. On submit, user message is appended to state
 * 3. streamChat() is called, which POSTs to the backend
 * 4. Each SSE chunk updates the last (assistant) message in-place
 * 5. On done, streaming state is cleared
 *
 * @module app/page
 */

"use client";

import { useRef, useState, useEffect } from "react";
import ChatMessage from "@/components/ChatMessage";
import ChatInput from "@/components/ChatInput";
import { streamChat, checkHealth, HealthResponse } from "@/lib/api";

/** A single message in the conversation. */
interface Message {
  role: "user" | "assistant";
  content: string;
}

export default function Home() {
  // ── State ──────────────────────────────────────────────────────
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [isStreaming, setIsStreaming] = useState(false);
  const [health, setHealth] = useState<HealthResponse | null>(null);

  /** Ref to the bottom of the message list for auto-scroll. */
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // ── Fetch health data on mount ─────────────────────────────
  useEffect(() => {
    checkHealth()
      .then(setHealth)
      .catch((err) => console.error("Health check failed:", err));
  }, []);

  // ── Auto-scroll on new messages or streaming updates ──────────
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  // ── Send handler ──────────────────────────────────────────────

  /**
   * Sends the current input to the backend and starts streaming
   * the assistant's response. Appends both the user message and
   * an empty assistant message (which gets filled token-by-token).
   */
  const handleSend = async () => {
    const userMessage = input.trim();
    if (!userMessage || isStreaming) return;

    // Append user message + empty assistant placeholder
    const newMessages: Message[] = [
      ...messages,
      { role: "user", content: userMessage },
      { role: "assistant", content: "" },
    ];
    setMessages(newMessages);
    setInput("");
    setIsStreaming(true);

    // Stream response from backend
    await streamChat(userMessage, {
      onChunk: (event) => {
        setMessages((prev) => {
          const updated = [...prev];
          const lastIdx = updated.length - 1;
          updated[lastIdx] = {
            ...updated[lastIdx],
            content: updated[lastIdx].content + event.content,
          };
          return updated;
        });
      },
      onDone: () => {
        setIsStreaming(false);
      },
      onError: (error) => {
        setMessages((prev) => {
          const updated = [...prev];
          const lastIdx = updated.length - 1;
          updated[lastIdx] = {
            ...updated[lastIdx],
            content: `⚠️ Error: ${error.message}`,
          };
          return updated;
        });
        setIsStreaming(false);
      },
    });
  };

  // ── Render ────────────────────────────────────────────────────
  return (
    <main id="chat-page" className="chat-container">
      {/* Header */}
      <header className="chat-header">
        <div className="header-content">
          <div className="logo-group">
            <div className="logo-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M9 10h.01" />
                <path d="M15 10h.01" />
                <path d="M12 2a8 8 0 0 0-8 8v12l3-3 2.5 2.5L12 19l2.5 2.5L17 19l3 3V10a8 8 0 0 0-8-8z" />
              </svg>
            </div>
            <h1 className="logo-text">GhostAI <span className="logo-lite">Lite</span></h1>
          </div>
          <div className="header-badges">
            {health && (
              <div className="token-badge">Limit: {health.maxTokens} tokens</div>
            )}
            <div className="header-badge">{health?.version || "v0.1.0"}</div>
          </div>
        </div>
      </header>

      {/* Messages area */}
      <div id="messages-container" className="messages-area">
        {messages.length === 0 ? (
          <div className="empty-state">
            <div className="empty-icon">👻</div>
            <h2>Welcome to GhostAI Lite</h2>
            <p>Start a conversation with your local AI assistant.</p>
            <div className="empty-hints">
              <button className="hint-chip" onClick={() => setInput("Explain how Go handles concurrency")}>
                Explain Go concurrency
              </button>
              <button className="hint-chip" onClick={() => setInput("Write a haiku about programming")}>
                Write a haiku
              </button>
              <button className="hint-chip" onClick={() => setInput("What is Docker and why should I use it?")}>
                What is Docker?
              </button>
            </div>
          </div>
        ) : (
          messages.map((msg, i) => (
            <ChatMessage
              key={i}
              role={msg.role}
              content={msg.content}
              isStreaming={isStreaming && i === messages.length - 1 && msg.role === "assistant"}
            />
          ))
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Input bar */}
      <div className="input-area">
        <ChatInput
          value={input}
          onChange={setInput}
          onSend={handleSend}
          disabled={isStreaming}
        />
        <p className="disclaimer">
          GhostAI Lite runs models locally via Ollama. Responses may vary.
        </p>
      </div>
    </main>
  );
}
