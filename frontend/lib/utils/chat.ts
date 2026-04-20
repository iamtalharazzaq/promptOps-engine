import { ChatEvent } from "../api";

/** A single message in the conversation. */
export interface Message {
  role: "user" | "assistant";
  content: string;
  status?: string;
  retryCount?: number;
}

/**
 * Updates an array of messages based on a new ChatEvent.
 * 
 * Handles different event types:
 * - content: Appends text (or replaces if status is valid/invalid)
 * - status: Updates validation status
 * - retryCount: Updates retry attempt count
 * - retrying: Clears content for the next attempt
 */
export function updateMessagesWithEvent(
  prevMessages: Message[],
  event: ChatEvent
): Message[] {
  const updated = [...prevMessages];
  const lastIdx = updated.length - 1;

  if (lastIdx < 0) return updated;

  // Clear content on retry to show fresh attempt
  if (event.status === "retrying") {
    updated[lastIdx] = {
      ...updated[lastIdx],
      content: "",
      status: event.status,
      retryCount: event.retryCount
    };
    return updated;
  }

  if (event.content) {
    // In schema mode, the final valid content comes with status="valid"
    if (event.status === "valid" || event.status === "invalid") {
      updated[lastIdx] = {
        ...updated[lastIdx],
        content: event.content
      };
    } else if (!event.status) {
      // Normal streaming mode
      updated[lastIdx] = {
        ...updated[lastIdx],
        content: updated[lastIdx].content + event.content
      };
    }
  }

  if (event.status) {
    updated[lastIdx].status = event.status;
  }
  if (event.retryCount !== undefined) {
    updated[lastIdx].retryCount = event.retryCount;
  }

  return updated;
}
