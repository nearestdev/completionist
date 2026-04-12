"use client";

import { useEffect, useState } from "react";
import api from "@/services/api";
import Button from "@/components/ui/Button";

interface Campaign {
  id: number;
  name: string;
  status: string;
  placement: string;
  createdAt: string;
}

export default function AdManager() {
  const [campaigns, setCampaigns] = useState<Campaign[]>([]);

  useEffect(() => {
    api.get<Campaign[]>("/admin/ads?limit=50").then((r) => setCampaigns(r.data || [])).catch(() => {});
  }, []);

  return (
    <div>
      {campaigns.length === 0 ? (
        <p className="text-muted">No ad campaigns yet.</p>
      ) : (
        <div className="space-y-3">
          {campaigns.map((c) => (
            <div key={c.id} className="bg-card rounded-xl border border-border p-4 flex items-center justify-between">
              <div>
                <h4 className="font-medium text-foreground">{c.name}</h4>
                <p className="text-xs text-muted">{c.placement} &middot; {c.status}</p>
              </div>
              <span className={`text-xs px-2 py-0.5 rounded-full ${
                c.status === "active" ? "bg-accent/10 text-accent" : "bg-muted/10 text-muted"
              }`}>{c.status}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
