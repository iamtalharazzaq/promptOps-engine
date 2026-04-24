"use client";

import { useRef, useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import ChatMessage from "@/components/ChatMessage";
import ChatInput from "@/components/ChatInput";
import Sidebar from "@/components/Sidebar";
import { streamChat, checkHealth, HealthResponse } from "@/lib/api";
import { Message, updateMessagesWithEvent } from "@/lib/utils/chat";
import { useAuth } from "@/context/AuthContext";

export default function Home() {
  const { token, isAuthenticated } = useAuth();
  const router = useRouter();

  // ── State ──────────────────────────────────────────────────────
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [isStreaming, setIsStreaming] = useState(false);
  const [health, setHealth] = useState<HealthResponse | null>(null);
  const [selectedModel, setSelectedModel] = useState("tinyllama");
  const [isStructuredMode, setIsStructuredMode] = useState(false);
  const [currentChatId, setCurrentChatId] = useState<string | null>(null);

  const availableModels = [
    { id: "tinyllama", name: "TinyLlama (Fast)" },
    { id: "llama3", name: "Llama 3" },
    { id: "mistral", name: "Mistral" },
    { id: "phi3", name: "Phi-3 Mini" }
  ];

  const messagesEndRef = useRef<HTMLDivElement>(null);

  // ── Auth Guard ──────────────────────────────────────────────
  useEffect(() => {
    if (!token && typeof window !== 'undefined') {
       const savedToken = localStorage.getItem('auth_token');
       if (!savedToken) {
         router.push('/login');
       }
    }
  }, [token, router]);

  // ── Fetch health data on mount ─────────────────────────────
  useEffect(() => {
    checkHealth()
      .then(setHealth)
      .catch((err) => console.error("Health check failed:", err));
  }, []);

  // ── Auto-scroll ───────────────────────────────────────────────
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  const handleSelectChat = (id: string) => {
    if (id === "") {
      setMessages([]);
      setCurrentChatId(null);
    } else {
      setCurrentChatId(id);
      // Fetch specific chat history if needed
      setMessages([{ role: "assistant", content: `Retrieved sequence: ${id.slice(0, 8)}... (Full history loading implemented in next phase)` }]);
    }
  };

  const handleSend = async () => {
    const userMessage = input.trim();
    if (!userMessage || isStreaming || !token) return;

    const newMessages: Message[] = [
      ...messages,
      { role: "user", content: userMessage },
      { role: "assistant", content: "" },
    ];
    setMessages(newMessages);
    setInput("");
    setIsStreaming(true);

    const schema = isStructuredMode ? `{
      "type": "object",
      "properties": {
        "name": { "type": "string" },
        "role": { "type": "string" }
      },
      "required": ["name", "role"]
    }` : undefined;

    await streamChat(userMessage, {
      onChunk: (event) => {
        if (event.chat_id && !currentChatId) {
          setCurrentChatId(event.chat_id);
        }
        setMessages((prev) => updateMessagesWithEvent(prev, event));
      },
      onDone: () => setIsStreaming(false),
      onError: (error) => {
        setMessages((prev) => {
          const updated = [...prev];
          updated[updated.length - 1] = { ...updated[updated.length - 1], content: `⚠️ Error: ${error.message}` };
          return updated;
        });
        setIsStreaming(false);
      },
    }, {
      model: selectedModel,
      schema,
      token,
      chat_id: currentChatId || undefined
    });
  };

  if (!isAuthenticated && !token) return null;

  return (
    <div className="flex h-screen bg-[#0a0a0f] overflow-hidden">
      <Sidebar onSelectChat={handleSelectChat} />

      <main className="flex-1 flex flex-col relative">
        {/* Header */}
        <header className="h-16 border-b border-emerald-500/10 flex items-center justify-between px-6 bg-[#0a0a0f]/80 backdrop-blur-md z-20">
          <div className="flex items-center gap-4">
            <h1 className="text-zinc-100 font-bold tracking-tight">
              PromptOps <span className="text-emerald-500 uppercase text-[10px] tracking-[0.2em] font-medium ml-1">Kernel</span>
            </h1>
            <div className="w-[1px] h-4 bg-zinc-800" />
            <select 
              className="bg-transparent text-xs text-zinc-400 font-mono focus:outline-none cursor-pointer hover:text-emerald-400 transition-colors"
              value={selectedModel}
              onChange={(e) => setSelectedModel(e.target.value)}
              disabled={isStreaming}
            >
              {availableModels.map(m => (
                <option key={m.id} value={m.id} className="bg-[#0a0a0f]">{m.name.toUpperCase()}</option>
              ))}
            </select>
          </div>
          
          <div className="flex items-center gap-4">
             <button 
                onClick={() => setIsStructuredMode(!isStructuredMode)}
                className={`flex items-center gap-2 px-3 py-1.5 rounded-lg border transition-all text-[10px] font-mono uppercase tracking-widest ${
                  isStructuredMode 
                  ? 'border-emerald-500/50 bg-emerald-500/10 text-emerald-400 shadow-[0_0_10px_rgba(16,185,129,0.2)]' 
                  : 'border-zinc-800 text-zinc-500 hover:border-zinc-700'
                }`}
              >
                <div className={`w-1.5 h-1.5 rounded-full ${isStructuredMode ? 'bg-emerald-500 animate-pulse' : 'bg-zinc-800'}`} />
                Schema Guard
              </button>
              
              {health && (
                <div className="flex items-center gap-2 px-2 py-1 bg-emerald-500/5 rounded border border-emerald-500/10">
                  <div className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                  <span className="text-[10px] font-mono text-emerald-500/80">{health.version}</span>
                </div>
              )}
          </div>
        </header>

        {/* Messages area */}
        <div className="flex-1 overflow-y-auto custom-scrollbar p-6 scroll-smooth">
          <div className="max-w-3xl mx-auto space-y-8 pb-32">
            {messages.length === 0 ? (
              <div className="mt-20 text-center space-y-6">
                <div className="w-16 h-16 bg-emerald-500/10 rounded-2xl border border-emerald-500/20 flex items-center justify-center mx-auto shadow-[0_0_30px_rgba(16,185,129,0.05)]">
                  <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" className="text-emerald-500/60">
                    <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                  </svg>
                </div>
                <div className="space-y-2">
                  <h2 className="text-xl font-semibold text-zinc-100">Initialize Prompt Sequence</h2>
                  <p className="text-zinc-500 text-sm max-w-sm mx-auto">
                    Awaiting instructions for the PromptOps Engine. Secure session established.
                  </p>
                </div>
                <div className="grid grid-cols-2 gap-3 max-w-lg mx-auto">
                  {["System Architecture Audit", "Dockerize Go Backend", "Deploy to Kubernetes", "Optimize CI/CD Pipe"].map(t => (
                    <button 
                      key={t}
                      onClick={() => setInput(t)}
                      className="p-4 bg-[#12121e] border border-emerald-500/5 rounded-xl text-left text-xs text-zinc-400 hover:border-emerald-500/30 hover:bg-emerald-500/[0.02] transition-all group"
                    >
                      <div className="text-emerald-500/40 mb-2 truncate group-hover:text-emerald-500/60 font-mono tracking-tighter">CMD_INPUT</div>
                      {t}
                    </button>
                  ))}
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
        </div>

        {/* Input area */}
        <div className="absolute bottom-0 left-0 right-0 p-6 bg-gradient-to-t from-[#0a0a0f] via-[#0a0a0f] to-transparent">
          <div className="max-w-3xl mx-auto">
            <ChatInput
              value={input}
              onChange={setInput}
              onSend={handleSend}
              disabled={isStreaming}
            />
            <p className="text-[10px] text-zinc-600 text-center mt-4 font-mono tracking-widest uppercase">
              Engine Status: <span className="text-emerald-500/60">Operational</span> • Latency: <span className="text-emerald-500/60">12ms</span> • Context: <span className="text-emerald-500/60">Secure</span>
            </p>
          </div>
        </div>
      </main>
    </div>
  );
}

