"use client";

import { useEffect, useState } from "react";
import { useAuth } from "@/hooks/useAuth";
import statsService from "@/services/statsService";
import { UserStats, HeatmapDay } from "@/types/stats";
import TimeDonut from "@/components/stats/TimeDonut";
import GenreRadar from "@/components/stats/GenreRadar";
import HeatmapCalendar from "@/components/stats/HeatmapCalendar";

export default function StatsPage() {
  const { user } = useAuth();
  const [stats, setStats] = useState<UserStats | null>(null);
  const [heatmap, setHeatmap] = useState<HeatmapDay[]>([]);

  useEffect(() => {
    if (!user) return;
    statsService.getUserStats(user.id).then((data) => setStats(data || null)).catch(() => {});
    statsService.getUserHeatmap(user.id).then((data) => setHeatmap(data || [])).catch(() => {});
  }, [user]);

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <h1 className="text-3xl font-bold text-foreground mb-8">My Stats</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <div className="bg-card rounded-xl border border-border p-6">
          <h2 className="text-lg font-semibold text-foreground mb-4">Time Breakdown</h2>
          <TimeDonut data={stats?.timeDonut || []} />
        </div>
        <div className="bg-card rounded-xl border border-border p-6">
          <h2 className="text-lg font-semibold text-foreground mb-4">Genre DNA</h2>
          <GenreRadar data={stats?.genres || []} />
        </div>
      </div>

      <div className="bg-card rounded-xl border border-border p-6">
        <h2 className="text-lg font-semibold text-foreground mb-4">Activity</h2>
        <HeatmapCalendar data={heatmap} />
      </div>
    </div>
  );
}
