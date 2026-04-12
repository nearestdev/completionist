"use client";

import { useEffect, useState } from "react";
import { Check, X } from "@phosphor-icons/react";
import api from "@/services/api";
import Button from "@/components/ui/Button";

interface ModerationItem {
  id: number;
  entityType: string;
  entityId: number;
  reason?: string;
  aiFlagged: boolean;
  status: string;
  createdAt: string;
}

export default function ModerationQueue() {
  const [items, setItems] = useState<ModerationItem[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchQueue = () => {
    api.get<ModerationItem[]>("/admin/moderation?status=pending&limit=50")
      .then((r) => setItems(r.data || []))
      .catch(() => {})
      .finally(() => setLoading(false));
  };

  useEffect(() => { fetchQueue(); }, []);

  const review = async (id: number, status: "approved" | "rejected") => {
    await api.put(`/admin/moderation/${id}`, { status });
    setItems(items.filter((i) => i.id !== id));
  };

  if (loading) return <p className="text-muted">Loading queue...</p>;
  if (items.length === 0) return <p className="text-muted">No pending items.</p>;

  return (
    <div className="space-y-3">
      {items.map((item) => (
        <div key={item.id} className="bg-card rounded-xl border border-border p-4 flex items-center justify-between">
          <div>
            <span className="text-sm font-medium text-foreground capitalize">{item.entityType} #{item.entityId}</span>
            {item.reason && <p className="text-xs text-muted mt-1">{item.reason}</p>}
            {item.aiFlagged && <span className="text-xs text-red-500 font-semibold ml-2">AI Flagged</span>}
          </div>
          <div className="flex gap-2">
            <Button size="sm" variant="outline" onClick={() => review(item.id, "approved")}>
              <Check size={14} className="mr-1" /> Approve
            </Button>
            <Button size="sm" variant="ghost" onClick={() => review(item.id, "rejected")}>
              <X size={14} className="mr-1" /> Reject
            </Button>
          </div>
        </div>
      ))}
    </div>
  );
}
