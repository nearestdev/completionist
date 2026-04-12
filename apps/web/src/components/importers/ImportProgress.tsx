"use client";

import { ImportJob } from "@/types/connected_account";

const statusColors: Record<string, string> = {
  pending: "bg-muted/10 text-muted",
  running: "bg-primary/10 text-primary",
  completed: "bg-accent/10 text-accent",
  failed: "bg-red-500/10 text-red-500",
  cancelled: "bg-muted/10 text-muted",
};

export default function ImportProgress({ job }: { job: ImportJob }) {
  const pct = job.totalItems > 0 ? Math.round((job.importedItems / job.totalItems) * 100) : 0;

  return (
    <div className="bg-card rounded-lg border border-border p-3">
      <div className="flex justify-between text-sm mb-1">
        <span className="font-medium text-foreground capitalize">{job.provider}</span>
        <span className={`text-xs px-2 py-0.5 rounded-full ${statusColors[job.status] || ""}`}>{job.status}</span>
      </div>
      {job.totalItems > 0 && (
        <div className="w-full h-1.5 bg-muted/20 rounded-full overflow-hidden mb-1">
          <div className="h-full bg-primary rounded-full transition-all" style={{ width: `${pct}%` }} />
        </div>
      )}
      <p className="text-xs text-muted">
        {job.importedItems} imported &middot; {job.skippedItems} skipped &middot; {job.failedItems} failed
      </p>
    </div>
  );
}
