"use client";

import React, { useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { login as loginApi } from '@/lib/api';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const { token } = await loginApi(email, password);
      login(token);
      router.push('/');
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#0a0a0f] flex items-center justify-center p-6 relative overflow-hidden">
      {/* Background Glow */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[500px] h-[500px] bg-emerald-500/10 rounded-full blur-[120px] pointer-events-none" />

      <div className="w-full max-w-md bg-[#12121e]/80 border border-emerald-500/20 rounded-2xl p-8 backdrop-blur-xl shadow-2xl relative z-10">
        <div className="flex flex-col items-center mb-10">
          <div className="w-16 h-16 bg-emerald-500/20 rounded-xl flex items-center justify-center border border-emerald-500/30 mb-4 group transition-all duration-500 hover:border-emerald-500/60 shadow-[0_0_20px_rgba(16,185,129,0.1)]">
             <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-emerald-400">
               <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
             </svg>
          </div>
          <h1 className="text-2xl font-bold text-white tracking-tight">Access Control</h1>
          <p className="text-zinc-400 text-sm mt-2">Initialize session to PromptOps Engine</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="space-y-1">
            <label className="text-xs font-medium text-emerald-500/70 uppercase tracking-widest ml-1">Identity</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full bg-[#0a0a0f]/50 border border-emerald-500/10 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-emerald-500/50 transition-all placeholder:text-zinc-600"
              placeholder="operator@promptops.sh"
              required
            />
          </div>

          <div className="space-y-1">
            <label className="text-xs font-medium text-emerald-500/70 uppercase tracking-widest ml-1">Access Key</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full bg-[#0a0a0f]/50 border border-emerald-500/10 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-emerald-500/50 transition-all placeholder:text-zinc-600"
              placeholder="••••••••"
              required
            />
          </div>

          {error && (
            <div className="bg-red-500/10 border border-red-500/20 text-red-400 text-sm p-3 rounded-lg flex items-center gap-2">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
              {error}
            </div>
          )}

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-emerald-500 hover:bg-emerald-400 text-black font-bold py-3 rounded-xl transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed shadow-[0_0_20px_rgba(16,185,129,0.2)] active:scale-[0.98]"
          >
            {loading ? 'Authenticating...' : 'Authorize Session'}
          </button>
        </form>

        <div className="mt-8 text-center">
          <p className="text-zinc-500 text-sm">
            Unregistered entity?{' '}
            <Link href="/register" className="text-emerald-400 hover:text-emerald-300 transition-colors font-medium">
              Create credentials
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
