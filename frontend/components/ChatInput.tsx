/**
 * @fileoverview ChatInput component for PromptOps Engine.
 *
 * Provides a styled text input and send button for composing messages.
 * Supports keyboard submission (Enter key) and is automatically disabled
 * while an assistant response is streaming to prevent concurrent requests.
 *
 * The component is a controlled input — the parent manages the value via
 * `value` and `onChange` props.
 *
 * @module components/ChatInput
 */

"use client";

import React, { FormEvent, KeyboardEvent } from "react";

/** Props accepted by the ChatInput component. */
interface ChatInputProps {
  /** Current input value (controlled) */
  value: string;
  /** Called when the input value changes */
  onChange: (value: string) => void;
  /** Called when the user submits the message (Enter or button click) */
  onSend: () => void;
  /** When true, the input and button are disabled (streaming in progress) */
  disabled?: boolean;
}

/**
 * ChatInput renders the message composition area at the bottom of the chat.
 *
 * Features:
 * - Enter key sends the message (Shift+Enter for newline — future)
 * - Send button with arrow icon
 * - Visual disabled state during streaming
 * - Focus ring with accent color
 *
 * @example
 * ```tsx
 * <ChatInput
 *   value={input}
 *   onChange={setInput}
 *   onSend={handleSend}
 *   disabled={isStreaming}
 * />
 * ```
 */
export default function ChatInput({
  value,
  onChange,
  onSend,
  disabled = false,
}: ChatInputProps) {
  /** Handle form submission (prevents page reload). */
  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (value.trim() && !disabled) {
      onSend();
    }
  };

  /** Send on Enter (without Shift). */
  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      if (value.trim() && !disabled) {
        onSend();
      }
    }
  };

  return (
    <form id="chat-input-form" className="chat-input-container" onSubmit={handleSubmit}>
      <input
        id="chat-input"
        type="text"
        className="chat-input"
        placeholder="Ask PromptOps Engine anything..."
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onKeyDown={handleKeyDown}
        disabled={disabled}
        autoComplete="off"
        autoFocus
      />
      <button
        id="send-button"
        type="submit"
        className="send-button"
        disabled={disabled || !value.trim()}
        aria-label="Send message"
      >
        <svg
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <line x1="22" y1="2" x2="11" y2="13" />
          <polygon points="22 2 15 22 11 13 2 9 22 2" />
        </svg>
      </button>
    </form>
  );
}
