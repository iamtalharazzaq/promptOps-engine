/**
 * @fileoverview Root layout for PromptOps Engine frontend.
 *
 * Sets up:
 * - Inter font via next/font/google
 * - Dark theme meta colour
 * - SEO meta tags (title, description)
 * - Global CSS import
 *
 * This layout wraps every page in the app and provides consistent
 * styling and metadata across the entire application.
 *
 * @module app/layout
 */

import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

/** Inter font with Latin subset for clean, modern typography. */
const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "PromptOps Engine",
  description:
    "LLM Orchestration Platform with Schema Validation, Metrics & CI/CD — built with Go, Next.js, and Ollama.",
};

/**
 * RootLayout wraps every page in the app.
 *
 * Responsibilities:
 * - Applies the Inter font via className
 * - Sets the color-scheme to dark (respects OS preference)
 * - Imports global CSS (design tokens, resets, component styles)
 */
export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head>
        <meta name="theme-color" content="#0a0a0f" />
      </head>
      <body className={inter.className}>{children}</body>
    </html>
  );
}
