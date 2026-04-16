/**
 * @fileoverview ChatMessage component for PromptOps Engine.
 *
 * Renders a single message bubble in the conversation thread. Supports
 * two roles:
 *   - **user**:      Right-aligned bubble with accent gradient
 *   - **assistant**: Left-aligned bubble with glass-morphism styling
 *
 * When `isStreaming` is true, a blinking cursor (▌) is appended to
 * the assistant's message to indicate that tokens are still arriving.
 *
 * @module components/ChatMessage
 */

import React from "react";

/** Props accepted by the ChatMessage component. */
interface ChatMessageProps {
  /** The role determines bubble alignment and styling */
  role: "user" | "assistant";
  /** The text content of the message */
  content: string;
  /** When true, shows a blinking cursor after the content (assistant only) */
  isStreaming?: boolean;
}

/**
 * ChatMessage renders a single chat bubble.
 *
 * @example
 * ```tsx
 * <ChatMessage role="user" content="Hello, what is Go?" />
 * <ChatMessage role="assistant" content="Go is a..." isStreaming={true} />
 * ```
 */
export default function ChatMessage({
  role,
  content,
  isStreaming = false,
}: ChatMessageProps) {
  return (
    <div className={`message-row ${role}`}>
      {/* Avatar indicator */}
      <div className={`avatar ${role}`}>
        {role === "user" ? "ME" : "PS"}
      </div>

      {/* Message bubble */}
      <div className={`message-bubble ${role}`}>
        <p className="message-text">
          {content}
          {isStreaming && <span className="streaming-cursor">▌</span>}
        </p>
      </div>
    </div>
  );
}
