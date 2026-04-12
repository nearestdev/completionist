"use client";

import { useEffect, useState } from "react";
import { Flame, Fire, FireSimple } from "@phosphor-icons/react";
import streakService from "@/services/streakService";
import { UserStreak, StreakTier } from "@/types/streak";

const tierConfig: Record<StreakTier, { icon: typeof Flame; color: string; label: string }> = {
  none: { icon: Flame, color: "text-muted", label: "" },
  spark: { icon: Flame, color: "text-orange-400", label: "Spark" },
  fire: { icon: Fire, color: "text-orange-500", label: "Fire" },
  inferno: { icon: FireSimple, color: "text-red-500", label: "Inferno" },
};

export default function StreakPill({ userId }: { userId: number }) {
  const [streak, setStreak] = useState<UserStreak | null>(null);

  useEffect(() => {
    streakService.getUserStreak(userId).then(setStreak).catch(() => {});
  }, [userId]);

  if (!streak || streak.currentStreak === 0) return null;

  const config = tierConfig[streak.tier] || tierConfig.spark;
  const Icon = config.icon;

  return (
    <div className={`flex items-center gap-1 px-2 py-1 rounded-full bg-card border border-border ${config.color}`}>
      <Icon size={16} weight="fill" />
      <span className="text-xs font-bold">{streak.currentStreak}</span>
    </div>
  );
}
