"use client";

import { BadgeWithEarned } from "@/types/badge";

export default function BadgeCard({ badge }: { badge: BadgeWithEarned }) {
  const earned = !!badge.earnedAt;

  return (
    <div
      className={`flex flex-col items-center gap-2 p-3 rounded-xl border transition-all ${
        earned
          ? "border-primary/30 bg-primary/5 shadow-sm"
          : "border-border bg-card opacity-40 grayscale"
      }`}
    >
      <div className={`w-12 h-12 rounded-full flex items-center justify-center text-2xl ${
        earned ? "bg-gradient-to-br from-primary/20 to-accent/20" : "bg-muted/20"
      }`}>
        {badge.icon ? "🏆" : "🔒"}
      </div>
      <span className="text-xs font-semibold text-center text-foreground">{badge.name}</span>
      {badge.description && (
        <span className="text-[10px] text-muted text-center leading-tight">{badge.description}</span>
      )}
    </div>
  );
}
