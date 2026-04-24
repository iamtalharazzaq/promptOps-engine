"use client";

import React, { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { getHistory } from '@/lib/api';

interface ChatHistoryItem {
  id: string;
  title: string;
  updated_at: string;
}

export default function Sidebar({ onSelectChat }: { onSelectChat: (id: string) => void }) {
  const [history, setHistory] = useState<ChatHistoryItem[]>([]);
  const { token, logout } = useAuth();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (token) {
      getHistory(token)
        .then(setHistory)
        .catch(console.error)
        .finally(() => setLoading(false));
    }
  }, [token]);

  return (
    <div className="w-72 bg-[#0a0a0f] border-r border-emerald-500/10 flex flex-col h-full overflow-hidden">
      <div className="p-4 border-b border-emerald-500/10 mb-2">
        <button 
          onClick={() => onSelectChat('')}
          className="w-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 hover:bg-emerald-500/20 py-2.5 rounded-xl flex items-center justify-center gap-2 transition-all font-medium text-sm"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
          New Sequence
        </button>
      </div>

      <div className="flex-1 overflow-y-auto px-3 space-y-1 custom-scrollbar">
        {loading ? (
          <div className="p-4 text-center text-zinc-600 text-xs animate-pulse font-mono uppercase tracking-widest">
            Scanning blocks...
          </div>
        ) : history.length === 0 ? (
          <div className="p-8 text-center text-zinc-600 text-xs font-mono italic">
            No history found.
          </div>
        ) : (
          history.map((chat) => (
            <button
              key={chat.id}
              onClick={() => onSelectChat(chat.id)}
              className="w-full group px-3 py-3 rounded-xl text-left hover:bg-emerald-500/5 transition-all border border-transparent hover:border-emerald-500/20"
            >
              <div className="text-zinc-100 text-sm font-medium truncate group-hover:text-emerald-400 transition-colors">
                {chat.title || "Untitled Sequence"}
              </div>
              <div className="text-[10px] text-zinc-600 mt-1 font-mono uppercase tracking-tighter">
                {new Date(chat.updated_at).toLocaleDateString()} • {chat.id.slice(0, 8)}
              </div>
            </button>
          ))
        )}
      </div>

      <div className="p-4 border-t border-emerald-500/10 bg-[#0d0d15]/50">
        <button 
          onClick={logout}
          className="w-full flex items-center justify-between px-3 py-2 rounded-lg hover:bg-red-500/10 text-zinc-500 hover:text-red-400 transition-all text-sm group"
        >
          <span className="font-medium">Terminate Session</span>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" className="opacity-50 group-hover:opacity-100">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" /><polyline points="16 17 21 12 16 7" /><line x1="21" y1="12" x2="9" y2="12" />
          </svg>
        </button>
      </div>
    </div>
  );
}
