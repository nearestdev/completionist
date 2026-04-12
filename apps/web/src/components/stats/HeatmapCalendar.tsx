"use client";

import { HeatmapDay } from "@/types/stats";

const intensityClass = (count: number): string => {
  if (count === 0) return "bg-muted/10";
  if (count <= 2) return "bg-primary/20";
  if (count <= 5) return "bg-primary/40";
  if (count <= 10) return "bg-primary/60";
  return "bg-primary/80";
};

export default function HeatmapCalendar({ data }: { data: HeatmapDay[] }) {
  const dayMap = new Map(data.map((d) => [d.date, d.count]));
  const today = new Date();
  const startDate = new Date(today);
  startDate.setFullYear(startDate.getFullYear() - 1);

  const days: { date: string; count: number }[] = [];
  const current = new Date(startDate);
  while (current <= today) {
    const dateStr = current.toISOString().split("T")[0];
    days.push({ date: dateStr, count: dayMap.get(dateStr) || 0 });
    current.setDate(current.getDate() + 1);
  }

  return (
    <div>
      <div className="flex flex-wrap gap-[2px]">
        {days.map((day) => (
          <div
            key={day.date}
            title={`${day.date}: ${day.count} completions`}
            className={`w-3 h-3 rounded-sm ${intensityClass(day.count)}`}
          />
        ))}
      </div>
      <div className="flex items-center gap-2 mt-2 text-xs text-muted">
        <span>Less</span>
        <div className="w-3 h-3 rounded-sm bg-muted/10" />
        <div className="w-3 h-3 rounded-sm bg-primary/20" />
        <div className="w-3 h-3 rounded-sm bg-primary/40" />
        <div className="w-3 h-3 rounded-sm bg-primary/60" />
        <div className="w-3 h-3 rounded-sm bg-primary/80" />
        <span>More</span>
      </div>
    </div>
  );
}
