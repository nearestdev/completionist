"use client";

import Widgets from "@/components/layout/Widgets";
import { PlusIcon } from "@phosphor-icons/react";

const MockFeed = () => (
  <div className="flex flex-col gap-6">
    <div className="flex items-center justify-between mb-2">
      <h2 className="font-heading font-bold text-2xl text-foreground">Your Feed</h2>
      <button className="flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium text-white bg-gradient-to-r from-primary to-primary-hover hover:opacity-90 transition-all shadow-sm hover:shadow-md">
        <PlusIcon size={18} weight="bold" />
        <span>Create Post</span>
      </button>
    </div>
    <div 
      className="bg-card p-6 rounded-xl border border-border transition-all duration-300 hover:shadow-lg" 
      style={{ boxShadow: 'var(--shadow-sm)' }}
    >
      <h3 className="font-heading font-bold text-lg mb-2">Welcome to your Feed</h3>
      <p className="text-muted">Start tracking your progress and see what your friends are up to!</p>
    </div>
  </div>
);

export default function HomePage() {
  return (
    <div className="flex flex-col lg:flex-row justify-center w-full">
      <main className="w-full max-w-[680px] px-4 py-8 mb-20 lg:mb-0 min-h-screen">
         <MockFeed />
      </main>

      <Widgets />
    </div>
  );
}