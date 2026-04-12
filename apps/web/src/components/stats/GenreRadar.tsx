"use client";

import { RadarChart, PolarGrid, PolarAngleAxis, PolarRadiusAxis, Radar, ResponsiveContainer, Tooltip } from "recharts";
import { GenreEntry } from "@/types/stats";

export default function GenreRadar({ data }: { data: GenreEntry[] }) {
  if (!data || data.length === 0) {
    return <p className="text-muted text-sm">No genre data yet.</p>;
  }

  return (
    <div className="w-full h-64">
      <ResponsiveContainer>
        <RadarChart data={data}>
          <PolarGrid stroke="var(--border)" />
          <PolarAngleAxis dataKey="genre" tick={{ fill: "var(--muted)", fontSize: 11 }} />
          <PolarRadiusAxis tick={false} axisLine={false} />
          <Radar dataKey="count" stroke="#4F46E5" fill="#4F46E5" fillOpacity={0.3} />
          <Tooltip />
        </RadarChart>
      </ResponsiveContainer>
    </div>
  );
}
