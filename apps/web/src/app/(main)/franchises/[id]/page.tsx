"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import franchiseService from "@/services/franchiseService";
import { FranchiseWithItems } from "@/types/franchise";
import FranchiseTree from "@/components/franchise/FranchiseTree";

export default function FranchiseDetailPage() {
  const params = useParams();
  const id = Number(params.id);
  const [franchise, setFranchise] = useState<FranchiseWithItems | null>(null);

  useEffect(() => {
    if (id) franchiseService.getById(id).then(setFranchise).catch(() => {});
  }, [id]);

  if (!franchise) return <div className="p-8 text-muted">Loading...</div>;

  const pct = franchise.totalItems > 0 ? Math.round((franchise.completed / franchise.totalItems) * 100) : 0;

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      <h1 className="text-3xl font-bold text-foreground mb-2">{franchise.name}</h1>
      {franchise.description && <p className="text-muted mb-6">{franchise.description}</p>}

      <div className="bg-card rounded-xl border border-border p-4 mb-6">
        <div className="flex justify-between text-sm text-muted mb-2">
          <span>{franchise.completed}/{franchise.totalItems} completed</span>
          <span>{pct}%</span>
        </div>
        <div className="w-full h-2 bg-muted/20 rounded-full overflow-hidden">
          <div className="h-full bg-gradient-to-r from-primary to-accent rounded-full transition-all" style={{ width: `${pct}%` }} />
        </div>
      </div>

      <FranchiseTree items={franchise.items} />
    </div>
  );
}
