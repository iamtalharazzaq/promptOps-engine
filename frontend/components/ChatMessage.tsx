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
  /** Optional validation status for Schema Guard */
  status?: string;
  /** Current retry attempt if validation failed */
  retryCount?: number;
}

/**
 * ChatMessage renders a single chat bubble.
 */
export default function ChatMessage({
  role,
  content,
  isStreaming = false,
  status,
  retryCount,
}: ChatMessageProps) {
  const getStatusDisplay = () => {
    switch (status) {
      case "validating":
        return { label: "Validating...", icon: "🔍", class: "validating" };
      case "valid":
        return { label: "Schema Valid", icon: "✅", class: "valid" };
      case "retrying":
        return { label: `Retrying (${retryCount}/3)...`, icon: "🔄", class: "retrying" };
      case "invalid":
        return { label: "Validation Failed", icon: "❌", class: "invalid" };
      default:
        return null;
    }
  };

  const statusInfo = getStatusDisplay();

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

        {/* Validation Status Indicator */}
        {role === "assistant" && statusInfo && (
          <div className="validation-status">
            <div className={`status-indicator ${statusInfo.class}`}>
              <span className="status-icon">{statusInfo.icon}</span>
              <span className="status-text">{statusInfo.label}</span>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
