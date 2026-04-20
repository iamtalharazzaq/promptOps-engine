/**
 * @fileoverview Main chat page for PromptOps Engine.
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
import { Message, updateMessagesWithEvent } from "@/lib/utils/chat";



export default function Home() {
  // ── State ──────────────────────────────────────────────────────
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [isStreaming, setIsStreaming] = useState(false);
  const [health, setHealth] = useState<HealthResponse | null>(null);
  const [selectedModel, setSelectedModel] = useState("tinyllama");
  const [isStructuredMode, setIsStructuredMode] = useState(false);

  const availableModels = [
    { id: "tinyllama", name: "TinyLlama (Fast)" },
    { id: "llama3", name: "Llama 3 (Smart)" },
    { id: "mistral", name: "Mistral 7B" },
    { id: "phi3", name: "Phi-3 Mini" }
  ];

  const sampleSchema = `{
    "type": "object",
    "properties": {
      "name": { "type": "string" },
      "role": { "type": "string" },
      "skills": { 
        "type": "array",
        "items": { "type": "string" }
      }
    },
    "required": ["name", "role", "skills"]
  }`;

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

    const schema = isStructuredMode ? sampleSchema : undefined;

    // Stream response from backend
    await streamChat(userMessage, {
      onChunk: (event) => {
        setMessages((prev) => updateMessagesWithEvent(prev, event));
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
    }, selectedModel, schema);
  };

  // ── Render ────────────────────────────────────────────────────
  return (
    <main id="chat-page" className="chat-container">
      {/* Header */}
      <header className="chat-header">
        <div className="header-content">
          <div className="logo-group">
            <div className="logo-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
                <path d="M12 8v8" />
                <path d="M8 12h8" />
              </svg>
            </div>
            <h1 className="logo-text">PromptOps <span className="logo-lite">Engine</span></h1>
          </div>
          
          <div className="header-controls">
            <div className="toggle-group">
              <button 
                className={`structured-toggle ${isStructuredMode ? 'active' : ''}`}
                onClick={() => setIsStructuredMode(!isStructuredMode)}
                title="Structured Output Mode (Schema Guard)"
                disabled={isStreaming}
              >
                <span className="toggle-icon">🛡️</span>
                <span className="toggle-label">Schema Guard</span>
              </button>
            </div>

            <select 
              className="model-selector"
              value={selectedModel}
              onChange={(e) => setSelectedModel(e.target.value)}
              disabled={isStreaming}
            >
              {availableModels.map(m => (
                <option key={m.id} value={m.id}>{m.name}</option>
              ))}
            </select>

            <div className="header-badges">
              {health && (
                <div className="token-badge">{health.maxTokens}T</div>
              )}
              <div className="header-badge">{health?.version || "OSS"}</div>
            </div>
          </div>
        </div>
      </header>

      {/* Messages area */}
      <div id="messages-container" className="messages-area">
        {messages.length === 0 ? (
          <div className="empty-state">
            <div className="empty-icon">⚙️</div>
            <h2>Welcome to PromptOps Engine</h2>
            <p>High-performance LLM orchestration with schema validation.</p>
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
              status={msg.status}
              retryCount={msg.retryCount}
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
        <div className="disclaimer-container">
          <svg className="ollama-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
          </svg>
          <p className="disclaimer">
            Powered by <b>Ollama</b> • Trusted Schema Validation
          </p>
        </div>
      </div>
    </main>
  );
}
